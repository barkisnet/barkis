//nolint
package types

import (
	sdk "github.com/barkisnet/barkis/types"
)

// Local code type
type CodeType = sdk.CodeType

const (
	// Default asset codespace
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeInvalidTokenName   CodeType = 101
	CodeInvalidTokenSymbol CodeType = 102
	CodeInvalidTotalSupply CodeType = 103
	CodeInvalidDecimal     CodeType = 104
	CodeInvalidMintAmount  CodeType = 105
)

func ErrNoInvalidTokenName(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTokenName, msg)
}

func ErrInvalidTokenSymbol(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTokenSymbol, msg)
}

func ErrInvalidTotalSupply(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTotalSupply, msg)
}

func ErrInvalidDecimal(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDecimal, msg)
}

func ErrInvalidMintAmount(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMintAmount, msg)
}
