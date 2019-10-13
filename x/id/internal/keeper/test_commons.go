package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

//This function create an environment to test modules
func SetupTestInput() (cdc *codec.Codec, ctx sdk.Context, govK government.Keeper, bk bank.Keeper, keeper Keeper) {

	memDB := db.NewMemDB()
	cdc = testCodec()

	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	govKey := sdk.NewKVStoreKey("government")

	// CommercioID
	storeKey := sdk.NewKVStoreKey("id")

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(govKey, sdk.StoreTypeIAVL, memDB)

	_ = ms.LoadLatestVersion()

	ctx = sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	govK = government.NewKeeper(cdc, govKey)
	_ = govK.SetGovernmentAddress(ctx, TestGovernment)

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, authKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk = bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, map[string]bool{})

	idk := NewKeeper(cdc, storeKey, bk)

	return cdc, ctx, govK, bk, idk

}

func testCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	types.RegisterCodec(cdc)

	cdc.Seal()

	return cdc
}

// ---------------------
// --- Test variables
// ---------------------

// Identities
var TestGovernment, _ = sdk.AccAddressFromBech32("cosmos15dnqp80tmkkkqdqx9ryky82cdasydtr5t9pgyx")
var TestOwnerAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestDidDocument = types.DidDocument{
	Uri:         "https://test.example.com/did-document#1",
	ContentHash: "6a40d9907d256795096b57d4bea23c0560aa3fe8f8a66c8207623f774c09c3a6",
}

// Deposit requests
var TestDepositor, _ = sdk.AccAddressFromBech32("cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6")
var TestPairwiseDid, _ = sdk.AccAddressFromBech32("cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa")
var TestDidDepositRequest = types.DidDepositRequest{
	Recipient:     TestPairwiseDid,
	Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
	Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
	EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
	FromAddress:   TestDepositor,
}

// Power up requests
var TestDidPowerUpRequest = types.DidPowerUpRequest{
	Claimant:      TestDepositor,
	Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
	Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
	EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
}
