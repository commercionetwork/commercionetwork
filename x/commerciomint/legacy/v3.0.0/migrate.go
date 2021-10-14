package v3_0_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	v220commerciomint "github.com/commercionetwork/commercionetwork/x/commerciomint/legacy/v2.2.0"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

// Migrate accepts exported genesis state from v2.2.0 and migrates it to v3.0.0
func Migrate(oldGenState v220commerciomint.GenesisState) *types.GenesisState {

	var postions []*types.Position
	for _, oldPosition := range oldGenState.Positions {
		var position types.Position

		position.Owner = oldPosition.Owner.String()
		position.Collateral = oldPosition.Collateral.Int64()
		position.CreatedAt = &oldPosition.CreatedAt
		position.Credits = &oldPosition.Credits
		position.ExchangeRate = oldPosition.ExchangeRate

		postions = append(postions, &position)
	}

	var coins []*sdk.Coin
	for _, coin := range oldGenState.LiquidityPoolAmount {
		coins = append(coins, &coin)
	}

	return &types.GenesisState{
		Positions:      postions,
		PoolAmount:     coins,
		CollateralRate: &sdk.DecProto{Dec: oldGenState.CollateralRate},
		FreezePeriod:   oldGenState.FreezePeriod.String(),
	}

}
