package tbr

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	PoolAmount sdk.DecCoins `json:"pool_amount"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets the initial Block Reward Pool amount for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	//Set the intial amount of Block Rewards Pool
	keeper.SetBlockRewardsPool(ctx, data.PoolAmount)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{PoolAmount: keeper.GetBlockRewardsPool(ctx)}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.PoolAmount.IsAnyNegative() {
		return fmt.Errorf("negative Funds in block reward pool, is %v", data.PoolAmount)
	}

	return nil
}
