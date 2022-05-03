package vbr

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/vbr/keeper"
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set the reward pool - Should never be nil as its validated inside the ValidateGenesis method
	k.SetTotalRewardPool(ctx, genState.PoolAmount)

	moduleAcc := k.VbrAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	coins := keeper.GetCoins(k, ctx, moduleAcc)
	if coins.Empty() {
		amount, _ := genState.PoolAmount.TruncateDecimal()
		err := k.MintVBRTokens(ctx, sdk.NewCoins(amount...))
		if err != nil {
			panic(err) // could not mint tokens on chain start, fatal!
		}
	}
	if err := k.SetParamSet(ctx, genState.Params); err != nil {
		panic(err)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		PoolAmount: k.GetTotalRewardPool(ctx),
		Params:     k.GetParamSet(ctx),
	}
}
