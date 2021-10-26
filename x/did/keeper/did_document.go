package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/did/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AppendId appends a didDocument in the store with given id
func (k Keeper) AppendId(
	ctx sdk.Context,
	didDocument types.DidDocument,
) string {
	// Create the Document
	store := ctx.KVStore(k.storeKey)
	store.Set(getIdentityStoreKey(sdk.AccAddress(didDocument.ID)), k.cdc.MustMarshalBinaryBare(&didDocument))
	return didDocument.ID
}

func getIdentityStoreKey(owner sdk.AccAddress) []byte {
	return append([]byte(types.IdentitiesStorePrefix), owner...)
}

func getDidPowerUpRequestStoreKey(id string) []byte {
	return []byte(types.DidPowerUpRequestStorePrefix + id)
}

func (k Keeper) HasIdentity(ctx sdk.Context, ID string) bool {
	store := ctx.KVStore(k.storeKey)

	identityKey := getIdentityStoreKey(sdk.AccAddress(ID))
	if !store.Has(identityKey) {
		return false
	}
	return true
}

// GetAllDidDocument returns all didDocument
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
