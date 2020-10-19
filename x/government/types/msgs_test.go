package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

// Test vars
var government, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var newTumbler, _ = sdk.AccAddressFromBech32("cosmos1h7tw92a66gr58pxgmf6cc336lgxadpjz5d5psf")

// ----------------------
// --- MsgSetTumblerAddress
// ----------------------

var msgSetTumblerAddress = NewMsgSetTumblerAddress(government, newTumbler)

func TestMsgSetTumblerAddress_Route(t *testing.T) {
	require.Equal(t, RouterKey, msgSetTumblerAddress.Route())
}

func TestMsgSetTumblerAddress_Type(t *testing.T) {
	require.Equal(t, MsgTypeSetTumblerAddress, msgSetTumblerAddress.Type())
}

func TestMsgSetTumblerAddress_ValidateBasic(t *testing.T) {
	tests := []struct {
		name  string
		msg   MsgSetTumblerAddress
		error error
	}{
		{
			name:  "Valid message returns no error",
			msg:   msgSetTumblerAddress,
			error: nil,
		},
		{
			name:  "Missing government returns error",
			msg:   MsgSetTumblerAddress{Government: nil, NewTumbler: newTumbler},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "invalid government address: "),
		},
		{
			name:  "Missing new tumbler returns error",
			msg:   MsgSetTumblerAddress{Government: government, NewTumbler: nil},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "invalid new tumbler address: "),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.error != nil {
				require.Equal(t, test.error.Error(), test.msg.ValidateBasic().Error())
			} else {
				require.NoError(t, test.msg.ValidateBasic())
			}
		})
	}
}

func TestMsgSetTumblerAddress_GetSignBytes(t *testing.T) {
	actual := msgSetTumblerAddress.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgSetTumblerAddress))
	require.Equal(t, expected, actual)
}

func TestMsgSetTumblerAddress_GetSigners(t *testing.T) {
	actual := msgSetTumblerAddress.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgSetTumblerAddress.Government, actual[0])
}

func TestMsgSetTumblerAddress_UnmarshalJson(t *testing.T) {
	json := `{"type":"commercio/MsgSetTumblerAddress","value":{"government":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","new_tumbler":"cosmos1h7tw92a66gr58pxgmf6cc336lgxadpjz5d5psf"}}`

	var msg MsgSetTumblerAddress
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	require.Equal(t, government, msg.Government)
	require.Equal(t, newTumbler, msg.NewTumbler)
}
