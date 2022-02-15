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

func Test_queryGetSentDocuments(t *testing.T) {

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

			path := []string{types.QuerySentDocuments, testingSender.String()}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})

			var got []types.Document
			legacyAmino.MustUnmarshalJSON(gotBz, &got)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)

		})
	}
}

func Test_queryGetReceivedDocsReceipts(t *testing.T) {

	tests := []struct {
		name      string
		docOrigin []types.Document
		want      []types.DocumentReceipt
		wantErr   bool
	}{
		{
			name:      "empty",
			docOrigin: []types.Document(nil),
			want:      []types.DocumentReceipt(nil),
		},
		{
			name:      "one",
			docOrigin: []types.Document{testingDocument},
			want:      []types.DocumentReceipt{testingDocumentReceipt},
		},
		{
			name:      "origin document not present",
			docOrigin: []types.Document{testingDocument},
			want:      []types.DocumentReceipt{testingDocumentReceiptNoDoc},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			for _, document := range tt.docOrigin {
				err := k.SaveDocument(ctx, document)
				require.NoError(t, err)
			}
			for _, receipt := range tt.want {
				err := k.SaveReceipt(ctx, receipt)
				if tt.wantErr {
					require.Error(t, err)
					return
				} else {
					require.NoError(t, err)
				}
			}

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()
			querier := NewQuerier(*k, legacyAmino)

			path := []string{types.QueryReceivedReceipts, testingRecipient.String()}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})

			var got []types.DocumentReceipt
			legacyAmino.MustUnmarshalJSON(gotBz, &got)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_queryGetSentDocsReceipts(t *testing.T) {
	tests := []struct {
		name      string
		docOrigin []types.Document
		want      []types.DocumentReceipt
		wantErr   bool
	}{
		{
			name:      "empty",
			docOrigin: []types.Document(nil),
			want:      []types.DocumentReceipt(nil),
		},
		{
			name:      "one",
			docOrigin: []types.Document{testingDocument},
			want:      []types.DocumentReceipt{testingDocumentReceipt},
		},
		{
			name:      "origin document not present",
			docOrigin: []types.Document{testingDocument},
			want:      []types.DocumentReceipt{testingDocumentReceiptNoDoc},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				k, ctx := setupKeeper(t)

				for _, document := range tt.docOrigin {
					err := k.SaveDocument(ctx, document)
					require.NoError(t, err)
				}
				for _, receipt := range tt.want {
					err := k.SaveReceipt(ctx, receipt)
					if tt.wantErr {
						require.Error(t, err)
						return
					} else {
						require.NoError(t, err)
					}
				}

				app := simapp.Setup(false)
				legacyAmino := app.LegacyAmino()
				querier := NewQuerier(*k, legacyAmino)

				path := []string{types.QuerySentReceipts, testingSender.String()}
				gotBz, err := querier(ctx, path, abci.RequestQuery{})

				var got []types.DocumentReceipt
				legacyAmino.MustUnmarshalJSON(gotBz, &got)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			})
		})
	}
}
