package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	accType "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

func (k Keeper) GetModuleAccount(ctx sdk.Context) accType.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k Keeper) GetModuleBalance(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	return k.bankKeeper.GetAllBalances(ctx, addr)
}

func (k Keeper) GetLiquidityPoolAmount(ctx sdk.Context) sdk.Coins {
	moduleAccount := k.GetModuleAccount(ctx)
	return k.bankKeeper.GetAllBalances(ctx, moduleAccount.GetAddress())

}

func (k Keeper) SetLiquidityPoolToAccount(ctx sdk.Context, coins sdk.Coins) error {
	moduleAccount := k.GetModuleAccount(ctx)
	if err := k.bankKeeper.SetBalances(ctx, moduleAccount.GetAddress(), coins); err != nil {
		return err
	}
	// TODO: check liquidity amount on migration
	supply := k.bankKeeper.GetSupply(ctx)
	supply.Inflate(coins)
	k.bankKeeper.SetSupply(ctx, supply)
	return nil
	//return k.bankKeeper.AddCoins(ctx, moduleAccount.GetAddress(), coins)
}

func (k Keeper) SetModuleAccount(ctx sdk.Context, acc accType.ModuleAccountI) {
	k.accountKeeper.SetModuleAccount(ctx, acc)
}
