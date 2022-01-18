package keeper

import (
	"reflect"
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
			k.UpdateIdentity(ctx, types.ValidIdentity)

			identity, err := k.GetLastIdentityOfAddress(ctx, types.ValidIdentity.DidDocument.ID)
			require.NoError(t, err)
			require.Equal(t, types.ValidIdentity, *identity)
		})
	}
}

var identitiesOnlyOne = []*types.Identity{&types.ValidIdentity}

func TestKeeper_GetIdentityHistoryOfAddress(t *testing.T) {
	type args struct {
		ID string
	}
	tests := []struct {
		name       string
		args       args
		identities []*types.Identity
	}{
		{
			name: "empty store",
			args: args{
				ID: types.ValidIdentity.DidDocument.ID,
			},
			identities: []*types.Identity{},
		},
		{
			name: "one",
			args: args{
				ID: types.ValidIdentity.DidDocument.ID,
			},
			identities: identitiesOnlyOne,
		},
		// {
		// 	name: "one among identities with different ID",
		// 	args: args{
		// 		ID: types.ValidIdentity.DidDocument.ID,
		// 	},
		// 	identities: ,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			expected := []*types.Identity{}
			for _, identity := range tt.identities {
				k.UpdateIdentity(ctx, *identity)
				if identity.DidDocument.ID == tt.args.ID {
					expected = append(expected, identity)
				}
			}

			result := k.GetIdentityHistoryOfAddress(ctx, types.ValidIdentity.DidDocument.ID)
			require.Equal(t, expected, result)

			// newIdentity := &types.ValidIdentity
			// newIdentity.Metadata.Updated = "abc"

			// k.UpdateIdentity(ctx, newIdentity)

			// result = k.GetIdentityHistoryOfAddress(ctx, types.ValidIdentity.DidDocument.ID)
			// require.Len(t, result, 2)

		})
	}
}

func TestKeeper_GetAllIdentities(t *testing.T) {

	tests := []struct {
		name string
		want []*types.Identity
	}{
		{
			name: "empty",
			want: []*types.Identity{},
		},
		{
			name: "one",
			want: []*types.Identity{&types.ValidIdentity},
		},
		{
			name: "more",
			want: []*types.Identity{&types.ValidIdentity},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			for _, identity := range tt.want {
				k.UpdateIdentity(ctx, *identity)
			}

			if got := k.GetAllIdentities(ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.GetAllIdentities() = %v, want %v", got, tt.want)
			}
		})
	}
}
