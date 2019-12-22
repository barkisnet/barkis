package types

import (
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/barkisnet/barkis/types"
)

const (
	IssueMsgType = "issueMsg"
	MintMsgType  = "mintMsg"

	TokenSymbolSuffixLen       = 3
	maxTokenNameLength         = 32
	maxTokenSymbolLength       = 10
	maxTokenDescription        = 128
	MaxTotalSupply       int64 = 9000000000000000000 // int64 max value: 9,223,372,036,854,775,807
)

var (
	isAlpha  = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
	isTxHash = regexp.MustCompile(`^[a-fA-F0-9]+$`).MatchString
)

var _ sdk.Msg = IssueMsg{}

type IssueMsg struct {
	From        sdk.AccAddress `json:"from"`
	Name        string         `json:"name"`
	Symbol      string         `json:"symbol"`
	TotalSupply int64          `json:"total_supply"`
	Mintable    bool           `json:"mintable"`
	Decimals    int8           `json:"decimals"`
	Description string         `json:"description"`
}

func NewIssueMsg(from sdk.AccAddress, name, symbol string, supply int64, mintable bool, decimals int8, description string) IssueMsg {
	return IssueMsg{
		From:        from,
		Name:        name,
		Symbol:      symbol,
		TotalSupply: supply,
		Mintable:    mintable,
		Decimals:    decimals,
		Description: description,
	}
}

func (msg IssueMsg) Route() string                { return RouterKey }
func (msg IssueMsg) Type() string                 { return IssueMsgType }
func (msg IssueMsg) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg IssueMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
func (msg IssueMsg) ValidateBasic() sdk.Error {
	if len(msg.From) != sdk.AddrLen {
		return sdk.ErrInvalidAddress(fmt.Sprintf("sender address length should be %d", sdk.AddrLen))
	}

	if len(msg.Name) == 0 || len(msg.Name) > maxTokenNameLength {
		return ErrNoInvalidTokenName(DefaultCodespace, fmt.Sprintf("token name length shoud be in (0, %d]", maxTokenNameLength))
	}
	if msg.Name == sdk.DefaultBondDenom {
		return ErrNoInvalidTokenName(DefaultCodespace, fmt.Sprintf("token name should be identical to native token %s", sdk.DefaultBondDenom))
	}

	if len(msg.Symbol) == 0 || len(msg.Symbol) > maxTokenSymbolLength {
		return ErrInvalidTokenSymbol(DefaultCodespace, fmt.Sprintf("token symbol length shoud be in (0, %d]", maxTokenSymbolLength))
	}
	if msg.Symbol == sdk.DefaultBondDenom {
		return ErrInvalidTokenSymbol(DefaultCodespace, fmt.Sprintf("token symbol should be identical to native token %s", sdk.DefaultBondDenom))
	}
	if !isAlpha(msg.Symbol) {
		return ErrInvalidTokenSymbol(DefaultCodespace, "token symbol should only contains alphabet")
	}

	if msg.TotalSupply < 0 || msg.TotalSupply > MaxTotalSupply {
		return ErrInvalidTotalSupply(DefaultCodespace, fmt.Sprintf("total supply should be in [0, %d]", MaxTotalSupply))
	}

	if msg.Decimals < 0 {
		return ErrInvalidDecimal(DefaultCodespace, fmt.Sprintf("token decimal %d is negative", msg.Decimals))
	}
	if len(msg.Description) > maxTokenDescription {
		return ErrInvalidTokenDescription(DefaultCodespace, fmt.Sprintf("token description length %d should be less than %d", len(msg.Description), maxTokenDescription))
	}

	return nil
}

type MintMsg struct {
	From   sdk.AccAddress `json:"from"`
	Symbol string         `json:"symbol"`
	Amount int64          `json:"amount"`
}

func NewMintMsg(from sdk.AccAddress, symbol string, amount int64) MintMsg {
	return MintMsg{
		From:   from,
		Symbol: symbol,
		Amount: amount,
	}
}

func (msg MintMsg) Route() string                { return RouterKey }
func (msg MintMsg) Type() string                 { return MintMsgType }
func (msg MintMsg) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.From} }
func (msg MintMsg) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
func (msg MintMsg) ValidateBasic() sdk.Error {
	if msg.From == nil {
		return sdk.ErrInvalidAddress("sender address cannot be empty")
	}
	symbolPaths := strings.Split(msg.Symbol, "-")
	if len(symbolPaths) != 2 {
		return ErrInvalidMintAmount(DefaultCodespace, fmt.Sprintf("valid token symbol should be XXX-YYY"))
	}

	originalSymbol := symbolPaths[0]
	if len(originalSymbol) == 0 || len(originalSymbol) > maxTokenSymbolLength {
		return ErrInvalidTokenSymbol(DefaultCodespace, fmt.Sprintf("token symbol length shoud be in (0, %d]", maxTokenSymbolLength))
	}
	if originalSymbol == sdk.DefaultBondDenom {
		return ErrInvalidTokenSymbol(DefaultCodespace, fmt.Sprintf("token symbol should be identical to native token %s", sdk.DefaultBondDenom))
	}
	if !isAlpha(originalSymbol) {
		return ErrInvalidTokenSymbol(DefaultCodespace, "token symbol should only contains alphabet")
	}

	symbolSuffix := symbolPaths[1]
	if len(symbolSuffix) != TokenSymbolSuffixLen {
		return ErrInvalidTokenSymbol(DefaultCodespace, fmt.Sprintf("token symbol suffix length should be %d", TokenSymbolSuffixLen))
	}
	if !isTxHash(symbolSuffix) {
		return ErrInvalidTokenSymbol(DefaultCodespace, "token symbol suffix should be transaction hash")
	}

	if msg.Amount <= 0 || msg.Amount > MaxTotalSupply {
		return ErrInvalidMintAmount(DefaultCodespace, fmt.Sprintf("mint amount should be in (0, %d]", MaxTotalSupply))
	}
	return nil
}
