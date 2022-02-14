package keeper

import (
	"reflect"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_SaveReceipt(t *testing.T) {

	tests := []struct {
		name           string
		storedDocument *types.Document
		storedReceipt  *types.DocumentReceipt
		receipt        func() types.DocumentReceipt
		wantErr        bool
	}{
		{
			name:           "ok",
			storedDocument: &types.ValidDocument,
			receipt: func() types.DocumentReceipt {
				return types.ValidDocumentReceiptRecipient1
			},
			wantErr: false,
		},
		{
			name: "invalid UUID",
			receipt: func() types.DocumentReceipt {
				rec := types.ValidDocumentReceiptRecipient1
				rec.UUID = rec.UUID + "$"
				return rec
			},
			wantErr: true,
		},
		{
			name: "no corresponding document",
			receipt: func() types.DocumentReceipt {
				return types.ValidDocumentReceiptRecipient1
			},
			wantErr: true,
		},
		{
			name:           "receipt already in store",
			storedDocument: &types.ValidDocument,
			storedReceipt:  &types.ValidDocumentReceiptRecipient1,
			receipt: func() types.DocumentReceipt {
				return types.ValidDocumentReceiptRecipient1
			},
			wantErr: true,
		},
		{
			name:           "a receipt for the document already been sent by the sender",
			storedDocument: &types.ValidDocument,
			storedReceipt:  &types.ValidDocumentReceiptRecipient1,
			receipt: func() types.DocumentReceipt {
				receipt := types.ValidDocumentReceiptRecipient1
				receipt.UUID = "4beff972-03a4-42da-9ebd-9303ae342be8"
				return receipt
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testReceipt := tt.receipt()

			keeper, ctx := setupKeeper(t)
			store := ctx.KVStore(keeper.storeKey)

			if tt.storedDocument != nil {
				err := keeper.SaveDocument(ctx, *tt.storedDocument)
				require.NoError(t, err)

				if tt.storedReceipt != nil {
					store.Set(getReceiptStoreKey(tt.storedReceipt.UUID), keeper.cdc.MustMarshalBinaryBare(tt.storedReceipt))

					marshaledReceiptID := []byte(tt.storedReceipt.UUID)

					sender, err := sdk.AccAddressFromBech32(tt.storedReceipt.Sender)
					require.NoError(t, err)
					store.Set(getSentReceiptsIdsUUIDStoreKey(sender, tt.storedReceipt.DocumentUUID), marshaledReceiptID)

					// TODO: remove check for (recipientAccAdrr, receipt.UUID)
					// recipient, err := sdk.AccAddressFromBech32(tt.storedReceipt.Recipient)
					// require.NoError(t, err)
					// store.Set(getReceivedReceiptsIdsUUIDStoreKey(recipient, tt.storedReceipt.UUID), marshaledReceiptID)

					store.Set(getDocumentReceiptsIdsUUIDStoreKey(tt.storedReceipt.DocumentUUID, tt.storedReceipt.UUID), marshaledReceiptID)
				}
			}

			if err := keeper.SaveReceipt(ctx, testReceipt); (err != nil) != tt.wantErr {
				t.Errorf("Keeper.SaveReceipt() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				var stored types.DocumentReceipt
				receiptBz := store.Get(getReceiptStoreKey(testReceipt.UUID))
				keeper.cdc.MustUnmarshalBinaryBare(receiptBz, &stored)
				require.Equal(t, stored, testReceipt)

				sender, err := sdk.AccAddressFromBech32(testReceipt.Sender)
				require.NoError(t, err)
				sentReceiptBz := store.Get(getSentReceiptsIdsUUIDStoreKey(sender, testReceipt.DocumentUUID))
				require.Equal(t, testReceipt.UUID, string(sentReceiptBz))

				recipient, err := sdk.AccAddressFromBech32(testReceipt.Recipient)
				require.NoError(t, err)
				receivedReceiptBz := store.Get(getReceivedReceiptsIdsUUIDStoreKey(recipient, testReceipt.UUID))
				require.Equal(t, testReceipt.UUID, string(receivedReceiptBz))

				documentsReceiptsBz := store.Get(getDocumentReceiptsIdsUUIDStoreKey(testReceipt.DocumentUUID, testReceipt.UUID))
				require.Equal(t, testReceipt.UUID, string(documentsReceiptsBz))
			}
		})
	}
}

func TestKeeper_GetReceiptByID(t *testing.T) {

	tests := []struct {
		name          string
		storedReceipt *types.DocumentReceipt
		ID            string
		want          types.DocumentReceipt
		wantErr       bool
	}{
		{
			name:          "empty store",
			storedReceipt: nil,
			ID:            types.ValidDocumentReceiptRecipient1.UUID,
			wantErr:       true,
		},
		{
			name:          "ok",
			storedReceipt: &types.ValidDocumentReceiptRecipient1,
			ID:            types.ValidDocumentReceiptRecipient1.UUID,
			want:          types.ValidDocumentReceiptRecipient1,
			wantErr:       false,
		},
		{
			name:          "store with another receipt",
			storedReceipt: &types.ValidDocumentReceiptRecipient1,
			ID:            anotherDocumentReceiptUUID,
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			keeper, ctx := setupKeeper(t)

			if tt.storedReceipt != nil {
				err := keeper.SaveDocument(ctx, types.ValidDocument)
				require.NoError(t, err)

				err = keeper.SaveReceipt(ctx, *tt.storedReceipt)
				require.NoError(t, err)
			}

			got, err := keeper.GetReceiptByID(ctx, tt.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keeper.GetReceiptByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.GetReceiptByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_UserSentReceiptsIterator(t *testing.T) {
	tests := []struct {
		name     string
		sender   string
		docs     []types.Document
		receipts []types.DocumentReceipt
	}{
		{
			name:   "empty",
			sender: types.ValidDocumentReceiptRecipient1.Sender,
		},
		{
			name:   "empty receipts",
			sender: types.ValidDocumentReceiptRecipient1.Sender,
			docs:   []types.Document{types.ValidDocument},
		},
		{
			name:   "one receipt",
			sender: types.ValidDocumentReceiptRecipient1.Sender,
			docs:   []types.Document{types.ValidDocument},
			receipts: []types.DocumentReceipt{
				types.ValidDocumentReceiptRecipient1,
			},
		},
		{
			name:   "two receipts",
			sender: types.ValidDocumentReceiptRecipient1.Sender,
			docs:   []types.Document{types.ValidDocument},
			receipts: []types.DocumentReceipt{
				types.ValidDocumentReceiptRecipient1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keeper, ctx := setupKeeper(t)

			for _, document := range tt.docs {
				keeper.SaveDocument(ctx, document)
			}

			for _, receipt := range tt.receipts {
				keeper.SaveReceipt(ctx, receipt)
			}

			senderAddr, err := sdk.AccAddressFromBech32(tt.sender)
			require.NoError(t, err)

			receipts := []types.DocumentReceipt{}
			di := keeper.UserSentReceiptsIterator(ctx, senderAddr)
			defer di.Close()

			for ; di.Valid(); di.Next() {
				id := string(di.Value())
				receipt, err := keeper.GetReceiptByID(ctx, id)
				require.NoError(t, err)

				receipts = append(receipts, receipt)
			}

			require.ElementsMatch(t, tt.receipts, receipts)
		})
	}
}

func TestKeeper_UserReceivedReceiptsIterator(t *testing.T) {
	tests := []struct {
		name           string
		recipient      string
		docs           []types.Document
		storedReceipts []types.DocumentReceipt
		receipts       []types.DocumentReceipt
	}{
		{
			name:      "empty",
			recipient: types.ValidDocumentReceiptRecipient1.Recipient,
		},
		{
			name:      "empty receipts",
			recipient: types.ValidDocumentReceiptRecipient1.Recipient,
			docs:      []types.Document{types.ValidDocument},
		},
		{
			name:      "one receipt",
			recipient: types.ValidDocumentReceiptRecipient1.Recipient,
			docs:      []types.Document{types.ValidDocument},
			receipts: []types.DocumentReceipt{
				types.ValidDocumentReceiptRecipient1,
			},
		},
		{
			name:      "two receipts",
			recipient: types.ValidDocumentReceiptRecipient1.Recipient,
			docs:      []types.Document{types.ValidDocument},
			storedReceipts: []types.DocumentReceipt{
				types.ValidDocumentReceiptRecipient1,
				types.ValidDocumentReceiptRecipient2,
			},
			receipts: []types.DocumentReceipt{
				types.ValidDocumentReceiptRecipient1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keeper, ctx := setupKeeper(t)

			for _, document := range tt.docs {
				keeper.SaveDocument(ctx, document)
			}

			for _, receipt := range tt.receipts {
				keeper.SaveReceipt(ctx, receipt)
			}

			recipientAddr, err := sdk.AccAddressFromBech32(tt.recipient)
			require.NoError(t, err)

			receipts := []types.DocumentReceipt{}
			di := keeper.UserReceivedReceiptsIterator(ctx, recipientAddr)
			defer di.Close()

			for ; di.Valid(); di.Next() {
				id := string(di.Value())
				receipt, err := keeper.GetReceiptByID(ctx, id)
				require.NoError(t, err)

				receipts = append(receipts, receipt)
			}

			require.ElementsMatch(t, tt.receipts, receipts)
		})
	}
}

func TestKeeper_UserDocumentsReceiptsIterator(t *testing.T) {
	tests := []struct {
		name           string
		documentUUID   string
		docs           []types.Document
		storedReceipts []types.DocumentReceipt
		receipts       []types.DocumentReceipt
	}{
		{
			name:         "empty",
			documentUUID: types.ValidDocumentReceiptRecipient1.DocumentUUID,
		},
		{
			name:         "empty receipts",
			documentUUID: types.ValidDocumentReceiptRecipient1.DocumentUUID,
			docs:         []types.Document{types.ValidDocument},
		},
		{
			name:         "one receipt",
			documentUUID: types.ValidDocumentReceiptRecipient1.DocumentUUID,
			docs:         []types.Document{types.ValidDocument},
			receipts: []types.DocumentReceipt{
				types.ValidDocumentReceiptRecipient1,
			},
		},
		{
			name:         "two receipts",
			documentUUID: types.ValidDocumentReceiptRecipient1.DocumentUUID,
			docs:         []types.Document{types.ValidDocument},
			storedReceipts: []types.DocumentReceipt{
				types.ValidDocumentReceiptRecipient1,
				types.ValidDocumentReceiptRecipient2,
			},
			receipts: []types.DocumentReceipt{
				types.ValidDocumentReceiptRecipient1,
				types.ValidDocumentReceiptRecipient2,
			},
		},
		{
			name:         "three receipts with one not concerning the document",
			documentUUID: types.ValidDocumentReceiptRecipient1.DocumentUUID,
			docs: []types.Document{
				types.ValidDocument,
				types.AnotherValidDocument,
			},
			storedReceipts: []types.DocumentReceipt{
				types.ValidDocumentReceiptRecipient1,
				types.ValidDocumentReceiptRecipient2,
				types.AnotherValidDocumentReceipt,
			},
			receipts: []types.DocumentReceipt{
				types.ValidDocumentReceiptRecipient1,
				types.ValidDocumentReceiptRecipient2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keeper, ctx := setupKeeper(t)

			for _, document := range tt.docs {
				keeper.SaveDocument(ctx, document)
			}

			for _, receipt := range tt.receipts {
				keeper.SaveReceipt(ctx, receipt)
			}

			receipts := []types.DocumentReceipt{}
			di := keeper.UserDocumentsReceiptsIterator(ctx, tt.documentUUID)
			defer di.Close()

			for ; di.Valid(); di.Next() {
				id := string(di.Value())
				receipt, err := keeper.GetReceiptByID(ctx, id)
				require.NoError(t, err)

				receipts = append(receipts, receipt)
			}

			require.ElementsMatch(t, tt.receipts, receipts)
		})
	}
}
