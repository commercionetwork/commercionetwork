package tbr

import (
	"github.com/cosmos/cosmos-sdk/x/staking"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/tbr/internal/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker retrieves all the active validators, and based on how many are of them, calculate
// the reward ONLY for the block proposer on every begin block.
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper, stakeKeeper staking.Keeper) {

	//Get the number of active validators
	activeValidators := stakeKeeper.GetLastValidators(ctx)
	valNumber := int64(len(activeValidators))

	//Get the validator who proposed the block
	previousProposer := k.DistributionKeeper.GetPreviousProposerConsAddr(ctx)

	//Get the block height
	if ctx.BlockHeight() > 1 {
		//retrieve validator from consesus address
		validator := stakeKeeper.ValidatorByConsAddr(ctx, previousProposer)

		//Compute the reward based on the number of validators, the validator's staked tokens and the total staked tokens
		reward := k.ComputeProposerReward(ctx, valNumber, validator, stakeKeeper.TotalBondedTokens(ctx))

		//Distribute the reward to the block proposer
		err := k.DistributeBlockRewards(ctx, validator, reward)
		if err != nil {
			panic(err)
		}

	}

	// record the proposer for when we payout on the next block
	consAddr := sdk.ConsAddress(req.Header.ProposerAddress)
	k.DistributionKeeper.SetPreviousProposerConsAddr(ctx, consAddr)
}

func EndBlocker(_ sdk.Context, _ keeper.Keeper) {}
