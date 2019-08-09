package keeper

import (
	"github.com/commercionetwork/commercionetwork/types"
	"github.com/commercionetwork/commercionetwork/utilities"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ----------------------------------
// --- Keeper definition
// ----------------------------------

const (
	SentDocumentsPrefix     = "sentBy:"
	ReceivedDocumentsPrefix = "received:"
)

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

// ShareDocument allows the sharing of a document
func (keeper Keeper) ShareDocument(ctx sdk.Context, document types.Document) {

	store := ctx.KVStore(keeper.StoreKey)

	sender := document.Sender.String()
	recipient := document.Recipient.String()

	var sentDocsList, receiverDocsList []types.Document

	// Get the existing received documents
	receivedDocs := store.Get([]byte(ReceivedDocumentsPrefix + recipient))
	keeper.cdc.MustUnmarshalBinaryBare(receivedDocs, &receiverDocsList)

	// Append the new received document
	receiverDocsList = utilities.AppendDocIfMissing(receiverDocsList, document)

	// Save the new list
	store.Set([]byte(ReceivedDocumentsPrefix+recipient), keeper.cdc.MustMarshalBinaryBare(&receiverDocsList))

	// Get the existing sent list
	sentDocs := store.Get([]byte(SentDocumentsPrefix + sender))
	if sentDocs != nil {
		keeper.cdc.MustUnmarshalBinaryBare(sentDocs, &sentDocsList)
	}

	// Append the new sent document
	sentDocsList = utilities.AppendDocIfMissing(sentDocsList, document)

	// Save the new list
	store.Set([]byte(SentDocumentsPrefix+sender), keeper.cdc.MustMarshalBinaryBare(&sentDocsList))
}

//Get all the received documents by user
func (keeper Keeper) GetUserReceivedDocuments(ctx sdk.Context, user sdk.AccAddress) []types.Document {

	store := ctx.KVStore(keeper.StoreKey)
	receivedDocs := store.Get([]byte(ReceivedDocumentsPrefix + user.String()))

	var receivedDocsList []types.Document
	keeper.cdc.MustUnmarshalBinaryBare(receivedDocs, &receivedDocsList)

	return receivedDocsList
}

//Get all the sent documents by user
func (keeper Keeper) GetUserSentDocuments(ctx sdk.Context, user sdk.AccAddress) []types.Document {
	store := ctx.KVStore(keeper.StoreKey)
	sentDocs := store.Get([]byte(SentDocumentsPrefix + user.String()))

	var sentDocsList []types.Document
	keeper.cdc.MustUnmarshalBinaryBare(sentDocs, &sentDocsList)

	return sentDocsList
}

//TODO Implement these functions when it useful

//Get Document associated with checksum given
func (keeper Keeper) GetDocument(ctx sdk.Context, checksumValue string) types.Document {
	return types.Document{}
}

// Get all the documents that given sender has shared with given receiver
func (keeper Keeper) GetSharedDocumentsWithUser(ctx sdk.Context, sender sdk.AccAddress, receiver sdk.AccAddress) []types.Document {
	return []types.Document{}
}
