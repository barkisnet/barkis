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
	Description string         `json:"description"`
	Owner       sdk.AccAddress `json:"owner"`
}

func NewToken(symbol, name string, decimal int8, totalSupply int64,
	mintable bool, description string, owner sdk.AccAddress) *Token {
	return &Token{
		Symbol:      symbol,
		Name:        name,
		Decimal:     decimal,
		TotalSupply: totalSupply,
		Mintable:    mintable,
		Description: description,
		Owner:       owner,
	}
}
