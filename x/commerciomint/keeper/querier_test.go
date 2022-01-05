package keeper

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func Test_queryConversionRate(t *testing.T) {
	t.Run("expected sdk.NewDec(2)", func(t *testing.T) {
		ctx, _, _, k := SetupTestInput()

		expected := sdk.NewDec(2)

		app := simapp.Setup(false)
		legacyAmino := app.LegacyAmino()

		gotBz, err := queryConversionRate(ctx, k, legacyAmino)
		require.NoError(t, err)

		var got sdk.Dec

		legacyAmino.MustUnmarshalJSON(gotBz, &got)

		require.Equal(t, expected, got)
	})
}
func Test_queryFreezePeriod(t *testing.T) {

	tests := []struct {
		name            string
		setFreezePeriod bool
	}{
		{
			name:            "ok",
			setFreezePeriod: true,
		},
		{
			name:            "empty",
			setFreezePeriod: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			var expected, got time.Duration

			if tt.setFreezePeriod {
				expected = time.Minute
				require.NoError(t, k.UpdateFreezePeriod(ctx, expected))
			}

			app := simapp.Setup(false)
			legacyAmino := app.LegacyAmino()

			gotBz, err := queryFreezePeriod(ctx, k, legacyAmino)
			require.NoError(t, err)

			legacyAmino.MustUnmarshalJSON(gotBz, &got)

			require.Equal(t, expected, got)
		})
	}
}
