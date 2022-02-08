package keeper

import (
	"fmt"
	"strings"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SaveReceipt allows to properly store the given receipt
func (keeper Keeper) SaveReceipt(ctx sdk.Context, receipt types.DocumentReceipt) error {
	// TODO: change to UUID validation
	// Check the id
	if len(strings.TrimSpace(receipt.UUID)) == 0 {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("invalid document receipt id: %s", receipt.UUID))
	}

	if _, err := keeper.GetDocumentByID(ctx, receipt.DocumentUUID); err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "receipt points to a non-existing document UUID")
	}

	store := ctx.KVStore(keeper.storeKey)
	senderAccadrr, _ := sdk.AccAddressFromBech32(receipt.Sender)
	sentReceiptsIdsStoreKey := getSentReceiptsIdsUUIDStoreKey(senderAccadrr, receipt.DocumentUUID)
	recipientAccAdrr, _ := sdk.AccAddressFromBech32(receipt.Recipient)
	receivedReceiptIdsStoreKey := getReceivedReceiptsIdsUUIDStoreKey(recipientAccAdrr, receipt.DocumentUUID)

	marshaledReceipt := keeper.cdc.MustMarshalBinaryBare(&receipt)
	marshaledReceiptID := []byte(receipt.UUID)

	// Store the receipt as sent
	if store.Has(sentReceiptsIdsStoreKey) {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("sent receipt for document with UUID %s already present: %s", receipt.DocumentUUID, receipt.UUID))
	}
	store.Set(sentReceiptsIdsStoreKey, marshaledReceiptID)

	// Store the receipt as received
	if store.Has(receivedReceiptIdsStoreKey) {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("received receipt for document with UUID %s already present: %s", receipt.DocumentUUID, receipt.UUID))
	}
	store.Set(receivedReceiptIdsStoreKey, marshaledReceiptID)

	// Store the receipt
	store.Set(getReceiptStoreKey(receipt.UUID), marshaledReceipt)

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

// TODO: change documentation or use directly GetReceiptByID instead of this method
// ExtractReceipt returns a DocumentReceipt slice instance and its UUID given an iterator byte stream value.
func (keeper Keeper) ExtractReceipt(ctx sdk.Context, iterVal []byte) (types.DocumentReceipt, string, error) {
	rid := string(iterVal)

	newReceipt, err := keeper.GetReceiptByID(ctx, rid)
	return newReceipt, rid, err
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

// ReceiptsIterator returns an Iterator for Sent and Received receipts.
func (keeper Keeper) ReceiptsIterators(ctx sdk.Context) (sdk.Iterator, sdk.Iterator) {
	store := ctx.KVStore(keeper.storeKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.SentDocumentsReceiptsPrefix)),
		sdk.KVStorePrefixIterator(store, []byte(types.ReceivedDocumentsReceiptsPrefix))
}

// UserSentReceiptsIterator returns an Iterator for all the Document Sent Receipts for a user.
func (keeper Keeper) UserSentReceiptsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)

	return sdk.KVStorePrefixIterator(store, getSentReceiptsIdsStoreKey(user))
}

// UserReceivedReceiptsIterator returns an Iterator for all the Document Received Receipts for a user.
func (keeper Keeper) UserReceivedReceiptsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)

	return sdk.KVStorePrefixIterator(store, getReceivedReceiptsIdsStoreKey(user))
}
