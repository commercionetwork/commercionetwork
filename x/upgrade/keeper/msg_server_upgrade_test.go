package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeeper_handlerMsgScheduleUpgrade(t *testing.T) {
	plan := upgradeTypes.Plan{
		Name:   "name",
		Info:   "info info info",
		Height: 100,
	}

	tests := []struct {
		name    string
		msg     *types.MsgScheduleUpgrade
		wantErr bool
	}{
		{
			"MsgSetTumblerAddress by government",
			types.NewMsgScheduleUpgrade(governmentTestAddress.String(), plan),
			false,
		},
		{
			"MsgSetTumblerAddress by fake government yields error",
			types.NewMsgScheduleUpgrade(notGovernmentAddress.String(), plan),
			true,
		},
		/*{
			"MsgSetTumblerAddress by government",
			types.NewMsgDeleteUpgrade(governmentTestAddress.String()),
			false,
		},
		{
			"MsgSetTumblerAddress by fake government yields error",
			types.NewMsgDeleteUpgrade(notGovernmentAddress.String()),
			true,
		},*/
	}
	
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			msgServer := NewMsgServerImpl(*k)
			//msgServer.ScheduleUpgrade(sdk.WrapSDKContext(ctx), tt.msg)
			//handler := NewHandler(k)

			msgScheduleUpgradeResponse, err := msgServer.ScheduleUpgrade(sdk.WrapSDKContext(ctx), tt.msg) //handler(ctx, tt.msg)

			if tt.wantErr {
				require.Nil(t, msgScheduleUpgradeResponse)
				require.Error(t, err)
			} else {
				require.NotNil(t, msgScheduleUpgradeResponse)
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeper_handlerMsgDeleteUpgrade(t *testing.T) {
	tests := []struct {
		name    string
		msg     *types.MsgDeleteUpgrade
		wantErr bool
	}{
		{
			"MsgSetTumblerAddress by government",
			types.NewMsgDeleteUpgrade(governmentTestAddress.String()),
			false,
		},
		{
			"MsgSetTumblerAddress by fake government yields error",
			types.NewMsgDeleteUpgrade(notGovernmentAddress.String()),
			true,
		},
	}
	
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			msgServer := NewMsgServerImpl(*k)
			msgDeleteUpgradeResponse, err := msgServer.DeleteUpgrade(sdk.WrapSDKContext(ctx), tt.msg)

			if tt.wantErr {
				require.Nil(t, msgDeleteUpgradeResponse)
				require.Error(t, err)
			} else {
				require.NotNil(t, msgDeleteUpgradeResponse)
				require.NoError(t, err)
			}
		})
	}
}
