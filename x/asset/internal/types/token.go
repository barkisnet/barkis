package types

import (
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/barkisnet/barkis/types"
)

const (
	TokenJoiner = "_"
)

var (
	isAlpha  = regexp.MustCompile(`^[a-z]+$`).MatchString
	isTxHash = regexp.MustCompile(`^[a-f0-9]+$`).MatchString
)

type Token struct {
	Symbol      string         `json:"symbol"`
	Name        string         `json:"name"`
	Decimal     int8           `json:"decimals"`
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

func (token *Token) String() string {
	return fmt.Sprintf(`Token:
  name:          %s
  symbol:      %s
  Decimal:      %d
  TotalSupply:    %d
  Mintable: %t
  Owner: %s
  Description:   %s`, token.Name, token.Symbol, token.Decimal,
		token.TotalSupply, token.Mintable, token.Owner.String(), token.Description)
}

type TokenList []*Token

func (tokenList TokenList) String() (out string) {
	for _, token := range tokenList {
		out += token.String() + "\n"
	}
	return strings.TrimSpace(out)
}

func ValidateToken(token *Token) error {
	if len(token.Owner) != sdk.AddrLen {
		return fmt.Errorf("sender address length should be %d", sdk.AddrLen)
	}

	if token.Name == sdk.DefaultBondDenom {
		return fmt.Errorf("token name should not be identical to native token name %s", sdk.DefaultBondDenom)
	}
	if len(token.Name) > MaxTokenNameLength {
		return fmt.Errorf("token name length should be less than %d", MaxTokenNameLength)
	}

	if len(token.Description) > MaxTokenDescription {
		return fmt.Errorf("token description length should be less than %d", MaxTokenDescription)
	}
	if len(token.Description) > MaxTokenDescription {
		return fmt.Errorf("token description length should be less than %d", MaxTokenDescription)
	}

	if err := validateTokenSymbol(token.Symbol); err != nil {
		return err
	}

	if token.Decimal < 0 {
		return fmt.Errorf("token decimal %d is negative", token.Decimal)
	}

	if token.TotalSupply <= 0 || token.TotalSupply > MaxTotalSupply {
		return fmt.Errorf("mint amount should be in (0, %d]", MaxTotalSupply)
	}
	return nil
}

func validateOriginalTokenSymbol(symbol string) error {
	if len(symbol) == 0 || len(symbol) > MaxTokenSymbolLength {
		return fmt.Errorf("token symbol length shoud be in (0, %d]", MaxTokenSymbolLength)
	}
	if symbol == sdk.DefaultBondDenom {
		return fmt.Errorf("token symbol should be identical to native token %s", sdk.DefaultBondDenom)
	}
	if !isAlpha(symbol) {
		return fmt.Errorf("token symbol should only contains alphabet")
	}
	return nil
}

func validateTokenSymbol(symbol string) error {
	symbolPaths := strings.Split(symbol, TokenJoiner)
	if len(symbolPaths) != 2 {
		return fmt.Errorf("valid token symbol should be XXX-YYY")
	}

	originalSymbol := symbolPaths[0]
	if err := validateOriginalTokenSymbol(originalSymbol); err != nil {
		return err
	}

	symbolSuffix := symbolPaths[1]
	if len(symbolSuffix) != TokenSymbolSuffixLen {
		return fmt.Errorf("token symbol suffix length should be %d", TokenSymbolSuffixLen)
	}
	if !isTxHash(symbolSuffix) {
		return fmt.Errorf("token symbol suffix should be transaction hash")
	}
	return nil
}
