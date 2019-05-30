package commerciodocs

import (
	"commercio-network/x/commercioid"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

type testInput struct {
	cdc        *codec.Codec
	ctx        sdk.Context
	accKeeper  auth.AccountKeeper
	bankKeeper bank.BaseKeeper
	docsKeeper Keeper
}

//commercioauth module initialisation
var input = setupTestInput()

//This function create an enviroment to test modules
func setupTestInput() testInput {

	db := dbm.NewMemDB()
	cdc := makeCodec()
	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	// CommercioID
	keyIDIdentities := sdk.NewKVStoreKey("id_identities")
	keyIDOwners := sdk.NewKVStoreKey("id_owners")
	keyIDConnections := sdk.NewKVStoreKey("id_connections")

	// CommercioDOCS
	keyDOCSOwners := sdk.NewKVStoreKey("docs_owners")
	keyDOCSMetadata := sdk.NewKVStoreKey("docs_metadata")
	keyDOCSSharing := sdk.NewKVStoreKey("docs_sharing")
	keyDOCSReaders := sdk.NewKVStoreKey("docs_readers")

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	ms.MountStoreWithDB(keyDOCSReaders, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDOCSOwners, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDOCSMetadata, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyDOCSSharing, sdk.StoreTypeIAVL, db)

	ms.LoadLatestVersion()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams)
	ak := auth.NewAccountKeeper(cdc, authKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	idk := commercioid.NewKeeper(keyIDIdentities, keyIDOwners, keyIDConnections, cdc)
	dck := NewKeeper(idk, keyDOCSOwners, keyDOCSMetadata, keyDOCSSharing, keyDOCSReaders, cdc)

	ak.SetParams(ctx, auth.DefaultParams())

	return testInput{
		cdc:        cdc,
		ctx:        ctx,
		accKeeper:  ak,
		bankKeeper: bk,
		docsKeeper: dck,
	}

}

func makeCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterInterface((*auth.Account)(nil), nil)

	cdc.Seal()

	return cdc
}
