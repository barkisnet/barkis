package types

import (
	"fmt"
	"regexp"

	sdk "github.com/barkisnet/barkis/types"
)

const (
	IssueMsgType = "issueMsg"
	MintMsgType  = "mintMsg"

	maxTokenNameLength         = 32
	maxTokenSymbolLength       = 10
	maxDecimals                = 10
	maxTotalSupply       int64 = 9000000000000000000 // int64 max value: 9,223,372,036,854,775,807
)

var (
	isAlpha = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
)

var _ sdk.Msg = IssueMsg{}

type IssueMsg struct {
	From        sdk.AccAddress `json:"from"`
	Name        string         `json:"name"`
	Symbol      string         `json:"symbol"`
	TotalSupply int64          `json:"total_supply"`
	Mintable    bool           `json:"mintable"`
	Decimals    int64          `json:"decimals"`
}

func NewIssueMsg(from sdk.AccAddress, name, symbol string, supply int64, mintable bool) IssueMsg {
	return IssueMsg{
		From:        from,
		Name:        name,
		Symbol:      symbol,
		TotalSupply: supply,
		Mintable:    mintable,
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

	if msg.TotalSupply < 0 || msg.TotalSupply > maxTotalSupply {
		return ErrInvalidTotalSupply(DefaultCodespace, fmt.Sprintf("total supply should be in [0, %d]", maxTotalSupply))
	}

	if msg.Decimals < 0 || msg.Decimals > maxDecimals {
		return ErrInvalidDecimal(DefaultCodespace, fmt.Sprintf("total supply should be in [0, %d]", maxDecimals))
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

	if len(msg.Symbol) == 0 || len(msg.Symbol) > maxTokenSymbolLength {
		return ErrInvalidTokenSymbol(DefaultCodespace, fmt.Sprintf("token symbol length shoud be in (0, %d]", maxTokenSymbolLength))
	}
	if msg.Symbol == sdk.DefaultBondDenom {
		return ErrInvalidTokenSymbol(DefaultCodespace, fmt.Sprintf("token symbol should be identical to native token %s", sdk.DefaultBondDenom))
	}
	if !isAlpha(msg.Symbol) {
		return ErrInvalidTokenSymbol(DefaultCodespace, "token symbol should only contains alphabet")
	}

	if msg.Amount <= 0 || msg.Amount > maxTotalSupply {
		return ErrInvalidMintAmount(DefaultCodespace, fmt.Sprintf("mint amount should be in (0, %d]", maxTotalSupply))
	}
	return nil
}
