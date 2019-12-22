package keeper

import (
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/asset/internal/types"
	"github.com/barkisnet/barkis/x/params"
)


var (
	ParamKeyMaxDecimal = []byte("paramMaxDecimal")
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = types.ModuleName
)

// issue new assets parameters
type Params struct {
	MaxDecimal int8 `json:"param_max_decimal"`
}

// ParamTable for issuing new assets
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{ParamKeyMaxDecimal, &p.MaxDecimal},
	}
}

// returns the current CommunityTax rate from the global param store
// nolint: errcheck
func (k Keeper) GetMaxDecimal(ctx sdk.Context) int8 {
	var maxDecimal int8
	k.paramSpace.Get(ctx, ParamKeyMaxDecimal, &maxDecimal)
	return maxDecimal
}

// nolint: errcheck
func (k Keeper) SetMaxDecimal(ctx sdk.Context, maxDecimal int8) {
	k.paramSpace.Set(ctx, ParamKeyMaxDecimal, &maxDecimal)
}