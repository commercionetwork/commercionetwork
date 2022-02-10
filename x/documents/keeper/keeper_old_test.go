package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// func TestKeeper_UserReceivedDocumentsIteratorPlus(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		recipient []byte
// 		docs      []types.Document
// 		receipts  []types.DocumentReceipt
// 	}{
// 		{
// 			"no document in store",
// 			nil,
// 			[]types.Document{},
// 			[]types.DocumentReceipt{},
// 		},
// 		{
// 			"one document in store",
// 			testingRecipient,
// 			[]types.Document{
// 				types.ValidDocument,
// 			},
// 			[]types.DocumentReceipt{
// 				types.ValidDocumentReceipt,
// 			},
// 		},
// 		{
// 			"multiple documents in store",
// 			testingRecipient,
// 			[]types.Document{
// 				testingDocument,
// 				{ // TestingDocument with different uuid
// 					UUID:           anotherValidDocumentUUID,
// 					Sender:         testingDocument.Sender,
// 					Recipients:     testingDocument.Recipients,
// 					Metadata:       testingDocument.Metadata,
// 					ContentURI:     testingDocument.ContentURI,
// 					Checksum:       testingDocument.Checksum,
// 					EncryptionData: testingDocument.EncryptionData,
// 				},
// 			},
// 			[]types.DocumentReceipt{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			keeper, ctx := setupKeeper(t)

// 			for _, document := range tt.docs {
// 				keeper.SaveDocument(ctx, document)
// 			}

// 			for _, receipt := range tt.receipts {
// 				keeper.SaveReceipt(ctx, receipt)
// 			}

// 			rdi := keeper.UserReceivedDocumentsIterator(ctx, tt.recipient)
// 			defer rdi.Close()

// 			documents := []types.Document{}
// 			for ; rdi.Valid(); rdi.Next() {
// 				id := string(rdi.Value())
// 				doc, err := keeper.GetDocumentByID(ctx, id)
// 				require.NoError(t, err)

// 				documents = append(documents, doc)
// 			}

// 			require.ElementsMatch(t, tt.docs, documents)
// 		})
// 	}
// }

func TestKeeper_UserSentDocumentsIterator(t *testing.T) {
	tests := []struct {
		name   string
		sender []byte
		docs   []types.Document
	}{
		{
			"no document in store",
			testingSender,
			[]types.Document{},
		},
		{
			"one document in store",
			testingSender,
			[]types.Document{
				testingDocument,
			},
		},
		{
			"multiple documents in store",
			testingSender,
			[]types.Document{
				testingDocument,
				{ // TestingDocument with different uuid
					UUID:           anotherDocumentReceiptUUID,
					Sender:         testingDocument.Sender,
					Recipients:     testingDocument.Recipients,
					Metadata:       testingDocument.Metadata,
					ContentURI:     testingDocument.ContentURI,
					Checksum:       testingDocument.Checksum,
					EncryptionData: testingDocument.EncryptionData,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			store := ctx.KVStore(k.storeKey)
			for _, document := range tt.docs {
				store.Set(getDocumentStoreKey(document.UUID), k.cdc.MustMarshalBinaryBare(&document))
				store.Set(getSentDocumentsIdsUUIDStoreKey(tt.sender, document.UUID), []byte(document.UUID))
			}

			documents := []types.Document{}
			di := k.UserSentDocumentsIterator(ctx, tt.sender)
			defer di.Close()

			for ; di.Valid(); di.Next() {
				id := string(di.Value())
				doc, err := k.GetDocumentByID(ctx, id)
				require.NoError(t, err)

				documents = append(documents, doc)
			}

			require.ElementsMatch(t, tt.docs, documents)
		})
	}
}

func TestKeeper_UserSentReceiptsIterator(t *testing.T) {
	tests := []struct {
		name       string
		document   types.Document
		receipt    types.DocumentReceipt
		newReceipt types.DocumentReceipt
	}{
		{
			"empty list",
			testingDocument,
			testingDocumentReceipt,
			types.DocumentReceipt{},
		},
		{
			"sent receipt already present",
			testingDocument,
			testingDocumentReceipt,
			types.DocumentReceipt{
				UUID:         anotherValidDocumentUUID,
				Sender:       testingSender.String(),
				Recipient:    testingDocumentReceipt.Recipient,
				TxHash:       testingDocumentReceipt.TxHash,
				DocumentUUID: testingDocument.UUID,
				Proof:        testingDocumentReceipt.Proof,
			},
		},
		{
			"received receipt already present",
			testingDocument,
			testingDocumentReceipt,
			types.DocumentReceipt{
				UUID:         anotherValidDocumentUUID,
				Sender:       anotherTestingSender.String(),
				Recipient:    testingDocumentReceipt.Recipient,
				TxHash:       testingDocumentReceipt.TxHash,
				DocumentUUID: testingDocument.UUID,
				Proof:        testingDocumentReceipt.Proof,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			require.NoError(t, k.SaveDocument(ctx, tt.document))

			tdr := tt.receipt
			tdr.DocumentUUID = tt.document.UUID
			require.NoError(t, k.SaveReceipt(ctx, tdr))

			store := ctx.KVStore(k.storeKey)

			senderAccadrr, _ := sdk.AccAddressFromBech32(tdr.Sender)
			docReceiptBz := store.Get(getSentReceiptsIdsUUIDStoreKey(senderAccadrr, tdr.DocumentUUID))
			storedID := string(docReceiptBz)

			stored, err := k.GetReceiptByID(ctx, storedID)
			require.NoError(t, err)

			require.Equal(t, stored, tdr)

			require.Error(t, k.SaveReceipt(ctx, tt.newReceipt))

			var storedSlice []types.DocumentReceipt
			senderAccadrr, _ = sdk.AccAddressFromBech32(tt.receipt.Sender)
			si := k.UserSentReceiptsIterator(ctx, senderAccadrr)

			defer si.Close()
			for ; si.Valid(); si.Next() {
				rid := string(si.Value())

				newReceipt, err := k.GetReceiptByID(ctx, rid)
				require.NoError(t, err)
				storedSlice = append(storedSlice, newReceipt)
			}

			require.Equal(t, 1, len(storedSlice))
			require.Contains(t, storedSlice, tdr)
			require.NotContains(t, storedSlice, tt.newReceipt)
		})
	}
}

func TestKeeper_UserReceivedReceiptsIterator(t *testing.T) {
	tests := []struct {
		name             string
		documentReceipts []types.DocumentReceipt
	}{
		{
			"Empty list",
			[]types.DocumentReceipt{},
		},
		{
			"Filled list",
			[]types.DocumentReceipt{testingDocumentReceipt},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			for _, tdr := range tt.documentReceipts {
				store := ctx.KVStore(k.storeKey)
				recipientAccAdrr, _ := sdk.AccAddressFromBech32(tdr.Recipient)

				store.Set(getReceivedReceiptsIdsUUIDStoreKey(recipientAccAdrr, tdr.DocumentUUID), []byte(tdr.UUID))

				store.Set(getReceiptStoreKey(tdr.UUID), k.cdc.MustMarshalBinaryBare(&tdr))
			}
			recipientAccAdrr, _ := sdk.AccAddressFromBech32(testingDocumentReceipt.Recipient)
			urri := k.UserReceivedReceiptsIterator(ctx, recipientAccAdrr)
			defer urri.Close()

			receipts := []types.DocumentReceipt{}
			for ; urri.Valid(); urri.Next() {
				rid := string(urri.Value())

				r, err := k.GetReceiptByID(ctx, rid)
				require.NoError(t, err)

				receipts = append(receipts, r)
			}

			require.Equal(t, tt.documentReceipts, receipts)

		})
	}
}
