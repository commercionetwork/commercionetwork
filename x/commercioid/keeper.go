package commercioid

import (
	// Provides tools to work with the Cosmos encoding format, Amino.
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	identitiesStoreKey sdk.StoreKey // Map containing DID -> DID Document references
	ownersStoresKey    sdk.StoreKey // Map containing DID -> Owner references

	// Pointer to the codec that is used by Amino to encode and decode binary structs.
	cdc *codec.Codec
}

// NewKeeper creates new instances of the nameservice Keeper
func NewKeeper(
	identitiesStoreKey sdk.StoreKey,
	ownersStoresKey sdk.StoreKey,
	cdc *codec.Codec) Keeper {
	return Keeper{
		identitiesStoreKey: identitiesStoreKey,
		ownersStoresKey:    ownersStoresKey,
		cdc:                cdc,
	}
}

func (keeper Keeper) SetIdentity(ctx sdk.Context, did string, ddoReference string) {
	// Get the store that contains the reference to the identities
	store := ctx.KVStore(keeper.identitiesStoreKey)

	// Insert the <DID, DID Document Reference> keypair into the store
	store.Set([]byte(did), []byte(ddoReference))
}

func (keeper Keeper) GetIdentity(ctx sdk.Context, did string) string {
	store := ctx.KVStore(keeper.identitiesStoreKey)
	result := store.Get([]byte(did))
	return string(result)
}

func (keeper Keeper) SetOwner(ctx sdk.Context, did string, owner sdk.AccAddress) {
	store := ctx.KVStore(keeper.ownersStoresKey)
	store.Set([]byte(did), owner)
}

func (keeper Keeper) HasOwner(ctx sdk.Context, did string) bool {
	store := ctx.KVStore(keeper.ownersStoresKey)
	result := store.Get([]byte(did))
	return result != nil
}

func (keeper Keeper) GetOwner(ctx sdk.Context, did string) sdk.AccAddress {
	store := ctx.KVStore(keeper.ownersStoresKey)
	result := store.Get([]byte(did))
	return result
}
