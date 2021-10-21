package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// this line is used by starport scaffolding # ibc/keeper/import
)

type (
	Keeper struct {
		cdc      codec.Marshaler
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey,
	memKey sdk.StoreKey,
) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
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
		return r.Status == nil
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
