package mint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

// GenesisState - docs genesis state
type GenesisState struct {
	Cdps                Cdps      `json:"cdps"`
	LiquidityPoolAmount sdk.Coins `json:"pool_amount"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets docs information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, supplyKeeper supply.Keeper, data GenesisState) {
	moduleAcc := keeper.GetMintModuleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", ModuleName))
	}

	if moduleAcc.GetCoins().IsZero() {
		if err := moduleAcc.SetCoins(data.LiquidityPoolAmount); err != nil {
			panic(err)
		}
		supplyKeeper.SetModuleAccount(ctx, moduleAcc)
	}

	for _, cdp := range data.Cdps {
		keeper.AddCdp(ctx, cdp)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	moduleAcc := keeper.GetMintModuleAccount(ctx)
	return GenesisState{
		Cdps:                keeper.GetTotalCdps(ctx),
		LiquidityPoolAmount: moduleAcc.GetCoins(),
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
	return nil
}
