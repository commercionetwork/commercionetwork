package types_test

import (
	"testing"
	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestDefaultGenesisState(t *testing.T) {
	expted := &types.GenesisState{}
	require.Equal(t, expted, types.DefaultGenesis())
}

func TestValidateGenesis(t *testing.T) {
	defGen := types.DefaultGenesis()
	err := defGen.Validate()
	require.NoError(t, err)

	// Test negative coins
	defGenNegativeLiquidity := types.DefaultGenesis()
	var coin sdk.Coin
	var coins sdk.Coins
	coin.Denom = "somecoin"
	coin.Amount = math.NewInt(-1)
	coins = append(coins, coin)
	defGenNegativeLiquidity.LiquidityPoolAmount = coins
	errNeg := defGenNegativeLiquidity.Validate()
	require.Error(t, errNeg)

}
