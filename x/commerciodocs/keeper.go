package commerciodocs

import (
	"commercio-network/types"
	"commercio-network/utilities"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/genproto/googleapis/type/date"
)

// ----------------------------------
// --- Keeper definition
// ----------------------------------

type Keeper struct {
	// Key of the map { Address => []DocReference }
	docsStoreKey sdk.StoreKey

	// Key of the map { DocumentReference => Metadata }
	metadataStoreKey sdk.StoreKey

	// Key of the map { Address => []Sharing }
	sharingStoreKey sdk.StoreKey

	cdc *codec.Codec
}

func NewKeeper(
	docsStoreKey sdk.StoreKey,
	metadataStoreKey sdk.StoreKey,
	sharingStoreKey sdk.StoreKey,
	cdc *codec.Codec) Keeper {
	return Keeper{
		docsStoreKey:     docsStoreKey,
		metadataStoreKey: metadataStoreKey,
		sharingStoreKey:  sharingStoreKey,
		cdc:              cdc,
	}
}

type Sharing struct {
	Sender   types.Did
	Receiver types.Did
	Date     date.Date
	Document string
}

// ----------------------------------
// --- Keeper methods
// ----------------------------------

// StoreDocument stores the given document reference assigning it to the given owner, so that we know who created it
func (keeper Keeper) StoreDocument(ctx sdk.Context, reference string, owner sdk.AccAddress) {
	store := ctx.KVStore(keeper.docsStoreKey)

	var documents []string

	// Read the existing documents for the given user, if there are any
	existingDocuments := store.Get(owner)
	if existingDocuments != nil {
		keeper.cdc.MustUnmarshalBinaryBare(existingDocuments, &documents)
	}

	// Append the given document to the list, if not already present
	documents = utilities.AppendStringIfMissing(documents, reference)

	// Save the documents back to the blockchain
	store.Set(owner, keeper.cdc.MustMarshalBinaryBare(documents))
}

// IsOwner tells whenever the given address is the owner of the document or not
func (keeper Keeper) IsOwner(ctx sdk.Context, owner sdk.AccAddress, reference string) bool {
	store := ctx.KVStore(keeper.docsStoreKey)

	existingDocuments := store.Get(owner)
	if existingDocuments == nil {
		return false
	}

	var documents []string
	keeper.cdc.MustUnmarshalBinaryBare(existingDocuments, &documents)

	return utilities.StringInSlice(reference, documents)
}

func getUserSharing(sharingStore sdk.KVStore, sender types.Did, keeper Keeper) []Sharing {
	var userSharing []Sharing
	existingSharing := sharingStore.Get([]byte(sender))
	if existingSharing != nil {
		keeper.cdc.MustUnmarshalBinaryBare(existingSharing, userSharing)
	}
	return userSharing
}

func isPresentInsideSharing(documentReference string, list []Sharing) bool {
	for _, ele := range list {
		if ele.Document == documentReference {
			return true
		}
	}
	return false
}

// ShareDocument allows the sharing of a document represented by the given reference, between the given sender and the
// given recipient
func (keeper Keeper) ShareDocument(ctx sdk.Context, reference string, sender types.Did, recipient types.Did) {
	sharing := Sharing{
		Sender:   sender,
		Receiver: recipient,
		Date:     date.Date{},
		Document: reference,
	}

	sharingStore := ctx.KVStore(keeper.sharingStoreKey)

	// Save the shared document for the sender
	if !keeper.CanReadDocument(ctx, sender, reference) {
		sharingStore.Set([]byte(sender), keeper.cdc.MustMarshalBinaryBare(sharing))
	}

	// Save the shared document for the recipient
	if !keeper.CanReadDocument(ctx, recipient, reference) {
		sharingStore.Set([]byte(recipient), keeper.cdc.MustMarshalBinaryBare(sharing))
	}
}

// CanReadDocument tells whenever a given user has access to a document or not.
// Returns true if the user has access to the document, false otherwise
func (keeper Keeper) CanReadDocument(ctx sdk.Context, user types.Did, reference string) bool {
	store := ctx.KVStore(keeper.sharingStoreKey)
	return isPresentInsideSharing(reference, getUserSharing(store, user, keeper))
}
