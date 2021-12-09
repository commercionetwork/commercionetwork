package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

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

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// AppendDidDocument appends a DID document in the store, retruning the ID contained in the DID document
func (k Keeper) AppendDidDocument(ctx sdk.Context, didDocument types.DidDocument) string {
	store := ctx.KVStore(k.storeKey)
	store.Set(getIdentityStoreKey(sdk.AccAddress(didDocument.ID)), k.cdc.MustMarshalBinaryBare(&didDocument))
	return didDocument.ID
}

// GetDidDocumentOfAddress returns the DID document reference associated to a given address.
// If the given address has no DID document associated, returns nil.
func (k Keeper) GetDidDocumentOfAddress(ctx sdk.Context, owner sdk.AccAddress) (types.DidDocument, error) {
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

// HasDidDocument returns true if there is a DID document associated to a given ID.
func (k Keeper) HasDidDocument(ctx sdk.Context, ID string) bool {
	store := ctx.KVStore(k.storeKey)
	identityKey := getIdentityStoreKey(sdk.AccAddress(ID))
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
