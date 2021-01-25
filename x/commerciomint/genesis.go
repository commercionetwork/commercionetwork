package commerciomint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/keeper"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

// GenesisState - documents genesis state
type GenesisState struct {
	Positions           []types.Position `json:"positions"`
	LiquidityPoolAmount sdk.Coins        `json:"pool_amount"`
	CollateralRate      sdk.Dec          `json:"collateral_rate"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Positions:           []types.Position{},
		LiquidityPoolAmount: sdk.Coins{},
		CollateralRate:      sdk.NewDec(2),
	}
}

// InitGenesis sets documents information for genesis.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, supplyKeeper supply.Keeper, data GenesisState) {

	// Get the module account
	moduleAcc := keeper.GetMintModuleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// Get the initial pool coins
	if moduleAcc.GetCoins().IsZero() {
		if err := moduleAcc.SetCoins(data.LiquidityPoolAmount); err != nil {
			panic(err)
		}
		supplyKeeper.SetModuleAccount(ctx, moduleAcc)
	}

	err := keeper.SetConversionRate(ctx, data.CollateralRate)
	if err != nil {
		panic(err)
	}

	// Add the existing ETPs
	for _, position := range data.Positions {
		keeper.SetPosition(ctx, position)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) GenesisState {
	return GenesisState{
		Positions:           keeper.GetAllPositions(ctx),
		LiquidityPoolAmount: keeper.GetLiquidityPoolAmount(ctx),
		CollateralRate:      keeper.GetConversionRate(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(state GenesisState) error {
	for _, position := range state.Positions {
		err := position.Validate()
		if err != nil {
			return err
		}
	}
	return types.ValidateConversionRate(state.CollateralRate)
}
