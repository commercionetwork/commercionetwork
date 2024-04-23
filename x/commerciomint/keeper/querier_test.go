package keeper

import (
	"testing"
	"time"

	//"cosmossdk.io/simapp"
	"cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	"github.com/cosmos/ibc-go/v4/testing/simapp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewQuerier_queryGetEtp(t *testing.T) {

	tests := []struct {
		name              string
		positionsToCreate []types.Position
		shouldFind        bool
	}{
		{
			name: "find among one",
			positionsToCreate: []types.Position{
				testEtp,
			},
			shouldFind: true,
		},
		{
			name: "find among many",
			positionsToCreate: []types.Position{
				testEtp1,
				testEtp,
				testEtp2,
				testEtpAnotherOwner,
			},
			shouldFind: true,
		},
		{
			name:              "empty",
			positionsToCreate: []types.Position{},
			shouldFind:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, p := range tt.positionsToCreate {
				require.NoError(t, k.SetPosition(ctx, p))
			}

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()
			querier := NewQuerier(k, legacyAmino)
			path := []string{types.QueryGetEtpRest, testEtp.ID}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})

			var got types.Position

			if tt.shouldFind {
				legacyAmino.MustUnmarshalJSON(gotBz, &got)
				require.NoError(t, err)
				require.Equal(t, testEtp, got)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func Test_NewQuerier_queryGetEtpsByOwner(t *testing.T) {

	t.Run("invalid address", func(t *testing.T) {
		ctx, _, _, k := SetupTestInput()

		app := simapp.Setup(false)
		legacyAmino := app.LegacyAmino()
		querier := NewQuerier(k, legacyAmino)
		path := []string{types.QueryGetEtpsByOwnerRest, ""}

		_, err := querier(ctx, path, abci.RequestQuery{})
		require.Error(t, err)

	})

	tests := []struct {
		name              string
		positionsToCreate []types.Position
		expected          []types.Position
	}{
		{
			name: "find none",
			positionsToCreate: []types.Position{
				testEtpAnotherOwner,
			},
			expected: []types.Position{},
		},
		{
			name: "find among one",
			positionsToCreate: []types.Position{
				testEtp,
			},
			expected: []types.Position{
				testEtp,
			},
		},
		{
			name: "find among it and another",
			positionsToCreate: []types.Position{
				testEtp,
				testEtpAnotherOwner,
			},
			expected: []types.Position{
				testEtp,
			},
		},
		{
			name: "find all by the same",
			positionsToCreate: []types.Position{
				testEtp,
				testEtp1,
				testEtp2,
				testEtpAnotherOwner,
			},
			expected: []types.Position{
				testEtp,
				testEtp1,
				testEtp2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, p := range tt.positionsToCreate {
				require.NoError(t, k.SetPosition(ctx, p))
			}

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()
			querier := NewQuerier(k, legacyAmino)
			path := []string{types.QueryGetEtpsByOwnerRest, testEtp.Owner}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})

			var got []types.Position
			legacyAmino.MustUnmarshalJSON(gotBz, &got)
			require.NoError(t, err)

			for _, etp := range tt.expected {
				require.Contains(t, got, etp)
			}

		})
	}
}

func Test_NewQuerier_queryGetAllEtp(t *testing.T) {

	tests := []struct {
		name              string
		positionsToCreate []types.Position
	}{
		{
			name:              "empty",
			positionsToCreate: []types.Position{},
		},
		{
			name: "find one",
			positionsToCreate: []types.Position{
				testEtp,
			},
		},
		{
			name: "find many",
			positionsToCreate: []types.Position{
				testEtp,
				testEtpAnotherOwner,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, p := range tt.positionsToCreate {
				require.NoError(t, k.SetPosition(ctx, p))
			}

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()
			querier := NewQuerier(k, legacyAmino)
			path := []string{types.QueryGetallEtpsRest}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})

			var got []types.Position

			legacyAmino.MustUnmarshalJSON(gotBz, &got)
			require.NoError(t, err)

			for _, etp := range tt.positionsToCreate {
				require.Contains(t, got, etp)
			}

		})
	}
}

func Test_NewQuerier_queryConversionRate(t *testing.T) {
	t.Run("expected math.LegacyNewDec(2)", func(t *testing.T) {
		ctx, _, _, k := SetupTestInput()

		expected := math.LegacyNewDec(2)

		app := simapp.Setup(false)
		legacyAmino := app.LegacyAmino()
		querier := NewQuerier(k, legacyAmino)
		path := []string{types.QueryConversionRateRest}
		gotBz, err := querier(ctx, path, abci.RequestQuery{})
		require.NoError(t, err)

		var got math.LegacyDec
		legacyAmino.MustUnmarshalJSON(gotBz, &got)
		require.Equal(t, expected, got)
	})
}
func Test_NewQuerier_queryFreezePeriod(t *testing.T) {

	tests := []struct {
		name            string
		setFreezePeriod bool
	}{
		{
			name:            "ok",
			setFreezePeriod: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			var expected, got time.Duration

			if tt.setFreezePeriod {
				expected = time.Minute

				params := validParams
				params.FreezePeriod = time.Minute
				assert.NotEqual(t, params.FreezePeriod, validParams.FreezePeriod)
				require.NoError(t, k.UpdateParams(ctx, params))
			}

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()
			querier := NewQuerier(k, legacyAmino)
			path := []string{types.QueryFreezePeriodRest}
			gotBz, err := querier(ctx, path, abci.RequestQuery{})
			require.NoError(t, err)

			legacyAmino.MustUnmarshalJSON(gotBz, &got)
			require.Equal(t, expected, got)
		})
	}
}

func Test_NewQuerier_queryGetParams(t *testing.T) {

	t.Run("ok", func(t *testing.T) {
		ctx, _, _, k := SetupTestInput()

		app := simapp.Setup(false)
		legacyAmino := app.LegacyAmino()
		querier := NewQuerier(k, legacyAmino)
		path := []string{types.QueryGetParamsRest}
		gotBz, err := querier(ctx, path, abci.RequestQuery{})
		require.NoError(t, err)

		var got types.Params
		legacyAmino.MustUnmarshalJSON(gotBz, &got)
		require.Equal(t, validParams, got)
	})

}

func Test_NewQuerier_default(t *testing.T) {

	t.Run("default request", func(t *testing.T) {
		ctx, _, _, k := SetupTestInput()

		app := simapp.Setup(false)
		legacyAmino := app.LegacyAmino()
		querier := NewQuerier(k, legacyAmino)
		path := []string{"abcd"}
		_, err := querier(ctx, path, abci.RequestQuery{})
		require.Error(t, err)
	})
}
