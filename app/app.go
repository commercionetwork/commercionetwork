package app

import (
	"commercio-network/x/commercioid/internal/keeper"
	"commercio-network/x/commercioid/internal/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/tendermint/tendermint/abci/version"
	"io"
	"os"

	"commercio-network/x/commercioauth"
	"commercio-network/x/commerciodocs"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"

	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"

	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
)

const (
	appName = "Commercio.network"
	Version = "1.0.1"

	DefaultBondDenom = "ucommercio"

	// DefaultKeyPass contains the default key password for genesis transactions
	DefaultKeyPass = "12345678"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32MainPrefix = "comnet"

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

	ModuleBasics module.BasicManager

	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
	}
)

func init() {
	ModuleBasics = module.NewBasicManager(
		genaccounts.AppModuleBasic{},
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(paramsclient.ProposalHandler, distrclient.ProposalHandler),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		supply.AppModuleBasic{},
		//todo Add this
		//commercioId.AppModuleBasic{}
		//commercioDocs.AppModuleBasic{}
	)
}

// custom tx codec
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	//TODO probably our modules should declared in init

	// CommercioAUTH
	commercioauth.RegisterCodec(cdc)

	// CommercioID
	types.RegisterCodec(cdc)

	// CommercioDOCS
	commerciodocs.RegisterCodec(cdc)

	return cdc
}

// Extended ABCI application
type commercioNetworkApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	// sdk keys to access the substores
	keyMain     *sdk.KVStoreKey
	keyAccount  *sdk.KVStoreKey
	keySupply   *sdk.KVStoreKey
	keyStaking  *sdk.KVStoreKey
	tkeyStaking *sdk.TransientStoreKey
	keySlashing *sdk.KVStoreKey
	keyMint     *sdk.KVStoreKey
	keyDistr    *sdk.KVStoreKey
	tkeyDistr   *sdk.TransientStoreKey
	keyGov      *sdk.KVStoreKey
	keyParams   *sdk.KVStoreKey
	tkeyParams  *sdk.TransientStoreKey

	// commercio-network keys to access the substores
	//CommercioID
	keyIDIdentities  *sdk.KVStoreKey
	keyIDOwners      *sdk.KVStoreKey
	keyIDConnections *sdk.KVStoreKey
	//CommercioDOCS
	keyDOCSOwners   *sdk.KVStoreKey
	keyDOCSMetadata *sdk.KVStoreKey
	keyDOCSSharing  *sdk.KVStoreKey
	keyDOCSReaders  *sdk.KVStoreKey

	// sdk keepers
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

	// commercio-network keepers
	// CommercioAUTH
	commercioAuthKeeper commercioauth.Keeper
	// CommercioID
	commercioIdKeeper keeper.Keeper
	// CommercioDOCS
	commercioDocsKeeper commerciodocs.Keeper

	mm *module.Manager
}

// NewCommercioNetworkApp returns a reference to an initialized CommercioNetworkApp.
func NewCommercioNetworkApp(logger log.Logger, db dbm.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint,
	baseAppOptions ...func(*bam.BaseApp)) *commercioNetworkApp {

	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)

	// Here you initialize your application with the store keys it requires
	var app = &commercioNetworkApp{
		BaseApp:     bApp,
		cdc:         cdc,
		keyMain:     sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:  sdk.NewKVStoreKey(auth.StoreKey),
		keySupply:   sdk.NewKVStoreKey(supply.StoreKey),
		keyStaking:  sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking: sdk.NewTransientStoreKey(staking.TStoreKey),
		keyMint:     sdk.NewKVStoreKey(mint.StoreKey),
		keyDistr:    sdk.NewKVStoreKey(distr.StoreKey),
		tkeyDistr:   sdk.NewTransientStoreKey(distr.TStoreKey),
		keySlashing: sdk.NewKVStoreKey(slashing.StoreKey),
		keyGov:      sdk.NewKVStoreKey(gov.StoreKey),
		keyParams:   sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams:  sdk.NewTransientStoreKey(params.TStoreKey),

		// CommercioID
		keyIDIdentities:  sdk.NewKVStoreKey("id_identities"),
		keyIDOwners:      sdk.NewKVStoreKey("id_owners"),
		keyIDConnections: sdk.NewKVStoreKey("id_connections"),

		// CommercioDOCS
		keyDOCSOwners:   sdk.NewKVStoreKey("docs_owners"),
		keyDOCSMetadata: sdk.NewKVStoreKey("docs_metadata"),
		keyDOCSSharing:  sdk.NewKVStoreKey("docs_sharing"),
		keyDOCSReaders:  sdk.NewKVStoreKey("docs_readers"),
	}

	// init params keeper and subspaces

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams, params.DefaultCodespace)

	//subspaces
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := app.paramsKeeper.Subspace(mint.DefaultParamspace)
	distrSubspace := app.paramsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := app.paramsKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := app.paramsKeeper.Subspace(gov.DefaultParamspace)
	crisisSubspace := app.paramsKeeper.Subspace(crisis.DefaultParamspace)

	// add keepers
	app.accountKeeper = auth.NewAccountKeeper(app.cdc, app.keyAccount, authSubspace, auth.ProtoBaseAccount)
	app.bankKeeper = bank.NewBaseKeeper(app.accountKeeper, bankSubspace, bank.DefaultCodespace)
	app.supplyKeeper = supply.NewKeeper(app.cdc, app.keySupply, app.accountKeeper, app.bankKeeper,
		supply.DefaultCodespace, maccPerms)

	stakingKeeper := staking.NewKeeper(app.cdc, app.keyStaking, app.tkeyStaking, app.supplyKeeper, stakingSubspace,
		staking.DefaultCodespace)

	app.mintKeeper = mint.NewKeeper(app.cdc, app.keyMint, mintSubspace, &stakingKeeper, app.supplyKeeper,
		auth.FeeCollectorName)

	app.distrKeeper = distr.NewKeeper(app.cdc, app.keyDistr, distrSubspace, &stakingKeeper, app.supplyKeeper,
		distr.DefaultCodespace, auth.FeeCollectorName)

	app.slashingKeeper = slashing.NewKeeper(app.cdc, app.keySlashing, &stakingKeeper, slashingSubspace,
		slashing.DefaultCodespace)

	app.crisisKeeper = crisis.NewKeeper(crisisSubspace, invCheckPeriod, app.supplyKeeper, auth.FeeCollectorName)

	// The CommercioAUTH keeper handles interactions for the CommercioAUTH module
	app.commercioAuthKeeper = commercioauth.NewKeeper(app.accountKeeper, app.cdc)

	// The CommercioID keeper handles interactions for the CommercioID module
	app.commercioIdKeeper = keeper.NewKeeper(app.keyIDIdentities, app.keyIDOwners, app.keyIDConnections, app.cdc)

	// The CommercioDOCS keeper handles interactions for the CommercioDOCS module
	app.commercioDocsKeeper = commerciodocs.NewKeeper(app.commercioIdKeeper, app.keyDOCSOwners, app.keyDOCSMetadata,
		app.keyDOCSSharing, app.keyDOCSReaders, app.cdc)

	// register the proposal types
	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, params.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(distr.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.distrKeeper))
	app.govKeeper = gov.NewKeeper(app.cdc, app.keyGov, app.paramsKeeper, govSubspace,
		app.supplyKeeper, &stakingKeeper, gov.DefaultCodespace, govRouter)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()))

	app.mm = module.NewManager(
		genaccounts.NewAppModule(app.accountKeeper),
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		crisis.NewAppModule(app.crisisKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		distr.NewAppModule(app.distrKeeper, app.supplyKeeper),
		gov.NewAppModule(app.govKeeper, app.supplyKeeper),
		mint.NewAppModule(app.mintKeeper),
		slashing.NewAppModule(app.slashingKeeper, app.stakingKeeper),
		staking.NewAppModule(app.stakingKeeper, app.distrKeeper, app.accountKeeper, app.supplyKeeper),
		//todo Finish this
		//commercioId.NewAppModule(app.commercioIdKeeper, ....
		//commercioDocs.NewAppModule(app.commercioDocsKeeper, ...
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	app.mm.SetOrderBeginBlockers(mint.ModuleName, distr.ModuleName, slashing.ModuleName)

	app.mm.SetOrderEndBlockers(gov.ModuleName, staking.ModuleName)

	// genutils must occur after staking so that pools are properly
	// initialized with tokens from genesis accounts.
	app.mm.SetOrderInitGenesis(genaccounts.ModuleName, distr.ModuleName,
		staking.ModuleName, auth.ModuleName, bank.ModuleName, slashing.ModuleName,
		gov.ModuleName, mint.ModuleName, supply.ModuleName, crisis.ModuleName, genutil.ModuleName)

	app.mm.RegisterInvariants(&app.crisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// initialize stores
	app.MountStores(app.keyMain, app.keyAccount, app.keySupply, app.keyStaking,
		app.keyMint, app.keyDistr, app.keySlashing, app.keyGov, app.keyParams,
		app.tkeyParams, app.tkeyStaking, app.tkeyDistr)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.supplyKeeper, auth.DefaultSigVerificationGasConsumer))
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		err := app.LoadLatestVersion(app.keyMain)
		if err != nil {
			cmn.Exit(err.Error())
		}
	}
	return app

	/*
		// The app.Router is the main transaction router where each module registers its routes
		app.Router().
			AddRoute(bank.RouterKey, bank.NewHandler(app.bankKeeper)).
			AddRoute(staking.RouterKey, staking.NewHandler(app.stakingKeeper)).
			AddRoute(distr.RouterKey, distr.NewHandler(app.distrKeeper)).
			AddRoute(slashing.RouterKey, slashing.NewHandler(app.slashingKeeper)).
			AddRoute(gov.RouterKey, gov.NewHandler(app.govKeeper)).
			AddRoute(crisis.RouterKey, crisis.NewHandler(app.crisisKeeper)).
			AddRoute("commercioauth", commercioauth.NewHandler(app.commercioAuthKeeper)).
			AddRoute("commercioid", commercioid.NewHandler(app.commercioIdKeeper)).
			AddRoute("commerciodocs", commerciodocs.NewHandler(app.commercioDocsKeeper))

		// The app.QueryRouter is the main query router where each module registers its routes
		app.QueryRouter().
			AddRoute(auth.QuerierRoute, auth.NewQuerier(app.accountKeeper)).
			AddRoute(distr.QuerierRoute, distr.NewQuerier(app.distrKeeper)).
			AddRoute(gov.QuerierRoute, gov.NewQuerier(app.govKeeper)).
			AddRoute(slashing.QuerierRoute, slashing.NewQuerier(app.slashingKeeper, app.cdc)).
			AddRoute(staking.QuerierRoute, staking.NewQuerier(app.stakingKeeper, app.cdc)).
			AddRoute("commercioauth", commercioauth.NewQuerier(app.commercioAuthKeeper)).
			AddRoute("commercioid", commercioid.NewQuerier(app.commercioIdKeeper)).
			AddRoute("commerciodocs", commerciodocs.NewQuerier(app.commercioDocsKeeper))

		// initialize BaseApp
		app.MountStores(app.keyMain, app.keyAccount, app.keyStaking, app.keyMint, app.keyDistr,
			app.keySlashing, app.keyGov, app.keyFeeCollection, app.keyParams,
			app.tkeyParams, app.tkeyStaking, app.tkeyDistr,

			// CommercioAUTH does not use any specific store as we base it on the auth module

			// CommercioID
			app.keyIDOwners,
			app.keyIDIdentities,
			app.keyIDConnections,

			// CommercioDOCS
			app.keyDOCSOwners,
			app.keyDOCSMetadata,
			app.keyDOCSSharing,
			app.keyDOCSReaders,
		)
		app.SetInitChainer(app.initChainer)
		app.SetBeginBlocker(app.BeginBlocker)
		// The AnteHandler handles signature verification and transaction pre-processing
		app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))
		app.SetEndBlocker(app.EndBlocker)

		if loadLatest {
			err := app.LoadLatestVersion(app.keyMain)
			if err != nil {
				cmn.Exit(err.Error())
			}
		}

		return app
	*/
}

// application updates every begin block
func (app *commercioNetworkApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

// application updates every end block
func (app *commercioNetworkApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// application update at chain initialization
func (app *commercioNetworkApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState simapp.GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.mm.InitGenesis(ctx, genesisState)
}

// load a particular height
func (app *commercioNetworkApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *commercioNetworkApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[app.supplyKeeper.GetModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// ______________________________________________________________________________________________
/*
var _ sdk.StakingHooks = StakingHooks{}

// StakingHooks contains combined distribution and slashing hooks needed for the
// staking module.
type StakingHooks struct {
	dh distr.Hooks
	sh slashing.Hooks
}

func NewStakingHooks(dh distr.Hooks, sh slashing.Hooks) StakingHooks {
	return StakingHooks{dh, sh}
}

// nolint
func (h StakingHooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.dh.AfterValidatorCreated(ctx, valAddr)
	h.sh.AfterValidatorCreated(ctx, valAddr)
}
func (h StakingHooks) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) {
	h.dh.BeforeValidatorModified(ctx, valAddr)
	h.sh.BeforeValidatorModified(ctx, valAddr)
}
func (h StakingHooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.AfterValidatorRemoved(ctx, consAddr, valAddr)
	h.sh.AfterValidatorRemoved(ctx, consAddr, valAddr)
}
func (h StakingHooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.AfterValidatorBonded(ctx, consAddr, valAddr)
	h.sh.AfterValidatorBonded(ctx, consAddr, valAddr)
}
func (h StakingHooks) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.dh.AfterValidatorBeginUnbonding(ctx, consAddr, valAddr)
	h.sh.AfterValidatorBeginUnbonding(ctx, consAddr, valAddr)
}
func (h StakingHooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.BeforeDelegationCreated(ctx, delAddr, valAddr)
	h.sh.BeforeDelegationCreated(ctx, delAddr, valAddr)
}
func (h StakingHooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
	h.sh.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
}
func (h StakingHooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.BeforeDelegationRemoved(ctx, delAddr, valAddr)
	h.sh.BeforeDelegationRemoved(ctx, delAddr, valAddr)
}
func (h StakingHooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) {
	h.dh.AfterDelegationModified(ctx, delAddr, valAddr)
	h.sh.AfterDelegationModified(ctx, delAddr, valAddr)
}
func (h StakingHooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) {
	h.dh.BeforeValidatorSlashed(ctx, valAddr, fraction)
	h.sh.BeforeValidatorSlashed(ctx, valAddr, fraction)
}
*/
