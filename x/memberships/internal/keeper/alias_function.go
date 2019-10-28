package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
)

func (k Keeper) GetPoolFunds(ctx sdk.Context) sdk.Coins {
	return k.GetMembershipModuleAccount(ctx).GetCoins()
}

func (k Keeper) GetMembershipModuleAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}
