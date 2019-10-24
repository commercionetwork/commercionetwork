package tbr

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	PoolAmount       sdk.DecCoins `json:"pool_amount"`
	YearlyPoolAmount sdk.DecCoins `json:"yearly_pool_amount"`
	YearlyRemains    sdk.DecCoins `json:"yearly_remains"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets the initial Block Reward Pool amount for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// Set the reward pool - Should never be nil as its validated inside the ValidateGenesis method
	keeper.SetTotalRewardPool(ctx, data.PoolAmount)

	// Default yearly reward pool amount
	if data.YearlyPoolAmount.Empty() {
		data.YearlyPoolAmount = data.PoolAmount.QuoDec(sdk.NewDec(5))
	}

	// Default yearly remains
	if data.YearlyRemains.Empty() {
		data.YearlyRemains = data.YearlyPoolAmount
	}

	// Set the yearly reward pool and the remains
	keeper.SetYearlyRewardPool(ctx, data.YearlyPoolAmount)
	keeper.SetYearlyPoolRemains(ctx, data.YearlyRemains)

	// Compute and set the year number
	yearNumber := keeper.ComputeYearFromBlockHeight(ctx.BlockHeight())
	keeper.SetYearNumber(ctx, yearNumber)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		PoolAmount:       keeper.GetTotalRewardPool(ctx),
		YearlyPoolAmount: keeper.GetYearlyRewardPool(ctx),
		YearlyRemains:    keeper.GetRemainingYearlyPool(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if data.PoolAmount == nil || data.PoolAmount.Empty() {
		return errors.New("transaction block reward pool cannot be empty")
	}

	if !data.PoolAmount.IsValid() {
		return errors.New(fmt.Sprintf("invalid transaction block reward pool: %s", data.PoolAmount.String()))
	}

	if !data.YearlyPoolAmount.IsValid() {
		return errors.New(fmt.Sprintf("invalid yearly transaction block reward pool: %s", data.YearlyPoolAmount.String()))
	}

	return nil
}
