package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

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
func (k Keeper) GetDidDocuments(ctx sdk.Context) []types.DidDocument {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.IdentitiesStorePrefix))

	var didDocuments []types.DidDocument
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var didDocument types.DidDocument
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &didDocument)
		didDocuments = append(didDocuments, didDocument)
	}

	return didDocuments
}

// ----------------------------
// --- Did power up requests
// ----------------------------

// StorePowerUpRequest allows to save the given request. Returns an error if a request with
// the same proof already exists
func (k Keeper) StorePowerUpRequest(ctx sdk.Context, request types.DidPowerUpRequest) error {
	store := ctx.KVStore(k.storeKey)

	requestStoreKey := getDidPowerUpRequestStoreKey(request.ID)
	if store.Has(requestStoreKey) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "PowerUp request with the same id already exists")
	}

	store.Set(requestStoreKey, k.cdc.MustMarshalBinaryBare(&request))

	return nil
}

// GetDidDepositRequestByProof returns the request having the same proof.
func (k Keeper) GetPowerUpRequestByID(ctx sdk.Context, id string) (types.DidPowerUpRequest, error) {
	store := ctx.KVStore(k.storeKey)

	requestStoreKey := getDidPowerUpRequestStoreKey(id)
	if !store.Has(requestStoreKey) {
		return types.DidPowerUpRequest{}, fmt.Errorf("power-up request with id %s not found", id)
	}

	request := types.DidPowerUpRequest{}
	k.cdc.MustUnmarshalBinaryBare(store.Get(requestStoreKey), &request)
	return request, nil
}

// ChangePowerUpRequestStatus changes the status of the request having the same proof, or returns an error
// if no request with the given proof could be found
func (k Keeper) ChangePowerUpRequestStatus(ctx sdk.Context, id string, status types.RequestStatus) error {
	store := ctx.KVStore(k.storeKey)

	request, err := k.GetPowerUpRequestByID(ctx, id)
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
	}

	// Update and store the request
	request.Status = &status
	store.Set(getDidPowerUpRequestStoreKey(id), k.cdc.MustMarshalBinaryBare(&request))

	return nil
}

// GetPowerUpRequests returns the list the requests saved inside the given context
func (k Keeper) GetPowerUpRequests(ctx sdk.Context) (requests []types.DidPowerUpRequest) {
	return k.iterRequestsWithFunc(ctx, func(r types.DidPowerUpRequest) bool {
		return true
	})
}

// GetApprovedPowerUpRequests returns the list of handled requests saved inside the given context
func (k Keeper) GetApprovedPowerUpRequests(ctx sdk.Context) (requests []types.DidPowerUpRequest) {
	return k.iterRequestsWithFunc(ctx, func(r types.DidPowerUpRequest) bool {
		if r.Status != nil {
			if r.Status.Type == types.StatusApproved {
				return true
			}
		}

		return false
	})
}

// GetRejectedPowerUpRequests returns the list of rejected (canceled, invalid) requests saved inside the given context
func (k Keeper) GetRejectedPowerUpRequests(ctx sdk.Context) (requests []types.DidPowerUpRequest) {
	return k.iterRequestsWithFunc(ctx, func(r types.DidPowerUpRequest) bool {
		if r.Status != nil {
			if r.Status.Type != types.StatusApproved {
				return true
			}
		}

		return false
	})
}

// GetPendingPowerUpRequests returns the list of pending requests saved inside the given context
func (k Keeper) GetPendingPowerUpRequests(ctx sdk.Context) (requests []types.DidPowerUpRequest) {
	return k.iterRequestsWithFunc(ctx, func(r types.DidPowerUpRequest) bool {
		if r.Status == nil {
			return true
		}

		return false
	})
}

// iterRequestsWithFunc returns a slice of requests, based on the logic of rationale.
// If rationale() returns true, r will be added to requests.
func (k Keeper) iterRequestsWithFunc(ctx sdk.Context, rationale func(r types.DidPowerUpRequest) bool) (requests []types.DidPowerUpRequest) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.DidPowerUpRequestStorePrefix))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var request types.DidPowerUpRequest
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &request)

		if rationale(request) {
			requests = append(requests, request)
		}
	}

	return requests
}
