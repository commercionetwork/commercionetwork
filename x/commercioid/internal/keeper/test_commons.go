package keeper

import (
	"github.com/commercionetwork/commercionetwork/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

type testInput struct {
	Cdc        *codec.Codec
	Ctx        sdk.Context
	accKeeper  auth.AccountKeeper
	bankKeeper bank.BaseKeeper
	IdKeeper   Keeper
}

//This function create an enviroment to test modules
func setupTestInput() testInput {

	memDB := db.NewMemDB()
	cdc := testCodec()
	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	// CommercioID
	keyIDIdentities := sdk.NewKVStoreKey("id_identities")
	keyIDOwners := sdk.NewKVStoreKey("id_owners")
	keyIDConnections := sdk.NewKVStoreKey("id_connections")

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(keyIDIdentities, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyIDOwners, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyIDConnections, sdk.StoreTypeIAVL, memDB)

	ms.LoadLatestVersion()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, authKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	idk := NewKeeper(keyIDIdentities, keyIDOwners, keyIDConnections, cdc)

	ak.SetParams(ctx, auth.DefaultParams())

	return testInput{
		Cdc:        cdc,
		Ctx:        ctx,
		accKeeper:  ak,
		bankKeeper: bk,
		IdKeeper:   idk,
	}

}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterInterface((*auth.Account)(nil), nil)

	cdc.Seal()

	return cdc
}

var TestUtils = setupTestInput()

//TEST VARS
var TestAddress = "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"
var TestOwner, _ = sdk.AccAddressFromBech32(TestAddress)
var TestOwnerIdentity = types.Did("newReader")
var TestIdentityRef = "ddo-reference"
var TestReference = "TestReference"
var TestMetadata = "TestMetadata"
var TestRecipient = types.Did("recipient")
