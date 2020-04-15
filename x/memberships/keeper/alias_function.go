package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"

	"github.com/commercionetwork/commercionetwork/x/memberships/types"
)

// GetPoolFunds returns the funds currently present inside the rewards pool
func (k Keeper) GetPoolFunds(ctx sdk.Context) sdk.Coins {
	return k.GetMembershipModuleAccount(ctx).GetCoins()
}

// GetMembershipModuleAccount returns the module account for the accreditations module
func (k Keeper) GetMembershipModuleAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.SupplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}
