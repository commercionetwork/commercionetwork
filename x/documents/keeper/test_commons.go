package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	uuid "github.com/satori/go.uuid"
)

// Testing variables

var testingSender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var testingSender2, _ = sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")
var testingRecipient, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

var testingDocument = types.Document{
	UUID:       "test-document-uuid",
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
}

var testingDocumentReceipt = types.DocumentReceipt{
	UUID:         "testing-document-receipt-uuid",
	Sender:       testingSender.String(),
	Recipient:    testingRecipient.String(),
	TxHash:       "txHash",
	DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	Proof:        "proof",
}

func createNDocument(keeper *Keeper, ctx sdk.Context, n int) []*types.Document {
	items := []*types.Document{}
	for i := 0; i < n; i++ {
		item := &types.Document{
			Sender: testingSender.String(),
			Recipients: []string{testingRecipient.String()},
			UUID: uuid.NewV4().String(),
		}
		items = append(items, item)

		_ = keeper.SaveDocument(ctx, *items[i])
	}
	return items
}

func createNDocumentReceipt(keeper *Keeper, ctx sdk.Context, n int) []*types.DocumentReceipt {
	docs := createNDocument(keeper, ctx, n)

	items := []*types.DocumentReceipt{}
	for i := range docs {
		item := &types.DocumentReceipt{
					Sender: docs[i].Recipients[0],
					DocumentUUID: docs[i].UUID,
					Recipient: docs[i].Sender,
					UUID: uuid.NewV4().String(),
				}
		items = append(items, item)

		_ = keeper.SaveReceipt(ctx, *items[i])
	}
	return items
}