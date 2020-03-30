package creditrisk_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/creditrisk"
	"github.com/commercionetwork/commercionetwork/x/creditrisk/types"
)

func TestValidateGenesis(t *testing.T) {
	require.NoError(t, creditrisk.ValidateGenesis(types.GenesisState{}))

	invalidCoins := sdk.NewCoins(sdk.NewInt64Coin("test", 10))
	invalidCoins[0].Amount = sdk.NewInt(-10)
	require.Error(t, creditrisk.ValidateGenesis(types.GenesisState{
		Pool: invalidCoins,
	}))
}

func TestInitGenesis(t *testing.T) {
	ctx, _, k := SetupTestInput()

	// empty genesis data
	require.NotPanics(t, func() { creditrisk.InitGenesis(ctx, k, types.GenesisState{}) })
	require.True(t, k.GetPoolFunds(ctx).IsZero())

	// uninitialized context/keeper
	require.Panics(t, func() { creditrisk.InitGenesis(sdk.Context{}, k, types.GenesisState{}) })
	require.Panics(t, func() { creditrisk.InitGenesis(sdk.Context{}, creditrisk.Keeper{}, types.GenesisState{}) })

	// non empty genesis data
	require.NotPanics(t, func() {
		creditrisk.InitGenesis(ctx, k, types.GenesisState{Pool: sdk.NewCoins(sdk.NewInt64Coin("test", 10))})
	})
	require.Equal(t, int64(10), k.GetPoolFunds(ctx).AmountOf("test").Int64())
}
