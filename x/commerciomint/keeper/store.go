package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	accType "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

// GetLiquidityPoolAmount returns the current liquidity pool amount
func (k Keeper) GetModuleAccount(ctx sdk.Context) accType.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k Keeper) GetLiquidityPoolAmount(ctx sdk.Context) sdk.Coins {
	moduleAccount := k.GetModuleAccount(ctx)
	return k.bankKeeper.GetAllBalances(ctx, moduleAccount.GetAddress())

}

func (k Keeper) GetLiquidityPoolAmountCoins(ctx sdk.Context) sdk.Coins {

	moduleAccount := k.GetModuleAccount(ctx)
	var coins sdk.Coins
	for _, coin := range k.bankKeeper.GetAllBalances(ctx, moduleAccount.GetAddress()) {
		coins = append(coins, coin)
	}
	return coins

}

func (k Keeper) SetLiquidityPoolToAccount(ctx sdk.Context, coins sdk.Coins) error {
	moduleAccount := k.GetModuleAccount(ctx)
	return k.bankKeeper.AddCoins(ctx, moduleAccount.GetAddress(), coins)
}

func (k Keeper) SetModuleAccount(ctx sdk.Context, acc accType.ModuleAccountI) {
	k.accountKeeper.SetModuleAccount(ctx, acc)
}
