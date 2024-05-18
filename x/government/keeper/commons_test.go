package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	storemetrics "cosmossdk.io/store/metrics"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	v300 "github.com/commercionetwork/commercionetwork/x/government/legacy/v3.0.0"
	"github.com/commercionetwork/commercionetwork/x/government/types"
	cometdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var governmentTestAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var notGovernmentAddress, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")

// This function creates an environment to test the government module
// if address is defined it will be used to add the government address
func setupKeeperWithGovernmentAddress(t testing.TB, address sdk.AccAddress) (*Keeper, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := cometdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), storemetrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	keeper := NewKeeper(
		codec.NewProtoCodec(registry), storeKey, memStoreKey,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	if address != nil {
		store := ctx.KVStore(keeper.storeKey)
		store.Set([]byte(types.GovernmentStoreKey), address)
	}

	return keeper, ctx
}

func setupKeeperWithV300Government(t testing.TB, address sdk.AccAddress) (*Keeper, sdk.Context) {
	k, ctx := setupKeeperWithGovernmentAddress(t, nil)

	require.NotNil(t, address)

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(v300.GovernmentStoreKey), address)

	return k, ctx
}
