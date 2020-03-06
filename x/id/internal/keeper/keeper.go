package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

type Keeper struct {
	storeKey     sdk.StoreKey
	cdc          *codec.Codec
	accKeeper    auth.AccountKeeper
	supplyKeeper supply.Keeper
}

// NewKeeper creates new instances of the CommercioID Keeper
func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, accKeeper auth.AccountKeeper, supplyKeeper supply.Keeper) Keeper {

	// ensure commerciomint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:     storeKey,
		cdc:          cdc,
		accKeeper:    accKeeper,
		supplyKeeper: supplyKeeper,
	}
}

// ------------------
// --- Identities
// ------------------

// SaveDidDocument saves the given didDocumentUri associating it with the given owner, replacing any existent one.
func (k Keeper) SaveDidDocument(ctx sdk.Context, document types.DidDocument) error {
	// validate the DidDocument before saving it
	if err := document.Validate(); err != nil {
		return err
	}

	// Set the Did Document into the store
	identitiesStore := ctx.KVStore(k.storeKey)
	identitiesStore.Set(getIdentityStoreKey(document.ID), k.cdc.MustMarshalBinaryBare(document))

	return nil
}

// GetIdentity returns the Did Document reference associated to a given Did.
// If the given Did has no Did Document reference associated, returns nil.
func (k Keeper) GetDidDocumentByOwner(ctx sdk.Context, owner sdk.AccAddress) (types.DidDocument, error) {
	store := ctx.KVStore(k.storeKey)

	identityKey := getIdentityStoreKey(owner)
	if !store.Has(identityKey) {
		return types.DidDocument{}, fmt.Errorf("did document with owner %s not found", owner.String())
	}

	var didDocument types.DidDocument
	k.cdc.MustUnmarshalBinaryBare(store.Get(identityKey), &didDocument)
	return didDocument, nil
}

// -------------------------
// --- Genesis utils
// -------------------------

// GetDidDocuments returns the list of all identities for the given context
func (k Keeper) GetDidDocuments(ctx sdk.Context) ([]types.DidDocument, error) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.IdentitiesStorePrefix))

	var didDocuments []types.DidDocument
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var didDocument types.DidDocument
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &didDocument)
		didDocuments = append(didDocuments, didDocument)
	}

	return didDocuments, nil
}

// ----------------------------
// --- Did power up requests
// ----------------------------

// StorePowerUpRequest allows to save the given request. Returns an error if a request with
// the same proof already exists
func (k Keeper) StorePowerUpRequest(ctx sdk.Context, request types.DidPowerUpRequest) error {
	store := ctx.KVStore(k.storeKey)

	requestStoreKey := getDidPowerUpRequestStoreKey(request.Proof)
	if store.Has(requestStoreKey) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "PowerUp request with the same proof already exists")
	}

	store.Set(requestStoreKey, k.cdc.MustMarshalBinaryBare(&request))

	return nil
}

// GetDidDepositRequestByProof returns the request having the same proof.
func (k Keeper) GetPowerUpRequestByProof(ctx sdk.Context, proof string) (types.DidPowerUpRequest, error) {
	store := ctx.KVStore(k.storeKey)

	requestStoreKey := getDidPowerUpRequestStoreKey(proof)
	if !store.Has(requestStoreKey) {
		return types.DidPowerUpRequest{}, fmt.Errorf("power-up request with proof %s not found", proof)
	}

	request := types.DidPowerUpRequest{}
	k.cdc.MustUnmarshalBinaryBare(store.Get(requestStoreKey), &request)
	return request, nil
}

// ChangePowerUpRequestStatus changes the status of the request having the same proof, or returns an error
// if no request with the given proof could be found
func (k Keeper) ChangePowerUpRequestStatus(ctx sdk.Context, proof string, status types.RequestStatus) error {
	store := ctx.KVStore(k.storeKey)

	request, err := k.GetPowerUpRequestByProof(ctx, proof)
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
	}

	// Update and store the request
	request.Status = &status
	store.Set(getDidPowerUpRequestStoreKey(proof), k.cdc.MustMarshalBinaryBare(&request))

	return nil
}

// GetPowerUpRequests returns the list the requests saved inside the given context
func (k Keeper) GetPowerUpRequests(ctx sdk.Context) (requests []types.DidPowerUpRequest) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.DidPowerUpRequestStorePrefix))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var request types.DidPowerUpRequest
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &request)
		requests = append(requests, request)
	}

	return requests
}

// ------------------------
// --- Deposits handling
// ------------------------

func (k Keeper) SetPowerUpRequestHandled(ctx sdk.Context, activationReference string) error {
	store := ctx.KVStore(k.storeKey)

	if store.Has(getHandledPowerUpRequestsReferenceStoreKey(activationReference)) {
		return fmt.Errorf("reference %s already marked as handled", activationReference)
	}

	k.SetHandledPowerUpRequestsReference(ctx, activationReference)
	return nil
}

func (k Keeper) GetHandledPowerUpRequestsReferences(ctx sdk.Context) ctypes.Strings {
	hi := k.HandledPowerUpRequestsIterator(ctx)
	defer hi.Close()

	data := ctypes.Strings{}
	for ; hi.Valid(); hi.Next() {
		data = append(data, string(hi.Value()))
	}

	return data
}

func (k Keeper) SetHandledPowerUpRequestsReference(ctx sdk.Context, reference string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(getHandledPowerUpRequestsReferenceStoreKey(reference), []byte(reference))
}

func (k Keeper) HandledPowerUpRequestsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)

	return sdk.KVStorePrefixIterator(store, []byte(types.HandledPowerUpRequestsReferenceStorePrefix))
}
