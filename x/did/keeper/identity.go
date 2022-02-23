package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func getIdentityStoreKey(address string, updated string) []byte {
	res := append(getIdentityPrefix(address), ":"...)
	return append(res, updated...)
}

func getIdentityPrefix(address string) []byte {
	return append([]byte(types.IdentitiesStorePrefix), address...)
}

// SetIdentity sets an Identity in the store
func (k Keeper) SetIdentity(ctx sdk.Context, identity types.Identity) {
	store := ctx.KVStore(k.storeKey)
	address := identity.DidDocument.ID
	timestamp := identity.Metadata.Updated
	store.Set(getIdentityStoreKey(address, timestamp), k.cdc.MustMarshalBinaryBare(&identity))
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		updateDDO,
		sdk.NewAttribute("address", identity.DidDocument.ID),
	))
	logger := k.Logger(ctx)
	logger.Info("Identity successfully updated")

}

// GetIdentity returns the Identity associated to a given address and timestamp
// If there is no Identity associated to a given address and timestamp, returns an error
func (k Keeper) GetIdentity(ctx sdk.Context, address string, timestamp string) (*types.Identity, error) {

	store := ctx.KVStore(k.storeKey)
	identityBz := store.Get(getIdentityStoreKey(address, timestamp))

	if identityBz == nil {
		return nil, fmt.Errorf("could not find the Identity associated to %s %s", address, timestamp)
	}

	var identity types.Identity
	k.cdc.MustUnmarshalBinaryBare(identityBz, &identity)
	return &identity, nil
}

// GetLastIdentityOfAddress returns the last Identity associated to a given address
// If the given address has no DID Identity associated, returns an error
func (k Keeper) GetLastIdentityOfAddress(ctx sdk.Context, address string) (*types.Identity, error) {

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStoreReversePrefixIterator(store, getIdentityPrefix(address))
	defer iter.Close()

	if !iter.Valid() {
		return nil, fmt.Errorf("could not find the last Identity associated to %s", address)
	}

	var identity types.Identity
	k.cdc.MustUnmarshalBinaryBare(iter.Value(), &identity)
	return &identity, nil
}

// GetIdentityHistoryOfAddress returns all the stored Identities associated to a given address, in ascending chronological order
func (k Keeper) GetIdentityHistoryOfAddress(ctx sdk.Context, address string) []*types.Identity {

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, getIdentityPrefix(address))

	defer iterator.Close()

	history := []*types.Identity{}

	for ; iterator.Valid(); iterator.Next() {
		var val types.Identity
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		history = append(history, &val)
	}

	return history
}

// GetAllIdentities returns all the stored Identities
func (k Keeper) GetAllIdentities(ctx sdk.Context) []*types.Identity {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.IdentitiesStorePrefix))

	defer iterator.Close()

	list := []*types.Identity{}

	for ; iterator.Valid(); iterator.Next() {
		var val types.Identity
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, &val)
	}

	return list
}
