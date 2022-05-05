package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

var governmentTestAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var notGovernmentAddress, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")

// This function creates an environment to test the government module
// if address is defined it will be used to add the government address
func setupKeeperWithGovernmentAddress(t testing.TB, address sdk.AccAddress) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	legacyAmino := codec.NewLegacyAmino()

	keys := sdk.NewKVStoreKeys(
		authTypes.StoreKey,
		bankTypes.StoreKey,
		paramsTypes.StoreKey,
		types.StoreKey,
	)
	tKeys := sdk.NewTransientStoreKeys(paramsTypes.TStoreKey)
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	app := simapp.Setup(false)
	cdc := app.AppCodec()

	pk := paramsKeeper.NewKeeper(cdc, legacyAmino, keys[paramsTypes.StoreKey], tKeys[paramsTypes.TStoreKey])
	ak := authKeeper.NewAccountKeeper(cdc, keys[authTypes.StoreKey], pk.Subspace(authTypes.DefaultParams().String()), authTypes.ProtoBaseAccount, nil)
	bk := bankKeeper.NewBaseKeeper(cdc, keys[bankTypes.StoreKey], ak, pk.Subspace(bankTypes.DefaultParams().String()), nil)

	registry := codectypes.NewInterfaceRegistry()
	keeper := NewKeeper(
		codec.NewProtoCodec(registry), storeKey, memStoreKey, bk,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	if address != nil {
		store := ctx.KVStore(keeper.storeKey)
		store.Set([]byte(types.GovernmentStoreKey), address)
	}

	return keeper, ctx
}
