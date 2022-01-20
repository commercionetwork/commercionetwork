package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestNewQuerier(t *testing.T) {
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
			path := []string{types.QueryResolveDid, types.ValidIdentity.DidDocument.ID}
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
