package types

import (
	"fmt"
	sdk "github.com/barkisnet/barkis/types"
	"strings"
)

const (
	//todo refactor name
	IssueMsgType = "issueMsg"
	MintMsgType  = "mintMsg"

	MaxTokenNameLength           = 32
	MaxTokenSymbolLength         = 12
	MinTokenSymbolLength         = 3
	MaxTokenDesLenLimit          = 128
	NewMaxTokenDesLenLimit       = 1024
	MaxTotalSupply         int64 = 9000000000000000000 // int64 max value: 9,223,372,036,854,775,807
)

var _ sdk.Msg = IssueMsg{}

type IssueMsg struct {
	From        sdk.AccAddress `json:"from"`
	Name        string         `json:"name"`
	Symbol      string         `json:"symbol"`
	TotalSupply int64          `json:"total_supply"`
	Mintable    bool           `json:"mintable"`
	Decimal     int8           `json:"decimal"`
	Description string         `json:"description"`
}

func NewIssueMsg(from sdk.AccAddress, name, symbol string, supply int64, mintable bool, decimal int8, description string) IssueMsg {
	return IssueMsg{
		From:        from,
		Name:        name,
		Symbol:      symbol,
		TotalSupply: supply,
		Mintable:    mintable,
		Decimal:     decimal,
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

	if len(msg.Name) == 0 || len(msg.Name) > MaxTokenNameLength {
		return ErrNoInvalidTokenName(DefaultCodespace, fmt.Sprintf("token name length shoud be in (0, %d]", MaxTokenNameLength))
	}
	if msg.Name == sdk.DefaultBondDenom || msg.Name == sdk.DefaultBondDenomName {
		return ErrNoInvalidTokenName(DefaultCodespace, fmt.Sprintf("token name should be identical to native token %s/%s", sdk.DefaultBondDenom, sdk.DefaultBondDenomName))
	}

	if err := validateTokenSymbol(strings.ToLower(msg.Symbol)); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}

	if msg.TotalSupply < 0 || msg.TotalSupply > MaxTotalSupply {
		return ErrInvalidTotalSupply(DefaultCodespace, fmt.Sprintf("total supply should be in [0, %d]", MaxTotalSupply))
	}

	if msg.Decimal < 0 {
		return ErrInvalidDecimal(DefaultCodespace, fmt.Sprintf("token decimal %d is negative", msg.Decimal))
	}
	desLenLimitation := MaxTokenDesLenLimit
	if sdk.GlobalUpgradeMgr.IsUpgradeApplied(sdk.TokenDesLenLimitUpgradeHeight) {
		desLenLimitation = NewMaxTokenDesLenLimit
	}
	if len(msg.Description) > desLenLimitation {
		return ErrInvalidTokenDescription(DefaultCodespace, fmt.Sprintf("token description length %d should be less than %d", len(msg.Description), desLenLimitation))
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
	if len(msg.From) != sdk.AddrLen {
		return sdk.ErrInvalidAddress(fmt.Sprintf("sender address length should be %d", sdk.AddrLen))
	}

	if err := validateTokenSymbol(msg.Symbol); err != nil {
		return ErrInvalidTokenSymbol(DefaultCodespace, err.Error())
	}

	if msg.Amount <= 0 || msg.Amount > MaxTotalSupply {
		return ErrInvalidMintAmount(DefaultCodespace, fmt.Sprintf("mint amount should be in (0, %d]", MaxTotalSupply))
	}
	return nil
}
