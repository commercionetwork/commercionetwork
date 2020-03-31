package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/types"
)

//This function create an environment to test modules
func SetupTestInput() (*codec.Codec, sdk.Context, government.Keeper, Keeper) {

	memDB := db.NewMemDB()
	cdc := testCodec()

	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	govKey := sdk.NewKVStoreKey("government")
	pricefeedKey := sdk.NewKVStoreKey("pricefeed")

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(govKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(pricefeedKey, sdk.StoreTypeIAVL, memDB)

	_ = ms.LoadLatestVersion()

	govkeeper := government.NewKeeper(cdc, govKey)

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	return cdc, ctx, govkeeper, NewKeeper(cdc, pricefeedKey, govkeeper)
}

func testCodec() *codec.Codec {
	var cdc = codec.New()
	government.RegisterCodec(cdc)

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	auth.RegisterCodec(cdc)
	cdc.Seal()

	return cdc
}

// Test variables
var TestPrice = types.Price{
	AssetName: "test",
	Value:     sdk.NewDec(10),
	Expiry:    sdk.NewInt(5000),
}

var testOracle, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var price = types.Price{AssetName: "test", Value: sdk.NewDec(10), Expiry: sdk.NewInt(5000)}
var testGovernment, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

var msgSetPrice = types.NewMsgSetPrice(price, testOracle)
var msgAddOracle = types.NewMsgAddOracle(testGovernment, testOracle)
