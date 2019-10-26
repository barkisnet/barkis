package types

import (
	sdk "github.com/barkisnet/barkis/types"
	authexported "github.com/barkisnet/barkis/x/auth/exported"
)

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	NewAccount(sdk.Context, authexported.Account) authexported.Account
	SetAccount(sdk.Context, authexported.Account)
	IterateAccounts(ctx sdk.Context, process func(authexported.Account) (stop bool))
}
