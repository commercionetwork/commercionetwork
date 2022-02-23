package v3_0_0

import (
	"time"

	v220commerciomint "github.com/commercionetwork/commercionetwork/x/commerciomint/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Migrate accepts exported genesis state from v2.2.0 and migrates it to v3.0.0
func Migrate(v220GenState v220commerciomint.GenesisState) *types.GenesisState {

	var positions []*types.Position
	for _, v220Position := range v220GenState.Positions {
		var timePosition time.Time
		timePosition = v220Position.CreatedAt
		var credits sdk.Coin
		credits = v220Position.Credits

		positions = append(positions, &types.Position{
			Owner:        v220Position.Owner.String(),
			Collateral:   v220Position.Collateral.Int64(),
			CreatedAt:    &timePosition,
			Credits:      &credits,
			ExchangeRate: v220Position.ExchangeRate,
			ID:           v220Position.ID,
		})
	}

	return &types.GenesisState{
		Positions:  positions,
		PoolAmount: v220GenState.LiquidityPoolAmount,
		Params: types.Params{
			ConversionRate: v220GenState.CollateralRate,
			FreezePeriod:   v220GenState.FreezePeriod,
		},
	}

}
