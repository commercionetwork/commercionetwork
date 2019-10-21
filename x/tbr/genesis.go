package tbr

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	PoolAmount       sdk.DecCoins `json:"pool_amount"`
	YearlyPoolAmount sdk.DecCoins `json:"yearly_pool_amount"`
	YearNumber       int64        `json:"year_number"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets the initial Block Reward Pool amount for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	if data.PoolAmount != nil {
		keeper.SetTotalRewardPool(ctx, data.PoolAmount)
		keeper.SetYearlyRewardPool(ctx, data.PoolAmount.QuoDec(sdk.NewDec(5)))
	}

	if data.YearlyPoolAmount != nil {
		keeper.SetYearlyRewardPool(ctx, data.YearlyPoolAmount)
	}

	keeper.SetYearNumber(ctx, data.YearNumber)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		PoolAmount: keeper.GetTotalRewardPool(ctx),
		YearNumber: keeper.GetYearNumber(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.PoolAmount.Empty() {
		return errors.New("transaction block reward pool cannot be empty")
	}

	if data.PoolAmount.IsAnyNegative() {
		return errors.New("transaction block reward pool cannot contain negative funds")
	}

	return nil
}
