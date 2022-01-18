package keeper

import (
	"fmt"

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

func getIdentityStoreKey(address string, updated string) []byte {
	res := append(getIdentityPrefix(address), ":"...)
	return append(res, updated...)
}

func getIdentityPrefix(address string) []byte {
	return append([]byte(types.IdentitiesStorePrefix), address...)
}

// UpdateIdentity appends an Identity in the store
func (k Keeper) UpdateIdentity(ctx sdk.Context, identity types.Identity) {
	store := ctx.KVStore(k.storeKey)
	store.Set(getIdentityStoreKey(identity.DidDocument.ID, identity.Metadata.Updated), k.cdc.MustMarshalBinaryBare(&identity))
}

// GetLastIdentityOfAddress returns the last Identity associated to a given address
// If the given address has no DID Identity associated, returns an error
func (k Keeper) GetLastIdentityOfAddress(ctx sdk.Context, address string) (*types.Identity, error) {

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStoreReversePrefixIterator(store, getIdentityPrefix(address))
	defer iter.Close()

	if !iter.Valid() {
		return nil, fmt.Errorf("could not find identity for %s", address)
	}

	var identity types.Identity
	k.cdc.MustUnmarshalBinaryBare(iter.Value(), &identity)
	return &identity, nil
}

// GetAllDidDocuments returns all the stored DID documents
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

// GetAllDidDocuments returns all the stored DID documents
func (k Keeper) GetAllIdentities(ctx sdk.Context) []*types.Identity {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	list := []*types.Identity{}

	for ; iterator.Valid(); iterator.Next() {
		var val types.Identity
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, &val)
	}

	return list
}
