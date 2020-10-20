package docs

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/docs/keeper"
	"github.com/commercionetwork/commercionetwork/x/docs/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// exportDocuments exports all the Documents in the store to a slice.
func exportDocuments(ctx sdk.Context, keeper keeper.Keeper) []types.Document {
	documents := []types.Document{}
	di := keeper.DocumentsIterator(ctx)
	defer di.Close()
	for ; di.Valid(); di.Next() {
		document, duuid, err := keeper.ExtractDocument(ctx, di.Key())
		if err != nil {
			panic(
				fmt.Sprintf(
					"could not find document with UUID %s even though the user has an associated received document",
					duuid,
				),
			)
		}

		documents = append(documents, document)
	}

	return documents
}

// exportReceipts exports all the Receipts in the store to a slice.
func exportReceipts(ctx sdk.Context, keeper keeper.Keeper) []types.DocumentReceipt {
	receipts := []types.DocumentReceipt{}
	sentDri, _ := keeper.ReceiptsIterators(ctx)
	defer sentDri.Close()

	// just iterate through sent receipt, received receipts are the same:
	// the per-user selection logic happens on the key-level
	for ; sentDri.Valid(); sentDri.Next() {
		receipt, uuid, err := keeper.ExtractReceipt(ctx, sentDri.Value())
		if err != nil {
			panic(
				fmt.Sprintf(
					"could not find document receipt with UUID %s even though the user has an associated document with it",
					uuid,
				),
			)
		}

		receipts = append(receipts, receipt)
	}

	return receipts
}

// exportMetadataSchema exports all the MetadataSchema in the store to a slice.
func exportMetadataSchemes(ctx sdk.Context, keeper keeper.Keeper) []types.MetadataSchema {
	sms := []types.MetadataSchema{}
	smsi := keeper.SupportedMetadataSchemesIterator(ctx)
	defer smsi.Close()

	for ; smsi.Valid(); smsi.Next() {
		sms = append(sms, keeper.ExtractMetadataSchema(smsi.Value()))
	}

	return sms
}

// exportTrustedSchemaProviders exports all the TrustedSchemaProviders in the store to a slice.
func exportTrustedSchemaProviders(ctx sdk.Context, keeper keeper.Keeper) []sdk.AccAddress {
	tsps := []sdk.AccAddress{}
	tspi := keeper.TrustedSchemaProposersIterator(ctx)
	defer tspi.Close()

	for ; tspi.Valid(); tspi.Next() {
		tsps = append(tsps, keeper.ExtractTrustedSchemaProposer(tspi.Value()))
	}

	return tsps
}
