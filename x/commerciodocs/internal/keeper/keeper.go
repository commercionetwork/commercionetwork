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
func (keeper Keeper) ShareDocument(ctx sdk.Context, document types.Document) sdk.Error {

	store := ctx.KVStore(keeper.StoreKey)

	receiver := document.Recipient.String()
	sender := document.Recipient.String()

	var receiverDocsList []types.Document
	var sentDocsList []types.Document

	//Update receiver documents list
	receivedDocs := store.Get([]byte(ReceivedDocumentsPrefix + receiver))
	keeper.cdc.MustUnmarshalBinaryBare(receivedDocs, &receiverDocsList)
	receiverDocsList = utilities.AppendDocIfMissing(receiverDocsList, document)
	store.Set([]byte(ReceivedDocumentsPrefix+receiver), keeper.cdc.MustMarshalBinaryBare(document))

	//Update sender documents list
	sentDocs := store.Get([]byte(SentDocumentsPrefix + sender))
	keeper.cdc.MustUnmarshalBinaryBare(sentDocs, &sentDocsList)
	sentDocsList = append(sentDocsList, document)
	store.Set([]byte(SentDocumentsPrefix+sender), keeper.cdc.MustMarshalBinaryBare(sentDocsList))

	return nil
}

//Get all the received documents by user
func (keeper Keeper) GetUserReceivedDocuments(ctx sdk.Context, user sdk.AccAddress) []types.Document {

	store := ctx.KVStore(keeper.StoreKey)

	var receivedDocsList []types.Document

	receivedDocs := store.Get([]byte(ReceivedDocumentsPrefix + user.String()))
	keeper.cdc.MustUnmarshalBinaryBare(receivedDocs, receivedDocsList)

	return receivedDocsList
}

//Get all the sent documents by user
func (keeper Keeper) GetUserSentDocuments(ctx sdk.Context, user sdk.AccAddress) []types.Document {
	store := ctx.KVStore(keeper.StoreKey)

	var sentDocsList []types.Document

	sentDocs := store.Get([]byte(SentDocumentsPrefix + user.String()))
	keeper.cdc.MustUnmarshalBinaryBare(sentDocs, sentDocsList)

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
