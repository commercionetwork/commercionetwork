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

	keyMain *sdk.KVStoreKey

	keyParams  *sdk.KVStoreKey
	tkeyParams *sdk.TransientStoreKey

	keyAccount    *sdk.KVStoreKey
	accountKeeper auth.AccountKeeper

	feeCollectionKeeper auth.FeeCollectionKeeper
	keyFeeCollection    *sdk.KVStoreKey

	bankKeeper   bank.Keeper
	paramsKeeper params.Keeper

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

		keyMain: sdk.NewKVStoreKey("main"),

		keyParams:  sdk.NewKVStoreKey("params"),
		tkeyParams: sdk.NewTransientStoreKey("transient_params"),

		keyAccount:       sdk.NewKVStoreKey("acc"),
		keyFeeCollection: sdk.NewKVStoreKey("fee_collection"),

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
	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams)

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
	app.feeCollectionKeeper = auth.NewFeeCollectionKeeper(cdc, app.keyFeeCollection)

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
		AddRoute("commercioauth", commercioauth.NewHandler(app.commercioAuthKeeper)).
		AddRoute("commercioid", commercioid.NewHandler(app.commercioIdKeeper)).
		AddRoute("commerciodocs", commerciodocs.NewHandler(app.commercioDocsKeeper))

	// The app.QueryRouter is the main query router where each module registers its routes
	app.QueryRouter().
		AddRoute("commercioauth", commercioauth.NewQuerier(app.commercioAuthKeeper)).
		AddRoute("commercioid", commercioid.NewQuerier(app.commercioIdKeeper)).
		AddRoute("commerciodocs", commerciodocs.NewQuerier(app.commercioDocsKeeper))

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.initChainer)

	app.MountStores(
		app.keyMain,
		app.keyAccount,

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

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState struct {
	AuthData auth.GenesisState   `json:"auth"`
	BankData bank.GenesisState   `json:"bank"`
	Accounts []*auth.BaseAccount `json:"accounts"`
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

	return abci.ResponseInitChain{}
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

	genState := GenesisState{Accounts: accounts}
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	return appState, validators, err
}

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
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
