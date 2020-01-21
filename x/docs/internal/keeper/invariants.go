package keeper

import (
	"fmt"
	"strings"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	docsReceiptsInvName string = "docs-receipts"
	docsSchemasInvName  string = "docs-schemas-valid"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, docsReceiptsInvName,
		DocsReceiptsInvariants(k))
	ir.RegisterRoute(types.ModuleName, docsSchemasInvName,
		DocsSchemasValidInvariant(k))
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
			receipt, _, err := k.ExtractReceipt(ctx, sentReceipts.Value())

			if err != nil {
				panic("could not extract sent receipt during invariant")
			}

			if _, found := docsLookup[receipt.DocumentUUID]; found {
				continue
			}

			_, err = k.GetDocumentByID(ctx, receipt.DocumentUUID)
			if err != nil {
				return sdk.FormatInvariant(
					types.ModuleName,
					docsReceiptsInvName,
					fmt.Sprintf(
						"found sent receipt %s which refers to non-existent document %s",
						receipt.UUID,
						receipt.DocumentUUID,
					),
				), true
			}

			docsLookup[receipt.DocumentUUID] = struct{}{}
		}

		// received receipts
		for ; receivedReceipts.Valid(); receivedReceipts.Next() {
			receipt, _, err := k.ExtractReceipt(ctx, receivedReceipts.Value())

			if err != nil {
				panic("could not extract received receipt during invariant")
			}

			if _, found := docsLookup[receipt.DocumentUUID]; found {
				continue
			}

			_, err = k.GetDocumentByID(ctx, receipt.DocumentUUID)
			if err != nil {
				return sdk.FormatInvariant(
					types.ModuleName,
					docsReceiptsInvName,
					fmt.Sprintf(
						"found received receipt %s which refers to non-existent document %s",
						receipt.UUID,
						receipt.DocumentUUID,
					),
				), true
			}

			docsLookup[receipt.DocumentUUID] = struct{}{}
		}

		return "", false
	}
}

// DocsSchemasValidInvariant checks that every Document SchemaType
// is a supported one.
func DocsSchemasValidInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		di := k.DocumentsIterator(ctx)
		defer di.Close()

		for ; di.Valid(); di.Next() {
			doc, _, err := k.ExtractDocument(ctx, di.Key())
			if err != nil {
				panic("could not extract document during invariant: " + err.Error())
			}

			if len(strings.TrimSpace(doc.Metadata.SchemaType)) != 0 {
				if !k.IsMetadataSchemeTypeSupported(ctx, doc.Metadata.SchemaType) {
					return sdk.FormatInvariant(
						types.ModuleName,
						docsSchemasInvName,
						fmt.Sprintf(
							"found document %s with invalid metadata schema type %s",
							doc.UUID,
							doc.Metadata.SchemaType,
						),
					), true
				}
			}
		}
		return "", false
	}
}
