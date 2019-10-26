package rpc

import (
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/barkisnet/barkis/client/context"
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/types/rest"
)

func TxSimulateRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashHexStr := vars["txBytes"]

		txBytes, err := hex.DecodeString(hashHexStr)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// run a simulation (via /app/simulate query) to
		// estimate gas and update TxBuilder accordingly
		rawRes, _, err := cliCtx.QueryWithData("/app/simulate", txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var simulationResult sdk.Result
		if err := cliCtx.Codec.UnmarshalBinaryLengthPrefixed(rawRes, &simulationResult); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		txResponse := sdk.TxResponse{
			Code:      uint32(simulationResult.Code),
			Data:      strings.ToUpper(hex.EncodeToString(simulationResult.Data)),
			Codespace: string(simulationResult.Codespace),
			RawLog:    simulationResult.Log,
			GasWanted: int64(simulationResult.GasWanted),
			GasUsed:   int64(simulationResult.GasUsed),
			Events:    sdk.StringifyEvents(simulationResult.Events.ToABCIEvents()),
		}

		rest.PostProcessResponseBare(w, cliCtx, txResponse)
	}
}
