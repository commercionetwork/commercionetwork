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
			ID:             anotherValidDocumentUUID,
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
		name string
		docs []types.Document
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
			k, ctx := setupKeeper(t)

			for _, document := range tt.docs {
				require.NoError(t, k.SaveDocument(ctx, document))
			}

			di := k.DocumentsIterator(ctx)
			defer di.Close()

			documents := []types.Document{}
			for ; di.Valid(); di.Next() {
				d := types.Document{}
				k.cdc.MustUnmarshalBinaryBare(di.Value(), &d)

				documents = append(documents, d)
			}

			require.ElementsMatch(t, tt.docs, documents)
		})
	}
}
