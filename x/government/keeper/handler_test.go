package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/vbr"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/government/types"
)

func Test_handleMsgSetTumblerAddress(t *testing.T) {
	tests := []struct {
		name        string
		govAddr     sdk.AccAddress
		tumblerAddr sdk.AccAddress
		wantErr     bool
	}{
		{
			"fake government tries to set a tumbler address",
			notGovernmentAddress,
			tumblerTestAddress,
			true,
		},
		{
			"real government sets a tumbler address",
			governmentTestAddress,
			tumblerTestAddress,
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput(true)

			// set tumbler to a fake value, just to make the addr equality work
			err := k.SetTumblerAddress(ctx, notGovernmentAddress)
			require.NoError(t, err)

			msg := types.NewMsgSetTumblerAddress(tt.govAddr, tt.tumblerAddr)

			handler := NewHandler(k)

			res, err := handler(ctx, msg)

			if tt.wantErr {
				require.False(t, k.GetTumblerAddress(ctx).Equals(tumblerTestAddress))
				require.Nil(t, res)
				require.Error(t, err)
			} else {
				require.True(t, k.GetTumblerAddress(ctx).Equals(tumblerTestAddress))
				require.NotNil(t, res)
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeper_handlerFunc(t *testing.T) {
	tests := []struct {
		name    string
		msg     sdk.Msg
		wantErr bool
	}{
		{
			"a message which is not MsgSetTumblerAddress",
			vbr.MsgIncrementsBlockRewardsPool{},
			true,
		},
		{
			"MsgSetTumblerAddress",
			types.NewMsgSetTumblerAddress(governmentTestAddress, tumblerTestAddress),
			false,
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
