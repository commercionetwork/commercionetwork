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
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

var TestUtils = setupTestInput()

type testInput struct {
	Cdc        *codec.Codec
	Ctx        sdk.Context
	accKeeper  auth.AccountKeeper
	bankKeeper bank.BaseKeeper
	DocsKeeper Keeper
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

	// CommercioDOCS
	keyDOCS := sdk.NewKVStoreKey("docs")

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(keyDOCS, sdk.StoreTypeIAVL, memDB)
	_ = ms.LoadLatestVersion()

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, authKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, map[string]bool{})

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	dck := NewKeeper(keyDOCS, cdc)

	ak.SetParams(ctx, auth.DefaultParams())

	return testInput{
		Cdc:        cdc,
		Ctx:        ctx,
		accKeeper:  ak,
		bankKeeper: bk,
		DocsKeeper: dck,
	}
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterInterface((*auth.Account)(nil), nil)

	cdc.Seal()

	return cdc
}

// Testing variables

var TestingSender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestingRecipient, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

var TestingDocument = types.Document{
	Sender:     TestingSender,
	Recipient:  TestingRecipient,
	ContentUri: "https://example.com/document",
	Metadata: types.DocumentMetadata{
		ContentUri: "",
		Schema: types.DocumentMetadataSchema{
			Uri:     "https://example.com/document/metadata/schema",
			Version: "1.0.0",
		},
		Proof: "73666c68676c7366676c7366676c6a6873666c6a6768",
	},
	Checksum: types.DocumentChecksum{
		Value:     "93dfcaf3d923ec47edb8580667473987",
		Algorithm: "md5",
	},
}
