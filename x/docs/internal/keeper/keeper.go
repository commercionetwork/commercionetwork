package keeper

import (
	"github.com/commercionetwork/commercionetwork/types"
	"github.com/commercionetwork/commercionetwork/utilities"
	keys "github.com/commercionetwork/commercionetwork/x/docs/internal/types"
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

// ----------------------------------
// --- Keeper methods
// ----------------------------------

// ----------------------------------
// --- ShareDocument
// ----------------------------------

// ShareDocument allows the sharing of a document
func (keeper Keeper) ShareDocument(ctx sdk.Context, document types.Document) {

	store := ctx.KVStore(keeper.StoreKey)

	sender := document.Sender.String()
	recipient := document.Recipient.String()

	var sentDocsList, recipientDocsList []types.Document

	// Get the existing received documents
	receivedDocs := store.Get([]byte(keys.ReceivedDocumentsPrefix + recipient))
	keeper.cdc.MustUnmarshalBinaryBare(receivedDocs, &recipientDocsList)

	// Append the new received document
	recipientDocsList = utilities.AppendDocIfMissing(recipientDocsList, document)

	// Save the new list
	store.Set([]byte(keys.ReceivedDocumentsPrefix+recipient), keeper.cdc.MustMarshalBinaryBare(&recipientDocsList))

	// Get the existing sent list
	sentDocs := store.Get([]byte(keys.SentDocumentsPrefix + sender))
	if sentDocs != nil {
		keeper.cdc.MustUnmarshalBinaryBare(sentDocs, &sentDocsList)
	}

	// Append the new sent document
	sentDocsList = utilities.AppendDocIfMissing(sentDocsList, document)

	// Save the new list
	store.Set([]byte(keys.SentDocumentsPrefix+sender), keeper.cdc.MustMarshalBinaryBare(&sentDocsList))
}

// GetUserReceivedDocuments returns a list of all the documents that has been received from a user
func (keeper Keeper) GetUserReceivedDocuments(ctx sdk.Context, user sdk.AccAddress) []types.Document {

	store := ctx.KVStore(keeper.StoreKey)
	receivedDocs := store.Get([]byte(keys.ReceivedDocumentsPrefix + user.String()))

	var receivedDocsList []types.Document
	keeper.cdc.MustUnmarshalBinaryBare(receivedDocs, &receivedDocsList)

	return receivedDocsList
}

//GetUserSentDocuments returns a list of all documents sent by user
func (keeper Keeper) GetUserSentDocuments(ctx sdk.Context, user sdk.AccAddress) []types.Document {
	store := ctx.KVStore(keeper.StoreKey)
	sentDocs := store.Get([]byte(keys.SentDocumentsPrefix + user.String()))

	var sentDocsList []types.Document
	keeper.cdc.MustUnmarshalBinaryBare(sentDocs, &sentDocsList)

	return sentDocsList
}

//TODO Implement these functions when it useful

//Get Document associated with UUID given
func (keeper Keeper) GetDocument(ctx sdk.Context, uuid string) types.Document {
	return types.Document{}
}

// Get all the documents that given sender has shared with given recipient
func (keeper Keeper) GetSharedDocumentsWithUser(ctx sdk.Context, sender sdk.AccAddress, recipient sdk.AccAddress) []types.Document {
	return []types.Document{}
}

// ----------------------------------
// --- DocumentReceipt
// ----------------------------------

//Share the receipt with the recipient inside it
func (keeper Keeper) ShareDocumentReceipt(ctx sdk.Context, receipt types.DocumentReceipt) {
	store := ctx.KVStore(keeper.StoreKey)

	//Check if the receipt is already in the store
	receiptBz := store.Get([]byte(keys.DocumentReceiptPrefix + receipt.Uuid + receipt.Recipient.String()))
	//if it's not, insert it
	if receiptBz == nil {
		store.Set([]byte(keys.DocumentReceiptPrefix+receipt.Uuid+receipt.Recipient.String()),
			keeper.cdc.MustMarshalBinaryBare(&receipt))
	}
}

func (keeper Keeper) GetUserReceivedReceipts(ctx sdk.Context, user sdk.AccAddress) []types.DocumentReceipt {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(keys.DocumentReceiptPrefix))

	var receivedReceipts []types.DocumentReceipt
	var receipt types.DocumentReceipt
	for ; iterator.Valid(); iterator.Next() {
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &receipt)
		if receipt.Recipient.Equals(user) {
			receivedReceipts = utilities.AppendReceiptIfMissing(receivedReceipts, receipt)
		}
	}
	return receivedReceipts
}

func (keeper Keeper) GetReceiptByDocumentUuid(ctx sdk.Context, recipient sdk.AccAddress, uuid string) types.DocumentReceipt {
	store := ctx.KVStore(keeper.StoreKey)

	var receipt types.DocumentReceipt
	receiptBz := store.Get([]byte(keys.DocumentReceiptPrefix + uuid + recipient.String()))
	keeper.cdc.MustUnmarshalBinaryBare(receiptBz, &receipt)

	return receipt
}
