package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	errorsmod "cosmossdk.io/errors"
	"github.com/gofrs/uuid"
)

// SaveDocument allows the sharing of a document
func (keeper Keeper) SaveDocument(ctx sdk.Context, document types.Document) error {
	// Check the id validity
	if _, err := uuid.FromString(document.UUID); err != nil {
		return errorsmod.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("invalid document UUID: %s", document.UUID))
	}

	store := ctx.KVStore(keeper.storeKey)

	// Check for an existing document
	if store.Has(getDocumentStoreKey(document.UUID)) {
		return errorsmod.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("document with uuid %s already present", document.UUID))
	}

	// Store the document instance
	store.Set(getDocumentStoreKey(document.UUID), keeper.cdc.MustMarshal(&document))

	// Store the document as sent by the sender
	senderAccAdrr, _ := sdk.AccAddressFromBech32(document.Sender)
	sentDocumentsStoreKey := getSentDocumentsIdsUUIDStoreKey(senderAccAdrr, document.UUID)
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
	logger := keeper.Logger(ctx)
	logger.Debug("Document successfully shared")
	return nil
}

// getDocumentStoreKey generates an unique store key for a Document UUID
func getDocumentStoreKey(uuid string) []byte {
	return []byte(types.DocumentStorePrefix + uuid)
}

// GetDocumentByID returns the document having the given id
func (k Keeper) GetDocumentByID(ctx sdk.Context, id string) (types.Document, error) {
	store := ctx.KVStore(k.storeKey)
	documentKey := getDocumentStoreKey(id)
	if !store.Has(documentKey) {
		return types.Document{}, fmt.Errorf("cannot find document with uuid %s", id)
	}

	var document types.Document
	k.cdc.MustUnmarshal(store.Get(documentKey), &document)
	return document, nil
}

// getSentDocumentsIdsUUIDStoreKey generates a SentDocumentID for a given user and document UUID
func getSentDocumentsIdsUUIDStoreKey(user sdk.AccAddress, documentUUID string) []byte {
	userPart := append(user, []byte(":"+documentUUID)...)
	return append([]byte(types.SentDocumentsPrefix), userPart...)
}

// getSentDocumentsIdsStoreKey generates a ReceivedDocumentsID store key for a given user
func getSentDocumentsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.SentDocumentsPrefix), user...)
}

// getReceivedDocumentsIdsStoreKey generates a ReceivedDocumentsID store key for a given user
func getReceivedDocumentsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.ReceivedDocumentsPrefix), user...)
}

// getReceivedDocumentsIdsUUIDStoreKey generates a ReceivedDocumentID for a given user and document UUID
func getReceivedDocumentsIdsUUIDStoreKey(user sdk.AccAddress, documentUUID string) []byte {
	userPart := append(user, []byte(":"+documentUUID)...)

	return append([]byte(types.ReceivedDocumentsPrefix), userPart...)
}

// DocumentsIterator returns an Iterator for all the Documents saved in the store.
func (keeper Keeper) DocumentsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.DocumentStorePrefix))
}

// UserSentDocumentsIterator returns an Iterator for all the sent Documents of a user.
func (keeper Keeper) UserSentDocumentsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)

	return sdk.KVStorePrefixIterator(store, getSentDocumentsIdsStoreKey(user))
}

// UserReceivedDocumentsIterator returns an Iterator for all the received Documents of a user.
func (keeper Keeper) UserReceivedDocumentsIterator(ctx sdk.Context, user sdk.AccAddress) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)

	return sdk.KVStorePrefixIterator(store, getReceivedDocumentsIdsStoreKey(user))
}
