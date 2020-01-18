package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/barkisnet/barkis/codec"
	"github.com/barkisnet/barkis/store"
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/x/asset/internal/types"
	"github.com/barkisnet/barkis/x/auth"
	"github.com/barkisnet/barkis/x/bank"
	distr "github.com/barkisnet/barkis/x/distribution"
	"github.com/barkisnet/barkis/x/gov"
	"github.com/barkisnet/barkis/x/mint"
	"github.com/barkisnet/barkis/x/params"
	"github.com/barkisnet/barkis/x/staking"
	"github.com/barkisnet/barkis/x/supply"
)

func SetupTestInput() (*codec.Codec, sdk.Context, Keeper, auth.AccountKeeper, bank.Keeper, supply.Keeper, params.Keeper) {
	db := dbm.NewMemDB()

	cdc := codec.New()
	codec.RegisterCrypto(cdc)

	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	params.RegisterCodec(cdc)


	authKey := sdk.NewKVStoreKey(auth.StoreKey)
	paramsKey := sdk.NewKVStoreKey(params.StoreKey)
	tParamsKey := sdk.NewTransientStoreKey(params.TStoreKey)
	supplyKey := sdk.NewKVStoreKey(supply.StoreKey)
	assetKey := sdk.NewKVStoreKey(types.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(paramsKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tParamsKey, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(supplyKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(assetKey, sdk.StoreTypeIAVL, db)

	_ = ms.LoadLatestVersion()

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[sdk.AccAddress([]byte("moduleAcc")).String()] = true

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())


	paramKeeper := params.NewKeeper(cdc, paramsKey, tParamsKey, params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(cdc, authKey, paramKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, paramKeeper.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, blacklistedAddrs)
	accountKeeper.SetParams(ctx, auth.DefaultParams())

	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		types.ModuleName:          {supply.Minter},
	}
	supplyKeeper := supply.NewKeeper(cdc, supplyKey, accountKeeper, bankKeeper, maccPerms)
	assetKeeper := NewKeeper(cdc, assetKey, paramKeeper.Subspace(DefaultParamspace), supplyKeeper, types.DefaultCodespace)
	assetKeeper.SetParams(ctx, types.DefaultParams())


	addr1 := sdk.AccAddress(crypto.AddressHash([]byte("addr1")))
	addr2 := sdk.AccAddress(crypto.AddressHash([]byte("addr2")))

	acc1 := accountKeeper.NewAccountWithAddress(ctx, addr1)
	acc2 := accountKeeper.NewAccountWithAddress(ctx, addr2)

	accountKeeper.SetAccount(ctx, acc1)
	accountKeeper.SetAccount(ctx, acc2)

	_ = bankKeeper.SetCoins(ctx, addr1, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000000000)))
	_ = bankKeeper.SetCoins(ctx, addr2, sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 10000000000)))

	supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.NewCoins(sdk.NewInt64Coin(sdk.DefaultBondDenom, 20000000000))))

	return cdc, ctx, assetKeeper, accountKeeper,  bankKeeper, supplyKeeper, paramKeeper
}
