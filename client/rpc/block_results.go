package rpc

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/barkisnet/barkis/client/context"
	"github.com/barkisnet/barkis/client/flags"
	"github.com/barkisnet/barkis/codec"
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/types/rest"
)

//BlockCommand returns the verified block data for a given heights
func BlockResultsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block-results [height]",
		Short: "Get the block results at given height",
		Args:  cobra.MaximumNArgs(1),
		RunE:  printBlockResults,
	}
	cmd.Flags().StringP(flags.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")
	viper.BindPFlag(flags.FlagNode, cmd.Flags().Lookup(flags.FlagNode))
	cmd.Flags().Bool(flags.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")
	viper.BindPFlag(flags.FlagTrustNode, cmd.Flags().Lookup(flags.FlagTrustNode))
	cmd.Flags().Bool(flags.FlagIndentResponse, false, "Add indent to JSON response")
	return cmd
}

type responseDeliverTx struct {
	Code      uint32           `json:"code,omitempty"`
	Data      []byte           `json:"data,omitempty"`
	Log       string           `json:"log,omitempty"`
	Info      string           `json:"info,omitempty"`
	GasWanted int64            `json:"gas_wanted,omitempty"`
	GasUsed   int64            `json:"gas_used,omitempty"`
	Events    sdk.StringEvents `json:"events,omitempty"`
	Codespace string           `json:"codespace,omitempty"`
}

func toResponseDeliverTx(deliverTxs []*abci.ResponseDeliverTx) []*responseDeliverTx {
	var responseDeliverTxs []*responseDeliverTx
	for _, deliverTx := range deliverTxs {
		responseDeliverTxs = append(responseDeliverTxs, &responseDeliverTx{
			Code:      deliverTx.Code,
			Data:      deliverTx.Data,
			Log:       deliverTx.Log,
			Info:      deliverTx.Info,
			GasWanted: deliverTx.GasWanted,
			GasUsed:   deliverTx.GasUsed,
			Events:    sdk.StringifyEvents(deliverTx.Events),
			Codespace: deliverTx.Codespace,
		})
	}
	return responseDeliverTxs
}

type responseEndBlock struct {
	ValidatorUpdates      []abci.ValidatorUpdate `json:"validator_updates"`
	ConsensusParamUpdates *abci.ConsensusParams  `json:"consensus_param_updates"`
	Events                sdk.StringEvents       `json:"events,omitempty"`
}

func toResponseEndBlock(endBlock *abci.ResponseEndBlock) *responseEndBlock {
	return &responseEndBlock{
		ValidatorUpdates:      endBlock.ValidatorUpdates,
		ConsensusParamUpdates: endBlock.ConsensusParamUpdates,
		Events:                sdk.StringifyEvents(endBlock.Events),
	}
}

type responseBeginBlock struct {
	Events sdk.StringEvents `json:"events,omitempty"`
}

func toResponseBeginBlock(beginBlock *abci.ResponseBeginBlock) *responseBeginBlock {
	return &responseBeginBlock{
		Events: sdk.StringifyEvents(beginBlock.Events),
	}
}

type abciResponses struct {
	DeliverTx  []*responseDeliverTx `json:"deliver_tx"`
	EndBlock   *responseEndBlock    `json:"end_block"`
	BeginBlock *responseBeginBlock  `json:"begin_block"`
}

type resultBlockResults struct {
	Height  int64          `json:"height"`
	Results *abciResponses `json:"results"`
}

func getBlockResults(cliCtx context.CLIContext, height *int64) ([]byte, error) {
	// get the node
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	res, err := node.BlockResults(height)
	if err != nil {
		return nil, err
	}

	blockResults := &resultBlockResults{
		Height: res.Height,
		Results: &abciResponses{
			DeliverTx:  toResponseDeliverTx(res.Results.DeliverTx),
			EndBlock:   toResponseEndBlock(res.Results.EndBlock),
			BeginBlock: toResponseBeginBlock(res.Results.BeginBlock),
		},
	}

	if cliCtx.Indent {
		return codec.Cdc.MarshalJSONIndent(blockResults, "", "  ")
	}

	return codec.Cdc.MarshalJSON(blockResults)
}

// CMD

func printBlockResults(cmd *cobra.Command, args []string) error {
	var height *int64
	// optional height
	if len(args) > 0 {
		h, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		if h > 0 {
			tmp := int64(h)
			height = &tmp
		}
	}

	output, err := getBlockResults(context.NewCLIContext(), height)
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}

// REST

// REST handler to get a block results
func BlockResultsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		height, err := strconv.ParseInt(vars["height"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest,
				"couldn't parse block height. Assumed format is '/block/{height}'.")
			return
		}

		chainHeight, err := GetChainHeight(cliCtx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "failed to parse chain height")
			return
		}

		if height > chainHeight {
			rest.WriteErrorResponse(w, http.StatusNotFound, "requested block height is bigger then the chain length")
			return
		}

		output, err := getBlockResults(cliCtx, &height)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, output)
	}
}

// REST handler to get the latest block
func LatestBlockResultsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		output, err := getBlockResults(cliCtx, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, output)
	}
}
