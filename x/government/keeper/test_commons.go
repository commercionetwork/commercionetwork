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
)

// This function create an environment to test modules
func SetupTestInput(setGovTumbAddresses bool) (cdc *codec.Codec, ctx sdk.Context, keeper Keeper) {

	memDB := db.NewMemDB()
	cdc = testCodec()
	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	// Store keys
	keyDocs := sdk.NewKVStoreKey("docs")
	keyGovernment := sdk.NewKVStoreKey("government")

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(keyDocs, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyGovernment, sdk.StoreTypeIAVL, memDB)
	_ = ms.LoadLatestVersion()

	ctx = sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	govk := NewKeeper(cdc, keyGovernment)

	if setGovTumbAddresses {
		_ = govk.SetGovernmentAddress(ctx, governmentTestAddress)
		_ = govk.SetTumblerAddress(ctx, tumblerTestAddress)
	}
	return cdc, ctx, govk
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	auth.RegisterCodec(cdc)
	cdc.Seal()

	return cdc
}

// Testing variables
var governmentTestAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var tumblerTestAddress, _ = sdk.AccAddressFromBech32("cosmos1h7tw92a66gr58pxgmf6cc336lgxadpjz5d5psf")
var notGovernmentAddress, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
