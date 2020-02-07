package keeper

import (
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/asset/internal/types"
	"github.com/barkisnet/barkis/x/params"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = types.ModuleName
)

// ParamTable for issuing new assets
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&types.Params{})
}

// nolint: errcheck
func (k Keeper) GetMaxDecimal(ctx sdk.Context) int8 {
	var maxDecimal int8
	k.paramSpace.Get(ctx, types.ParamKeyMaxDecimal, &maxDecimal)
	return maxDecimal
}

// nolint: errcheck
func (k Keeper) SetMaxDecimal(ctx sdk.Context, maxDecimal int8) {
	k.paramSpace.Set(ctx, types.ParamKeyMaxDecimal, &maxDecimal)
}

// nolint: errcheck
func (k Keeper) GetIssueFee(ctx sdk.Context) sdk.Coins {
	var issueFee sdk.Coins
	k.paramSpace.Get(ctx, types.ParamKeyIssueFee, &issueFee)
	return issueFee
}

// nolint: errcheck
func (k Keeper) SetIssueFee(ctx sdk.Context, issueFee sdk.Coins) {
	k.paramSpace.Set(ctx, types.ParamKeyIssueFee, &issueFee)
}

// nolint: errcheck
func (k Keeper) GetMintFee(ctx sdk.Context) sdk.Coins {
	var mintFee sdk.Coins
	k.paramSpace.Get(ctx, types.ParamKeyMintFee, &mintFee)
	return mintFee
}

// nolint: errcheck
func (k Keeper) SetMintFee(ctx sdk.Context, mintFee sdk.Coins) {
	k.paramSpace.Set(ctx, types.ParamKeyMintFee, &mintFee)
}

// Get all parameteras as Params
func (k Keeper) GetParams(ctx sdk.Context) *types.Params {
	return types.NewParams(k.GetMaxDecimal(ctx), k.GetIssueFee(ctx), k.GetMintFee(ctx))
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params *types.Params) {
	k.paramSpace.SetParamSet(ctx, params)
}
