package asset

import (
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/asset/internal/types"
)

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	Params *types.Params  `json:"params" yaml:"params"`
	Tokens []*types.Token `json:"tokens" yaml:"tokens"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		Params: types.DefaultParams(),
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
	keeper.SetParams(ctx, data.Params)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	iter := keeper.ListToken(ctx)
	defer iter.Close()

	var tokens []*types.Token
	for ; iter.Valid(); iter.Next() {
		token := keeper.DecodeToken(iter.Value())
		tokens = append(tokens, token)
	}

	return GenesisState{
		Params: keeper.GetParams(ctx),
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
	if err := data.Params.Validate(); err != nil {
		return err
	}
	return nil
}
