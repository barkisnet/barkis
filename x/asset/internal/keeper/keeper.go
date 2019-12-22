package keeper

import (
	"github.com/barkisnet/barkis/codec"
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/asset/internal/types"
	"github.com/barkisnet/barkis/x/params"
)

// Keeper of the distribution store
type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	paramSpace   params.Subspace
	SupplyKeeper types.SupplyKeeper
	codespace    sdk.CodespaceType
}

// NewKeeper creates a new distribution Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, supplyKeeper types.SupplyKeeper, codespace sdk.CodespaceType) Keeper {

	return Keeper{
		storeKey:     key,
		cdc:          cdc,
		paramSpace:   paramSpace.WithKeyTable(ParamKeyTable()),
		SupplyKeeper: supplyKeeper,
		codespace:    codespace,
	}
}

func (k *Keeper) SetToken(ctx sdk.Context, token *types.Token) {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.BuildTokenKey(token.Symbol)
	store.Set(tokenKey, k.serializeToken(token))
}

func (k *Keeper) GetToken(ctx sdk.Context, symbol string) *types.Token {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.BuildTokenKey(symbol)
	bz := store.Get(tokenKey)
	if bz == nil {
		return nil
	}
	return k.decodeToToken(bz)
}

func (k *Keeper) IsTokenExist(ctx sdk.Context, symbol string) bool {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.BuildTokenKey(symbol)
	return store.Has(tokenKey)
}

func (k *Keeper) serializeToken(token *types.Token) []byte {
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(*token)
	if err != nil {
		panic(err)
	}
	return bz
}

func (k *Keeper) decodeToToken(bz []byte) *types.Token {
	var token types.Token
	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &token)
	if err != nil {
		panic(err)
	}
	return &token
}
