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
				sentReceiptBz := store.Get(getSentReceiptsIdsUUIDStoreKey(sender, testReceipt.UUID))
				require.Equal(t, testReceipt.UUID, string(sentReceiptBz))

				recipient, err := sdk.AccAddressFromBech32(testReceipt.Recipient)
				require.NoError(t, err)
				receivedReceiptBz := store.Get(getReceivedReceiptsIdsUUIDStoreKey(recipient, testReceipt.UUID))
				require.Equal(t, testReceipt.UUID, string(receivedReceiptBz))
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
		// {
		// 	name:   "empty",
		// 	sender: types.ValidDocument.Sender,
		// },
		// {
		// 	name:   "empty receipts",
		// 	sender: types.ValidDocument.Sender,
		// 	docs:   []types.Document{types.ValidDocument},
		// },
		// {
		// 	name:   "one receipt",
		// 	sender: types.ValidDocument.Sender,
		// 	docs:   []types.Document{types.ValidDocument},
		// 	receipts: []types.DocumentReceipt{
		// 		types.ValidDocumentReceiptRecipient1,
		// 	},
		// },
		{
			name:   "two receipts",
			sender: types.ValidDocument.Sender,
			docs:   []types.Document{types.ValidDocument},
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
