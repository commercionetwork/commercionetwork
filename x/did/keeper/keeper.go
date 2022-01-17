package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/did/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
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

// UpdateDidDocument appends a DID document in the store, returning the ID contained in the DID document
func (k Keeper) UpdateDidDocument(ctx sdk.Context, didDocument types.DidDocument) string {
	store := ctx.KVStore(k.storeKey)
	store.Set(getIdentityStoreKey(didDocument.ID), k.cdc.MustMarshalBinaryBare(&didDocument))
	return didDocument.ID
}

// GetDidDocumentOfAddress returns the DID document reference associated to a given address.
// If the given address has no DID document associated, returns an error.
func (k Keeper) GetDidDocumentOfAddress(ctx sdk.Context, address string) (types.DidDocument, error) {
	if !k.HasDidDocument(ctx, address) {
		return types.DidDocument{}, fmt.Errorf("DID document for %s not found", address)
	}

	identityKey := getIdentityStoreKey(address)
	store := ctx.KVStore(k.storeKey)
	var didDocument types.DidDocument
	k.cdc.MustUnmarshalBinaryBare(store.Get(identityKey), &didDocument)
	return didDocument, nil
}

func getIdentityStoreKey(owner string) []byte {
	return append([]byte(types.IdentitiesStorePrefix), owner...)
}

// HasDidDocument returns true if there is a DID document associated to a given ID.
func (k Keeper) HasDidDocument(ctx sdk.Context, ID string) bool {
	store := ctx.KVStore(k.storeKey)
	identityKey := getIdentityStoreKey(ID)
	return store.Has(identityKey)
}

// GetAllDidDocuments returns all the stored DID documents
func (k Keeper) GetAllDidDocuments(ctx sdk.Context) (list []types.DidDocument) {
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
