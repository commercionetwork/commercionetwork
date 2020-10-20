package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"

	"github.com/commercionetwork/commercionetwork/x/id/types"
)

// GetIdentitiesModuleAccount returns the id ModuleAccount
func (k Keeper) GetIdentitiesModuleAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}
