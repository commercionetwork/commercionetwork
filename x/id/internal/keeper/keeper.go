package keeper

import (
	"strings"

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

func (keeper Keeper) getIdentiyStoreKey(owner sdk.AccAddress) []byte {
	return []byte(types.IdentitiesStorePrefix + owner.String())
}

// SaveIdentity saves the given didDocumentUri associating it with the given owner, replacing any existent one.
func (keeper Keeper) SaveIdentity(ctx sdk.Context, owner sdk.AccAddress, didDocumentUri string) {
	identitiesStore := ctx.KVStore(keeper.StoreKey)
	identitiesStore.Set(keeper.getIdentiyStoreKey(owner), []byte(didDocumentUri))
}

// GetIdentity returns the Did Document reference associated to a given Did.
// If the given Did has no Did Document reference associated, returns nil.
func (keeper Keeper) GetDidDocumentUriByDid(ctx sdk.Context, owner sdk.AccAddress) string {
	store := ctx.KVStore(keeper.StoreKey)
	result := store.Get(keeper.getIdentiyStoreKey(owner))
	return string(result)
}

// -------------------------
// --- Genesis utils
// -------------------------

// GetIdentities returns the list of all identities for the given context
func (keeper Keeper) GetIdentities(ctx sdk.Context) ([]types.Identity, error) {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.IdentitiesStorePrefix))

	var identities []types.Identity
	for ; iterator.Valid(); iterator.Next() {
		addressString := strings.ReplaceAll(string(iterator.Key()), types.IdentitiesStorePrefix, "")
		address, err := sdk.AccAddressFromBech32(addressString)
		if err != nil {
			return nil, err
		}

		identity := types.Identity{
			Owner:       address,
			DidDocument: string(iterator.Value()),
		}
		identities = append(identities, identity)
	}

	return identities, nil
}

// SetIdentities allows to bulk save a bunch of identities inside the store
func (keeper Keeper) SetIdentities(ctx sdk.Context, identities []types.Identity) {
	for _, identity := range identities {
		keeper.SaveIdentity(ctx, identity.Owner, identity.DidDocument)
	}
}
