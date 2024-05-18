package keeper

import (
	"fmt"

	errors "cosmossdk.io/errors"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	storetypes "cosmossdk.io/store/types"
	"github.com/gofrs/uuid"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SaveReceipt allows to properly store the given receipt
func (keeper Keeper) SaveReceipt(ctx sdk.Context, receipt types.DocumentReceipt) error {

	// Check the id validity
	if _, err := uuid.FromString(receipt.UUID); err != nil {
		return errors.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("invalid document receipt UUID: %s", receipt.UUID))
	}
	store := ctx.KVStore(keeper.storeKey)

	// Check for an existing document
	if !store.Has(getDocumentStoreKey(receipt.DocumentUUID)) {
		return errors.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("receipt points to a non-existing document UUID: %s", receipt.DocumentUUID))
	}

	marshaledReceiptID := []byte(receipt.UUID)
	receiptStoreKey := getReceiptStoreKey(receipt.UUID)

	// Check for usage of same ID
	if store.Has(receiptStoreKey) {
		return errors.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("receipt for document with UUID %s already present: %s", receipt.DocumentUUID, receipt.UUID))
	}

	senderAccAdrr, _ := sdk.AccAddressFromBech32(receipt.Sender)

	// check if Sender is included among the recipients of the document
	if !store.Has(getReceivedDocumentsIdsUUIDStoreKey(senderAccAdrr, receipt.DocumentUUID)) {
		return errors.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("sender for document receipt with address %s not among the recipients of the document: %s", receipt.Sender, receipt.UUID))
	}

	sentReceiptsIdsStoreKey := getSentReceiptsIdsUUIDStoreKey(senderAccAdrr, receipt.DocumentUUID)
	if store.Has(sentReceiptsIdsStoreKey) {
		return errors.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("sent receipt for document with UUID %s already present: %s", receipt.DocumentUUID, receipt.UUID))
	}

	recipientAccAdrr, _ := sdk.AccAddressFromBech32(receipt.Recipient)
	receivedReceiptIdsStoreKey := getReceivedReceiptsIdsUUIDStoreKey(recipientAccAdrr, receipt.UUID)

	documentsReceiptsIdsStoreKey := getDocumentReceiptsIdsUUIDStoreKey(receipt.DocumentUUID, receipt.UUID)

	// Store the receipt
	marshaledReceipt := keeper.cdc.MustMarshal(&receipt)
	store.Set(receiptStoreKey, marshaledReceipt)
	// Store the receipt ID as sent
	store.Set(sentReceiptsIdsStoreKey, marshaledReceiptID)
	// Store the receipt ID as received
	store.Set(receivedReceiptIdsStoreKey, marshaledReceiptID)
	// Store the receipt ID along the receipts for the document
	store.Set(documentsReceiptsIdsStoreKey, marshaledReceiptID)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		eventSavedReceipt,
		sdk.NewAttribute("receipt_id", receipt.UUID),
		sdk.NewAttribute("document_id", receipt.DocumentUUID),
		sdk.NewAttribute("sender", receipt.Sender),
		sdk.NewAttribute("recipient", receipt.Recipient),
	))
	logger := keeper.Logger(ctx)
	logger.Debug("Receipt Document successfully sent")
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
	keeper.cdc.MustUnmarshal(store.Get(key), &receipt)
	return receipt, nil
}

// getReceiptStoreKey generates a store key for a document UUID
func getReceiptStoreKey(id string) []byte {
	return []byte(types.ReceiptsStorePrefix + id)
}

// getSentReceiptsIdsStoreKey generates a SentReceiptsID store key for a given user
func getSentReceiptsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.SentDocumentsReceiptsPrefix), user...)
}

// getSentReceiptsIdsUUIDStoreKey generates a SentReceiptsID store key for a given user
func getSentReceiptsIdsUUIDStoreKey(user sdk.AccAddress, receiptUUID string) []byte {
	receiptPart := append(user, []byte(":"+receiptUUID)...)

	return append([]byte(types.SentDocumentsReceiptsPrefix), receiptPart...)
}

// getReceivedReceiptsIdsStoreKey generates a ReceivedReceiptsID store key for a given user
func getReceivedReceiptsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.ReceivedDocumentsReceiptsPrefix), user...)
}

// getReceivedReceiptsIdsUUIDStoreKey generates a ReceivedReceiptsID store key for a given user
func getReceivedReceiptsIdsUUIDStoreKey(user sdk.AccAddress, receiptUUID string) []byte {
	receiptPart := append(user, []byte(":"+receiptUUID)...)

	return append([]byte(types.ReceivedDocumentsReceiptsPrefix), receiptPart...)
}

// getDocumentReceiptsIdsStoreKey generates a ReceivedReceiptsID store key for a given user
func getDocumentReceiptsIdsStoreKey(documentUUID string) []byte {
	return append([]byte(types.DocumentsReceiptsPrefix), []byte(documentUUID)...)
}

// getDocumentReceiptsIdsUUIDStoreKey generates a DocumentReceiptsID store key for a given document
func getDocumentReceiptsIdsUUIDStoreKey(documentUUID string, receiptUUID string) []byte {
	receiptPart := append([]byte(documentUUID), []byte(":"+receiptUUID)...)

	return append([]byte(types.DocumentsReceiptsPrefix), receiptPart...)
}

// DocumentReceiptsIterator returns an Iterator for all Receipts saved in the store.
func (keeper Keeper) DocumentReceiptsIterator(ctx sdk.Context) storetypes.Iterator {
	store := ctx.KVStore(keeper.storeKey)

	return storetypes.KVStorePrefixIterator(store, []byte(types.ReceiptsStorePrefix))
}

// UserSentReceiptsIterator returns an Iterator for all the Document Sent Receipts for a user.
func (keeper Keeper) UserSentReceiptsIterator(ctx sdk.Context, user sdk.AccAddress) storetypes.Iterator {
	store := ctx.KVStore(keeper.storeKey)

	return storetypes.KVStorePrefixIterator(store, getSentReceiptsIdsStoreKey(user))
}

// UserReceivedReceiptsIterator returns an Iterator for all the Document Received Receipts for a user.
func (keeper Keeper) UserReceivedReceiptsIterator(ctx sdk.Context, user sdk.AccAddress) storetypes.Iterator {
	store := ctx.KVStore(keeper.storeKey)

	return storetypes.KVStorePrefixIterator(store, getReceivedReceiptsIdsStoreKey(user))
}

// UserReceivedReceiptsIterator returns an Iterator for all the Receipts for a Document.
func (keeper Keeper) UUIDDocumentsReceiptsIterator(ctx sdk.Context, documentUUID string) storetypes.Iterator {
	store := ctx.KVStore(keeper.storeKey)

	return storetypes.KVStorePrefixIterator(store, getDocumentReceiptsIdsStoreKey(documentUUID))
}
