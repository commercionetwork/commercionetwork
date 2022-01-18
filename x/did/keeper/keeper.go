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
func (k Keeper) UpdateIdentity(ctx sdk.Context, identity *types.Identity) {
	store := ctx.KVStore(k.storeKey)
	store.Set(getIdentityStoreKey(identity.DidDocument.ID, identity.Metadata.Updated), k.cdc.MustMarshalBinaryBare(identity))
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
func (k Keeper) GetIdentityHistoryOfAddress(ctx sdk.Context, address string) []types.Identity {

	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, getIdentityPrefix(address))

	defer iterator.Close()

	history := []types.Identity{}

	for ; iterator.Valid(); iterator.Next() {
		var val types.Identity
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		history = append(history, val)
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

// UpdateDidDocument appends a DID document in the store, returning the ID contained in the DID document
// func (k Keeper) UpdateDidDocument(ctx sdk.Context, didDocument types.DidDocument) string {
// 	store := ctx.KVStore(k.storeKey)
// 	store.Set(getDidDocumentStoreKey(didDocument.ID), k.cdc.MustMarshalBinaryBare(&didDocument))
// 	return didDocument.ID
// }

// TODO change using GetIdentityOfAddress
// GetDidDocumentOfAddress returns the DID document reference associated to a given address.
// If the given address has no DID document associated, returns an error.
// func (k Keeper) GetDidDocumentOfAddress(ctx sdk.Context, address string) (types.DidDocument, error) {
// 	identity, err := k.GetLastIdentityOfAddress(ctx, address)
// 	if err != nil {
// 		return types.DidDocument{}, err
// 	}
// 	didDocument := *identity.DidDocument

// 	return didDocument, nil
// }

// func getDidDocumentStoreKey(owner string) []byte {
// 	return append([]byte(types.IdentitiesStorePrefix), owner...)
// }

//  TODO use HasIdentity
// HasDidDocument returns true if there is a DID document associated to a given ID.
// func (k Keeper) HasDidDocument(ctx sdk.Context, ID string) bool {
// 	store := ctx.KVStore(k.storeKey)
// 	identityKey := getDidDocumentStoreKey(ID)
// 	return store.Has(identityKey)
// }

// GetAllDidDocuments returns all the stored DID documents
// func (k Keeper) GetAllDidDocuments(ctx sdk.Context) (list []types.DidDocument) {
// 	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.IdentitiesStorePrefix))
// 	iterator := sdk.KVStorePrefixIterator(store, []byte{})

// 	defer iterator.Close()

// 	for ; iterator.Valid(); iterator.Next() {
// 		var val types.DidDocument
// 		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
// 		list = append(list, val)
// 	}

// 	return
// }
