package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AppendDid appends a DID document in the store with given id
func (k Keeper) AppendDid(ctx sdk.Context, didDocument types.DidDocument) string {
	// Create the Document
	store := ctx.KVStore(k.storeKey)
	store.Set(getIdentityStoreKey(sdk.AccAddress(didDocument.ID)), k.cdc.MustMarshalBinaryBare(&didDocument))
	return didDocument.ID
}

// GetDdoByOwner returns the DID document reference associated to a given DID.
// If the given DID has no DID document reference associated, returns nil.
func (k Keeper) GetDdoByOwner(ctx sdk.Context, owner sdk.AccAddress) (types.DidDocument, error) {
	store := ctx.KVStore(k.storeKey)

	identityKey := getIdentityStoreKey(owner)
	if !store.Has(identityKey) {
		return types.DidDocument{}, fmt.Errorf("DID document with owner %s not found", owner.String())
	}

	var DidDocument types.DidDocument
	k.cdc.MustUnmarshalBinaryBare(store.Get(identityKey), &DidDocument)
	return DidDocument, nil
}

func getIdentityStoreKey(owner sdk.AccAddress) []byte {
	return append([]byte(types.IdentitiesStorePrefix), owner...)
}

func (k Keeper) HasIdentity(ctx sdk.Context, ID string) bool {
	store := ctx.KVStore(k.storeKey)

	identityKey := getIdentityStoreKey(sdk.AccAddress(ID))
	return store.Has(identityKey)
}

// GetAllDidDocument returns all DID document
func (k Keeper) GetAllDidDocument(ctx sdk.Context) (list []types.DidDocument) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.IdentitiesStorePrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DidDocument
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
