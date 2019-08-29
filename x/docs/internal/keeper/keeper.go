package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ----------------------------------
// --- Keeper definition
// ----------------------------------

type Keeper struct {
	StoreKey sdk.StoreKey
	cdc      *codec.Codec
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		StoreKey: storeKey,
		cdc:      cdc,
	}
}

// ShareDocument allows the sharing of a document
func (keeper Keeper) ShareDocument(ctx sdk.Context, document types.Document) {

	store := ctx.KVStore(keeper.StoreKey)

	sender := document.Sender.String()
	recipient := document.Recipient.String()

	var sentDocsList, recipientDocsList types.Documents

	// Get the existing received documents
	receivedDocs := store.Get([]byte(types.ReceivedDocumentsPrefix + recipient))
	keeper.cdc.MustUnmarshalBinaryBare(receivedDocs, &recipientDocsList)

	// Append the new received document
	recipientDocsList = recipientDocsList.AppendIfMissing(document)

	// Save the new list
	store.Set([]byte(types.ReceivedDocumentsPrefix+recipient), keeper.cdc.MustMarshalBinaryBare(&recipientDocsList))

	// Get the existing sent list
	sentDocs := store.Get([]byte(types.SentDocumentsPrefix + sender))
	if sentDocs != nil {
		keeper.cdc.MustUnmarshalBinaryBare(sentDocs, &sentDocsList)
	}

	// Append the new sent document
	sentDocsList = sentDocsList.AppendIfMissing(document)

	// Save the new list
	store.Set([]byte(types.SentDocumentsPrefix+sender), keeper.cdc.MustMarshalBinaryBare(&sentDocsList))
}

// GetUserReceivedDocuments returns a list of all the documents that has been received from a user
func (keeper Keeper) GetUserReceivedDocuments(ctx sdk.Context, user sdk.AccAddress) types.Documents {

	store := ctx.KVStore(keeper.StoreKey)
	receivedDocs := store.Get([]byte(types.ReceivedDocumentsPrefix + user.String()))

	var receivedDocsList types.Documents
	keeper.cdc.MustUnmarshalBinaryBare(receivedDocs, &receivedDocsList)

	return receivedDocsList
}

//GetUserSentDocuments returns a list of all documents sent by user
func (keeper Keeper) GetUserSentDocuments(ctx sdk.Context, user sdk.AccAddress) types.Documents {
	store := ctx.KVStore(keeper.StoreKey)
	sentDocs := store.Get([]byte(types.SentDocumentsPrefix + user.String()))

	var sentDocsList types.Documents
	keeper.cdc.MustUnmarshalBinaryBare(sentDocs, &sentDocsList)

	return sentDocsList
}

//TODO Implement these functions when it useful

//Get Document associated with UUID given
func (keeper Keeper) GetDocument(ctx sdk.Context, uuid string) types.Document {
	return types.Document{}
}

// Get all the documents that given sender has shared with given recipient
func (keeper Keeper) GetSharedDocumentsWithUser(ctx sdk.Context, sender sdk.AccAddress, recipient sdk.AccAddress) types.Documents {
	return types.Documents{}
}

// ShareDocumentReceipt allows to properly store the given receipt
func (keeper Keeper) ShareDocumentReceipt(ctx sdk.Context, receipt types.DocumentReceipt) {
	store := ctx.KVStore(keeper.StoreKey)

	//Check if the receipt is already in the store
	receiptBz := store.Get([]byte(types.DocumentReceiptPrefix + receipt.Uuid + receipt.Recipient.String()))
	//if it's not, insert it
	if receiptBz == nil {
		store.Set([]byte(types.DocumentReceiptPrefix+receipt.Uuid+receipt.Recipient.String()),
			keeper.cdc.MustMarshalBinaryBare(&receipt))
	}
}

// GetUserReceivedReceipts returns the list of all the receipts that the given user has received
func (keeper Keeper) GetUserReceivedReceipts(ctx sdk.Context, user sdk.AccAddress) types.DocumentReceipts {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.DocumentReceiptPrefix))

	receivedReceipts := types.DocumentReceipts{}
	var receipt types.DocumentReceipt
	for ; iterator.Valid(); iterator.Next() {
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &receipt)
		if receipt.Recipient.Equals(user) {
			receivedReceipts = receivedReceipts.AppendReceiptIfMissing(receipt)
		}
	}
	return receivedReceipts
}

// GetReceiptByDocumentUuid returns the receipts that the given recipient has received for the document having the
// given uuid
func (keeper Keeper) GetReceiptByDocumentUuid(ctx sdk.Context, recipient sdk.AccAddress, uuid string) types.DocumentReceipt {
	store := ctx.KVStore(keeper.StoreKey)

	var receipt types.DocumentReceipt
	receiptBz := store.Get([]byte(types.DocumentReceiptPrefix + uuid + recipient.String()))
	keeper.cdc.MustUnmarshalBinaryBare(receiptBz, &receipt)

	return receipt
}
