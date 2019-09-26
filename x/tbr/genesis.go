package tbr

import (
	"errors"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	RewardDenom string       `json:"reward_denom"`
	PoolAmount  sdk.DecCoins `json:"pool_amount"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState(rewardDenom string) GenesisState {
	return GenesisState{
		RewardDenom: rewardDenom,
		PoolAmount:  sdk.DecCoins{},
	}
}

// InitGenesis sets the initial Block Reward Pool amount for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	if data.PoolAmount != nil {
		keeper.SetBlockRewardsPool(ctx, data.PoolAmount)
	}

	keeper.SetRewardDenom(ctx, data.RewardDenom)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		RewardDenom: keeper.GetRewardDenom(ctx),
		PoolAmount:  keeper.GetBlockRewardsPool(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if len(strings.TrimSpace(data.RewardDenom)) == 0 {
		return errors.New("transaction block reward reward denom cannot be empty")
	}

	if data.PoolAmount.Empty() {
		return errors.New("transaction block reward pool cannot be empty")
	}

	if data.PoolAmount.IsAnyNegative() {
		return errors.New("transaction block reward pool cannot contain negative funds")
	}

	return nil
}
