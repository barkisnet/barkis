package app

import (
	"io"
	"os"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/barkisnet/barkis/app/config"
	bam "github.com/barkisnet/barkis/baseapp"
	"github.com/barkisnet/barkis/codec"
	"github.com/barkisnet/barkis/simapp"
	sdk "github.com/barkisnet/barkis/types"
	"github.com/barkisnet/barkis/types/module"
	"github.com/barkisnet/barkis/version"
	"github.com/barkisnet/barkis/x/asset"
	"github.com/barkisnet/barkis/x/auth"
	"github.com/barkisnet/barkis/x/bank"
	"github.com/barkisnet/barkis/x/crisis"
	distr "github.com/barkisnet/barkis/x/distribution"
	"github.com/barkisnet/barkis/x/genaccounts"
	"github.com/barkisnet/barkis/x/genutil"
	"github.com/barkisnet/barkis/x/gov"
	"github.com/barkisnet/barkis/x/mint"
	"github.com/barkisnet/barkis/x/params"
	paramsclient "github.com/barkisnet/barkis/x/params/client"
	"github.com/barkisnet/barkis/x/slashing"
	"github.com/barkisnet/barkis/x/staking"
	"github.com/barkisnet/barkis/x/supply"
)

const appName = "BarkisApp"

var (
	// default home directories for barkiscli
	DefaultCLIHome = os.ExpandEnv("$HOME/.barkiscli")

	// default home directories for barkisd
	DefaultNodeHome = os.ExpandEnv("$HOME/.barkisd")

	// The module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		asset.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
		asset.ModuleName:          {supply.Minter},
	}

	BarkisContext = config.NewDefaultContext()
)

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()

	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	codec.RegisterEvidences(cdc)

	return cdc
}

// Extended ABCI application
type BarkisApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// keepers
	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	supplyKeeper   supply.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	mintKeeper     mint.Keeper
	distrKeeper    distr.Keeper
	govKeeper      gov.Keeper
	crisisKeeper   crisis.Keeper
	paramsKeeper   params.Keeper
	assetKeeper    asset.Keeper

	// the module manager
	mm *module.Manager
}

// NewBarkisApp returns a reference to an initialized BarkisApp.
func NewBarkisApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp)) *BarkisApp {

	cdc := MakeCodec()

	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
		supply.StoreKey, mint.StoreKey, distr.StoreKey, slashing.StoreKey,
		gov.StoreKey, params.StoreKey, asset.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	app := &BarkisApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tkeys:          tkeys,
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey], params.DefaultCodespace)
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := app.paramsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := app.paramsKeeper.Subspace(gov.DefaultParamspace)
	crisisSubspace := app.paramsKeeper.Subspace(crisis.DefaultParamspace)
	assetSubspace := app.paramsKeeper.Subspace(asset.DefaultParamspace)

	// add keepers
	app.accountKeeper = auth.NewAccountKeeper(app.cdc, keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper, bankSubspace, bank.DefaultCodespace, app.ModuleAccountAddrs())
	app.supplyKeeper = supply.NewKeeper(app.cdc, keys[supply.StoreKey], app.accountKeeper, app.bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(
		app.cdc, keys[staking.StoreKey], tkeys[staking.TStoreKey],
		app.supplyKeeper, stakingSubspace, staking.DefaultCodespace,
	)
	app.mintKeeper = mint.NewKeeper(app.cdc, keys[mint.StoreKey], mintSubspace, app.supplyKeeper, auth.FeeCollectorName)
	app.distrKeeper = distr.NewKeeper(app.cdc, keys[distr.StoreKey], distrSubspace, &stakingKeeper,
		app.supplyKeeper, distr.DefaultCodespace, auth.FeeCollectorName, app.ModuleAccountAddrs())
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc, keys[slashing.StoreKey], &stakingKeeper, slashingSubspace, slashing.DefaultCodespace,
	)
	app.crisisKeeper = crisis.NewKeeper(crisisSubspace, invCheckPeriod, app.supplyKeeper, auth.FeeCollectorName)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper))
	app.govKeeper = gov.NewKeeper(
		app.cdc, keys[gov.StoreKey], app.paramsKeeper, govSubspace,
		app.supplyKeeper, &stakingKeeper, gov.DefaultCodespace, govRouter,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	app.assetKeeper = asset.NewKeeper(cdc, keys[asset.StoreKey], assetSubspace, app.supplyKeeper, asset.DefaultCodespace)

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.
	app.mm = module.NewManager(
		genaccounts.NewAppModule(app.accountKeeper),
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		crisis.NewAppModule(&app.crisisKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.supplyKeeper),
		gov.NewAppModule(app.govKeeper, app.supplyKeeper),
		mint.NewAppModule(app.mintKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.distrKeeper, app.accountKeeper, app.supplyKeeper),
		asset.NewAppModule(app.assetKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName, slashing.ModuleName)

	app.mm.SetOrderEndBlockers(crisis.ModuleName, gov.ModuleName, staking.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		genaccounts.ModuleName, distr.ModuleName, staking.ModuleName,
		auth.ModuleName, bank.ModuleName, slashing.ModuleName, gov.ModuleName,
		mint.ModuleName, supply.ModuleName, crisis.ModuleName, genutil.ModuleName, asset.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.supplyKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	app.registerUpgrade()
	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			cmn.Exit(err.Error())
		}
	}

	return app
}

func (app *BarkisApp) registerUpgrade() {
	//Register upgrade height
	sdk.GlobalUpgradeMgr.RegisterUpgradeHeight(sdk.RewardUpgrade , BarkisContext.UpgradeConfig.RewardUpgrade)

	sdk.GlobalUpgradeMgr.RegisterBeginBlockerFirst(sdk.RewardUpgrade, func(ctx sdk.Context) {
		app.govKeeper.SetVotingParams(ctx, gov.NewVotingParams( 604800000000000)) // one week
		bonusProposerReward, err := sdk.NewDecFromStr("0.1838")
		if err != nil {
			panic(err)
		}
		mintSubspace , ok := app.paramsKeeper.GetSubspace(mint.DefaultParamspace)
		if ! ok {
			panic(fmt.Errorf("failed to get mint params subspace"))
		}
		mintSubspace.UpdateKeyTable(mint.UpdatedParamKeyTable())
		app.distrKeeper.SetBonusProposerReward(ctx, bonusProposerReward)
		app.mintKeeper.SetUnfreezeAmountPerBlock(ctx, 431000)
	})

	sdk.GlobalUpgradeMgr.RegisterUpgradeHeight(sdk.TokenIssueUpgrade, BarkisContext.UpgradeConfig.TokenIssueHeight)

	//Register new store if necessary
	sdk.GlobalUpgradeMgr.RegisterNewStore(sdk.TokenIssueUpgrade, asset.StoreKey)

	//Register new msg types if necessary
	sdk.GlobalUpgradeMgr.RegisterNewMsg(sdk.TokenIssueUpgrade, asset.IssueMsg{}.Type(), asset.MintMsg{}.Type())

	//Register BeginBlocker first for upgrade
	sdk.GlobalUpgradeMgr.RegisterBeginBlockerFirst(sdk.TokenIssueUpgrade, func(ctx sdk.Context) {
		app.assetKeeper.SetParams(ctx, asset.NewParams(10))
	})
}

// application updates every begin block
func (app *BarkisApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	sdk.GlobalUpgradeMgr.SetBlockHeight(ctx.BlockHeight())

	sdk.GlobalUpgradeMgr.BeginBlockersFirst(ctx)
	response := app.mm.BeginBlock(ctx, req)
	sdk.GlobalUpgradeMgr.BeginBlockersLast(ctx)
	return response
}

// application updates every end block
func (app *BarkisApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	sdk.GlobalUpgradeMgr.EndBlockersFirst(ctx)
	response := app.mm.EndBlock(ctx, req)
	sdk.GlobalUpgradeMgr.EndBlockersLast(ctx)
	return response
}

// application update at chain initialization
func (app *BarkisApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)

	return app.mm.InitGenesis(ctx, genesisState)
}

// load a particular height
func (app *BarkisApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *BarkisApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[app.supplyKeeper.GetModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}
