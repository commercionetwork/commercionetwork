package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// GetIdentitiesModuleAccount returns the id ModuleAccount
func (k Keeper) GetIdentitiesModuleAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetPoolAmount returns the current pool amount
func (k Keeper) GetPoolAmount(ctx sdk.Context) (pool sdk.Coins) {
	return k.GetIdentitiesModuleAccount(ctx).GetCoins()
}
