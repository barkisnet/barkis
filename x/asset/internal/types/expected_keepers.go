package types

import (
	sdk "github.com/barkisnet/barkis/types"
	supplyexported "github.com/barkisnet/barkis/x/supply/exported"
)

// SupplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
	GetModuleAccountAndPermissions(ctx sdk.Context, moduleName string) (supplyexported.ModuleAccountI, []string)
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error
}
