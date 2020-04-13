package commerciomint_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciomint"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/keeper"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

func TestInitGenesis(t *testing.T) {
	got := commerciomint.DefaultGenesisState("test")
	require.NoError(t, commerciomint.ValidateGenesis(got))
	ctx, _, _, _, sk, k := keeper.SetupTestInput()
	require.Equal(t, commerciomint.GenesisState{Positions: []types.Position{}, LiquidityPoolAmount: sdk.NewCoins(), CreditsDenom: "test", CollateralRate: sdk.NewDec(2)}, got)
	commerciomint.InitGenesis(ctx, k, sk, got)
	export := commerciomint.ExportGenesis(ctx, k)
	require.Equal(t, commerciomint.GenesisState{Positions: []types.Position{}, LiquidityPoolAmount: sdk.Coins(nil), CreditsDenom: "test", CollateralRate: sdk.NewDec(2)}, export)

	credits, err := sdk.ParseCoin("5test")
	require.NoError(t, err)
	testCdp := types.Position{Owner: []byte("test"), CreatedAt: 10, Deposit: sdk.NewCoins(sdk.NewInt64Coin("ucommercio", 10)), Credits: credits}
	k.SetPosition(ctx, testCdp)
	export = commerciomint.ExportGenesis(ctx, k)
	require.Equal(t, commerciomint.GenesisState{Positions: []types.Position{testCdp}, LiquidityPoolAmount: sdk.Coins(nil), CreditsDenom: "test", CollateralRate: sdk.NewDec(2)}, export)
}
