package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestNewQuerier(t *testing.T) {

	_, _, addr := testdata.KeyTestPubAddr()
	didDocument := types.DidDocument{
		ID: addr.String(),
		Context: []string{
			types.ContextDidV1,
			"https://w3id.org/security/suites/ed25519-2018/v1",
			"https://w3id.org/security/suites/x25519-2019/v1",
		}}

	tests := []struct {
		name       string
		want       types.QueryResolveDidDocumentResponse
		shouldFind bool
	}{
		// {
		// 	name:       "empty",
		// 	want:       types.QueryResolveDidDocumentResponse{},
		// 	shouldFind: false,
		// },
		// {
		// 	name: "ok",
		// 	want: types.QueryResolveDidDocumentResponse{
		// 		DidDocument: &didDocument,
		// 	},
		// 	shouldFind: true,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()
			querier := NewQuerier(*k, legacyAmino)

			id := k.UpdateDidDocument(ctx, didDocument)

			path := []string{types.QueryResolveDid, id}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})

			var got *types.DidDocument

			if tt.shouldFind {
				// t.Log(string(gotBz))
				// TODO check legacyAmino problem, this cannot unmarshal the DDO
				legacyAmino.MustUnmarshalJSON(gotBz, &got)
				t.Log(got)
				require.NoError(t, err)
				require.Equal(t, tt.want.DidDocument, got)
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
