package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	docsReceiptsInvName string = "docs-receipts"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, docsReceiptsInvName,
		DocsReceiptsInvariants(k))
}

// DocsReceiptsInvariants checks that every receipt points to an
// existing Document.
func DocsReceiptsInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		// Keep a map of all documents we've seen already
		docsLookup := make(map[string]struct{})

		sentReceipts, receivedReceipts := k.ReceiptsIterators(ctx)
		defer sentReceipts.Close()
		defer receivedReceipts.Close()

		// sent receipts
		for ; sentReceipts.Valid(); sentReceipts.Next() {
			receipt, uuid, err := k.ExtractReceipt(ctx, sentReceipts.Value())

			if err != nil {
				panic("could not extract receipt during invariant")
			}

			if _, found := docsLookup[uuid]; found {
				continue
			}

			_, err = k.GetDocumentByID(ctx, uuid)
			if err != nil {
				return sdk.FormatInvariant(
					types.ModuleName,
					docsReceiptsInvName,
					fmt.Sprintf(
						receipt.UUID,
						uuid,
					),
				), true
			}

			docsLookup[uuid] = struct{}{}
		}

		// received receipts
		for ; receivedReceipts.Valid(); receivedReceipts.Next() {
			receipt, uuid, err := k.ExtractReceipt(ctx, receivedReceipts.Value())

			if err != nil {
				panic("could not extract receipt during invariant")
			}

			if _, found := docsLookup[uuid]; found {
				continue
			}

			_, err = k.GetDocumentByID(ctx, uuid)
			if err != nil {
				return sdk.FormatInvariant(
					types.ModuleName,
					docsReceiptsInvName,
					fmt.Sprintf(
						receipt.UUID,
						uuid,
					),
				), true
			}

			docsLookup[uuid] = struct{}{}
		}

		return "", false
	}
}
