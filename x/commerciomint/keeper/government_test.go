package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetConversionRate(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	rate := sdk.NewDec(3)
	require.NoError(t, k.UpdateConversionRate(ctx, rate))
	require.Equal(t, rate, k.GetConversionRate(ctx))
}

func TestKeeper_UpdateConversionRate(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	require.Error(t, k.UpdateConversionRate(ctx, sdk.NewDec(0)))
	require.Error(t, k.UpdateConversionRate(ctx, sdk.NewDec(-1)))
	require.NoError(t, k.UpdateConversionRate(ctx, sdk.NewDec(2)))
	rate := sdk.NewDec(3)
	require.NoError(t, k.UpdateConversionRate(ctx, rate))

	got := k.GetConversionRate(ctx)
	require.Equal(t, rate, got)
}

func TestKeeper_SetFreezePeriod(t *testing.T) {

	tests := []struct {
		name         string
		freezePeriod time.Duration
		wantErr      bool
	}{
		{"correctly set freeze period", time.Minute, false},
		{"invalid freeze period", -time.Second, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			if err := k.UpdateFreezePeriod(ctx, tt.freezePeriod); (err != nil) != tt.wantErr {
				t.Errorf("Keeper.SetFreezePeriod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				require.Equal(t, tt.freezePeriod, k.GetFreezePeriod(ctx))
			}
		})
	}
}

func TestKeeper_GetFreezePeriod(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	freeze := time.Minute
	require.NoError(t, k.UpdateFreezePeriod(ctx, freeze))
	require.Equal(t, freeze, k.GetFreezePeriod(ctx))
}
