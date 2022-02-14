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
		name      string
		recipient string
		want      []types.Document
		wantErr   bool
	}{
		{
			name:      "invalid recipient",
			recipient: "",
			wantErr:   true,
		},
		{
			name:      "empty",
			recipient: types.ValidDocumentReceiptRecipient1.Sender,
			want:      []types.Document(nil),
		},
		{
			name:      "one",
			recipient: types.ValidDocumentReceiptRecipient1.Sender,
			want:      []types.Document{types.ValidDocument},
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

			path := []string{types.QueryReceivedDocuments, tt.recipient}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				var got []types.Document
				legacyAmino.MustUnmarshalJSON(gotBz, &got)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_queryGetSentDocuments(t *testing.T) {

	tests := []struct {
		name    string
		sender  string
		want    []types.Document
		wantErr bool
	}{
		{
			name:    "invalid sender",
			sender:  "",
			wantErr: true,
		},
		{
			name:   "empty",
			sender: types.ValidDocument.Sender,
			want:   []types.Document(nil),
		},
		{
			name:   "one",
			sender: types.ValidDocument.Sender,
			want:   []types.Document{types.ValidDocument},
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

			path := []string{types.QuerySentDocuments, tt.sender}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				var got []types.Document
				legacyAmino.MustUnmarshalJSON(gotBz, &got)
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func Test_queryGetReceivedDocsReceipts(t *testing.T) {

	tests := []struct {
		name            string
		receiver        string
		storedDocuments []types.Document
		storedReceipts  []types.DocumentReceipt
		wantedReceipts  []types.DocumentReceipt
		wantErr         bool
	}{
		{
			name:     "invalid receiver",
			receiver: "",
			wantErr:  true,
		},
		{
			name:     "empty store",
			receiver: types.ValidDocumentReceiptRecipient1.Recipient,
		},
		{
			name:            "receipt not stored",
			receiver:        types.ValidDocumentReceiptRecipient1.Recipient,
			storedDocuments: []types.Document{types.ValidDocument},
		},
		{
			name:            "one",
			receiver:        types.ValidDocumentReceiptRecipient1.Recipient,
			storedDocuments: []types.Document{types.ValidDocument},
			storedReceipts:  []types.DocumentReceipt{types.ValidDocumentReceiptRecipient1},
			wantedReceipts:  []types.DocumentReceipt{types.ValidDocumentReceiptRecipient1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			for _, document := range tt.storedDocuments {
				err := k.SaveDocument(ctx, document)
				require.NoError(t, err)
			}
			for _, receipt := range tt.storedReceipts {
				err := k.SaveReceipt(ctx, receipt)
				require.NoError(t, err)
			}

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()
			querier := NewQuerier(*k, legacyAmino)

			path := []string{types.QueryReceivedReceipts, tt.receiver}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				var got []types.DocumentReceipt
				legacyAmino.MustUnmarshalJSON(gotBz, &got)
				require.NoError(t, err)
				require.Equal(t, tt.wantedReceipts, got)
			}

		})
	}
}

func Test_queryGetSentDocsReceipts(t *testing.T) {
	tests := []struct {
		name            string
		sender          string
		storedDocuments []types.Document
		storedReceipts  []types.DocumentReceipt
		wantedReceipts  []types.DocumentReceipt
		wantErr         bool
	}{
		{
			name:   "empty store",
			sender: types.ValidDocumentReceiptRecipient1.Sender,
		},

		{
			name:            "one",
			sender:          types.ValidDocumentReceiptRecipient1.Sender,
			storedDocuments: []types.Document{types.ValidDocument},
			storedReceipts:  []types.DocumentReceipt{types.ValidDocumentReceiptRecipient1},
			wantedReceipts:  []types.DocumentReceipt{types.ValidDocumentReceiptRecipient1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				k, ctx := setupKeeper(t)

				for _, document := range tt.storedDocuments {
					err := k.SaveDocument(ctx, document)
					require.NoError(t, err)
				}
				for _, receipt := range tt.storedReceipts {
					err := k.SaveReceipt(ctx, receipt)
					require.NoError(t, err)
				}

				app := simapp.Setup(false)
				legacyAmino := app.LegacyAmino()
				querier := NewQuerier(*k, legacyAmino)

				path := []string{types.QuerySentReceipts, tt.sender}
				gotBz, err := querier(ctx, path, abci.RequestQuery{})

				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)

					var got []types.DocumentReceipt
					legacyAmino.MustUnmarshalJSON(gotBz, &got)
					require.NoError(t, err)
					require.Equal(t, tt.wantedReceipts, got)
				}
			})
		})
	}
}
