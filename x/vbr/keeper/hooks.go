package keeper

import (
	epochstypes "github.com/commercionetwork/commercionetwork/x/epochs/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
}

func (k Keeper) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	params := k.GetParamSet(ctx)
	if epochIdentifier == params.DistrEpochIdentifier {
		// Get the number of active validators
		activeValidators := k.stakingKeeper.GetLastValidators(ctx)
		valNumber := int64(len(activeValidators))

		for _, validator := range activeValidators {
			// Compute the reward based on the number of validators, the validator's staked tokens and the total staked tokens
			reward := k.ComputeProposerReward(ctx, valNumber, validator, k.stakingKeeper.BondDenom(ctx), params)

			// Distribute the reward to the block proposer
			// TODO: Don't panic if pool is empty or not enough to distribute something
			// _ = k.DistributeBlockRewards(ctx, validator, reward)
			if err := k.DistributeBlockRewards(ctx, validator, reward); err != nil {
				panic(err)
			}
		}
	}
}

// ___________________________________________________________________________________________________

// Hooks wrapper struct for incentives keeper
type Hooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// epochs hooks
func (h Hooks) BeforeEpochStart(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h Hooks) AfterEpochEnd(ctx sdk.Context, epochIdentifier string, epochNumber int64) {
	h.k.AfterEpochEnd(ctx, epochIdentifier, epochNumber)
}
