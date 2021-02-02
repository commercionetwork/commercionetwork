package keeper

import (
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQuerier_queryCurrentUpgrade(t *testing.T) {

	tests := []struct {
		name       string
		planExists bool
	}{
		{"a correctly scheduled plan exists",
			true,
		},
		{"no plan scheduled",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc, ctx, k := SetupTestInput(true)

			var plan upgrade.Plan
			if tt.planExists {
				plan = upgrade.Plan{
					Name:   "name",
					Info:   "info info info",
					Height: 100,
				}
				err := k.UpgradeKeeper.ScheduleUpgrade(ctx, plan)
				require.NoError(t, err)

				currentUpgrade, err := queryCurrentUpgrade(ctx, k)

				require.NoError(t, err)

				if tt.planExists {
					var p upgrade.Plan

					require.NotPanics(t, func() {
						cdc.MustUnmarshalJSON(currentUpgrade, &p)
					})
					require.Equal(t, plan, p)
				} else {
					require.Nil(t, currentUpgrade)
				}

			}
		})
	}
}
