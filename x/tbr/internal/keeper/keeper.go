package keeper

import (
	"errors"

	"github.com/commercionetwork/commercionetwork/x/tbr/internal/types"
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
	StakingKeeper      staking.Keeper
	DistributionKeeper distribution.Keeper

	Cdc *codec.Codec
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, bk bank.Keeper, sk staking.Keeper, dk distribution.Keeper) Keeper {
	return Keeper{
		StoreKey:           storeKey,
		BankKeeper:         bk,
		DistributionKeeper: dk,
		StakingKeeper:      sk,
		Cdc:                cdc,
	}
}

func (k Keeper) SetRewardDenom(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.RewardDenomStoreKey), []byte(denom))
}

func (k Keeper) GetRewardDenom(ctx sdk.Context) string {
	store := ctx.KVStore(k.StoreKey)
	return string(store.Get([]byte(types.RewardDenomStoreKey)))
}

// Utility method to set Block Reward Pool
func (k Keeper) SetBlockRewardsPool(ctx sdk.Context, updatedPool sdk.DecCoins) {
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.PoolStoreKey), k.Cdc.MustMarshalBinaryBare(&updatedPool))
}

//Return the Block Rewards Pool if exists
func (k Keeper) GetBlockRewardsPool(ctx sdk.Context) sdk.DecCoins {
	var brPool sdk.DecCoins
	store := ctx.KVStore(k.StoreKey)
	brpBz := store.Get([]byte(types.PoolStoreKey))
	k.Cdc.MustUnmarshalBinaryBare(brpBz, &brPool)
	return brPool
}

// IncrementBlockRewardsPool increases the block rewards pool with the specified coin amount
func (k Keeper) IncrementBlockRewardsPool(ctx sdk.Context, funder sdk.AccAddress, amount sdk.Coins) error {
	bk := k.BankKeeper
	brPool := sdk.DecCoins{}

	if bk.HasCoins(ctx, funder, amount) {
		brPool = k.GetBlockRewardsPool(ctx)
		_, err := bk.SubtractCoins(ctx, funder, amount)
		if err != nil {
			return err
		}
		brPool = brPool.Add(sdk.NewDecCoins(amount))
		k.SetBlockRewardsPool(ctx, brPool)
	}
	return errors.New("not enough tokens to supply the given funds' amount")

}

// Compute the reward of proposer validator as following:
//
// With 100 or less active validators, we calculate the reward like this:
//
// TPY = Tokens Per Year
// Divider =  (365 * 24 * 60 * 12)
//
// Reward100 = TPY / Divider
//
// Instead, if the active validators will be more than 100, we calculate the reward like this:
//
// V = Validators Number (assuming it's greater than 100)
//
// RewardVN = (Reward100 / V) * 100
//
// Summarizing these formulas we obtain:
//
// Reward(n, V) = TPY / Divider * 100 / V * STAKE / TOTALSTAKE
//
// Divider's members:
// 365 years' days
// 24 hours per day
// 60 minutes per hour
// 12 blocks per minutes
//
// where:
// Reward(n, V) indicates the reward for the n validator considering a set of V validators
// V 			 indicates the Validators Number
// STAKE 		 staked token's amount of n-esim validator
// TOTALSTAKE	 indicates all staked token's amount of all validators
var (
	TPY = sdk.NewDecWithPrec(25000, 0) //Tokens Per Year
	DPY = sdk.NewDecWithPrec(36524, 2) // Days Per Year
	HPD = sdk.NewDecWithPrec(24, 0)    //  Hours Per Day
	MPH = sdk.NewDecWithPrec(60, 0)    //  Minutes Per Hour
	BPM = sdk.NewDecWithPrec(12, 0)    // Blocks Per Minutes
)

// computeRawReward computes the raw reward for proposer, assuming that raw reward is the Reward(n, V)
// equation's result without the last multiplication between the (Stake/TotalStake) value
func computeRawReward(validatorsNumber int64) sdk.Dec {
	tokensPerYear := TPY.Mul(sdk.NewDecWithPrec(100000, 0))
	divider := DPY.Mul(HPD).Mul(MPH).Mul(BPM)

	firstMember := tokensPerYear.Quo(divider)

	averageValidatorsNumber := sdk.NewDecWithPrec(100, 0)
	vNumber := sdk.NewDecWithPrec(validatorsNumber, 0)

	secondMember := averageValidatorsNumber.Quo(vNumber)

	firstMember = firstMember.Mul(secondMember)

	return firstMember
}

// ComputeProposerReward computes the final reward for the validator block's proposer
func (k Keeper) ComputeProposerReward(ctx sdk.Context, validatorNumber int64, proposer exported.ValidatorI,
	totalStakedTokens sdk.Int) sdk.DecCoins {

	// Get the raw reward for proposer
	rawReward := computeRawReward(validatorNumber)

	// Retrieve staked tokens by proposer
	validatorStakedTokens := proposer.GetBondedTokens().ToDec()

	// Compute his validation power
	validatorPower := validatorStakedTokens.Quo(totalStakedTokens.ToDec())

	// Calculate the final reward for this proposer
	concreteReward := rawReward.Mul(validatorPower)

	coinReward := sdk.DecCoin{Denom: k.GetRewardDenom(ctx), Amount: concreteReward}
	return sdk.DecCoins{coinReward}
}

// DistributeBlockRewards distributes the computed reward to the block proposer
func (k Keeper) DistributeBlockRewards(ctx sdk.Context, validator exported.ValidatorI, reward sdk.DecCoins) error {
	brPool := k.GetBlockRewardsPool(ctx)

	// Check if the pool has enough funds

	if brPool.AmountOf(k.GetRewardDenom(ctx)).GTE(reward.AmountOf(k.GetRewardDenom(ctx))) {
		brPool = brPool.Sub(reward)
		k.SetBlockRewardsPool(ctx, brPool)

		// Get his current reward and then add the new one
		currentRewards := k.DistributionKeeper.GetValidatorCurrentRewards(ctx, validator.GetOperator())
		currentRewards.Rewards = currentRewards.Rewards.Add(reward)

		// Set the just earned reward
		k.DistributionKeeper.SetValidatorCurrentRewards(ctx, validator.GetOperator(), currentRewards)
	} else {
		return sdk.ErrInsufficientFunds("Pool hasn't got enough funds to supply validator's rewards")
	}
	return nil
}
