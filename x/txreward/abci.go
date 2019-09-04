package txreward

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/txreward/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/**
Retrieve all the active validators, and based on how many are them, calculate the reward ONLY for the block proposer.
The reward is distributed every time a new block is created (one every 5 sec).
*/
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {

	//Get the number of active validators
	activeValidators := k.StakeKeeper.GetLastValidators(ctx)
	valNumber := sdk.NewInt(int64(len(activeValidators)))

	//Get the validator who proposed the block
	previousProposer := k.DistributionKeeper.GetPreviousProposerConsAddr(ctx)
	validator := k.StakeKeeper.ValidatorByConsAddr(ctx, previousProposer)

	//Compute the reward based on the number of validators, the validator's staked tokens and the total staked tokens
	reward := k.ComputeProposerReward(ctx, valNumber, validator, k.StakeKeeper.TotalBondedTokens(ctx))

	//Get the block proposer
	if ctx.BlockHeight() > 1 {

		//Distribute the reward to the block proposer
		k.DistributeBlockRewards(ctx, validator, reward)
	}

	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.DistributionKeeper.SetPreviousProposerConsAddr(ctx, consAddr)
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	//FILL WITH OPERATIONS TO PERFORM AT BLOCK'S END
}
