package mint

import (
	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState - docs genesis state
type GenesisState struct {
	Cdps                types.Cdps `json:"cdps"`
	LiquidityPoolAmount sdk.Coins  `json:"pool_amount"`
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// InitGenesis sets docs information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetLiquidityPool(ctx, data.LiquidityPoolAmount)

	for _, cdp := range data.Cdps {
		keeper.AddCdp(ctx, cdp)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return GenesisState{
		Cdps:                keeper.GetTotalCdps(ctx),
		LiquidityPoolAmount: keeper.GetLiquidityPool(ctx),
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
