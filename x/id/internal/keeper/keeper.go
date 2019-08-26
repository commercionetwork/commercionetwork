package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	StoreKey sdk.StoreKey

	// Pointer to the codec that is used by Amino to encode and decode binary structs.
	Cdc *codec.Codec
}

// NewKeeper creates new instances of the CommercioID Keeper
func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		StoreKey: storeKey,
		Cdc:      cdc,
	}
}

// SaveIdentity saves the given didDocumentUri associating it with the given owner, replacing any existent one.
func (keeper Keeper) SaveIdentity(ctx sdk.Context, owner sdk.AccAddress, didDocumentUri string) {
	identitiesStore := ctx.KVStore(keeper.StoreKey)
	identitiesStore.Set([]byte(types.IdentitiesStorePrefix+owner.String()), []byte(didDocumentUri))
}

// GetIdentity returns the Did Document reference associated to a given Did.
// If the given Did has no Did Document reference associated, returns nil.
func (keeper Keeper) GetDidDocumentUriByDid(ctx sdk.Context, owner sdk.AccAddress) string {
	store := ctx.KVStore(keeper.StoreKey)
	result := store.Get([]byte(types.IdentitiesStorePrefix + owner.String()))
	return string(result)
}
