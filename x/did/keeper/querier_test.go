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
		name       string
		want       types.QueryResolveIdentityResponse
		shouldFind bool
	}{
		// {
		// 	name:       "empty",
		// 	want:       types.QueryResolveIdentityResponse{},
		// 	shouldFind: false,
		// },
		{
			name: "ok",
			want: types.QueryResolveIdentityResponse{
				Identity: &types.ValidIdentity,
			},
			shouldFind: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()
			querier := NewQuerier(*k, legacyAmino)

			k.UpdateIdentity(ctx, &types.ValidIdentity)

			path := []string{types.QueryResolveDid, types.ValidIdentity.DidDocument.ID}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})

			var got types.QueryResolveIdentityResponse

			if tt.shouldFind {
				// TODO check legacyAmino problem, this cannot unmarshal the DDO

				legacyAmino.MustUnmarshalJSON(gotBz, &got)
				t.Log(got)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			} else {
				require.Error(t, err)
			}

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
