package keeper

import (
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/asset/internal/types"
	"github.com/barkisnet/barkis/x/params"
)

var (
	ParamKeyMaxDecimal = []byte("paramMaxDecimal")
	ParamKeyIssueFee   = []byte("paramIssueFee")
	ParamKeyMintFee    = []byte("paramMintFee")
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = types.ModuleName
)

// issue new assets parameters
type Params struct {
	MaxDecimal int8      `json:"param_max_decimal"`
	IssueFee   sdk.Coins `json:"issue_fee"`
	MintFee    sdk.Coins `json:"mint_fee"`
}

func NewParams(decimal int8, issueFee, mintFee sdk.Coins) *Params {
	return &Params{
		MaxDecimal: decimal,
		IssueFee:   issueFee,
		MintFee:    mintFee,
	}
}

// ParamTable for issuing new assets
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{ParamKeyMaxDecimal, &p.MaxDecimal},
		{ParamKeyIssueFee, &p.IssueFee},
		{ParamKeyMintFee, &p.MintFee},
	}
}

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

// nolint: errcheck
func (k Keeper) GetIssueFee(ctx sdk.Context) sdk.Coins {
	var issueFee sdk.Coins
	k.paramSpace.Get(ctx, ParamKeyMaxDecimal, &issueFee)
	return issueFee
}

// nolint: errcheck
func (k Keeper) SetIssueFee(ctx sdk.Context, issueFee sdk.Coins) {
	k.paramSpace.Set(ctx, ParamKeyMaxDecimal, &issueFee)
}

// nolint: errcheck
func (k Keeper) GetMintFee(ctx sdk.Context) sdk.Coins {
	var mintFee sdk.Coins
	k.paramSpace.Get(ctx, ParamKeyMaxDecimal, &mintFee)
	return mintFee
}

// nolint: errcheck
func (k Keeper) SetMintFee(ctx sdk.Context, mintFee sdk.Coins) {
	k.paramSpace.Set(ctx, ParamKeyMaxDecimal, &mintFee)
}

// Get all parameteras as Params
func (k Keeper) GetParams(ctx sdk.Context) *Params {
	return NewParams(k.GetMaxDecimal(ctx), k.GetIssueFee(ctx), k.GetMintFee(ctx))
}

// set the params
func (k Keeper) SetParams(ctx sdk.Context, params *Params) {
	k.paramSpace.SetParamSet(ctx, params)
}
