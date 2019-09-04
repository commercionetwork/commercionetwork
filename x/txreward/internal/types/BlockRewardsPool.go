package types

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/txreward/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BlockRewardsPool struct {
	Funds sdk.DecCoins `json:"funds"`
}

func InitBlockRewardsPool() BlockRewardsPool {
	return BlockRewardsPool{
		Funds: sdk.DecCoins{},
	}
}

func (brp BlockRewardsPool) ValidateGenesis() error {
	if brp.Funds.IsAnyNegative() {
		return fmt.Errorf("negative Funds in block reward pool, is %v", brp.Funds)
	}

	return nil
}

//Utility function to set Block Reward Pool
func (brp BlockRewardsPool) SetBlockRewardsPool(ctx sdk.Context, keeper keeper.Keeper, updatedPool *BlockRewardsPool) {
	store := ctx.KVStore(keeper.StoreKey)
	store.Set([]byte(BlockRewardsPoolPrefix), keeper.Cdc.MustMarshalBinaryBare(&updatedPool))
}

//Utility function to get Block Reward Pool
func (brp BlockRewardsPool) GetBlockRewardsPool(ctx sdk.Context, keeper keeper.Keeper) BlockRewardsPool {
	store := ctx.KVStore(keeper.StoreKey)
	brpBz := store.Get([]byte(BlockRewardsPoolPrefix))
	if brpBz == nil {
		return InitBlockRewardsPool()
	}
	keeper.Cdc.MustUnmarshalBinaryBare(brpBz, &brp)
	return brp
}
