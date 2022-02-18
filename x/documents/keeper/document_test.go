package keeper

import (
	"reflect"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_SaveDocument(t *testing.T) {

	tests := []struct {
		name           string
		storedDocument *types.Document
		testDocument   func() types.Document
		wantErr        bool
	}{
		{
			name: "ok",
			testDocument: func() types.Document {
				return types.ValidDocument
			},
		},
		{
			name: "invalid UUID",
			testDocument: func() types.Document {
				doc := types.ValidDocument
				doc.UUID = doc.UUID + "$"
				return doc
			},
			wantErr: true,
		},
		{
			name:           "document already in store",
			storedDocument: &types.ValidDocument,
			testDocument: func() types.Document {
				return types.ValidDocument
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDocument := tt.testDocument()

			keeper, ctx := setupKeeper(t)
			store := ctx.KVStore(keeper.storeKey)

			if tt.storedDocument != nil {
				store.Set(getDocumentStoreKey(tt.storedDocument.UUID), keeper.cdc.MustMarshalBinaryBare(tt.storedDocument))
			}

			if err := keeper.SaveDocument(ctx, testDocument); (err != nil) != tt.wantErr {
				t.Errorf("Keeper.SaveDocument() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				var stored types.Document
				documentBz := store.Get(getDocumentStoreKey(testDocument.UUID))
				keeper.cdc.MustUnmarshalBinaryBare(documentBz, &stored)
				require.Equal(t, stored, testDocument)

				sender, err := sdk.AccAddressFromBech32(testDocument.Sender)
				require.NoError(t, err)
				sentDocumentBz := store.Get(getSentDocumentsIdsUUIDStoreKey(sender, testDocument.UUID))
				require.Equal(t, testDocument.UUID, string(sentDocumentBz))

				for _, recipientAddr := range testDocument.Recipients {
					recipient, err := sdk.AccAddressFromBech32(recipientAddr)
					require.NoError(t, err)
					receivedDocumentBz := store.Get(getReceivedDocumentsIdsUUIDStoreKey(recipient, testDocument.UUID))
					require.Equal(t, testDocument.UUID, string(receivedDocumentBz))
				}
			}
		})
	}
}

func TestKeeper_GetDocumentById(t *testing.T) {
	tests := []struct {
		name           string
		storedDocument *types.Document
		ID             string
		want           types.Document
		wantErr        bool
	}{
		{
			name:    "empty store",
			ID:      types.ValidDocument.UUID,
			wantErr: true,
		},
		{
			name:           "ok",
			storedDocument: &types.ValidDocument,
			ID:             types.ValidDocument.UUID,
			want:           types.ValidDocument,
		},
		{
			name:           "store with another document",
			storedDocument: &types.ValidDocument,
			ID:             types.AnotherValidDocument.UUID,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keeper, ctx := setupKeeper(t)

			if tt.storedDocument != nil {
				err := keeper.SaveDocument(ctx, *tt.storedDocument)
				require.NoError(t, err)
			}

			got, err := keeper.GetDocumentByID(ctx, tt.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keeper.GetDocumentByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.GetDocumentByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_DocumentsIterator(t *testing.T) {
	tests := []struct {
		name      string
		documents []types.Document
	}{
		{
			"empty",
			[]types.Document{},
		},
		{
			"one",
			[]types.Document{
				types.ValidDocument,
			},
		},
		{
			"multiple",
			[]types.Document{
				types.ValidDocument,
				types.AnotherValidDocument,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keeper, ctx := setupKeeper(t)

			for _, document := range tt.documents {
				require.NoError(t, keeper.SaveDocument(ctx, document))
			}

			di := keeper.DocumentsIterator(ctx)
			defer di.Close()

			documents := []types.Document{}
			for ; di.Valid(); di.Next() {
				d := types.Document{}
				keeper.cdc.MustUnmarshalBinaryBare(di.Value(), &d)

				documents = append(documents, d)
			}

			require.ElementsMatch(t, tt.documents, documents)
		})
	}
}

func TestKeeper_UserSentDocumentsIterator(t *testing.T) {
	tests := []struct {
		name       string
		sender     string
		storedDocs []types.Document
		wantedDocs []types.Document
	}{
		{
			name:   "empty",
			sender: types.ValidDocument.Sender,
		},
		{
			name:       "one document",
			sender:     types.ValidDocument.Sender,
			storedDocs: []types.Document{types.ValidDocument},
			wantedDocs: []types.Document{types.ValidDocument},
		},
		{
			name:       "two documents",
			sender:     types.ValidDocument.Sender,
			storedDocs: []types.Document{types.ValidDocument, types.AnotherValidDocument},
			wantedDocs: []types.Document{types.ValidDocument, types.AnotherValidDocument},
		},
		{
			name:       "two documents, one by another sender",
			sender:     types.ValidDocument.Sender,
			storedDocs: []types.Document{types.ValidDocument, types.ValidDocumentDifferentSenderRecipients},
			wantedDocs: []types.Document{types.ValidDocument},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keeper, ctx := setupKeeper(t)

			for _, document := range tt.storedDocs {
				keeper.SaveDocument(ctx, document)
			}

			senderAddr, err := sdk.AccAddressFromBech32(tt.sender)
			require.NoError(t, err)

			documents := []types.Document{}
			di := keeper.UserSentDocumentsIterator(ctx, senderAddr)
			defer di.Close()

			for ; di.Valid(); di.Next() {
				id := string(di.Value())
				doc, err := keeper.GetDocumentByID(ctx, id)
				require.NoError(t, err)

				documents = append(documents, doc)
			}

			require.ElementsMatch(t, tt.wantedDocs, documents)
		})
	}
}

func TestKeeper_UserReceivedDocumentsIterator(t *testing.T) {
	tests := []struct {
		name            string
		recipient       string
		storedDocuments []types.Document
		wantedDocuments []types.Document
	}{
		{
			name:      "empty",
			recipient: types.ValidDocumentReceiptRecipient1.Recipient,
		},
		{
			name:            "one document",
			recipient:       types.ValidDocumentReceiptRecipient1.Sender,
			storedDocuments: []types.Document{types.ValidDocument},
			wantedDocuments: []types.Document{types.ValidDocument},
		},
		{
			name:            "two documents",
			recipient:       types.ValidDocumentReceiptRecipient1.Sender,
			storedDocuments: []types.Document{types.ValidDocument, types.AnotherValidDocument},
			wantedDocuments: []types.Document{types.ValidDocument, types.AnotherValidDocument},
		},
		{
			name:            "two documents with one with different recipient",
			recipient:       types.ValidDocumentReceiptRecipient1.Sender,
			storedDocuments: []types.Document{types.ValidDocument, types.ValidDocumentDifferentSenderRecipients},
			wantedDocuments: []types.Document{types.ValidDocument},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keeper, ctx := setupKeeper(t)

			for _, document := range tt.storedDocuments {
				keeper.SaveDocument(ctx, document)
			}

			recipientAddr, err := sdk.AccAddressFromBech32(tt.recipient)
			require.NoError(t, err)

			rdi := keeper.UserReceivedDocumentsIterator(ctx, recipientAddr)
			defer rdi.Close()

			documents := []types.Document{}
			for ; rdi.Valid(); rdi.Next() {
				id := string(rdi.Value())
				doc, err := keeper.GetDocumentByID(ctx, id)
				require.NoError(t, err)

				documents = append(documents, doc)
			}

			require.ElementsMatch(t, tt.wantedDocuments, documents)
		})
	}
}
