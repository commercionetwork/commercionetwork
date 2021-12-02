package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

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

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetIdentity returns the Did Document reference associated to a given Did.
// If the given Did has no Did Document reference associated, returns nil.
func (k Keeper) GetDidDocumentByOwner(ctx sdk.Context, owner sdk.AccAddress) (types.DidDocument, error) {
	store := ctx.KVStore(k.storeKey)

	identityKey := getIdentityStoreKey(owner)
	if !store.Has(identityKey) {
		return types.DidDocument{}, fmt.Errorf("did document with owner %s not found", owner.String())
	}

	var didDocument types.DidDocument
	k.cdc.MustUnmarshalBinaryBare(store.Get(identityKey), &didDocument)
	return didDocument, nil
}
