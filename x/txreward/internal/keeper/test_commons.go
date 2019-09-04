package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/txreward"
	"github.com/commercionetwork/commercionetwork/x/txreward/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

type testInput struct {
	Cdc       *codec.Codec
	Ctx       sdk.Context
	TBRKeeper Keeper
}

func setupTestInput() testInput {
	memDB := db.NewMemDB()
	cdc := testCodec()
	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	//TBR
	storeKey := sdk.NewKVStoreKey(types.BlockRewardsPoolPrefix)

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, memDB)

	_ = ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	var ak auth.AccountKeeper
	var bk bank.BaseKeeper
	var dk distribution.Keeper
	var sk staking.Keeper

	ak.SetParams(ctx, auth.DefaultParams())
	bk.se
	tbrK := NewKeeper(storeKey)

}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	cdc.RegisterConcrete(txreward.MsgIncrementsBlockRewardsPool{}, "commercio/incrementBlockRewardsPool", nil)
	cdc.Seal()

	return cdc
}

var TestUtils = setupTestInput()

var TestFunder, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestAmount = sdk.Coin{
	Denom:  "ucommercio",
	Amount: sdk.NewInt(10),
}
