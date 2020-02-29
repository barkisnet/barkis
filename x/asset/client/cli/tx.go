package cli

import (
	"fmt"
	"math"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/barkisnet/barkis/client"
	"github.com/barkisnet/barkis/client/context"
	"github.com/barkisnet/barkis/codec"
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/asset/internal/types"
	"github.com/barkisnet/barkis/x/auth"
	"github.com/barkisnet/barkis/x/auth/client/utils"
)

const (
	flagSymbol       = "token-symbol"
	flagTotalSupply  = "total-supply"
	flagTokenName    = "token-name"
	flagTokenDesc    = "token-desc"
	flagTokenDecimal = "token-decimal"
	flagMintable     = "mintable"
	flagAmount       = "amount"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Transaction commands for the asset module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(client.PostCommands(
		IssueTokenCmd(cdc),
		MintTokenCmd(cdc),
	)...)
	return txCmd
}

// IssueTokenCmd will create a issue token tx and sign it with the given key.
func IssueTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "issue",
		Short: "Create and sign a issue token tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			issuerAddr := cliCtx.GetFromAddress()
			supply := viper.GetInt64(flagTotalSupply)
			decimalInt := viper.GetInt(flagTokenDecimal)
			if decimalInt > math.MaxInt8 {
				return fmt.Errorf("token decimal overflow int8")
			}
			decimal := int8(decimalInt)
			mintable := viper.GetBool(flagMintable)
			name := viper.GetString(flagTokenName)
			symbol := viper.GetString(flagSymbol)
			desc := viper.GetString(flagTokenDesc)

			msgs := []sdk.Msg{types.NewIssueMsg(issuerAddr, name, symbol, supply, mintable, decimal, desc)}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs)
		},
	}
	cmd.Flags().String(flagTokenName, "", "token name")
	cmd.Flags().String(flagSymbol, "", "token symbol")
	cmd.Flags().String(flagTokenDesc, "", "token description")
	cmd.Flags().Int8(flagTokenDecimal, 6, "token decimal")
	cmd.Flags().Int64(flagTotalSupply, 0, "total supply of the new token")
	cmd.Flags().Bool(flagMintable, false, "whether the token can be minted")
	return cmd
}

// MintTokenCmd will create a mint token tx and sign it with the given key.
func MintTokenCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mint",
		Short: "Create and sign a issue token tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			issuerAddr := cliCtx.GetFromAddress()
			symbol := viper.GetString(flagSymbol)
			amount := viper.GetInt64(flagAmount)

			msgs := []sdk.Msg{types.NewMintMsg(issuerAddr, symbol, amount)}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, msgs)
		},
	}
	cmd.Flags().String(flagSymbol, "", "token symbol")
	cmd.Flags().Int64(flagAmount, 0, "mint amount")
	return cmd
}

// DelayedTransferCmd will create a mint token tx and sign it with the given key.
func DelayedTransferCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delayed-transfer [from_key_or_address] [to_address] [amount] [delayed_time]",
		Short: "Create and sign a delayed transfer tx",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			to, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			coins, err := sdk.ParseCoins(args[2])
			if err != nil {
				return err
			}

			delayedTime, err := strconv.Atoi(args[3])
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewDelayedTransferMsg(cliCtx.GetFromAddress(), to, coins, int64(delayedTime))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	return cmd
}
