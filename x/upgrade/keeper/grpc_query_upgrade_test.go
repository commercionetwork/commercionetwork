package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/require"
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
			k, ctx := setupKeeper(t)

			var plan upgradeTypes.Plan
			if tt.planExists {
				plan = upgradeTypes.Plan{
					Name:   "name",
					Info:   "info info info",
					Height: 100,
				}
				err := k.UpgradeKeeper.ScheduleUpgrade(ctx, plan)
				require.NoError(t, err)
				
				queryCurrentUpgradeResponse, err := k.CurrentUpgrade(sdk.WrapSDKContext(ctx), &types.QueryCurrentUpgradeRequest{})

				require.NoError(t, err)

				if tt.planExists {
					var p upgradeTypes.Plan

					require.NotPanics(t, func() {
						types.ModuleCdc.UnmarshalBinaryBare(queryCurrentUpgradeResponse.CurrentUpgrade, &p)
					})
					require.Equal(t, plan, p)
				} else {
					require.Nil(t, queryCurrentUpgradeResponse)
				}

			}
		})
	}
}
