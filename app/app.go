package app

import (
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/commercionetwork/commercionetwork/x/ante"
	"github.com/commercionetwork/commercionetwork/x/commerciomint"
	commerciomintkeeper "github.com/commercionetwork/commercionetwork/x/commerciomint/keeper"
	commerciominttypes "github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	"github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/creditrisk"
	creditriskTypes "github.com/commercionetwork/commercionetwork/x/creditrisk/types"
	"github.com/commercionetwork/commercionetwork/x/docs"
	custombank "github.com/commercionetwork/commercionetwork/x/encapsulated/bank"
	customcrisis "github.com/commercionetwork/commercionetwork/x/encapsulated/crisis"
	customstaking "github.com/commercionetwork/commercionetwork/x/encapsulated/staking"
	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/id"
	idkeeper "github.com/commercionetwork/commercionetwork/x/id/keeper"
	idtypes "github.com/commercionetwork/commercionetwork/x/id/types"

	"github.com/commercionetwork/commercionetwork/x/memberships"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	"github.com/commercionetwork/commercionetwork/x/vbr"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"

	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
	tmos "github.com/tendermint/tendermint/libs/os"
	dbm "github.com/tendermint/tm-db"
)

const (
	appName = "Commercio.network"

	DefaultBondDenom   = "ucommercio"
	StableCreditsDenom = "uccc"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32MainPrefix = "did:com:"

	// PrefixValidator is the prefix for validator keys
	PrefixValidator = "val"
	// PrefixConsensus is the prefix for consensus keys
	PrefixConsensus = "cons"
	// PrefixPublic is the prefix for public keys
	PrefixPublic = "pub"
	// PrefixOperator is the prefix for operator keys
	PrefixOperator = "oper"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32MainPrefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32MainPrefix + PrefixPublic
	// Bech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	Bech32PrefixValAddr = Bech32MainPrefix + PrefixValidator + PrefixOperator
	// Bech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	Bech32PrefixValPub = Bech32MainPrefix + PrefixValidator + PrefixOperator + PrefixPublic
	// Bech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	Bech32PrefixConsAddr = Bech32MainPrefix + PrefixValidator + PrefixConsensus
	// Bech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	Bech32PrefixConsPub = Bech32MainPrefix + PrefixValidator + PrefixConsensus + PrefixPublic
)

// default home directories for expected binaries
var (
	DefaultCLIHome  = os.ExpandEnv("$HOME/.cncli")
	DefaultNodeHome = os.ExpandEnv("$HOME/.cnd")

	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		staking.AppModuleBasic{},
		distr.AppModuleBasic{},
		//gov.NewAppModuleBasic(paramsclient.ProposalHandler, distr.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},

		// Encapsulated modules
		customcrisis.NewAppModuleBasic(DefaultBondDenom),
		customstaking.NewAppModuleBasic(DefaultBondDenom),
		custombank.NewAppModuleBasic(bank.AppModuleBasic{}),

		// Custom modules
		docs.AppModuleBasic{},
		government.AppModuleBasic{},
		id.AppModuleBasic{},
		memberships.NewAppModuleBasic(StableCreditsDenom),
		commerciomint.NewAppModuleBasic(),
		pricefeed.AppModuleBasic{},
		vbr.AppModuleBasic{},
		creditrisk.AppModuleBasic{},
	)

	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		//gov.ModuleName:            {supply.Burner},

		// Custom modules
		commerciominttypes.ModuleName: {supply.Minter, supply.Burner},
		memberships.ModuleName:        {supply.Burner},
		idtypes.ModuleName:            nil,
		vbr.ModuleName:                {supply.Minter},
		creditriskTypes.ModuleName:    nil,
	}

	allowedModuleReceivers = types.Strings{
		commerciominttypes.ModuleName,
		memberships.ModuleName,
		vbr.ModuleName,
		creditriskTypes.ModuleName,
	}
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

func SetAddressPrefixes() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)
	config.Seal()
}

// Extended ABCI application
type CommercioNetworkApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// sdk keys to access the substores
	keys  map[string]*sdk.KVStoreKey
	tkeys map[string]*sdk.TransientStoreKey

	// sdk keepers
	accountKeeper  auth.AccountKeeper
	bankKeeper     bank.Keeper
	supplyKeeper   supply.Keeper
	stakingKeeper  staking.Keeper
	slashingKeeper slashing.Keeper
	distrKeeper    distr.Keeper
	//govKeeper      gov.Keeper
	crisisKeeper crisis.Keeper
	paramsKeeper params.Keeper

	// Encapsulated modules
	customBankKeeper custombank.Keeper

	// Custom modules
	docsKeeper       docs.Keeper
	governmentKeeper government.Keeper
	idKeeper         idkeeper.Keeper
	membershipKeeper memberships.Keeper
	mintKeeper       commerciomintkeeper.Keeper
	priceFeedKeeper  pricefeed.Keeper
	vbrKeeper        vbr.Keeper
	creditriskKeeper creditrisk.Keeper

	mm *module.Manager
}

// NewCommercioNetworkApp returns a reference to an initialized CommercioNetworkApp.
func NewCommercioNetworkApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*bam.BaseApp)) *CommercioNetworkApp {

	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	keys := sdk.NewKVStoreKeys(
		// Basics
		bam.MainStoreKey, auth.StoreKey, staking.StoreKey,
		supply.StoreKey, distr.StoreKey, slashing.StoreKey,
		//gov.StoreKey,
		params.StoreKey,

		// Encapsulated modules
		custombank.StoreKey,

		// Custom modules
		docs.StoreKey,
		government.StoreKey,
		idtypes.StoreKey,
		memberships.StoreKey,
		commerciominttypes.StoreKey,
		pricefeed.StoreKey,
		vbr.StoreKey,
		creditriskTypes.StoreKey,
	)
	tkeys := sdk.NewTransientStoreKeys(staking.TStoreKey, params.TStoreKey)

	// Here you initialize your application with the store keys it requires
	var app = &CommercioNetworkApp{
		BaseApp:        bApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		tkeys:          tkeys,
	}

	// init params keeper and subspaces
	app.paramsKeeper = params.NewKeeper(app.cdc, keys[params.StoreKey], tkeys[params.TStoreKey])
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	//govSubspace := app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	crisisSubspace := app.paramsKeeper.Subspace(crisis.DefaultParamspace)

	// add keepers
	app.accountKeeper = auth.NewAccountKeeper(app.cdc, keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper, bankSubspace, app.BlacklistedModuleAccAddrs())
	app.supplyKeeper = supply.NewKeeper(app.cdc, keys[supply.StoreKey], app.accountKeeper, app.bankKeeper, maccPerms)
	stakingKeeper := staking.NewKeeper(
		app.cdc, keys[staking.StoreKey],
		app.supplyKeeper, stakingSubspace,
	)
	app.distrKeeper = distr.NewKeeper(app.cdc, keys[distr.StoreKey], distrSubspace, &stakingKeeper,
		app.supplyKeeper, auth.FeeCollectorName, app.ModuleAccountAddrs())
	app.slashingKeeper = slashing.NewKeeper(
		app.cdc, keys[slashing.StoreKey], &stakingKeeper, slashingSubspace,
	)
	app.crisisKeeper = crisis.NewKeeper(crisisSubspace, invCheckPeriod, app.supplyKeeper, auth.FeeCollectorName)

	// Encapsulated modules
	app.customBankKeeper = custombank.NewKeeper(app.cdc, app.keys[custombank.StoreKey], app.bankKeeper)

	// Custom modules
	app.governmentKeeper = government.NewKeeper(app.cdc, app.keys[government.StoreKey])
	app.membershipKeeper = memberships.NewKeeper(app.cdc, app.keys[memberships.StoreKey], app.supplyKeeper, app.governmentKeeper, app.accountKeeper)
	app.docsKeeper = docs.NewKeeper(app.keys[docs.StoreKey], app.governmentKeeper, app.cdc)
	app.idKeeper = idkeeper.NewKeeper(app.cdc, app.keys[idtypes.StoreKey], app.accountKeeper, app.supplyKeeper)
	app.priceFeedKeeper = pricefeed.NewKeeper(app.cdc, app.keys[pricefeed.StoreKey], app.governmentKeeper)
	app.vbrKeeper = vbr.NewKeeper(app.cdc, app.keys[vbr.StoreKey], app.distrKeeper, app.supplyKeeper)
	app.mintKeeper = commerciomintkeeper.NewKeeper(app.cdc, app.keys[commerciominttypes.StoreKey], app.supplyKeeper, app.priceFeedKeeper, app.governmentKeeper)
	app.creditriskKeeper = creditrisk.NewKeeper(cdc, app.keys[creditriskTypes.StoreKey], app.supplyKeeper)

	// register the proposal types
	// govRouter := gov.NewRouter()
	// govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler)
	// app.govKeeper = gov.NewKeeper(
	// 	app.cdc, keys[gov.StoreKey], govSubspace,
	// 	app.supplyKeeper, &stakingKeeper, govRouter,
	// )

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	// Create default modules to be used from customs during encapsulation
	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.accountKeeper, app.supplyKeeper, app.stakingKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.accountKeeper, app.stakingKeeper),
		customstaking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		//gov.NewAppModule(app.govKeeper, app.accountKeeper, app.supplyKeeper),
		crisis.NewAppModule(&app.crisisKeeper),

		// Encapsulating modules
		custombank.NewAppModule(bank.NewAppModule(app.bankKeeper, app.accountKeeper), app.customBankKeeper, app.governmentKeeper),

		// Custom modules
		docs.NewAppModule(app.docsKeeper),
		government.NewAppModule(app.governmentKeeper),
		id.NewAppModule(app.idKeeper, app.governmentKeeper, app.supplyKeeper),
		memberships.NewAppModule(app.membershipKeeper, app.supplyKeeper, app.governmentKeeper, app.accountKeeper),
		commerciomint.NewAppModule(app.mintKeeper, app.supplyKeeper),
		pricefeed.NewAppModule(app.priceFeedKeeper, app.governmentKeeper),
		vbr.NewAppModule(app.vbrKeeper, app.stakingKeeper),
		creditrisk.NewAppModule(app.creditriskKeeper),
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(
		distr.ModuleName, slashing.ModuleName,

		// Custom modules
		vbr.ModuleName,
	)

	app.mm.SetOrderEndBlockers(
		crisis.ModuleName,
		//gov.ModuleName,
		staking.ModuleName,

		// Custom modules
		pricefeed.ModuleName,
		memberships.ModuleName,
		commerciominttypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(
		distr.ModuleName, staking.ModuleName, auth.ModuleName, bank.ModuleName,
		slashing.ModuleName, supply.ModuleName,
		//gov.ModuleName,
		crisis.ModuleName, genutil.ModuleName,

		// Custom modules
		government.ModuleName,
		docs.ModuleName,
		idtypes.ModuleName,
		memberships.ModuleName,
		commerciominttypes.ModuleName,
		pricefeed.ModuleName,
		vbr.ModuleName,
		creditriskTypes.ModuleName,
	)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(
		ante.NewAnteHandler(
			app.accountKeeper, app.supplyKeeper, app.priceFeedKeeper, app.governmentKeeper,
			auth.DefaultSigVerificationGasConsumer, StableCreditsDenom,
		),
	)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keys[bam.MainStoreKey])
		if err != nil {
			tmos.Exit(err.Error())
		}
	}

	return app
}

// application updates every begin block
func (app *CommercioNetworkApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// application updates every end block
func (app *CommercioNetworkApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// application update at chain initialization
func (app *CommercioNetworkApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, genesisState)
}

// load a particular height
func (app *CommercioNetworkApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keys[bam.MainStoreKey])
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *CommercioNetworkApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[app.supplyKeeper.GetModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlacklistedModuleAccAddrs returns all the app's module account addresses that
// are black listed from received tokens from the users.
func (app *CommercioNetworkApp) BlacklistedModuleAccAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[app.supplyKeeper.GetModuleAddress(acc).String()] = allowedModuleReceivers.Contains(acc)
	}

	return modAccAddrs
}
