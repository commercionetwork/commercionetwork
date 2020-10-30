package commerciomint_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciomint"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/keeper"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

func TestInitGenesis(t *testing.T) {
	got := commerciomint.DefaultGenesisState()
	require.NoError(t, commerciomint.ValidateGenesis(got))
	ctx, _, _, _, sk, k := keeper.SetupTestInput()
	require.Equal(t, commerciomint.GenesisState{Positions: []types.Position{}, LiquidityPoolAmount: sdk.NewCoins(), CollateralRate: sdk.NewInt(2)}, got)
	commerciomint.InitGenesis(ctx, k, sk, got)
	export := commerciomint.ExportGenesis(ctx, k)
	require.Equal(t, commerciomint.GenesisState{Positions: []types.Position(nil), LiquidityPoolAmount: sdk.Coins(nil), CollateralRate: sdk.NewInt(2)}, export)

	credits, err := sdk.ParseCoin("5test")
	require.NoError(t, err)
	testCdp := types.Position{Owner: []byte("test"), CreatedAt: time.Now(), Collateral: sdk.NewInt(10), Credits: credits, ExchangeRate: sdk.NewInt(2)}
	k.SetPosition(ctx, testCdp)
	export = commerciomint.ExportGenesis(ctx, k)

	require.True(t, export.Positions[0].Equals(testCdp))
	require.Equal(t, export.LiquidityPoolAmount, sdk.Coins(nil))
	require.Equal(t, export.CollateralRate, sdk.NewInt(2))
}
