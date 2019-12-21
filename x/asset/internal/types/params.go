package types

import (
	"github.com/barkisnet/barkis/x/params"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = ModuleName
)

// Parameter store keys
var (
)

// issue new assets parameters
type Params struct {
}

// ParamTable for issuing new assets
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
	}
}