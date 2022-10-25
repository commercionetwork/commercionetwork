package commerciomint_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

func TestInitGenesis(t *testing.T) {
	got := types.DefaultGenesis()
	require.NoError(t, got.Validate())
	/*ctx, _, _, sk, k := keeper.SetupTestInput()
	require.Equal(t, types.GenesisState{Positions: []*types.Position{}, PoolAmount: sdk.NewCoins(), CollateralRate: sdk.NewDec(1), FreezePeriod: types.DefaultFreezePeriod}, got)
	commerciomint.InitGenesis(ctx, k, sk, got)
	export := commerciomint.ExportGenesis(ctx, k)
	require.Equal(t, types.GenesisState{Positions: []*types.Position(nil), PoolAmount: sdk.Coins(nil), CollateralRate: sdk.NewDec(1), FreezePeriod: types.DefaultFreezePeriod}, export)

	credits, err := sdk.ParseCoin("5test")
	require.NoError(t, err)
	testEtp := types.Position{Owner: []byte("test"), CreatedAt: time.Now(), Collateral: sdk.NewInt(10), Credits: credits, ExchangeRate: sdk.NewDec(1)}
	k.SetPosition(ctx, testEtp)
	export = commerciomint.ExportGenesis(ctx, k)

	require.True(t, export.Positions[0].Equals(testEtp))
	require.Equal(t, export.PoolAmount, sdk.Coins(nil))
	require.Equal(t, export.CollateralRate, sdk.NewDec(1))*/
}
