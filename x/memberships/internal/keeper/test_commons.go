package keeper

import (
	"time"

	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

//This function create an environment to test modules
func GetTestInput() (*codec.Codec, sdk.Context, bank.Keeper, government.Keeper, Keeper) {

	memDB := db.NewMemDB()
	cdc := testCodec()

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	nftKey := sdk.NewKVStoreKey("nft")
	governmentKey := sdk.NewKVStoreKey("government")
	storeKey := sdk.NewKVStoreKey("accreditations")

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(nftKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(governmentKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, memDB)
	_ = ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankK := bank.NewBaseKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, map[string]bool{})

	govK := government.NewKeeper(cdc, governmentKey)
	nftK := nft.NewKeeper(cdc, nftKey)
	accK := NewKeeper(cdc, storeKey, nftK, bankK)
	accK.SetStableCreditsDenom(ctx, TestStableCreditsDenom)

	return cdc, ctx, bankK, govK, accK
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	nft.RegisterCodec(cdc)

	types.RegisterCodec(cdc)

	cdc.Seal()

	return cdc
}

// Testing variables
var TestUser, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
var TestUser2, _ = sdk.AccAddressFromBech32("cosmos1h7tw92a66gr58pxgmf6cc336lgxadpjz5d5psf")

var TestTsp, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

var zone, _ = time.LoadLocation("UTC")
var TestTimestamp = time.Date(1990, 10, 10, 20, 20, 0, 0, zone)

var TestStableCreditsDenom = "uccc"
var TestMembershipType = "bronze"
