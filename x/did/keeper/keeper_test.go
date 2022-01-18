package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

// func Test_GetDidDocumentOfAddress(t *testing.T) {
// 	keeper, ctx := setupKeeper(t)
// 	ddos := createNDidDocuments(keeper, ctx, 10)
// 	for _, item := range ddos {
// 		a, err := keeper.GetDidDocumentOfAddress(ctx, item.ID)
// 		require.NoError(t, err)
// 		assert.Equal(t, item, a)
// 	}
// }

// func Test_NewDocumentExist(t *testing.T) {
// 	keeper, ctx := setupKeeper(t)
// 	ddos := createNDidDocuments(keeper, ctx, 10)
// 	for _, item := range ddos {
// 		assert.True(t, keeper.HasDidDocument(ctx, item.ID))
// 	}
// }

// func Test_UpdateDidDocument(t *testing.T) {
// 	keeper, ctx := setupKeeper(t)
// 	ddos := createNDidDocuments(keeper, ctx, 10)
// 	for _, item := range ddos {
// 		ID := keeper.UpdateDidDocument(ctx, item)

// 		require.True(t, keeper.HasDidDocument(ctx, ID))

// 		created, err := keeper.GetDidDocumentOfAddress(ctx, ID)
// 		require.NoError(t, err)
// 		require.Equal(t, created, item)
// 	}
// }

// func Test_GetAllDidDocuments(t *testing.T) {
// 	keeper, ctx := setupKeeper(t)
// 	ddos := createNDidDocuments(keeper, ctx, 10)
// 	for _, item := range ddos {
// 		ID := keeper.UpdateDidDocument(ctx, item)
// 		require.True(t, keeper.HasDidDocument(ctx, ID))
// 	}
// 	all := keeper.GetAllDidDocuments(ctx)
// 	for _, item := range ddos {
// 		var found bool
// 		for _, a := range all {
// 			if a.ID == item.ID {
// 				found = true
// 			}
// 		}
// 		require.True(t, found)
// 	}
// }

// func createNDidDocuments(keeper *Keeper, ctx sdk.Context, n int) []types.DidDocument {
// 	ddos := make([]types.DidDocument, n)
// 	for i := range ddos {
// 		_, _, addr := testdata.KeyTestPubAddr()
// 		ddos[i].ID = string(addr)
// 		_ = keeper.UpdateDidDocument(ctx, ddos[i])
// 	}
// 	return ddos
// }

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

	return keeper, ctx
}

func TestKeeper_UpdateIdentity(t *testing.T) {

	type args struct {
		identity types.Identity
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "empty store",
			args: args{},
		},
		// {
		// 	name: "append over existing",
		// 	args: args{},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)
			k.UpdateIdentity(ctx, &types.ValidIdentity)

			identity, err := k.GetLastIdentityOfAddress(ctx, types.ValidIdentity.DidDocument.ID)
			require.NoError(t, err)
			require.Equal(t, types.ValidIdentity, *identity)
		})
	}
}

func TestKeeper_GetIdentityHistoryOfAddress(t *testing.T) {
	type args struct {
		identity types.Identity
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "empty store",
			args: args{},
		},
		// {
		// 	name: "append over existing",
		// 	args: args{},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)
			k.UpdateIdentity(ctx, &types.ValidIdentity)

			identity, err := k.GetLastIdentityOfAddress(ctx, types.ValidIdentity.DidDocument.ID)
			require.NoError(t, err)
			require.Equal(t, types.ValidIdentity, *identity)

			expected := []types.Identity{*identity}

			result := k.GetIdentityHistoryOfAddress(ctx, types.ValidIdentity.DidDocument.ID)
			require.Equal(t, expected, result)

			newIdentity := &types.ValidIdentity
			newIdentity.Metadata.Updated = "abc"

			k.UpdateIdentity(ctx, newIdentity)

			result = k.GetIdentityHistoryOfAddress(ctx, types.ValidIdentity.DidDocument.ID)
			require.Len(t, result, 2)

		})
	}
}
