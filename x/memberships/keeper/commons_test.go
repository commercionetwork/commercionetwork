package keeper_test

import (
	governmentKeeper "github.com/commercionetwork/commercionetwork/x/government/keeper"
	governmentTypes "github.com/commercionetwork/commercionetwork/x/government/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	creditrisk "github.com/commercionetwork/commercionetwork/x/creditrisk/types"
	"github.com/commercionetwork/commercionetwork/x/memberships/keeper"
	"github.com/commercionetwork/commercionetwork/x/memberships/types"
)

//This function create an environment to test modules
func SetupTestInput() (sdk.Context, bank.Keeper, governmentKeeper.Keeper, keeper.Keeper) {

	memDB := db.NewMemDB()
	cdc := testCodec()

	keys := sdk.NewKVStoreKeys(
		auth.StoreKey,
		params.StoreKey,
		supply.StoreKey,
		governmentTypes.StoreKey,

		types.StoreKey,
	)
	tKeys := sdk.NewTransientStoreKeys(params.TStoreKey)

	ms := store.NewCommitMultiStore(memDB)
	for _, key := range keys {
		ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, memDB)
	}
	for _, tkey := range tKeys {
		ms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, memDB)
	}
	_ = ms.LoadLatestVersion()

	ctx := sdk.NewContext(ms, abci.Header{ChainID: "test-chain-id"}, false, log.NewNopLogger())

	pk := params.NewKeeper(cdc, keys[params.StoreKey], tKeys[params.TStoreKey])
	ak := auth.NewAccountKeeper(cdc, keys[auth.StoreKey], pk.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bk := bank.NewBaseKeeper(ak, pk.Subspace(bank.DefaultParamspace), nil)
	maccPerms := map[string][]string{
		types.ModuleName:      {supply.Minter, supply.Burner},
		creditrisk.ModuleName: nil,
	}
	sk := supply.NewKeeper(cdc, keys[supply.StoreKey], ak, bk, maccPerms)
	sk.SetSupply(ctx, supply.NewSupply(sdk.NewCoins(sdk.NewInt64Coin("stake", 1))))

	govk := governmentKeeper.NewKeeper(cdc, keys[governmentTypes.StoreKey])

	k := keeper.NewKeeper(cdc, keys[types.StoreKey], sk, govk, ak)

	// Set module accounts
	memAcc := supply.NewEmptyModuleAccount(types.ModuleName, supply.Minter, supply.Burner)
	k.SupplyKeeper.SetModuleAccount(ctx, memAcc)

	// Set the stable credits denom
	k.SetStableCreditsDenom(ctx, "uccc")

	return ctx, bk, govk, k
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)

	types.RegisterCodec(cdc)

	cdc.Seal()

	return cdc
}

// Testing variables
var testUser, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
var testUser2, _ = sdk.AccAddressFromBech32("cosmos1h7tw92a66gr58pxgmf6cc336lgxadpjz5d5psf")
var testTsp, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var testDenom = "ucommercio"
var stableCreditDenom = "uccc"
