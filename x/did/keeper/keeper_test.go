package keeper

import (
	"testing"
	"time"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func TestIdentityGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNIdentityNew(keeper, ctx, 10)
	for _, item := range items {
		a, err := keeper.GetDdoByOwner(ctx, sdk.AccAddress(item.ID))
		require.NoError(t, err)
		assert.Equal(t, item, a)
	}
}

func TestNewDocumentExist(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNIdentityNew(keeper, ctx, 10)
	for _, item := range items {
		assert.True(t, keeper.HasIdentity(ctx, item.ID))
	}
}

func createNIdentityNew(keeper *Keeper, ctx sdk.Context, n int) []types.DidDocument {
	items := make([]types.DidDocument, n)
	for i := range items {
		_, _, addr := testdata.KeyTestPubAddr()
		items[i].ID = string(addr)
		_ = keeper.AppendDid(ctx, items[i])
	}
	return items
}

func setupKeeper(t testing.TB) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	keeper := NewKeeper(
		codec.NewProtoCodec(registry), storeKey, memStoreKey,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	ctx = ctx.WithBlockTime(time.Now())

	return keeper, ctx
}

func Test_getTimestamp(t *testing.T) {

	_, ctx := setupKeeper(t)

	timestamp, err := getTimestamp(ctx)
	assert.NoError(t, err)
	t.Log(timestamp)
	t.FailNow()

}
