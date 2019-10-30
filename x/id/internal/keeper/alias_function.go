package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

// GetIdModuleAccount returns the id ModuleAccount
func (k Keeper) GetIdModuleAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetPoolAmount returns the current pool amount
func (k Keeper) GetPoolAmount(ctx sdk.Context) (pool sdk.Coins) {
	return k.GetIdModuleAccount(ctx).GetCoins()
}
