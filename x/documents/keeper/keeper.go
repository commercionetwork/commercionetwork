package keeper

import (
	"fmt"
	"strings"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	eventNewMetadataScheme = "new_metadata_scheme"
	eventNewTMSP           = "new_tmsp"
	eventSavedDocument     = "new_saved_document"
	eventSavedReceipt      = "new_saved_receipt"
)

type (
	Keeper struct {
		cdc      codec.Marshaler
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// ----------------------
// --- Documents
// ----------------------

// SaveDocument allows the sharing of a document
func (keeper Keeper) SaveDocument(ctx sdk.Context, document types.Document) error {
	// Check the id validity
	if len(strings.TrimSpace(document.UUID)) == 0 {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid document id: %s", document.UUID))
	}

	// Check any existing document
	// when err is nil, we found a document with said document.UUID
	if _, err := keeper.GetDocumentByID(ctx, document.UUID); err == nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Document with uuid %s already present", document.UUID))
	}

	// Store the document object
	store := ctx.KVStore(keeper.storeKey)
	store.Set(getDocumentStoreKey(document.UUID), keeper.cdc.MustMarshalBinaryBare(&document))

	// Store the document as sent by the sender

	// Idea: SentDocumentsPrefix + address + document.UUID -> document.UUID
	senderAccadrr, _ := sdk.AccAddressFromBech32(document.Sender)
	sentDocumentsStoreKey := getSentDocumentsIdsUUIDStoreKey(senderAccadrr, document.UUID)

	store.Set(sentDocumentsStoreKey, []byte(document.UUID))

	// Store the documents as received for all the recipients
	for _, recipient := range document.Recipients {
		recipientAccAdrr, _ := sdk.AccAddressFromBech32(recipient)
		receivedDocumentsStoreKey := getReceivedDocumentsIdsUUIDStoreKey(recipientAccAdrr, document.UUID)

		store.Set(receivedDocumentsStoreKey, []byte(document.UUID))
	}

	attributes := make([]sdk.Attribute, 0, len(document.Recipients)+2)

	attributes = append(attributes,
		sdk.NewAttribute("sender", document.Sender),
		sdk.NewAttribute("doc_id", document.UUID),
	)

	for i, r := range document.Recipients {
		attributes = append(attributes, sdk.NewAttribute(
			fmt.Sprintf("receiver_%d", i),
			r,
		))
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventSavedDocument,
		attributes...,
	))

	return nil
}

// ExtractDocument returns a Document slice instance and its UUID given an iterator byte stream value.
func (keeper Keeper) ExtractDocument(ctx sdk.Context, keyVal []byte) (types.Document, string, error) {
	documentUUID := string(keyVal[len(types.DocumentStorePrefix):])

	document, err := keeper.GetDocumentByID(ctx, documentUUID)
	return document, documentUUID, err
}

// ----------------------
// --- Receipts
// ----------------------

// SaveReceipt allows to properly store the given receipt
func (keeper Keeper) SaveReceipt(ctx sdk.Context, receipt types.DocumentReceipt) error {
	// Check the id
	if len(strings.TrimSpace(receipt.UUID)) == 0 {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid document receipt id: %s", receipt.UUID))
	}

	if _, err := keeper.GetDocumentByID(ctx, receipt.DocumentUUID); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "recepit points to a non-existing document UUID")
	}

	store := ctx.KVStore(keeper.storeKey)
	senderAccadrr, _ := sdk.AccAddressFromBech32(receipt.Sender)
	sentReceiptsIdsStoreKey := getSentReceiptsIdsUUIDStoreKey(senderAccadrr, receipt.DocumentUUID)
	recipientAccAdrr, _ := sdk.AccAddressFromBech32(receipt.Recipient)
	receivedReceiptIdsStoreKey := getReceivedReceiptsIdsUUIDStoreKey(recipientAccAdrr, receipt.DocumentUUID)

	marshaledRecepit := keeper.cdc.MustMarshalBinaryBare(&receipt)
	marshaledRecepitID := []byte(receipt.UUID)

	// Store the receipt as sent
	if store.Has(sentReceiptsIdsStoreKey) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("sent receipt for document with UUID %s already present: %s", receipt.DocumentUUID, receipt.UUID))
	}
	store.Set(sentReceiptsIdsStoreKey, marshaledRecepitID)

	// Store the receipt as received
	if store.Has(receivedReceiptIdsStoreKey) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("received receipt for document with UUID %s already present: %s", receipt.DocumentUUID, receipt.UUID))
	}
	store.Set(receivedReceiptIdsStoreKey, marshaledRecepitID)

	// Store the receipt
	store.Set(getReceiptStoreKey(receipt.UUID), marshaledRecepit)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventSavedReceipt,
		sdk.NewAttribute("receipt_id", receipt.UUID),
		sdk.NewAttribute("document_id", receipt.DocumentUUID),
		sdk.NewAttribute("sender", receipt.Sender),
		sdk.NewAttribute("recipient", receipt.Recipient),
	))

	return nil
}

// GetReceiptByID returns the document receipt having the given id, or false if such receipt could not be found
func (keeper Keeper) GetReceiptByID(ctx sdk.Context, id string) (types.DocumentReceipt, error) {
	store := ctx.KVStore(keeper.storeKey)
	key := getReceiptStoreKey(id)

	if !store.Has(key) {
		return types.DocumentReceipt{}, fmt.Errorf("cannot find receipt with uuid %s", id)
	}

	var receipt types.DocumentReceipt
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(key), &receipt)
	return receipt, nil
}

// ExtractReceipt returns a DocumentReceipt slice instance and its UUID given an iterator byte stream value.
func (keeper Keeper) ExtractReceipt(ctx sdk.Context, iterVal []byte) (types.DocumentReceipt, string, error) {
	rid := string(iterVal)

	newReceipt, err := keeper.GetReceiptByID(ctx, rid)
	return newReceipt, rid, err
}
