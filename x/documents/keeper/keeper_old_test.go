package keeper

// func TestKeeper_UserSentReceiptsIterator(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		document   types.Document
// 		receipt    types.DocumentReceipt
// 		newReceipt types.DocumentReceipt
// 	}{
// 		{
// 			"empty list",
// 			testingDocument,
// 			testingDocumentReceipt,
// 			types.DocumentReceipt{},
// 		},
// 		{
// 			"sent receipt already present",
// 			testingDocument,
// 			testingDocumentReceipt,
// 			types.DocumentReceipt{
// 				UUID:         anotherValidDocumentUUID,
// 				Sender:       testingSender.String(),
// 				Recipient:    testingDocumentReceipt.Recipient,
// 				TxHash:       testingDocumentReceipt.TxHash,
// 				DocumentUUID: testingDocument.UUID,
// 				Proof:        testingDocumentReceipt.Proof,
// 			},
// 		},
// 		{
// 			"received receipt already present",
// 			testingDocument,
// 			testingDocumentReceipt,
// 			types.DocumentReceipt{
// 				UUID:         anotherValidDocumentUUID,
// 				Sender:       anotherTestingSender.String(),
// 				Recipient:    testingDocumentReceipt.Recipient,
// 				TxHash:       testingDocumentReceipt.TxHash,
// 				DocumentUUID: testingDocument.UUID,
// 				Proof:        testingDocumentReceipt.Proof,
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			k, ctx := setupKeeper(t)

// 			require.NoError(t, k.SaveDocument(ctx, tt.document))

// 			tdr := tt.receipt
// 			tdr.DocumentUUID = tt.document.UUID
// 			require.NoError(t, k.SaveReceipt(ctx, tdr))

// 			store := ctx.KVStore(k.storeKey)

// 			senderAccadrr, _ := sdk.AccAddressFromBech32(tdr.Sender)
// 			docReceiptBz := store.Get(getSentReceiptsIdsUUIDStoreKey(senderAccadrr, tdr.DocumentUUID))
// 			storedID := string(docReceiptBz)

// 			stored, err := k.GetReceiptByID(ctx, storedID)
// 			require.NoError(t, err)

// 			require.Equal(t, stored, tdr)

// 			require.Error(t, k.SaveReceipt(ctx, tt.newReceipt))

// 			var storedSlice []types.DocumentReceipt
// 			senderAccadrr, _ = sdk.AccAddressFromBech32(tt.receipt.Sender)
// 			si := k.UserSentReceiptsIterator(ctx, senderAccadrr)

// 			defer si.Close()
// 			for ; si.Valid(); si.Next() {
// 				rid := string(si.Value())

// 				newReceipt, err := k.GetReceiptByID(ctx, rid)
// 				require.NoError(t, err)
// 				storedSlice = append(storedSlice, newReceipt)
// 			}

// 			require.Equal(t, 1, len(storedSlice))
// 			require.Contains(t, storedSlice, tdr)
// 			require.NotContains(t, storedSlice, tt.newReceipt)
// 		})
// 	}
// }
