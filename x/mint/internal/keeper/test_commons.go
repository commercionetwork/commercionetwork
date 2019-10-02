package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/government"
	comMint "github.com/commercionetwork/commercionetwork/x/mint"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"

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

func SetupTestInput() (cdc *codec.Codec, ctx sdk.Context, bankKeeper bank.Keeper, pricefeedKeeper pricefeed.Keeper,
	keeper Keeper) {
	memDB := db.NewMemDB()
	cdc = testCodec()

	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	distrKey := sdk.NewKVStoreKey("distrKey")

	//custom modules keys
	govKey := sdk.NewKVStoreKey(government.StoreKey)
	pricefeedKey := sdk.NewKVStoreKey(pricefeed.StoreKey)
	cMintKey := sdk.NewKVStoreKey(comMint.StoreKey)

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(distrKey, sdk.StoreTypeIAVL, memDB)

	//custom modules loading
	ms.MountStoreWithDB(pricefeedKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(cMintKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(govKey, sdk.StoreTypeIAVL, memDB)

	_ = ms.LoadLatestVersion()

	ctx = sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	ak := auth.NewAccountKeeper(cdc, authKey, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, nil)

	govkeeper := government.NewKeeper(govKey, cdc)
	pricefeedK := pricefeed.NewKeeper(pricefeedKey, govkeeper, cdc)
	mintK := comMint.NewKeeper(cMintKey, bankKeeper, pricefeedK, cdc)

	return cdc, ctx, bk, pricefeedK, mintK
}

func testCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)

	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)

	//custom modules codecs
	pricefeed.RegisterCodec(cdc)
	government.RegisterCodec(cdc)

	codec.RegisterCrypto(cdc)

	cdc.Seal()

	return cdc
}
