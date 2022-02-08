package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentGet(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNDocument(keeper, ctx, 10)
	for _, item := range items {
		actual, _ := keeper.GetDocumentByID(ctx, item.UUID)
		assert.Equal(t, *item, actual)
	}
}

func TestDocumentExist(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	items := createNDocument(keeper, ctx, 10)
	for _, item := range items {
		assert.True(t, keeper.HasDocument(ctx, item.UUID))
	}
}

func TestKeeper_SaveDocumentPlus(t *testing.T) {
	tests := []struct {
		name           string
		storedDocument types.Document
		document       types.Document
		newRecipient   string
	}{
		{
			"No document in store",
			types.Document{},
			testingDocument,
			"",
		},
		{
			"One document in store, different recipient",
			testingDocument,
			testingDocument,
			anotherTestingRecipient.String(),
		},
		{
			"One document in store, different uuid",
			testingDocument,
			types.Document{
				UUID:       anotherValidDocumentUUID,
				ContentURI: testingDocument.ContentURI,
				Metadata:   testingDocument.Metadata,
				Checksum:   testingDocument.Checksum,
				Sender:     testingDocument.Sender,
				Recipients: testingDocument.Recipients,
			},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			store := ctx.KVStore(k.storeKey)

			if !tt.storedDocument.Equals(types.Document{}) {
				store.Set(getSentDocumentsIdsUUIDStoreKey(testingSender, tt.storedDocument.UUID), []byte(tt.storedDocument.UUID))
				store.Set(getReceivedDocumentsIdsUUIDStoreKey(testingRecipient, tt.storedDocument.UUID), []byte(tt.storedDocument.UUID))
			}

			if tt.newRecipient != "" {
				tt.document.Recipients = []string{tt.newRecipient}
			}

			err := k.SaveDocument(ctx, tt.document)
			require.NoError(t, err)

			docsBz := store.Get(getDocumentStoreKey(tt.document.UUID))
			sentDocsBz := store.Get(getSentDocumentsIdsUUIDStoreKey(testingSender, tt.document.UUID))
			receivedDocsBz := store.Get(getReceivedDocumentsIdsUUIDStoreKey(testingRecipient, tt.document.UUID))

			if tt.newRecipient != "" {
				recipientAddr, _ := sdk.AccAddressFromBech32(tt.newRecipient)
				newReceivedDocsBz := store.Get(getReceivedDocumentsIdsUUIDStoreKey(recipientAddr, tt.document.UUID))

				newReceivedDocs := string(newReceivedDocsBz)
				require.Equal(t, tt.document.UUID, newReceivedDocs)
			}

			var stored types.Document
			k.cdc.MustUnmarshalBinaryBare(docsBz, &stored)
			require.Equal(t, stored, tt.document)

			sentDocs := string(sentDocsBz)
			receivedDocs := string(receivedDocsBz)
			require.Equal(t, tt.document.UUID, sentDocs)
			require.Equal(t, tt.document.UUID, receivedDocs)

		})
	}
}

// SaveDocument duplicated tests
func TestKeeper_SaveDocument(t *testing.T) {
	tests := []struct {
		name     string
		document types.Document
		wantErr  bool
	}{
		{
			"duplicated document",
			testingDocument,
			true,
		},
		{
			"UUID already in store",
			types.Document{
				UUID: "test-document-uuid",
			},
			true,
		},
		{
			"document's UUID not in store",
			types.Document{
				UUID:       anotherValidDocumentUUID,
				ContentURI: "https://example.com/document",
				Metadata: &types.DocumentMetadata{
					ContentURI: "https://example.com/document/metadata",
					Schema: &types.DocumentMetadataSchema{
						URI:     "https://example.com/document/metadata/schema",
						Version: "1.0.0",
					},
				},
				Checksum: &types.DocumentChecksum{
					Value:     "93dfcaf3d923ec47edb8580667473987",
					Algorithm: "md5",
				},
				Sender:     testingSender.String(),
				Recipients: append([]string{}, testingRecipient.String()),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			store := ctx.KVStore(k.storeKey)
			store.Set(getDocumentStoreKey(testingDocument.UUID), k.cdc.MustMarshalBinaryBare(&testingDocument))

			if tt.wantErr {
				require.Error(t, k.SaveDocument(ctx, tt.document))
			} else {
				require.NoError(t, k.SaveDocument(ctx, tt.document))
			}
		})
	}
}
