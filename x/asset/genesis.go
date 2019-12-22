package asset

import (
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/asset/internal/types"
)

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	Assets []*types.Token `json:"assets"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		Assets: nil,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, token := range data.Assets {
		keeper.SetToken(ctx, token)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	iter := keeper.ListAllToken(ctx)
	defer iter.Close()

	var assets []*types.Token
	for ; iter.Valid(); iter.Next() {
		token := keeper.DecodeToToken(iter.Value())
		assets = append(assets, token)
	}

	return GenesisState{
		Assets: assets,
	}
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	for _, token := range data.Assets {
		err := types.ValidateToken(token)
		if err != nil {
			return err
		}
	}
	return nil
}
