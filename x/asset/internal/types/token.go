package types

import (
	sdk "github.com/barkisnet/barkis/types"
)

type Token struct {
	Symbol      string         `json:"symbol"`
	Name        string         `json:"name"`
	Decimal     int8           `json:"decimal"`
	TotalSupply int64          `json:"total_supply"`
	Mintable    bool           `json:"mintable"`
	Owner       sdk.AccAddress `json:"owner"`
}
