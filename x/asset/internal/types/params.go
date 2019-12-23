package types

import (
	"fmt"

	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/params"
)

var (
	ParamKeyMaxDecimal = []byte("paramMaxDecimal")
	ParamKeyIssueFee   = []byte("paramIssueFee")
	ParamKeyMintFee    = []byte("paramMintFee")
)

// issue new assets parameters
type Params struct {
	MaxDecimal int8      `json:"param_max_decimal"`
	IssueFee   sdk.Coins `json:"param_issue_fee"`
	MintFee    sdk.Coins `json:"param_mint_fee"`
}

func (params Params) String() string {
	return fmt.Sprintf(`Asset parameters:
  MaxDecimal:   %d
  IssueFee:     %s
  MintFee:      %s`, params.MaxDecimal, params.IssueFee.String(), params.MintFee.String())
}

func NewParams(decimal int8, issueFee, mintFee sdk.Coins) *Params {
	return &Params{
		MaxDecimal: decimal,
		IssueFee:   issueFee,
		MintFee:    mintFee,
	}
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{ParamKeyMaxDecimal, &p.MaxDecimal},
		{ParamKeyIssueFee, &p.IssueFee},
		{ParamKeyMintFee, &p.MintFee},
	}
}

// validate a set of params
func (p Params) Validate() error {
	if p.MaxDecimal < 0 {
		return fmt.Errorf("token decimal must not negative")
	}
	if !p.IssueFee.IsAllPositive() {
		return fmt.Errorf("issue fee must be positive")
	}
	if !p.MintFee.IsAllPositive() {
		return fmt.Errorf("mint fee must be positive")
	}
	return nil
}
