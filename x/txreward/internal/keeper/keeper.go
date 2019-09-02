package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/txreward/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

type Keeper struct {
	StoreKey sdk.StoreKey

	BankKeeper  bank.Keeper
	StakeKeeper staking.Keeper

	Cdc *codec.Codec
}

func NewKeeper(storeKey sdk.StoreKey, bk bank.Keeper, cdc *codec.Codec) Keeper {
	return Keeper{
		StoreKey:   storeKey,
		BankKeeper: bk,
		Cdc:        cdc,
	}
}

func (keeper Keeper) IncrementBlockRewardsPool(ctx sdk.Context, funder sdk.AccAddress, amount sdk.Coin) {
	store := ctx.KVStore(keeper.StoreKey)
	bk := keeper.BankKeeper
	brAmount := sdk.Coins{amount}
	var brPool sdk.Coins

	if bk.HasCoins(ctx, funder, brAmount) {
		poolBz := store.Get([]byte(types.BlockRewardsPoolPrefix))
		if poolBz == nil {
			store.Set([]byte(types.BlockRewardsPoolPrefix), keeper.Cdc.MustMarshalBinaryBare(&brAmount))
		} else {
			poolBz := store.Get([]byte(types.BlockRewardsPoolPrefix))
			keeper.Cdc.MustUnmarshalBinaryBare(poolBz, &brPool)
			brPool.Add(brAmount)
			store.Set([]byte(types.BlockRewardsPoolPrefix), keeper.Cdc.MustMarshalBinaryBare(&brPool))
		}
	}
}

func (keeper Keeper) DistributeBlockRewards(ctx sdk.Context, validators []staking.Validator) {

}

func (keeper Keeper) ComputeValidatorsReward(ctx sdk.Context, validatorNumber sdk.Int)
