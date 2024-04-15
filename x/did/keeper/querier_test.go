package keeper

import (
	"testing"

	//"cosmossdk.io/simapp"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/stretchr/testify/require"
)

func TestNewQuerier_queryGetLastIdentityOfAddress(t *testing.T) {
	tests := []struct {
		name string
		want *types.Identity
	}{
		{
			name: "empty",
			want: nil,
		},
		{
			name: "ok",
			want: &types.ValidIdentity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			if tt.want != nil {
				k.SetIdentity(ctx, types.ValidIdentity)
			}

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()
			querier := NewQuerier(*k, legacyAmino)
			path := []string{types.QueryResolveIdentity, types.ValidIdentity.DidDocument.ID}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})

			if tt.want == nil {
				require.Error(t, err)
				return
			}

			var got types.Identity

			legacyAmino.MustUnmarshalJSON(gotBz, &got)
			require.NoError(t, err)
			require.Equal(t, *tt.want, got)
		})
	}
}

func TestNewQuerier_GetIdentityHistoryOfAddress(t *testing.T) {
	tests := []struct {
		name string
		want []*types.Identity
	}{
		{
			name: "empty",
			want: nil,
		},
		{
			name: "ok",
			want: identitiesAtIncreasingMoments(2, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			for _, identity := range tt.want {
				k.SetIdentity(ctx, *identity)
			}

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()
			querier := NewQuerier(*k, legacyAmino)
			path := []string{types.QueryResolveIdentityHistory, types.ValidIdentity.DidDocument.ID}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})
			require.NoError(t, err)

			var got []*types.Identity

			legacyAmino.MustUnmarshalJSON(gotBz, &got)

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
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
