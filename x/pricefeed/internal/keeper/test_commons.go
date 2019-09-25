package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

//This function create an environment to test modules
func SetupTestInput() (cdc *codec.Codec, ctx sdk.Context, keeper Keeper) {

	memDB := db.NewMemDB()
	cdc = testCodec()

	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	govKey := sdk.NewKVStoreKey("government")

	storeKey := sdk.NewKVStoreKey("pricefeed")

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(govKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, memDB)

	_ = ms.LoadLatestVersion()

	ctx = sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	govkeeper := government.NewKeeper(govKey, cdc)
	pfk := NewKeeper(storeKey, govkeeper, cdc)

	return cdc, ctx, pfk
}

func testCodec() *codec.Codec {
	var cdc = codec.New()
	government.RegisterCodec(cdc)

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterInterface((*auth.Account)(nil), nil)
	cdc.Seal()

	return cdc
}

// Test variables
var TestOracle1, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestOracle2, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
var TestGovernment, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")
var TestPriceInfo = types.CurrentPrice{
	AssetName: "test",
	AssetCode: "0000",
	Price:     sdk.NewInt(10),
	Expiry:    sdk.NewInt(5000),
}
var TestPriceInfo2 = types.CurrentPrice{
	AssetName: "test2",
	AssetCode: "0001",
	Price:     sdk.NewInt(8),
	Expiry:    sdk.NewInt(4000),
}

var TestPriceInfo3 = types.CurrentPrice{
	AssetName: TestPriceInfo.AssetName,
	AssetCode: TestPriceInfo.AssetCode,
	Price:     sdk.NewInt(20),
	Expiry:    sdk.NewInt(7000),
}

var TestPriceInfoE = types.CurrentPrice{
	AssetName: "test",
	AssetCode: "0000",
	Price:     sdk.NewInt(0),
	Expiry:    sdk.NewInt(-1),
}
var TestRawPriceE = types.RawPrice{
	Oracle:    TestOracle1,
	PriceInfo: TestPriceInfoE,
}

var TestRawPrice1 = types.RawPrice{
	PriceInfo: TestPriceInfo,
	Oracle:    TestOracle1,
}

var TestRawPrice2 = types.RawPrice{
	PriceInfo: TestPriceInfo2,
	Oracle:    TestOracle2,
}

var TestRawPrice3 = types.RawPrice{
	Oracle:    TestOracle2,
	PriceInfo: TestPriceInfo3,
}

var TestAsset = types.Asset{
	Name: "ucommercio",
	Code: "0001",
}

var TestAsset2 = types.Asset{
	Name: "ucommerciocredits",
	Code: "0002",
}
