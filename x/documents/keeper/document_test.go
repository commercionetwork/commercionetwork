package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_SaveDocument(t *testing.T) {

	tests := []struct {
		name           string
		storedDocument *types.Document
		document       func() types.Document
		wantErr        bool
	}{
		{
			name: "ok",
			document: func() types.Document {
				return types.ValidDocument
			},
		},
		{
			name: "invalid UUID",
			document: func() types.Document {
				doc := types.ValidDocument
				doc.UUID = doc.UUID + "$"
				return doc
			},
			wantErr: true,
		},
		{
			name:           "document already in store",
			storedDocument: &types.ValidDocument,
			document: func() types.Document {
				return types.ValidDocument
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDocument := tt.document()

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
				docBz := store.Get(getDocumentStoreKey(testDocument.UUID))
				keeper.cdc.MustUnmarshalBinaryBare(docBz, &stored)
				require.Equal(t, stored, testDocument)

				sender, err := sdk.AccAddressFromBech32(testDocument.Sender)
				require.NoError(t, err)
				sentDocBz := store.Get(getSentDocumentsIdsUUIDStoreKey(sender, testDocument.UUID))
				require.Equal(t, testDocument.UUID, string(sentDocBz))

				for _, recipientAddr := range testDocument.Recipients {
					recipient, err := sdk.AccAddressFromBech32(recipientAddr)
					require.NoError(t, err)
					receivedDocsBz := store.Get(getReceivedDocumentsIdsUUIDStoreKey(recipient, testDocument.UUID))
					require.Equal(t, testDocument.UUID, string(receivedDocsBz))
				}
			}
		})
	}
}
