package commerciomint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/keeper"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

// GenesisState - docs genesis state
type GenesisState struct {
	Cdps                []types.Position `json:"cdps"`
	LiquidityPoolAmount sdk.Coins        `json:"pool_amount"`
	CreditsDenom        string           `json:"credits_denom"`
	CollateralRate      sdk.Dec          `json:"collateral_rate"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState(creditsDenom string) GenesisState {
	return GenesisState{
		Cdps:                []types.Position{},
		LiquidityPoolAmount: sdk.Coins{},
		CreditsDenom:        creditsDenom,
		CollateralRate:      sdk.NewDec(2),
	}
}

// InitGenesis sets docs information for genesis.
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

	// Add the existing CDPs
	for _, cdp := range data.Cdps {
		keeper.SetPosition(ctx, cdp)
	}

	// Set the stable credits denom
	keeper.SetCreditsDenom(ctx, data.CreditsDenom)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) GenesisState {
	return GenesisState{
		Cdps:                keeper.GetAllPositions(ctx),
		LiquidityPoolAmount: keeper.GetLiquidityPoolAmount(ctx),
		CreditsDenom:        keeper.GetCreditsDenom(ctx),
		CollateralRate:      keeper.GetCollateralRate(ctx),
	}
}

// ValidateGenesis performs basic validation of genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(state GenesisState) error {
	for _, cdp := range state.Cdps {
		err := cdp.Validate()
		if err != nil {
			return err
		}
	}
	return types.ValidateCollateralRate(state.CollateralRate)
}
