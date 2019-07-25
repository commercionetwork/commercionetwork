package keeper

import (
	"commercio-network/types"
	"commercio-network/utilities"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ----------------------------------
// --- Keeper definition
// ----------------------------------

type Keeper struct {
	// Key of the map { Did => Ddo }
	identitiesStoreKey sdk.StoreKey

	// Key of the map { Address => []Did }
	ownersStoresKey sdk.StoreKey

	// Key of the map { Did => []Did }
	connectionsStoreKey sdk.StoreKey

	// Pointer to the codec that is used by Amino to encode and decode binary structs.
	Cdc *codec.Codec
}

// NewKeeper creates new instances of the CommercioID Keeper
func NewKeeper(
	identitiesStoreKey sdk.StoreKey,
	ownersStoresKey sdk.StoreKey,
	connectionsStoreKey sdk.StoreKey,
	cdc *codec.Codec) Keeper {
	return Keeper{
		identitiesStoreKey:  identitiesStoreKey,
		ownersStoresKey:     ownersStoresKey,
		connectionsStoreKey: connectionsStoreKey,
		Cdc:                 cdc,
	}
}

// ----------------------------------
// --- Keeper methods
// ----------------------------------

// CreateIdentity creates a new identity for the given owner.
// The identity is represented as an association between the given Did and the respective Did Document reference.
// The data related to the identity should be placed inside the Did Document that can be reached following the proper
// reference.
func (keeper Keeper) CreateIdentity(ctx sdk.Context, owner sdk.AccAddress, did types.Did, ddoReference string) {
	// Store the Did => Ddo entry
	identitiesStore := ctx.KVStore(keeper.identitiesStoreKey)
	identitiesStore.Set([]byte(did), []byte(ddoReference))

	// --- Store the Address => IdentityReference entry ---
	// Get the store
	ownersStore := ctx.KVStore(keeper.ownersStoresKey)

	// Read the list of the dids associated to a given address
	existingReferences := ownersStore.Get(owner)

	// If the value exists, read it. Otherwise create an empty array
	var dids []types.Did
	if existingReferences != nil {
		keeper.Cdc.MustUnmarshalBinaryBare(existingReferences, &dids)
	}

	// Save the new Did inside the array
	dids = utilities.AppendDidIfMissing(dids, did)

	// Store the array back to the blockchain
	ownersStore.Set(owner, keeper.Cdc.MustMarshalBinaryBare(&dids))
}

// GetIdentity returns the Did Document reference associated to a given Did.
// If the given Did has no Did Document reference associated, returns nil.
func (keeper Keeper) GetDdoReferenceByDid(ctx sdk.Context, did types.Did) string {
	store := ctx.KVStore(keeper.identitiesStoreKey)
	result := store.Get([]byte(did))
	return string(result)
}

// IsOwner tells whenever the given AccAddress owns the identity associated with the given Did.
// If the address owns the identity returns true, otherwise returns false.
func (keeper Keeper) CanBeUsedBy(ctx sdk.Context, owner sdk.AccAddress, did types.Did) bool {
	identitiesStore := ctx.KVStore(keeper.identitiesStoreKey)
	existingIdentity := identitiesStore.Get([]byte(did))

	// If the identity has not yet been registered, everyone can use it
	if existingIdentity == nil {
		return true
	}

	ownersStore := ctx.KVStore(keeper.ownersStoresKey)
	result := ownersStore.Get(owner)

	// If the owner has no identities, he can't use this one
	if result == nil {
		return false
	}

	var dids []types.Did
	keeper.Cdc.MustUnmarshalBinaryBare(result, &dids)

	// If the owner has some identities, check if this one is inside the ones he has registered
	return utilities.DidInSlice(did, dids)
}

func addConnectionToUser(user types.Did, connection types.Did, store sdk.KVStore, codec *codec.Codec) {
	var userConnections []types.Did

	existingUserConnections := store.Get([]byte(user))
	if existingUserConnections != nil {
		codec.MustUnmarshalBinaryBare(existingUserConnections, &userConnections)
	}

	userConnections = utilities.AppendDidIfMissing(userConnections, connection)

	store.Set([]byte(user), codec.MustMarshalBinaryBare(&userConnections))
}

// AddConnection adds a connection between the first given user and the second given user.
// If a connection between the two users already exists, nothing is added
func (keeper Keeper) AddConnection(ctx sdk.Context, firstDid types.Did, secondDid types.Did) {
	connectionsStore := ctx.KVStore(keeper.connectionsStoreKey)

	addConnectionToUser(firstDid, secondDid, connectionsStore, keeper.Cdc)
	addConnectionToUser(secondDid, firstDid, connectionsStore, keeper.Cdc)
}

// GetConnections returns all the connections that a given did has.
// If no connection is found, an empty list is returned.
func (keeper Keeper) GetConnections(ctx sdk.Context, did types.Did) []types.Did {
	store := ctx.KVStore(keeper.connectionsStoreKey)

	var connections []types.Did

	result := store.Get([]byte(did))
	if result != nil {
		keeper.Cdc.MustUnmarshalBinaryBare(result, &connections)
	}

	return connections
}
