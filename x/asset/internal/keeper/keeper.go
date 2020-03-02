package keeper

import (
	"encoding/binary"
	"fmt"

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
	if store.Has(tokenKey) {
		panic(fmt.Errorf("duplicated token symbol"))
	}
	store.Set(tokenKey, k.EncodeToken(token))
}

func (k *Keeper) UpdateToken(ctx sdk.Context, token *types.Token) {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.BuildTokenKey(token.Symbol)
	if !store.Has(tokenKey) {
		panic(fmt.Errorf("non-exist token"))
	}
	store.Set(tokenKey, k.EncodeToken(token))
}

func (k *Keeper) GetToken(ctx sdk.Context, symbol string) *types.Token {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.BuildTokenKey(symbol)
	bz := store.Get(tokenKey)
	if bz == nil {
		return nil
	}
	return k.DecodeToken(bz)
}

func (k *Keeper) ListToken(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.TokenKeyPrefix)
}

func (k *Keeper) IsTokenExist(ctx sdk.Context, symbol string) bool {
	store := ctx.KVStore(k.storeKey)
	tokenKey := types.BuildTokenKey(symbol)
	return store.Has(tokenKey)
}

func (k *Keeper) EncodeToken(token *types.Token) []byte {
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(*token)
	if err != nil {
		panic(err)
	}
	return bz
}

func (k *Keeper) DecodeToken(bz []byte) *types.Token {
	var token types.Token
	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &token)
	if err != nil {
		panic(err)
	}
	return &token
}

func (k *Keeper) GetSequence(ctx sdk.Context) int64 {
	kvStore := ctx.KVStore(k.storeKey)
	bz := kvStore.Get(types.SequenceKey)
	if bz == nil {
		return 0
	}
	return int64(binary.BigEndian.Uint64(bz))
}

func (k *Keeper) IncSequence(ctx sdk.Context) {
	sequence := k.GetSequence(ctx)
	sequence++
	sequenceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(sequenceBytes, uint64(sequence))

	kvStore := ctx.KVStore(k.storeKey)
	kvStore.Set(types.SequenceKey, sequenceBytes)
}

func (k *Keeper) InsertDelayedTransfer(ctx sdk.Context, delayedTransfer *types.DelayedTransfer) {
	store := ctx.KVStore(k.storeKey)
	delayedTransferKey := types.BuildDelayedTransferKey(delayedTransfer.Sequence)
	if store.Has(delayedTransferKey) {
		panic(fmt.Errorf("duplicated delayedTransfer sequence"))
	}
	store.Set(delayedTransferKey, k.EncodeDelayedTransfer(delayedTransfer))

	sequenceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(sequenceBytes, uint64(delayedTransfer.Sequence))

	delayedTransferMatureTimeKey := types.BuildDelayedTransferMatureTimeKey(delayedTransfer.MaturedTime, delayedTransfer.Sequence)
	store.Set(delayedTransferMatureTimeKey, sequenceBytes)

	delayedTransferFromKey := types.BuildDelayedTransferFromKey(delayedTransfer.From, delayedTransfer.Sequence)
	store.Set(delayedTransferFromKey, sequenceBytes)

	delayedTransferToKey := types.BuildDelayedTransferToKey(delayedTransfer.To, delayedTransfer.Sequence)
	store.Set(delayedTransferToKey, sequenceBytes)

	k.IncSequence(ctx)
}

func (k *Keeper) DeleteDelayedTransfer(ctx sdk.Context, delayedTransfer *types.DelayedTransfer) {
	store := ctx.KVStore(k.storeKey)

	delayedTransferKey := types.BuildDelayedTransferKey(delayedTransfer.Sequence)
	store.Delete(delayedTransferKey)

	delayedTransferMatureTimeKey := types.BuildDelayedTransferMatureTimeKey(delayedTransfer.MaturedTime, delayedTransfer.Sequence)
	store.Delete(delayedTransferMatureTimeKey)

	delayedTransferFromKey := types.BuildDelayedTransferFromKey(delayedTransfer.From, delayedTransfer.Sequence)
	store.Delete(delayedTransferFromKey)

	delayedTransferToKey := types.BuildDelayedTransferToKey(delayedTransfer.To, delayedTransfer.Sequence)
	store.Delete(delayedTransferToKey)
}

func (k *Keeper) GetDelayedTransfer(ctx sdk.Context, sequence int64) *types.DelayedTransfer {
	store := ctx.KVStore(k.storeKey)
	delayedTransferKey := types.BuildDelayedTransferKey(sequence)
	bz := store.Get(delayedTransferKey)
	if bz == nil {
		return nil
	}
	return k.DecodeDelayedTransfer(bz)
}

func (k *Keeper) EncodeDelayedTransfer(delayedTransfer *types.DelayedTransfer) []byte {
	bz, err := k.cdc.MarshalBinaryLengthPrefixed(*delayedTransfer)
	if err != nil {
		panic(err)
	}
	return bz
}

func (k *Keeper) DecodeDelayedTransfer(bz []byte) *types.DelayedTransfer {
	var delayedTransfer types.DelayedTransfer
	err := k.cdc.UnmarshalBinaryLengthPrefixed(bz, &delayedTransfer)
	if err != nil {
		panic(err)
	}
	return &delayedTransfer
}

func (k *Keeper) ListDelayedTransfer(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.DelayedTransferPrefix)
}

func (k *Keeper) ListDelayedTransferMaturedTime(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.DelayedTransferMatureTimePrefix)
}

func (k *Keeper) ListDelayedTransferFrom(ctx sdk.Context, addr sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, append(types.DelayedTransferFromPrefix, []byte(addr)...))
}

func (k *Keeper) ListDelayedTransferTo(ctx sdk.Context, addr sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, append(types.DelayedTransferToPrefix, []byte(addr)...))
}
