package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"
)

//This function create an environment to test modules
func getTestInput() (sdk.Context, *codec.Codec, auth.AccountKeeper, bank.Keeper, Keeper) {

	memDB := db.NewMemDB()
	cdc := testCodec()

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	authKey := sdk.NewKVStoreKey("authCapKey")
	ibcKey := sdk.NewKVStoreKey("ibcCapKey")
	fckCapKey := sdk.NewKVStoreKey("fckCapKey")
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)

	keyAccreditations := sdk.NewKVStoreKey("accreditations")

	ms := store.NewCommitMultiStore(memDB)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(ibcKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(authKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(fckCapKey, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, memDB)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, memDB)
	ms.MountStoreWithDB(keyAccreditations, sdk.StoreTypeIAVL, memDB)
	_ = ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)
	accountKeeper := auth.NewAccountKeeper(cdc, keyAcc, pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, pk.Subspace(bank.DefaultParamspace), bank.DefaultCodespace, map[string]bool{})
	accreditationKeeper := NewKeeper(cdc, keyAccreditations, bankKeeper)

	return ctx, cdc, accountKeeper, bankKeeper, accreditationKeeper
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	bank.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)

	types.RegisterCodec(cdc)

	cdc.Seal()

	return cdc
}

// Testing variables
var TestUser, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
var TestAccrediter, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")
var TestSigner, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
