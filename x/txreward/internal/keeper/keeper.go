package keeper

import (
	"github.com/commercionetwork/commercionetwork/app"
	"github.com/commercionetwork/commercionetwork/x/txreward/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
)

type Keeper struct {
	StoreKey sdk.StoreKey

	BankKeeper         bank.Keeper
	StakeKeeper        staking.Keeper
	DistributionKeeper distribution.Keeper

	Cdc *codec.Codec
}

func NewKeeper(storeKey sdk.StoreKey, bk bank.Keeper, sk staking.Keeper, dk distribution.Keeper, cdc *codec.Codec) Keeper {
	return Keeper{
		StoreKey:           storeKey,
		BankKeeper:         bk,
		StakeKeeper:        sk,
		DistributionKeeper: dk,
		Cdc:                cdc,
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

func (keeper Keeper) GetPoolFunds(ctx sdk.Context) sdk.DecCoins {
	store := ctx.KVStore(keeper.StoreKey)
	var brPool sdk.Coins

	poolBz := store.Get([]byte(types.BlockRewardsPoolPrefix))
	keeper.Cdc.MustUnmarshalBinaryBare(poolBz, brPool)

	return sdk.NewDecCoins(brPool)
}

func (keeper Keeper) setPoolFunds(ctx sdk.Context, updatedPool sdk.Coins) {
	store := ctx.KVStore(keeper.StoreKey)
	store.Set([]byte(types.BlockRewardsPoolPrefix), keeper.Cdc.MustMarshalBinaryBare(&updatedPool))
}

/*
Compute the reward of proposer validator as following:

With 100 or less active validators, we calculate the reward like this:

TPY = Tokens Per Year
Reward100 = TPY / (365 * 24 * 60 * 12)

Instead, if the active validators will be more than 100, we calculate the reward like this:

V = Validators Number (assuming it's greater than 100)

RewardVN = (Reward100 / V) * 100

Summarizing these formulas we obtain:

Reward(n, V) = TPY / (365 * 24 * 60 * 12) * 100 / V * STAKE / TOTALSTAKE
365 years' days
24 hours per day
60 minutes per hour
12 blocks per minutes

where:
Reward(n, V) indicates the reward for the n validator considering a set of V validators
V 			 indicates the Validators Number
STAKE 		 staked token's amount of n-esim validator
TOTALSTAKE	 indicates all staked token's amount of all validators

*/

// constants representing Tokens Per Year, Days Per Year, Hours Per Day, Minutes Per Days, Blocks Per Minutes
const (
	TPY = 25000
	DPY = 365
	HPD = 24
	MPD = 60
	BPM = 12
)

func computeRawReward(validatorNumber sdk.Int) sdk.Dec {
	rawReward := (TPY * 1000000) / (DPY * HPD * MPD * BPM) * (100 / validatorNumber.Int64())
	return sdk.NewDec(rawReward)
}

func (keeper Keeper) ComputeValidatorReward(ctx sdk.Context, validatorNumber sdk.Int, proposer exported.ValidatorI,
	totalStakedTokens sdk.Int) sdk.DecCoins {

	//Get the raw reward for proposer
	rawReward := computeRawReward(validatorNumber)

	//Retrieve staked tokens by proposer
	validatorStakedTokens := proposer.GetBondedTokens()
	//compute his validation power
	validatorPower := validatorStakedTokens.Quo(totalStakedTokens).ToDec()

	//calculate the final reward for this proposer
	concreteReward := rawReward.Mul(validatorPower)

	coinReward := sdk.DecCoin{Denom: app.DefaultBondDenom, Amount: concreteReward}

	return append(sdk.DecCoins{}, coinReward)
}

//Distribute the computed reward to the block proposer
func (keeper Keeper) DistributeBlockRewards(ctx sdk.Context, validator exported.ValidatorI, reward sdk.DecCoins) {

	brPoolFunds := keeper.GetPoolFunds(ctx)
	if brPoolFunds.AmountOf(app.DefaultBondDenom).GTE(reward.AmountOf(app.DefaultBondDenom)) {
		brPoolFunds = brPoolFunds.Sub(reward)
		keeper.setPoolFunds(ctx)
	}

	//Get his current reward and then add the new one
	currentRewards := keeper.DistributionKeeper.GetValidatorCurrentRewards(ctx, validator.GetOperator())

	currentRewards.Rewards = currentRewards.Rewards.Add(reward)

	//Set the just earned reward
	keeper.DistributionKeeper.SetValidatorCurrentRewards(ctx, validator.GetOperator(), currentRewards)
}
