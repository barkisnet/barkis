package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/barkisnet/barkis/client"
	"github.com/barkisnet/barkis/client/context"
	"github.com/barkisnet/barkis/codec"
	"github.com/barkisnet/barkis/version"
	"github.com/barkisnet/barkis/x/asset/internal/types"
)

const (
	flagPage  = "page"
	flagLimit = "limit"
)

func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	distQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the asset module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	distQueryCmd.AddCommand(client.GetCommands(
		QueryParamsCmd(queryRoute, cdc),
		GetTokenCmd(queryRoute, cdc),
		ListTokenCmd(queryRoute, cdc),
		GetDelayedTransferCmd(queryRoute, cdc),
		ListDelayedTransferCmd(queryRoute, cdc),
		ListDelayedTransferFromCmd(queryRoute, cdc),
		ListDelayedTransferToCmd(queryRoute, cdc),
	)...)

	return distQueryCmd
}

// QueryParamsCmd implements the params query command.
func QueryParamsCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "Query the current asset parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as asset parameters.
Example:
$ %s query asset params
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParams)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(bz, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}

func GetTokenCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [symbol]",
		Short: "Get token information",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			symbol := args[0]

			resp, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.GetToken, symbol))
			if err != nil {
				return err
			}

			var token types.Token
			if err := cdc.UnmarshalJSON(resp, &token); err != nil {
				return err
			}

			return cliCtx.PrintOutput(&token)
		},
	}
	return cmd
}

func ListTokenCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List token",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			page := viper.GetInt(flagPage)
			limit := viper.GetInt(flagLimit)

			params := types.QueryTokensParams{
				Page:  page,
				Limit: limit,
			}

			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			resp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.ListToken), bz)
			if err != nil {
				return err
			}

			var tokenList types.TokenList
			if err := cdc.UnmarshalJSON(resp, &tokenList); err != nil {
				return err
			}

			return cliCtx.PrintOutput(tokenList)
		},
	}
	cmd.Flags().Int(flagPage, 1, "Query a specific page of paginated results")
	cmd.Flags().Int(flagLimit, 30, "Query number of transactions results per page returned")
	return cmd
}

func GetDelayedTransferCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-delayed-transfer [sequence]",
		Short: "Get delayed transfer by sequence",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			sequence := args[0]

			resp, _, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.GetDelayedTranfer, sequence))
			if err != nil {
				return err
			}

			var delayedTransfer types.DelayedTransfer
			if err := cdc.UnmarshalJSON(resp, &delayedTransfer); err != nil {
				return err
			}

			return cliCtx.PrintOutput(&delayedTransfer)
		},
	}

	return cmd
}

func ListDelayedTransferCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-delayed-transfer",
		Short: "List delayed transfer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			page := viper.GetInt(flagPage)
			limit := viper.GetInt(flagLimit)

			params := types.QueryDelayedTranferParams{
				Page:  page,
				Limit: limit,
			}

			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			resp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", queryRoute, types.ListDelayedTranfer), bz)
			if err != nil {
				return err
			}

			var delayedTransferList types.DelayedTransferList
			if err := cdc.UnmarshalJSON(resp, &delayedTransferList); err != nil {
				return err
			}

			return cliCtx.PrintOutput(delayedTransferList)
		},
	}
	cmd.Flags().Int(flagPage, 1, "Query a specific page of paginated results")
	cmd.Flags().Int(flagLimit, 30, "Query number of transactions results per page returned")
	return cmd
}

func ListDelayedTransferFromCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-delayed-transfer-from [from]",
		Short: "List delayed transfer from the specified address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			page := viper.GetInt(flagPage)
			limit := viper.GetInt(flagLimit)

			params := types.QueryDelayedTranferParams{
				Page:  page,
				Limit: limit,
			}

			fromAddr := args[0]

			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			resp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.ListDelayedTranferFrom, fromAddr), bz)
			if err != nil {
				return err
			}

			var delayedTransferList types.DelayedTransferList
			if err := cdc.UnmarshalJSON(resp, &delayedTransferList); err != nil {
				return err
			}

			return cliCtx.PrintOutput(delayedTransferList)
		},
	}
	cmd.Flags().Int(flagPage, 1, "Query a specific page of paginated results")
	cmd.Flags().Int(flagLimit, 30, "Query number of transactions results per page returned")
	return cmd
}

func ListDelayedTransferToCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-delayed-transfer-to [to]",
		Short: "List delayed transfer to the specified address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			page := viper.GetInt(flagPage)
			limit := viper.GetInt(flagLimit)

			params := types.QueryDelayedTranferParams{
				Page:  page,
				Limit: limit,
			}

			toAddr := args[0]

			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			resp, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.ListDelayedTranferTo, toAddr), bz)
			if err != nil {
				return err
			}

			var delayedTransferList types.DelayedTransferList
			if err := cdc.UnmarshalJSON(resp, &delayedTransferList); err != nil {
				return err
			}

			return cliCtx.PrintOutput(delayedTransferList)
		},
	}
	cmd.Flags().Int(flagPage, 1, "Query a specific page of paginated results")
	cmd.Flags().Int(flagLimit, 30, "Query number of transactions results per page returned")
	return cmd
}
