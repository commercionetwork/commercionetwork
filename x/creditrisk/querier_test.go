package creditrisk_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/commercionetwork/commercionetwork/x/creditrisk"
	"github.com/commercionetwork/commercionetwork/x/creditrisk/types"
	"github.com/commercionetwork/commercionetwork/x/government"
)

func TestWrongPath(t *testing.T) {
	ctx, _, k := SetupTestInput()
	querier := creditrisk.NewQuerier(k)
	path := []string{"invalid"}
	var req abci.RequestQuery
	_, err := querier(ctx, path, req)
	require.Error(t, err)
}

func TestGetPoolFunds(t *testing.T) {
	ctx, sk, k := SetupTestInput()
	querier := creditrisk.NewQuerier(k)
	path := []string{"pool"}
	var req abci.RequestQuery
	res, err := querier(ctx, path, req)
	require.NoError(t, err)
	var coins sdk.Coins
	require.NoError(t, creditrisk.ModuleCdc.UnmarshalJSON(res, &coins))
	require.True(t, coins.IsZero())

	modAcc := sk.GetModuleAccount(ctx, types.ModuleName)
	newcoins := modAcc.GetCoins().Add(sdk.NewInt64Coin("coin", 10))
	modAcc.SetCoins(newcoins)
	sk.SetModuleAccount(ctx, modAcc)

	res, err = querier(ctx, path, req)
	require.NoError(t, err)
	require.NoError(t, creditrisk.ModuleCdc.UnmarshalJSON(res, &coins))
	require.True(t, coins.IsEqual(newcoins))
}

func SetupTestInput() (sdk.Context, supply.Keeper, creditrisk.Keeper) {

	memDB := db.NewMemDB()
	cdc := testCodec()

	keys := sdk.NewKVStoreKeys(
		auth.StoreKey,
		params.StoreKey,
		supply.StoreKey,
		government.StoreKey,

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
		types.ModuleName: nil,
	}
	sk := supply.NewKeeper(cdc, keys[supply.StoreKey], ak, bk, maccPerms)
	sk.SetSupply(ctx, supply.NewSupply(sdk.NewCoins(sdk.NewInt64Coin("stake", 1))))

	k := creditrisk.NewKeeper(cdc, keys[types.StoreKey], sk)

	// Set module accounts
	// memAcc := supply.NewEmptyModuleAccount(ModuleName, supply.Minter, supply.Burner)
	// k.supplyKeeper.SetModuleAccount(ctx, memAcc)

	return ctx, sk, k
}

func testCodec() *codec.Codec {
	var cdc = codec.New()

	auth.RegisterCodec(cdc)
	bank.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	sdk.RegisterCodec(cdc)
	supply.RegisterCodec(cdc)

	creditrisk.RegisterCodec(cdc)

	cdc.Seal()

	return cdc
}
