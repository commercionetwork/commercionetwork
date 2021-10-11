package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/x/upgrade/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Test vars
var governmentTestAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var msgScheduleUpgrade = MsgScheduleUpgrade{
	Proposer: governmentTestAddress.String() ,
	Plan: &types.Plan{
		Name:   "name",
		Info:   "info info info",
		Height: 100,
	},
}
var msgDeleteUpgrade = MsgDeleteUpgrade{
	Proposer: governmentTestAddress.String(),
}
// ----------------------------------
// --- MsgScheduleUpgrade
// ----------------------------------

func TestMsgScheduleUpgrade_Route(t *testing.T) {
	actual := msgScheduleUpgrade.Route()
	require.Equal(t, QuerierRoute, actual)
}

func TestMsgScheduleUpgrade_Type(t *testing.T) {
	actual := msgScheduleUpgrade.Type()
	require.Equal(t, MsgTypeScheduleUpgrade, actual)
}

func TestMsgScheduleUpgrade_GetSignBytes(t *testing.T) {
	actual := msgScheduleUpgrade.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msgScheduleUpgrade))
	require.Equal(t, expected, actual)
}

func TestMsgScheduleUpgrade_GetSigners(t *testing.T) {
	actual := msgScheduleUpgrade.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgScheduleUpgrade.Proposer, actual[0].String())
}

func TestMsgScheduleUpgrade_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		sdr     MsgScheduleUpgrade
		haveErr bool
	}{
		{
			"valid MsgScheduleUpgrade",
			msgScheduleUpgrade,
			false,
		},
		{
			"MsgScheduleUpgrade with empty proposer",
			MsgScheduleUpgrade{Proposer: "",
				Plan: &types.Plan{
					Name:   "name",
					Info:   "info info info",
					Height: 100,
				}},
			true,
		},
		{
			"MsgScheduleUpgrade with invalid plan",
			MsgScheduleUpgrade{Proposer: governmentTestAddress.String(),
				Plan: &types.Plan{
					Name:   "",
					Info:   "info info info",
					Height: 100,
				}},
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sdr.ValidateBasic()
			if tt.haveErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// ----------------------------------
// --- MsgDeleteUpgrade
// ----------------------------------

func TestMsgDeleteUpgrade_Route(t *testing.T) {
	actual := msgDeleteUpgrade.Route()
	require.Equal(t, QuerierRoute, actual)
}

func TestMsgDeleteUpgrade_Type(t *testing.T) {
	actual := msgDeleteUpgrade.Type()
	require.Equal(t, MsgTypeDeleteUpgrade, actual)
}

func TestMsgDeleteUpgrade_GetSignBytes(t *testing.T) {
	actual := msgDeleteUpgrade.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msgDeleteUpgrade))
	require.Equal(t, expected, actual)
}

func TestMsgDeleteUpgrade_GetSigners(t *testing.T) {
	actual := msgDeleteUpgrade.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgDeleteUpgrade.Proposer, actual[0].String())
}

func TestMsgDeleteUpgrade_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		sdr     MsgDeleteUpgrade
		haveErr bool
	}{
		{
			"valid MsgDeleteUpgrade",
			msgDeleteUpgrade,
			false,
		},
		{
			"MsgDeleteUpgrade with empty proposer",
			MsgDeleteUpgrade{Proposer: ""},
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sdr.ValidateBasic()
			if tt.haveErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}