package keeper

import (
	"fmt"
	"strings"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	StoreKey sdk.StoreKey

	// Pointer to the codec that is used by Amino to encode and decode binary structs.
	Cdc *codec.Codec
}

// NewKeeper creates new instances of the CommercioID Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		StoreKey: storeKey,
		Cdc:      cdc,
	}
}

// ------------------
// --- Identities
// ------------------

func (keeper Keeper) getIdentityStoreKey(owner sdk.AccAddress) []byte {
	return []byte(types.IdentitiesStorePrefix + owner.String())
}

// SaveIdentity saves the given didDocumentUri associating it with the given owner, replacing any existent one.
func (keeper Keeper) SaveIdentity(ctx sdk.Context, owner sdk.AccAddress, document types.DidDocument) {
	identitiesStore := ctx.KVStore(keeper.StoreKey)
	identitiesStore.Set(keeper.getIdentityStoreKey(owner), keeper.Cdc.MustMarshalBinaryBare(document))
}

// GetDidDocumentByOwner returns the Did Document reference associated to a given Did
func (keeper Keeper) GetDidDocumentByOwner(ctx sdk.Context, owner sdk.AccAddress) (didDocument types.DidDocument, found bool) {
	store := ctx.KVStore(keeper.StoreKey)

	identityKey := keeper.getIdentityStoreKey(owner)
	if !store.Has(identityKey) {
		return types.DidDocument{}, false
	}

	keeper.Cdc.MustUnmarshalBinaryBare(store.Get(identityKey), &didDocument)
	return didDocument, true
}

// GetIdentities returns the list of all identities for the given context
func (keeper Keeper) GetIdentities(ctx sdk.Context) ([]types.Identity, error) {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.IdentitiesStorePrefix))

	var identities []types.Identity
	for ; iterator.Valid(); iterator.Next() {
		addressString := strings.ReplaceAll(string(iterator.Key()), types.IdentitiesStorePrefix, "")
		address, err := sdk.AccAddressFromBech32(addressString)
		if err != nil {
			return nil, err
		}

		var didDocument types.DidDocument
		keeper.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &didDocument)

		identity := types.Identity{
			Owner:       address,
			DidDocument: didDocument,
		}
		identities = append(identities, identity)
	}

	return identities, nil
}

// SetIdentities allows to bulk save a bunch of identities inside the store
func (keeper Keeper) SetIdentities(ctx sdk.Context, identities []types.Identity) {
	for _, identity := range identities {
		keeper.SaveIdentity(ctx, identity.Owner, identity.DidDocument)
	}
}

// ----------------------------
// --- Did deposit requests
// ----------------------------

func (keeper Keeper) getDepositRequestStoreKey(proof string) []byte {
	return []byte(types.DidDepositRequestStorePrefix + proof)
}

func (keeper Keeper) StoreDidDepositRequest(ctx sdk.Context, request types.DidDepositRequest) sdk.Error {
	store := ctx.KVStore(keeper.StoreKey)

	requestKey := keeper.getDepositRequestStoreKey(request.Proof)
	if store.Has(requestKey) {
		return sdk.ErrUnknownRequest("Did deposit request with the same proof already exists")
	}

	// Set the initial status to null and store the request
	request.Status = nil
	store.Set(requestKey, keeper.Cdc.MustMarshalBinaryBare(&request))

	return nil
}

func (keeper Keeper) GetDidDepositRequestByProof(ctx sdk.Context, proof string) (request types.DidDepositRequest, found bool) {
	store := ctx.KVStore(keeper.StoreKey)

	requestKey := keeper.getDepositRequestStoreKey(proof)
	if !store.Has(requestKey) {
		return types.DidDepositRequest{}, false
	}

	keeper.Cdc.MustUnmarshalBinaryBare(store.Get(requestKey), &request)
	return request, true
}

func (keeper Keeper) ChangeDepositRequestStatus(ctx sdk.Context, proof string, status types.DidDepositRequestStatus) sdk.Error {
	store := ctx.KVStore(keeper.StoreKey)

	request, found := keeper.GetDidDepositRequestByProof(ctx, proof)
	if !found {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Did deposit request with proof %s not fond", proof))
	}

	// Update and store the request
	request.Status = &status
	store.Set(keeper.getDepositRequestStoreKey(request.Proof), keeper.Cdc.MustMarshalBinaryBare(request))

	return nil
}

// GetDepositRequests returns the list of the deposit requests existing inside the given context
func (keeper Keeper) GetDepositRequests(ctx sdk.Context) (requests []types.DidDepositRequest) {
	store := ctx.KVStore(keeper.StoreKey)

	iterator := sdk.KVStorePrefixIterator(store, []byte(types.DidDepositRequestStorePrefix))
	for ; iterator.Valid(); iterator.Next() {
		var request types.DidDepositRequest
		keeper.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &request)
		requests = append(requests, request)
	}

	return requests
}

// ----------------------------
// --- Did PowerUp requests
// ----------------------------

func (keeper Keeper) getDidPowerUpRequestStoreKey(proof string) []byte {
	return []byte(types.DidPowerUpRequestStorePrefix + proof)
}

// StorePowerUpRequest allows to save the given request. Returns an error if a request with
// the same proof already exists
func (keeper Keeper) StorePowerUpRequest(ctx sdk.Context, request types.DidPowerUpRequest) sdk.Error {
	store := ctx.KVStore(keeper.StoreKey)

	requestStoreKey := keeper.getDidPowerUpRequestStoreKey(request.Proof)
	if store.Has(requestStoreKey) {
		return sdk.ErrUnknownRequest("PowerUp request with the same proof already exists")
	}

	// Set the initial status to null and store the request
	request.Status = nil
	store.Set(requestStoreKey, keeper.Cdc.MustMarshalBinaryBare(&request))

	return nil
}

func (keeper Keeper) GetPowerUpRequestByProof(ctx sdk.Context, proof string) (request types.DidPowerUpRequest, found bool) {
	store := ctx.KVStore(keeper.StoreKey)

	requestStoreKey := keeper.getDidPowerUpRequestStoreKey(proof)
	if !store.Has(requestStoreKey) {
		return types.DidPowerUpRequest{}, false
	}

	keeper.Cdc.MustUnmarshalBinaryBare(store.Get(requestStoreKey), &request)
	return request, true
}

// ChangePowerUpRequestStatus changes the status of the request having the same proof, or returns an error
// if no request with the given proof could be found
func (keeper Keeper) ChangePowerUpRequestStatus(ctx sdk.Context, proof string, status types.DidPowerUpRequestStatus) sdk.Error {
	store := ctx.KVStore(keeper.StoreKey)

	request, found := keeper.GetPowerUpRequestByProof(ctx, proof)
	if !found {
		return sdk.ErrUnknownRequest(fmt.Sprintf("PowerUp request with proof %s not found", proof))
	}

	// Update and store the request
	request.Status = &status
	store.Set(keeper.getDidPowerUpRequestStoreKey(proof), keeper.Cdc.MustMarshalBinaryBare(&request))

	return nil
}

// GetPowerUpRequests returns the list the requests saved inside the given context
func (keeper Keeper) GetPowerUpRequests(ctx sdk.Context) (requests []types.DidPowerUpRequest) {
	store := ctx.KVStore(keeper.StoreKey)

	iterator := sdk.KVStorePrefixIterator(store, []byte(types.DidPowerUpRequestStorePrefix))
	for ; iterator.Valid(); iterator.Next() {
		var request types.DidPowerUpRequest
		keeper.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &request)
		requests = append(requests, request)
	}

	return requests
}
