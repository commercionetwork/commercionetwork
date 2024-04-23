package keeper

import (
	"testing"

	//"cosmossdk.io/simapp"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v4/testing/simapp"
	"github.com/stretchr/testify/require"
)

func TestNewQuerier_default(t *testing.T) {
	t.Run("default request", func(t *testing.T) {
		k, ctx := SetupKeeper(t)

		app := simapp.Setup(false)
		legacyAmino := app.LegacyAmino()
		querier := NewQuerier(*k, legacyAmino)
		path := []string{"abcd"}
		_, err := querier(ctx, path, abci.RequestQuery{})
		require.Error(t, err)
	})
}

func TestNewQuerier_queryGetBlockRewardsPoolFunds(t *testing.T) {
	t.Run("queryGetBlockRewardsPoolFunds", func(t *testing.T) {
		k, ctx := SetupKeeper(t)

		expected := sdk.NewDecCoinsFromCoins(types.ValidMsgIncrementBlockRewardsPool.Amount...)
		k.SetTotalRewardPool(ctx, expected)
		amount, _ := expected.TruncateDecimal()
		err := k.MintVBRTokens(ctx, sdk.NewCoins(amount...))
		require.NoError(t, err)

		app := simapp.Setup(false)
		legacyAmino := app.LegacyAmino()

		querier := NewQuerier(*k, legacyAmino)
		path := []string{types.QueryBlockRewardsPoolFunds}
		gotBz, err := querier(ctx, path, abci.RequestQuery{})

		var got sdk.DecCoins
		legacyAmino.MustUnmarshalJSON(gotBz, &got)

		require.NoError(t, err)
		require.Equal(t, expected, got)
	})
}

func TestNewQuerier_queryParams(t *testing.T) {
	t.Run("queryParams", func(t *testing.T) {
		k, ctx := SetupKeeper(t)

		expected := types.NewParams(types.ValidMsgSetParams.DistrEpochIdentifier, types.ValidMsgSetParams.EarnRate)
		err := k.SetParamSet(ctx, expected)
		require.NoError(t, err)

		app := simapp.Setup(false)
		legacyAmino := app.LegacyAmino()
		querier := NewQuerier(*k, legacyAmino)
		path := []string{types.QueryParams}
		gotBz, err := querier(ctx, path, abci.RequestQuery{})

		var got types.Params
		legacyAmino.MustUnmarshalJSON(gotBz, &got)
		require.NoError(t, err)
		require.Equal(t, expected, got)
	})
}
