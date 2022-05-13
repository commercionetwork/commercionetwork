package keeper

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"errors"

	"github.com/commercionetwork/commercionetwork/x/government/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	// this line is used by starport scaffolding # ibc/keeper/import
)

type (
	Keeper struct {
		cdc      codec.Codec
		storeKey sdk.StoreKey
		memKey   sdk.StoreKey
	}
)

func NewKeeper(
	cdc codec.Codec,
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

// SetGovernmentAddress allows to set the given address as the one that
// the government will use later
func (k Keeper) SetGovernmentAddress(ctx sdk.Context, address sdk.AccAddress) error {
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
