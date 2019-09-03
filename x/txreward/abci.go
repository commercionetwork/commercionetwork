package txreward

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/txreward/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/**
Retrieve all the active validators, and based on how many are them, calculate the reward ONLY for the block proposer.
The reward is distributed every time a new block is created (one every 5 sec).
In particular,
with 100 or less active validators, we calculate the reward like this:

TPY = Tokens Per Year
Reward100 = TPY / (365 * 24 * 60 * 12)

Instead, if the active validators will be more than 100, we calculate the reward like this:

V = Validators Number (assuming it's greater than 100)

RewardVN = (Reward100 / V) * 100

Summarizing these formulas we obtain:

Reward(n, V) = TPY / (365 * 24 * 60 * 12) * 100 / V * STAKE / TOTALSTAKE

where:
Reward(n, V) indicates the reward for the n validator considering a set of V validators
V 			 indicates the Validators Number
STAKE 		 staked token's amount of n-esim validator
TOTALSTAKE	 indicates all staked token's amount of all validators
*/

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {

	//Get the number of active validators
	activeValidators := k.StakeKeeper.GetLastValidators(ctx)
	valNumber := sdk.NewInt(int64(len(activeValidators)))

	//Compute the reward based on the number of validators
	reward := k.ComputeValidatorsReward(ctx, valNumber)

	//Get the block proposer
	if ctx.BlockHeight() > 1 {
		//retrieve the validator who proposed the block
		previousProposer := k.DistributionKeeper.GetPreviousProposerConsAddr(ctx)

		//Distribute the reward to the block proposer
		k.DistributeBlockRewards(ctx, previousProposer, reward)
	}

	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.DistributionKeeper.SetPreviousProposerConsAddr(ctx, consAddr)
}

func EndBlocker(ctx sdk.Context, k keeper.Keeper) {
	//FILL WITH OPERATIONS TO PERFORM AT BLOCK'S END
}
