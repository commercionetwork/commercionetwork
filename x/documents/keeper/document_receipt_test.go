package keeper

import (
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
				return types.ValidDocumentReceipt
			},
			wantErr: false,
		},
		{
			name: "invalid UUID",
			receipt: func() types.DocumentReceipt {
				rec := types.ValidDocumentReceipt
				rec.UUID = rec.UUID + "$"
				return rec
			},
			wantErr: true,
		},
		{
			name: "no corresponding document",
			receipt: func() types.DocumentReceipt {
				return types.ValidDocumentReceipt
			},
			wantErr: true,
		},
		{
			name:           "receipt already in store",
			storedDocument: &types.ValidDocument,
			storedReceipt:  &types.ValidDocumentReceipt,
			receipt: func() types.DocumentReceipt {
				return types.ValidDocumentReceipt
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
				sentReceiptBz := store.Get(getSentReceiptsIdsUUIDStoreKey(sender, testReceipt.DocumentUUID))
				require.Equal(t, testReceipt.UUID, string(sentReceiptBz))

				recipient, err := sdk.AccAddressFromBech32(testReceipt.Recipient)
				require.NoError(t, err)
				receivedReceiptBz := store.Get(getReceivedReceiptsIdsUUIDStoreKey(recipient, testReceipt.DocumentUUID))
				require.Equal(t, testReceipt.UUID, string(receivedReceiptBz))
			}
		})
	}
}
