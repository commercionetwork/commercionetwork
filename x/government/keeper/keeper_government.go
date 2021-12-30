package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/government/types"
)

// SetGovernmentAddress allows to set the given address as the one that
// the government will use later
func (k Keeper) SetGovernmentAddress(ctx sdk.Context, address sdk.AccAddress) error {
	if address == nil {
		return errors.New("address is not defined")
	}

	store := ctx.KVStore(k.storeKey)

	if store.Has([]byte(types.GovernmentStoreKey)) {
		return errors.New("government address already set")
	}

	store.Set([]byte(types.GovernmentStoreKey), address)
	return nil
}

// GetGovernmentAddress returns the address that the government has currently
func (k Keeper) GetGovernmentAddress(ctx sdk.Context) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	return store.Get([]byte(types.GovernmentStoreKey))
}
