package documents

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/documents/keeper"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis sets documents and receipts information for genesis.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	for _, doc := range data.Documents {
		if err := keeper.SaveDocument(ctx, *doc); err != nil {
			panic(err)
		}
	}

	for _, receipt := range data.Receipts {
		if err := keeper.SaveReceipt(ctx, *receipt); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis returns the genesis state for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) types.GenesisState {
	documents := []*types.Document{}
	receipts := []*types.DocumentReceipt{}

	documentsIterator := keeper.DocumentsIterator(ctx)
	defer documentsIterator.Close()
	for ; documentsIterator.Valid(); documentsIterator.Next() {
		keyDocumentUUIDVal := documentsIterator.Key()
		documentUUID := string(keyDocumentUUIDVal[len(types.DocumentStorePrefix):])
		document, err := keeper.GetDocumentByID(ctx, documentUUID)
		if err != nil {
			panic(fmt.Sprintf("could not find document with UUID %s", documentUUID))
		}

		documents = append(documents, &document)

		receiptsIterator := keeper.UUIDDocumentsReceiptsIterator(ctx, documentUUID)
		for ; receiptsIterator.Valid(); receiptsIterator.Next() {
			keyReceiptUUIDVal := receiptsIterator.Key()
			receiptUUID := string(keyReceiptUUIDVal[len(types.DocumentsReceiptsPrefix+documentUUID+":"):])
			receipt, err := keeper.GetReceiptByID(ctx, receiptUUID)
			if err != nil {
				panic(fmt.Sprintf("could not find receipt with UUID %s", receiptUUID))
			}
			receipts = append(receipts, &receipt)
		}
	}

	return types.GenesisState{
		Documents: documents,
		Receipts:  receipts,
	}
}
