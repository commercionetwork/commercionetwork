package vbr

import (
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/vbr/keeper"
)

// BeginBlocker retrieves all the active validators, and based on how many are of them, calculate
// the reward ONLY for the block proposer on every begin block.
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper, stakeKeeper staking.Keeper) {

	// determine the total power signing the block
	// This calculate the real power in previus block
	// It consider the validators that were present in precommit process
	// sum all voting power and all voting power present
	/*var previousTotalPower, sumPreviousPrecommitPower int64
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		previousTotalPower += voteInfo.Validator.Power
		if voteInfo.SignedLastBlock {
			sumPreviousPrecommitPower += voteInfo.Validator.Power
		}
	}*/

	// Get the number of active validators
	activeValidators := stakeKeeper.GetLastValidators(ctx)
	valNumber := int64(len(activeValidators))

	// Get the block height
	if ctx.BlockHeight() > 1 {
		// Get the validator who proposed the block
		previousProposer := k.GetPreviousProposerConsAddr(ctx)

		// Retrieve the validator from its consensus address
		validator := stakeKeeper.ValidatorByConsAddr(ctx, previousProposer)

		// Compute the reward based on the number of validators, the validator's staked tokens and the total staked tokens
		reward := k.ComputeProposerReward(ctx, valNumber, validator, stakeKeeper.TotalBondedTokens(ctx))

		// Distribute the reward to the block proposer
		if err := k.DistributeBlockRewards(ctx, validator, reward); err != nil {
			panic(err)
		}
	}

	// Reward all
	if k.IsDailyWighDrawBlock(ctx.BlockHeight()) {
		k.WithdrawAllRewards(ctx, stakeKeeper)
	}

}
