package keeper

import (
	cmtypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/encapsulated/bank/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type Keeper struct {
	bank.Keeper
	storeKey sdk.StoreKey
	cdc      *codec.Codec
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, bankKeeper bank.Keeper) Keeper {
	return Keeper{
		Keeper:   bankKeeper,
		storeKey: storeKey,
		cdc:      cdc,
	}
}

// AddBlockedAddresses allows to add the given address as a blocked one
func (keeper Keeper) AddBlockedAddresses(ctx sdk.Context, address sdk.AccAddress) {
	addresses := keeper.GetBlockedAddresses(ctx)

	if newBlocked, edited := addresses.AppendIfMissing(address); edited {
		keeper.SetBlockedAddresses(ctx, newBlocked)
	}
}

// IsAddressBlocked returns true iff
func (keeper Keeper) IsAddressBlocked(ctx sdk.Context, address sdk.AccAddress) bool {
	return keeper.GetBlockedAddresses(ctx).Contains(address)
}

// GetBlockedAddresses returns the list of all the blocked addresses
func (keeper Keeper) GetBlockedAddresses(ctx sdk.Context) cmtypes.Addresses {
	store := ctx.KVStore(keeper.storeKey)

	var addresses cmtypes.Addresses
	keeper.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.BlockedAddressesStoreKey)), &addresses)
	return addresses
}

// SetBlockedAddresses allows to set the given addresses as the blocked ones
func (keeper Keeper) SetBlockedAddresses(ctx sdk.Context, addresses cmtypes.Addresses) {
	store := ctx.KVStore(keeper.storeKey)
	store.Set([]byte(types.BlockedAddressesStoreKey), keeper.cdc.MustMarshalBinaryBare(&addresses))
}
