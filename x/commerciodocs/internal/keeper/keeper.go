package keeper

import (
	"bytes"
	"fmt"
	"github.com/commercionetwork/commercionetwork/types"
	"github.com/commercionetwork/commercionetwork/x/commercioid"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/genproto/googleapis/type/date"
)

// ----------------------------------
// --- Keeper definition
// ----------------------------------

type Keeper struct {
	CommercioIdKeeper commercioid.Keeper

	// Key of the map { DocumentReference => Address }
	ownersStoreKey sdk.StoreKey

	// Key of the map { DocumentReference => Metadata }
	metadataStoreKey sdk.StoreKey

	// Key of the map { Did => []Sharing }
	sharingStoreKey sdk.StoreKey

	// Key of the map { DocumentReference => []Did }
	readersStoreKey sdk.StoreKey

	cdc *codec.Codec
}

func NewKeeper(
	commercioIdKeeper commercioid.Keeper,
	ownersStoreKey sdk.StoreKey,
	metadataStoreKey sdk.StoreKey,
	sharingStoreKey sdk.StoreKey,
	readersStoreKey sdk.StoreKey,
	cdc *codec.Codec) Keeper {
	return Keeper{
		CommercioIdKeeper: commercioIdKeeper,
		ownersStoreKey:    ownersStoreKey,
		metadataStoreKey:  metadataStoreKey,
		sharingStoreKey:   sharingStoreKey,
		readersStoreKey:   readersStoreKey,
		cdc:               cdc,
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

func (keeper Keeper) addReaderForDocument(ctx sdk.Context, user types.Did, reference string) {
	store := ctx.KVStore(keeper.readersStoreKey)

	// Get the readers for the document by reading the store
	var readers []types.Did
	existingReaders := store.Get([]byte(reference))
	if existingReaders != nil {
		keeper.cdc.MustUnmarshalBinaryBare(existingReaders, &readers)
	}

	// Append the user to the reading list
	//readers = utilities.AppendAddressIfMissing(readers, user)

	// Save the result into the store
	store.Set([]byte(reference), keeper.cdc.MustMarshalBinaryBare(&readers))
}

// StoreDocument stores the given document Reference assigning it to the given Owner, so that we know who created it.
func (keeper Keeper) StoreDocument(ctx sdk.Context, owner sdk.AccAddress, identity types.Did, reference string, metadata string) {
	// Save the document Owner
	ownersStore := ctx.KVStore(keeper.ownersStoreKey)
	ownersStore.Set([]byte(reference), owner)

	// Save the Metadata
	metadataStore := ctx.KVStore(keeper.metadataStoreKey)
	metadataStore.Set([]byte(reference), []byte(metadata))

	// Set the creator as one of the readers
	keeper.addReaderForDocument(ctx, identity, reference)
}

// HasOwners tells whenever the document with the given Reference as an Owner or not.
// Returns true iff the document already has an Owner, false otherwise.
func (keeper Keeper) HasOwner(ctx sdk.Context, reference string) bool {
	store := ctx.KVStore(keeper.ownersStoreKey)
	result := store.Get([]byte(reference))
	return result != nil
}

// IsOwner tells whenever the given Address is the Owner of the document or not.
func (keeper Keeper) IsOwner(ctx sdk.Context, owner sdk.AccAddress, reference string) bool {
	store := ctx.KVStore(keeper.ownersStoreKey)
	existingOwner := store.Get([]byte(reference))
	return bytes.Equal(existingOwner, owner)
}

// GetMetadata returns the Metadata Reference for the document with the given Reference.
func (keeper Keeper) GetMetadata(ctx sdk.Context, reference string) string {
	store := ctx.KVStore(keeper.metadataStoreKey)
	result := store.Get([]byte(reference))
	return string(result)
}

// ShareDocument allows the sharing of a document represented by the given Reference, between the given sender and the
// given recipient.
func (keeper Keeper) ShareDocument(ctx sdk.Context, reference string, sender types.Did, recipient types.Did) sdk.Error {
	sharing := Sharing{
		Sender:   sender,
		Receiver: recipient,
		Date:     date.Date{},
		Document: reference,
	}

	sharingStore := ctx.KVStore(keeper.sharingStoreKey)

	// Save the shared document for the sender
	if keeper.CanReadDocument(ctx, sender, reference) {
		sharingStore.Set([]byte(sender), keeper.cdc.MustMarshalBinaryBare(sharing))
	} else {
		return sdk.ErrUnauthorized(fmt.Sprintf("The sender with TestAddress %s doesnt have the rights on this document", sender))
	}

	// Save the shared document for the recipient
	if !keeper.CanReadDocument(ctx, recipient, reference) {
		sharingStore.Set([]byte(recipient), keeper.cdc.MustMarshalBinaryBare(sharing))
		keeper.addReaderForDocument(ctx, recipient, reference)
	}

	return nil
}

// CanReadDocument tells whenever a given user has access to a document or not.
// Returns true if the user has access to the document, false otherwise.
func (keeper Keeper) CanReadDocument(ctx sdk.Context, user types.Did, reference string) bool {
	store := ctx.KVStore(keeper.readersStoreKey)

	var readers []types.Did

	existingReaders := store.Get([]byte(reference))
	if existingReaders != nil {
		keeper.cdc.MustUnmarshalBinaryBare(existingReaders, &readers)
	}

	//return utilities.DidInSlice(user, readers)
	return false
}

// GetAuthorizedReaders lists all the users, represented by their identity, that have access to the document.
// This includes the creator of the document as well all the users to which the creator has shared the document itself.
func (keeper Keeper) GetAuthorizedReaders(ctx sdk.Context, reference string) []types.Did {
	store := ctx.KVStore(keeper.readersStoreKey)

	var readers []types.Did
	existingReaders := store.Get([]byte(reference))
	if existingReaders != nil {
		keeper.cdc.MustUnmarshalBinaryBare(existingReaders, &readers)
	}

	return readers
}
