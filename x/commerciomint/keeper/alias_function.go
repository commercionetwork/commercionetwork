package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

// GetMintModuleAccount returns the commerciomint ModuleAccount
func (k Keeper) GetMintModuleAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

// GetLiquidityPoolAmount returns the current liquidity pool amount
func (k Keeper) GetLiquidityPoolAmount(ctx sdk.Context) sdk.Coins {
	return k.GetMintModuleAccount(ctx).GetCoins()
}
