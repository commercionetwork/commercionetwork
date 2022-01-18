package v3_0_0

import (
	v220commerciomint "github.com/commercionetwork/commercionetwork/x/commerciomint/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

// Migrate accepts exported genesis state from v2.2.0 and migrates it to v3.0.0
func Migrate(oldGenState v220commerciomint.GenesisState) *types.GenesisState {

	var positions []*types.Position
	for _, oldPosition := range oldGenState.Positions {
		var position types.Position

		position.Owner = oldPosition.Owner.String()
		position.Collateral = oldPosition.Collateral.Int64()
		position.CreatedAt = &oldPosition.CreatedAt
		position.Credits = &oldPosition.Credits
		position.ExchangeRate = oldPosition.ExchangeRate
		position.ID = oldPosition.ID

		positions = append(positions, &position)
	}

	/*var coins []*sdk.Coin
	for _, coin := range oldGenState.LiquidityPoolAmount {
		coins = append(coins, &coin)
	}*/

	return &types.GenesisState{
		Positions:  positions,
		PoolAmount: oldGenState.LiquidityPoolAmount,
		Params: types.Params{
			ConversionRate: oldGenState.CollateralRate,
			FreezePeriod:   oldGenState.FreezePeriod,
		},
	}

}
