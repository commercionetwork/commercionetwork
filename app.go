package app

import (
	"commercio-network/x/commercioauth"
	"commercio-network/x/commerciodocs"
	"commercio-network/x/commercioid"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/x/params"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"	
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/slashing"

	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	appName = "Commercio.network"
)

type commercioNetworkApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyStaking       *sdk.KVStoreKey
	tkeyStaking      *sdk.TransientStoreKey
	keySlashing      *sdk.KVStoreKey
	keyMint          *sdk.KVStoreKey
	keyDistr         *sdk.KVStoreKey
	tkeyDistr        *sdk.TransientStoreKey
	keyGov           *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyParams        *sdk.KVStoreKey
	tkeyParams       *sdk.TransientStoreKey


	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          bank.Keeper
	stakingKeeper       staking.Keeper
	slashingKeeper      slashing.Keeper
	mintKeeper          mint.Keeper
	distrKeeper         distr.Keeper
	govKeeper           gov.Keeper
	paramsKeeper        params.Keeper


	// CommercioAUTH
	commercioAuthKeeper commercioauth.Keeper

	// CommercioID
	commercioIdKeeper commercioid.Keeper
	keyIDIdentities   *sdk.KVStoreKey
	keyIDOwners       *sdk.KVStoreKey
	keyIDConnections  *sdk.KVStoreKey

	// CommercioDOCS
	commercioDocsKeeper commerciodocs.Keeper
	keyDOCSOwners       *sdk.KVStoreKey
	keyDOCSMetadata     *sdk.KVStoreKey
	keyDOCSSharing      *sdk.KVStoreKey
	keyDOCSReaders      *sdk.KVStoreKey
}

func NewCommercioNetworkApp(logger log.Logger, db dbm.DB) *commercioNetworkApp {

	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc))

	// Here you initialize your application with the store keys it requires
	var app = &commercioNetworkApp{
		BaseApp: bApp,
		cdc:     cdc,
		keyMain:          sdk.NewKVStoreKey("main"),
		keyAccount:       sdk.NewKVStoreKey("acc"),
		keyStaking:       sdk.NewKVStoreKey("staking"),
		tkeyStaking:      sdk.NewTransientStoreKey("transient_staking"),
		keyMint:          sdk.NewKVStoreKey("mint"),
		keyDistr:         sdk.NewKVStoreKey("distr"),
		tkeyDistr:        sdk.NewTransientStoreKey("transient_distr"),
		keySlashing:      sdk.NewKVStoreKey("slashing"),
		keyGov:           sdk.NewKVStoreKey("gov"),
		keyFeeCollection: sdk.NewKVStoreKey("fee_collection"),
		keyParams:        sdk.NewKVStoreKey("params"),
		tkeyParams:       sdk.NewTransientStoreKey("transient_params"),

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

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(
		app.cdc, 
		app.keyParams, 
		app.tkeyParams,
	)

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		app.paramsKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		app.paramsKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)

	// The FeeCollectionKeeper collects transaction fees and renders them to the fee distribution module
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(
		cdc, 
		app.keyFeeCollection,
	)

	// The StakingKeeper 
	stakingKeeper := staking.NewKeeper(
		app.cdc,
		app.keyStaking,
		app.tkeyStaking,
		app.bankKeeper, 
		app.paramsKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)

	app.mintKeeper = mint.NewKeeper(
		app.cdc,
		app.keyMint,
		app.paramsKeeper.Subspace(mint.DefaultParamspace),
		app.stakingKeeper, 
		app.feeCollectionKeeper,
	)

	app.distrKeeper = distr.NewKeeper(
		app.cdc,
		app.keyDistr,
		app.paramsKeeper.Subspace(distr.DefaultParamspace),
		app.bankKeeper, 
		app.stakingKeeper, 
		app.feeCollectionKeeper,
		distr.DefaultCodespace,
	)

	app.slashingKeeper = slashing.NewKeeper(
		app.cdc,
		app.keySlashing,
		app.stakingKeeper, 
		app.paramsKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)
	app.govKeeper = gov.NewKeeper(
		app.cdc,
		app.keyGov,
		app.paramsKeeper, 
		app.paramsKeeper.Subspace(gov.DefaultParamspace), 
		app.bankKeeper, 
		app.stakingKeeper,
		gov.DefaultCodespace,
	)

	app.stakingKeeper = *stakingKeeper.SetHooks(
		NewStakingHooks(app.distrKeeper.Hooks(), app.slashingKeeper.Hooks()),
	)

	// The CommercioAUTH keeper handles interactions for the CommercioAUTH module
	app.commercioAuthKeeper = commercioauth.NewKeeper(
		app.accountKeeper,
		app.cdc)

	// The CommercioID keeper handles interactions for the CommercioID module
	app.commercioIdKeeper = commercioid.NewKeeper(
		app.keyIDIdentities,
		app.keyIDOwners,
		app.keyIDConnections,
		app.cdc)

	// The CommercioDOCS keeper handles interactions for the CommercioDOCS module
	app.commercioDocsKeeper = commerciodocs.NewKeeper(
		app.commercioIdKeeper,
		app.keyDOCSOwners,
		app.keyDOCSMetadata,
		app.keyDOCSSharing,
		app.keyDOCSReaders,
		app.cdc)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))

	// The app.Router is the main transaction router where each module registers its routes
	app.Router().
		AddRoute("bank", bank.NewHandler(app.bankKeeper)).
		AddRoute("staking", staking.NewHandler(app.stakingKeeper)).
		AddRoute("distr", distr.NewHandler(app.distrKeeper)).
		AddRoute("slashing", slashing.NewHandler(app.slashingKeeper)).
		AddRoute("gov", gov.NewHandler(app.govKeeper)).
		AddRoute("commercioauth", commercioauth.NewHandler(app.commercioAuthKeeper)).
		AddRoute("commercioid", commercioid.NewHandler(app.commercioIdKeeper)).
		AddRoute("commerciodocs", commerciodocs.NewHandler(app.commercioDocsKeeper))

	// The app.QueryRouter is the main query router where each module registers its routes
	app.QueryRouter().
		AddRoute("distr", distr.NewQuerier(app.distrKeeper)).
		AddRoute("gov", gov.NewQuerier(app.govKeeper)).
		AddRoute("slashing", slashing.NewQuerier(app.slashingKeeper, app.cdc)).
		AddRoute("staking", staking.NewQuerier(app.stakingKeeper, app.cdc)).
		AddRoute("commercioauth", commercioauth.NewQuerier(app.commercioAuthKeeper)).
		AddRoute("commercioid", commercioid.NewQuerier(app.commercioIdKeeper)).
		AddRoute("commerciodocs", commerciodocs.NewQuerier(app.commercioDocsKeeper))

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.initChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(auth.NewAnteHandler(app.accountKeeper, app.feeCollectionKeeper))
	app.SetEndBlocker(app.EndBlocker)



	app.MountStores(
		app.keyMain,
		app.keyAccount,

		app.keyStaking, 
		app.keyMint, 
		app.keyDistr,
		app.keySlashing, 
		app.keyGov, 
		app.keyFeeCollection, 
		app.keyParams,
		app.tkeyParams, 
		app.tkeyStaking, 
		app.tkeyDistr,
	
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

	err := app.LoadLatestVersion(app.keyMain)
	if err != nil {
		cmn.Exit(err.Error())
	}

	return app
}


// application updates every end block
func (app *commercioNetworkApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	// mint new tokens for the previous block
	mint.BeginBlocker(ctx, app.mintKeeper)

	// distribute rewards for the previous block
	distr.BeginBlocker(ctx, req, app.distrKeeper)

	// slash anyone who double signed.
	// NOTE: This should happen after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool,
	// so as to keep the CanWithdrawInvariant invariant.
	// TODO: This should really happen at EndBlocker.
	tags := slashing.BeginBlocker(ctx, req, app.slashingKeeper)

	return abci.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
}

// application updates every end block
// nolint: unparam
func (app *commercioNetworkApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	tags := gov.EndBlocker(ctx, app.govKeeper)
	validatorUpdates, endBlockerTags := staking.EndBlocker(ctx, app.stakingKeeper)
	tags = append(tags, endBlockerTags...)

	//app.assertRuntimeInvariants()

	return abci.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Tags:             tags,
	}
}


// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState struct {
	AuthData auth.GenesisState   `json:"auth"`
	BankData bank.GenesisState   `json:"bank"`
	StakingData  staking.GenesisState  `json:"staking"`
	MintData     mint.GenesisState     `json:"mint"`
	DistrData    distr.GenesisState    `json:"distr"`
	GovData      gov.GenesisState      `json:"gov"`
	SlashingData slashing.GenesisState `json:"slashing"`
	Accounts []*auth.BaseAccount `json:"accounts"`
	GenTxs       []json.RawMessage     `json:"gentxs"`
}

func (app *commercioNetworkApp) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err)
	}

	for _, acc := range genesisState.Accounts {
		acc.AccountNumber = app.accountKeeper.GetNextAccountNumber(ctx)
		app.accountKeeper.SetAccount(ctx, acc)
	}
	// initialize distribution (must happen before staking)
	distr.InitGenesis(ctx, app.distrKeeper, genesisState.DistrData)
	// load the initial staking information
	validators, err := staking.InitGenesis(ctx, app.stakingKeeper, genesisState.StakingData)
	if err != nil {
		panic(err) // TODO find a way to do this w/o panics
	}
	// initialize module-specific stores
	auth.InitGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper, genesisState.AuthData)
	bank.InitGenesis(ctx, app.bankKeeper, genesisState.BankData)
	slashing.InitGenesis(ctx, app.slashingKeeper, genesisState.SlashingData, genesisState.StakingData.Validators.ToSDKValidators())
	gov.InitGenesis(ctx, app.govKeeper, genesisState.GovData)
	mint.InitGenesis(ctx, app.mintKeeper, genesisState.MintData)

	if len(genesisState.GenTxs) > 0 {
		for _, genTx := range genesisState.GenTxs {
			var tx auth.StdTx
			err = app.cdc.UnmarshalJSON(genTx, &tx)
			if err != nil {
				panic(err)
			}
			bz := app.cdc.MustMarshalBinaryLengthPrefixed(tx)
			res := app.BaseApp.DeliverTx(bz)
			if !res.IsOK() {
				panic(res.Log)
			}
		}

		validators = app.stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}
	return abci.ResponseInitChain{
		Validators: validators,		
	}
}

// ExportAppStateAndValidators does the things
func (app *commercioNetworkApp) ExportAppStateAndValidators() (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {
	ctx := app.NewContext(true, abci.Header{})
	var accounts []*auth.BaseAccount


	appendAccountsFn := func(acc auth.Account) bool {
		account := &auth.BaseAccount{
			Address: acc.GetAddress(),
			Coins:   acc.GetCoins(),
		}

		accounts = append(accounts, account)
		return false
	}

	app.accountKeeper.IterateAccounts(ctx, appendAccountsFn)

	genState := NewGenesisState(
		accounts,
		auth.ExportGenesis(ctx, app.accountKeeper, app.feeCollectionKeeper),
		bank.ExportGenesis(ctx, app.bankKeeper),
		staking.ExportGenesis(ctx, app.stakingKeeper),
		mint.ExportGenesis(ctx, app.mintKeeper),
		distr.ExportGenesis(ctx, app.distrKeeper),
		gov.ExportGenesis(ctx, app.govKeeper),
		slashing.ExportGenesis(ctx, app.slashingKeeper),
	)
	
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	return appState, validators, err
}

func NewGenesisState(accounts []*auth.BaseAccount, 
	authData auth.GenesisState,
	bankData bank.GenesisState,
	stakingData staking.GenesisState, 
	mintData mint.GenesisState,
	distrData distr.GenesisState, 
	govData gov.GenesisState,
	slashingData slashing.GenesisState) GenesisState {

	return GenesisState{
		Accounts:     accounts,
		AuthData:     authData,
		BankData:     bankData,
		StakingData:  stakingData,
		MintData:     mintData,
		DistrData:    distrData,
		GovData:      govData,
		SlashingData: slashingData,
	}
}



// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	distr.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	// CommercioAUTH
	commercioauth.RegisterCodec(cdc)

	// CommercioID
	commercioid.RegisterCodec(cdc)

	// CommercioDOCS
	commerciodocs.RegisterCodec(cdc)

	return cdc
}


//______________________________________________________________________________________________

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
