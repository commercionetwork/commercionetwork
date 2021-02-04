package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
	vbrTypes "github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestKeeper_handlerFunc(t *testing.T) {
	plan := upgrade.Plan{
		Name:   "name",
		Info:   "info info info",
		Height: 100,
	}

	tests := []struct {
		name    string
		msg     sdk.Msg
		wantErr bool
	}{
		{
			"a message which is not compatible with upgrade yields error",
			vbrTypes.MsgIncrementBlockRewardsPool{},
			true,
		},
		{
			"MsgSetTumblerAddress by government",
			types.NewMsgScheduleUpgrade(governmentTestAddress, plan),
			false,
		},
		{
			"MsgSetTumblerAddress by fake government yields error",
			types.NewMsgScheduleUpgrade(notGovernmentAddress, plan),
			true,
		},
		{
			"MsgSetTumblerAddress by government",
			types.NewMsgDeleteUpgrade(governmentTestAddress),
			false,
		},
		{
			"MsgSetTumblerAddress by fake government yields error",
			types.NewMsgDeleteUpgrade(notGovernmentAddress),
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput(true)

			handler := NewHandler(k)

			res, err := handler(ctx, tt.msg)

			if tt.wantErr {
				require.Nil(t, res)
				require.Error(t, err)
			} else {
				require.NotNil(t, res)
				require.NoError(t, err)
			}
		})
	}
}
