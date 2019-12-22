package asset

import (
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/asset/internal/types"
)

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	Tokens []*types.Token `json:"tokens"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		Tokens: nil,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	for _, token := range data.Tokens {
		keeper.SetToken(ctx, token)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	iter := keeper.ListToken(ctx)
	defer iter.Close()

	var tokens []*types.Token
	for ; iter.Valid(); iter.Next() {
		token := keeper.DecodeToToken(iter.Value())
		tokens = append(tokens, token)
	}

	return GenesisState{
		Tokens: tokens,
	}
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	for _, token := range data.Tokens {
		err := types.ValidateToken(token)
		if err != nil {
			return err
		}
	}
	return nil
}
