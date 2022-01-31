package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

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

func Test_queryGetReceivedDocuments(t *testing.T) {

	tests := []struct {
		name    string
		want    []types.Document
		wantErr bool
	}{
		{
			name: "empty",
			want: []types.Document(nil),
		},
		{
			name: "one",
			want: []types.Document{testingDocument},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			for _, document := range tt.want {
				err := k.SaveDocument(ctx, document)
				require.NoError(t, err)
			}

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()
			querier := NewQuerier(*k, legacyAmino)

			path := []string{types.QueryReceivedDocuments, testingRecipient.String()}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})

			var got []types.Document
			legacyAmino.MustUnmarshalJSON(gotBz, &got)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
