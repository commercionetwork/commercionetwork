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

	// Record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.SetPreviousProposerConsAddr(ctx, consAddr)

}
