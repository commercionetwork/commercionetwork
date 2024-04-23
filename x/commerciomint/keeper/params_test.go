package keeper

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeeper_GetConversionRate(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	params := validParams
	params.ConversionRate = math.LegacyNewDec(3)
	assert.NotEqual(t, params.ConversionRate, validParams.ConversionRate)
	require.NoError(t, k.UpdateParams(ctx, params))
	require.Equal(t, params.ConversionRate, k.GetConversionRate(ctx))
}

func TestKeeper_GetFreezePeriod(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	params := validParams
	params.FreezePeriod = time.Minute
	assert.NotEqual(t, params.FreezePeriod, validParams.FreezePeriod)

	require.NoError(t, k.UpdateParams(ctx, params))
	require.Equal(t, params.FreezePeriod, k.GetFreezePeriod(ctx))
}

func TestKeeper_GetParams(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	k.UpdateParams(ctx, types.Params{
		ConversionRate: types.DefaultConversionRate,
		FreezePeriod:   types.DefaultFreezePeriod,
	})
	require.Equal(t, types.Params{ConversionRate: types.DefaultConversionRate,
		FreezePeriod: types.DefaultFreezePeriod}, k.GetParams(ctx))
}

func TestKeeper_UpdateParams(t *testing.T) {

	type args struct {
		params types.Params
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				params: types.Params{
					ConversionRate: math.LegacyNewDec(5),
					FreezePeriod:   time.Minute,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid conversion rate",
			args: args{
				params: types.Params{
					ConversionRate: invalidConversionRate,
					FreezePeriod:   validFreezePeriod,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid freeze period",
			args: args{
				params: types.Params{
					ConversionRate: validConversionRate,
					FreezePeriod:   invalidFreezePeriod,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			oldParams := k.GetParams(ctx)
			assert.NotEqual(t, oldParams, tt.args.params)

			if err := k.UpdateParams(ctx, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("Keeper.UpdateParams() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				require.Equal(t, oldParams, k.GetParams(ctx))
			} else {
				require.Equal(t, tt.args.params, k.GetParams(ctx))
			}

		})
	}
}
