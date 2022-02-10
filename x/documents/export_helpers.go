package documents

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/documents/keeper"
	"github.com/commercionetwork/commercionetwork/x/documents/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// exportDocuments exports all the Documents in the store to a slice.
func exportDocuments(ctx sdk.Context, keeper keeper.Keeper) []*types.Document {
	documents := []*types.Document{}
	di := keeper.DocumentsIterator(ctx)
	defer di.Close()
	for ; di.Valid(); di.Next() {
		keyVal := di.Key()
		uuid := string(keyVal[len(types.DocumentStorePrefix):])
		document, err := keeper.GetDocumentByID(ctx, uuid)
		if err != nil {
			panic(fmt.Sprintf("could not find document with UUID %s", uuid))
		}

		documents = append(documents, &document)
	}

	return documents
}

// exportReceipts exports all the Receipts in the store to a slice.
func exportReceipts(ctx sdk.Context, keeper keeper.Keeper) []*types.DocumentReceipt {
	receipts := []*types.DocumentReceipt{}
	sentDri, _ := keeper.ReceiptsIterators(ctx)
	defer sentDri.Close()

	// just iterate through sent receipt, received receipts are the same:
	// the per-user selection logic happens on the key-level
	for ; sentDri.Valid(); sentDri.Next() {
		receipt, uuid, err := keeper.ExtractReceipt(ctx, sentDri.Value())
		if err != nil {
			panic(
				fmt.Sprintf(
					"could not find document receipt with UUID %s",
					uuid,
				),
			)
		}

		receipts = append(receipts, &receipt)
	}

	return receipts
}
