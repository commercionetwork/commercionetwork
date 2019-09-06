package keeper

import (
	"math/big"

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

//Utility method to set Block Reward Pool
func (k Keeper) setBlockRewardsPool(ctx sdk.Context, updatedPool types.BlockRewardsPool) {
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.BlockRewardsPoolPrefix), k.Cdc.MustMarshalBinaryBare(&updatedPool))
}

//Utility method to get Block Reward Pool
func (k Keeper) getBrPool(ctx sdk.Context) types.BlockRewardsPool {
	var brPool types.BlockRewardsPool
	store := ctx.KVStore(k.StoreKey)
	brpBz := store.Get([]byte(types.BlockRewardsPoolPrefix))
	if brpBz == nil {
		return types.InitBlockRewardsPool()
	}
	k.Cdc.MustUnmarshalBinaryBare(brpBz, &brPool)
	return brPool
}

func (k Keeper) getFunders(ctx sdk.Context) types.Funders {
	var funders types.Funders
	store := ctx.KVStore(k.StoreKey)
	fundersBz := store.Get([]byte(types.BlockRewardsPoolFundersPrefix))
	k.Cdc.MustUnmarshalBinaryBare(fundersBz, &funders)
	return funders
}

func (k Keeper) setFunders(ctx sdk.Context, updatedFunders types.Funders) {
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.BlockRewardsPoolFundersPrefix), k.Cdc.MustMarshalBinaryBare(&updatedFunders))
}

//Add new funder to the list of the authorized funders of the block rewards pool
func (k Keeper) AddBlockRewardsPoolFunder(ctx sdk.Context, funder types.Funder) {
	funders := k.getFunders(ctx)
	funders.AppendFunderIfMissing(funder)
	k.setFunders(ctx, funders)
}

//Return all the block rewards pool's funders
func (k Keeper) GetBlockRewardsPoolFunders(ctx sdk.Context) types.Funders {
	return k.getFunders(ctx)
}

//Increase the Block Rewards Pool with the specified coin amount
func (k Keeper) IncrementBlockRewardsPool(ctx sdk.Context, funder sdk.AccAddress, amount sdk.Coin) {
	bk := k.BankKeeper
	brAmount := sdk.Coins{amount}
	brPool := types.InitBlockRewardsPool()

	if bk.HasCoins(ctx, funder, brAmount) {
		brPool = k.getBrPool(ctx)
		if brPool.Funds.IsZero() {
			brPool.Funds.Add(sdk.NewDecCoins(brAmount))
			k.setBlockRewardsPool(ctx, brPool)
		} else {
			brPool.Funds.Add(sdk.NewDecCoins(brAmount))
			k.setBlockRewardsPool(ctx, brPool)
		}
	}
}

//Return the Block Rewards Pool if exists
func (k Keeper) GetBlockRewardsPool(ctx sdk.Context) types.BlockRewardsPool {
	return k.getBrPool(ctx)
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
//TODO I used big.Int to mantain the precision in these operations and in prevision of very large numbers,
// I don't know if it's the right choise, need to discuss
var (
	TPY = big.NewInt(25000) //Tokens Per Year
	//TODO when sdk.Dec get an implementation with big.Float, switch this to 365.24 to count even the bissextile year
	DPY = big.NewInt(365) // Days Per Year
	HPD = big.NewInt(24)  // 	Hours Per Day
	MPD = big.NewInt(60)  // 	Minutes Per Days
	BPM = big.NewInt(12)  // 	Blocks Per Minutes
)

/*
Compute the Raw Reward for proposer, assuming that Raw Reward is the Reward(n, V) equation's result
without the last multiplication between the (Stake/TotalStake) value
*/
func computeRawReward(validatorsNumber sdk.Int) sdk.Dec {
	var firstMember big.Int
	var secondMember big.Int
	var thirdMember big.Int

	averageValidatorsNumber := big.NewInt(100)
	vNumber := big.NewInt(validatorsNumber.Int64())

	firstMember.Mul(TPY, big.NewInt(1000000))

	secondMember.Mul(DPY, HPD)
	secondMember.Mul(&secondMember, MPD)
	secondMember.Mul(&secondMember, BPM)

	thirdMember.Quo(averageValidatorsNumber, vNumber)

	firstMember.Quo(&firstMember, &secondMember)

	rawReward := firstMember.Mul(&firstMember, &thirdMember)

	return sdk.NewDecFromBigInt(rawReward)
}

//Compute the final reward for the validator block's proposer
func (k Keeper) ComputeProposerReward(ctx sdk.Context, validatorNumber sdk.Int, proposer exported.ValidatorI,
	totalStakedTokens sdk.Int) sdk.DecCoins {

	//Get the raw reward for proposer
	rawReward := computeRawReward(validatorNumber)

	//Retrieve staked tokens by proposer
	validatorStakedTokens := proposer.GetBondedTokens()
	//compute his validation power
	validatorPower := validatorStakedTokens.Quo(totalStakedTokens).ToDec()

	//calculate the final reward for this proposer
	concreteReward := rawReward.Mul(validatorPower)

	coinReward := sdk.DecCoin{Denom: types.DefaultBondDenom, Amount: concreteReward}

	return append(sdk.DecCoins{}, coinReward)
}

//Distribute the computed reward to the block proposer
func (k Keeper) DistributeBlockRewards(ctx sdk.Context, validator exported.ValidatorI, reward sdk.DecCoins) {

	var brPool types.BlockRewardsPool

	k.GetBlockRewardsPool(ctx)

	//Check if the pool has enough funds
	if brPool.Funds.AmountOf(types.DefaultBondDenom).GTE(reward.AmountOf(types.DefaultBondDenom)) {
		brPool.Funds = brPool.Funds.Sub(reward)
		k.setBlockRewardsPool(ctx, brPool)
	}

	//Get his current reward and then add the new one
	currentRewards := k.DistributionKeeper.GetValidatorCurrentRewards(ctx, validator.GetOperator())
	currentRewards.Rewards = currentRewards.Rewards.Add(reward)

	//Set the just earned reward
	k.DistributionKeeper.SetValidatorCurrentRewards(ctx, validator.GetOperator(), currentRewards)
}
