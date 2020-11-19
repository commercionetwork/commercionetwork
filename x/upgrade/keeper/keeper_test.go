package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestKeeper_GetUpgradePlan(t *testing.T) {

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
			_, ctx, k := SetupTestInput(true)

			var plan upgrade.Plan
			if tt.planExists {
				plan = upgrade.Plan{
					Name:   "name",
					Info:   "info info info",
					Height: 100,
				}
				err := k.UpgradeKeeper.ScheduleUpgrade(ctx, plan)
				require.NoError(t, err)
			} else {
				plan = upgrade.Plan{Name: "", Time: time.Time{}, Height: 0, Info: ""}
			}

			upgradePlan, havePlan := k.GetUpgradePlan(ctx)

			require.Equal(t, tt.planExists, havePlan)
			require.Equal(t, plan, upgradePlan)
		})
	}
}

func TestKeeper_ScheduleUpgrade(t *testing.T) {

	timePast := time.Now()
	timeFuture := time.Now().Add(time.Duration(123456789))

	tests := []struct {
		name                 string
		proposedByGovernment bool
		time                 *time.Time
		height               int64
		wantErr              bool
	}{
		{"government schedules with future block height",
			true,
			nil,
			100,
			false,
		},
		{"government schedules with future time",
			true,
			&timeFuture,
			-1,
			false,
		},
		{"government schedules with time that has already passed raises error",
			true,
			&timePast,
			-1,
			true,
		},
		{"non-government address trying to schedule raises error",
			false,
			nil,
			100,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput(true)

			plan := upgrade.Plan{
				Name: "name",
				Info: "info info info",
			}

			if tt.time != nil {
				plan.Time = *tt.time
			} else {
				plan.Height = tt.height
			}

			var proposer sdk.AccAddress
			if tt.proposedByGovernment {
				proposer = governmentTestAddress
			} else {
				proposer = notGovernmentAddress
			}

			err := k.ScheduleUpgrade(ctx, proposer, plan)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				upgradePlan, havePlan := k.GetUpgradePlan(ctx)
				require.True(t, havePlan)
				require.NotNil(t, upgradePlan)
			}
		})
	}
}

func TestKeeper_DeleteUpgrade(t *testing.T) {

	tests := []struct {
		name                 string
		proposedByGovernment bool
		planExists           bool
		wantErr              bool
	}{
		{"government deletes the scheduled plan",
			true,
			true,
			false,
		},
		{"government deletes nothing",
			true,
			false,
			false,
		},
		{"non-government address trying to delete raises error",
			false,
			true,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput(true)

			var plan upgrade.Plan
			if tt.planExists {
				plan = upgrade.Plan{
					Name:   "name",
					Info:   "info info info",
					Height: 100,
				}
				err := k.ScheduleUpgrade(ctx, governmentTestAddress, plan)
				require.NoError(t, err)
			}
			nilPlan := upgrade.Plan{Name: "", Time: time.Time{}, Height: 0, Info: ""}

			var proposer sdk.AccAddress
			if tt.proposedByGovernment {
				proposer = governmentTestAddress
			} else {
				proposer = notGovernmentAddress
			}

			err := k.DeleteUpgrade(ctx, proposer)

			if tt.wantErr {
				require.Error(t, err)

				upgradePlan, havePlan := k.GetUpgradePlan(ctx)
				require.True(t, havePlan)
				require.Equal(t, plan, upgradePlan)

			} else {
				require.NoError(t, err)

				upgradePlan, havePlan := k.GetUpgradePlan(ctx)
				require.False(t, havePlan)
				require.Equal(t, nilPlan, upgradePlan)
			}
		})
	}
}
