package cli

import (
	"github.com/spf13/cobra"

	"github.com/barkisnet/barkis/client"
	"github.com/barkisnet/barkis/codec"
	"github.com/barkisnet/barkis/x/asset/internal/types"
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

	)...)

	return distQueryCmd
}
