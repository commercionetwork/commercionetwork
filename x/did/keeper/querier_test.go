package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestNewQuerier_queryGetLastIdentityOfAddress(t *testing.T) {
	tests := []struct {
		name string
		want types.QueryResolveIdentityResponse
	}{
		{
			name: "empty",
			want: types.QueryResolveIdentityResponse{},
		},
		{
			name: "ok",
			want: types.QueryResolveIdentityResponse{
				Identity: &types.ValidIdentity,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			if tt.want.Identity != nil {
				k.SetIdentity(ctx, types.ValidIdentity)
			}

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()
			querier := NewQuerier(*k, legacyAmino)
			path := []string{types.QueryResolveIdentity, types.ValidIdentity.DidDocument.ID}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})

			if tt.want.Identity == nil {
				require.Error(t, err)
				return
			}

			var got types.QueryResolveIdentityResponse

			legacyAmino.MustUnmarshalJSON(gotBz, &got)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNewQuerier_GetIdentityHistoryOfAddress(t *testing.T) {
	tests := []struct {
		name string
		want types.QueryResolveIdentityHistoryResponse
	}{
		// TODO: consider changing proto of query and keeper methods to return a slice of values, no pointers
		// then, want.Identities can be []types.Identity
		{
			name: "empty",
			want: types.QueryResolveIdentityHistoryResponse{
				Identities: nil,
			},
		},
		{
			name: "ok",
			want: types.QueryResolveIdentityHistoryResponse{
				Identities: identitiesAtIncreasingMoments(2, 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			for _, identity := range tt.want.Identities {
				k.SetIdentity(ctx, *identity)
			}

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()
			querier := NewQuerier(*k, legacyAmino)
			path := []string{types.QueryResolveIdentityHistory, types.ValidIdentity.DidDocument.ID}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})
			require.NoError(t, err)

			var got types.QueryResolveIdentityHistoryResponse

			legacyAmino.MustUnmarshalJSON(gotBz, &got)

			expected := types.QueryResolveIdentityHistoryResponse{
				Identities: tt.want.Identities,
			}
			require.NoError(t, err)
			require.Equal(t, expected, got)
		})
	}
}

func Test_NewQuerier_default(t *testing.T) {

	t.Run("default request", func(t *testing.T) {
		k, ctx := setupKeeper(t)

		app := simapp.Setup(false)
		legacyAmino := app.LegacyAmino()
		querier := NewQuerier(*k, legacyAmino)
		path := []string{"abcd"}
		_, err := querier(ctx, path, abci.RequestQuery{})
		require.Error(t, err)
	})
}