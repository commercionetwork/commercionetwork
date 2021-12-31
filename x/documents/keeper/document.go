package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	//	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AppendDocument appends a document in the store with a new id and update the count
func (k Keeper) AppendDocument(ctx sdk.Context, document types.Document) string {
	// Create the Document
	store := ctx.KVStore(k.storeKey)
	//store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DocumentKey))
	store.Set(getDocumentStoreKey(document.UUID), k.cdc.MustMarshalBinaryBare(&document))
	return document.UUID
}

// getDocumentStoreKey generates an unique store key for a Document UUID
func getDocumentStoreKey(uuid string) []byte {
	return []byte(types.DocumentStorePrefix + uuid)
}

/*
// SetDocument set a specific document in the store
func (k Keeper) SetDocument(ctx sdk.Context, document types.Document) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DocumentKey))
	b := k.cdc.MustMarshalBinaryBare(&document)
	store.Set(getDocumentStoreKey(document.UUID), b)
}*/

// GetDocumentByID returns the document having the given id
func (k Keeper) GetDocumentByID(ctx sdk.Context, id string) (types.Document, error) {
	store := ctx.KVStore(k.storeKey)
	//store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DocumentKey))
	documentKey := getDocumentStoreKey(id)
	if !store.Has(documentKey) {
		return types.Document{}, fmt.Errorf("cannot find document with uuid %s", id)
	}

	var document types.Document
	k.cdc.MustUnmarshalBinaryBare(store.Get(documentKey), &document)
	return document, nil
}

// HasDocument checks if the Document exists in the store
func (k Keeper) HasDocument(ctx sdk.Context, id string) bool {
	//store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DocumentKey))
	store := ctx.KVStore(k.storeKey)
	documentKey := getDocumentStoreKey(id)
	return store.Has(documentKey)
}

// GetDocumentOwner returns the creator of the document
func (k Keeper) GetDocumentOwner(ctx sdk.Context, id string) string {
	document, _ := k.GetDocumentByID(ctx, id)
	return document.Sender
}

// getSentDocumentsIdsUUIDStoreKey generates a SentDocumentID for a given user and document UUID
func getSentDocumentsIdsUUIDStoreKey(user sdk.AccAddress, documentUUID string) []byte {
	userPart := append(user, []byte(":"+documentUUID)...)
	return append([]byte(types.SentDocumentsPrefix), userPart...)
}

// getReceivedDocumentsIdsUUIDStoreKey generates a ReceivedDocumentID for a given user and document UUID
func getReceivedDocumentsIdsUUIDStoreKey(user sdk.AccAddress, documentUUID string) []byte {
	userPart := append(user, []byte(":"+documentUUID)...)

	return append([]byte(types.ReceivedDocumentsPrefix), userPart...)
}

// getReceivedDocumentsIdsStoreKey generates a ReceivedDocumentsID store key for a given user
func getReceivedDocumentsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.ReceivedDocumentsPrefix), user...)
}

// getSentDocumentsIdsStoreKey generates a ReceivedDocumentsID store key for a given user
func getSentDocumentsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.SentDocumentsPrefix), user...)
}

// getReceiptStoreKey generates a store key for a document UUID
func getReceiptStoreKey(id string) []byte {
	return []byte(types.ReceiptsStorePrefix + id)
}

// getSentReceiptsIdsStoreKey generates a SentReceiptsID store key for a given user
func getSentReceiptsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.SentDocumentsReceiptsPrefix), user...)
}

// getReceivedReceiptsIdsStoreKey generates a ReceivedReceiptsID store key for a given user
func getReceivedReceiptsIdsStoreKey(user sdk.AccAddress) []byte {
	return append([]byte(types.ReceivedDocumentsReceiptsPrefix), user...)
}

// getSentReceiptsIdsUUIDStoreKey generates a SentReceiptsID store key for a given user
func getSentReceiptsIdsUUIDStoreKey(user sdk.AccAddress, recepitUUID string) []byte {
	recepitPart := append(user, []byte(":"+recepitUUID)...)

	return append([]byte(types.SentDocumentsReceiptsPrefix), recepitPart...)
}

// getReceivedReceiptsIdsUUIDStoreKey generates a ReceivedReceiptsID store key for a given user
func getReceivedReceiptsIdsUUIDStoreKey(user sdk.AccAddress, recepitUUID string) []byte {
	recepitPart := append(user, []byte(":"+recepitUUID)...)

	return append([]byte(types.ReceivedDocumentsReceiptsPrefix), recepitPart...)
}
