package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var testMsgOpenCDP = MsgOpenCDP{Request: TestCdpRequest}
var testMsgCloseCDP = MsgCloseCDP{
	Signer:    TestOwner,
	Timestamp: TestTimestamp,
}

func TestMsgOpenCDP_Route(t *testing.T) {
	actual := testMsgOpenCDP.Route()
	assert.Equal(t, RouterKey, actual)
}

func TestMsgOpenCDP_Type(t *testing.T) {
	actual := testMsgOpenCDP.Type()
	assert.Equal(t, MsgTypeOpenCDP, actual)
}

func TestMsgOpenCDP_ValidateBasic_Valid(t *testing.T) {
	actual := testMsgOpenCDP.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgOpenCDP_ValidateBasic_InvalidOwnerAddr(t *testing.T) {
	invalidMsg := MsgOpenCDP{Request: CDPRequest{
		Signer:          nil,
		DepositedAmount: nil,
		Timestamp:       "",
	}}
	err := invalidMsg.ValidateBasic()
	assert.Error(t, err)
	assert.Equal(t, sdk.ErrInvalidAddress(invalidMsg.Request.Signer.String()), err)
}

func TestMsgOpenCDP_ValidateBasic_InvalidDepositedAmount(t *testing.T) {
	invalidMsg := MsgOpenCDP{Request: CDPRequest{
		Signer:          TestOwner,
		DepositedAmount: nil,
		Timestamp:       "",
	}}
	err := invalidMsg.ValidateBasic()
	assert.Error(t, err)
	assert.Equal(t, sdk.ErrInvalidCoins(invalidMsg.Request.DepositedAmount.String()), err)
}

func TestMsgOpenCDP_ValidateBasic_InvalidTimestamp(t *testing.T) {
	invalidMsg := MsgOpenCDP{Request: CDPRequest{
		Signer:          TestOwner,
		DepositedAmount: TestDepositedAmount,
		Timestamp:       "  ",
	}}
	err := invalidMsg.ValidateBasic()
	assert.Error(t, err)
	assert.Equal(t, sdk.ErrUnknownRequest("cdp request's timestamp can't be empty"), err)
}

func TestMsgOpenCDP_GetSignBytes(t *testing.T) {
	actual := testMsgOpenCDP.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(testMsgOpenCDP))
	assert.Equal(t, expected, actual)
}

func TestMsgOpenCDP_GetSigners(t *testing.T) {
	actual := testMsgOpenCDP.GetSigners()
	expected := []sdk.AccAddress{testMsgOpenCDP.Request.Signer}
	assert.Equal(t, expected, actual)
}

///////////////////
///MsgCloseCDP////
/////////////////
func TestMsgCloseCDP_Route(t *testing.T) {
	actual := testMsgCloseCDP.Route()
	assert.Equal(t, RouterKey, actual)
}

func TestMsgCloseCDP_Type(t *testing.T) {
	actual := testMsgCloseCDP.Type()
	assert.Equal(t, MsgTypeCloseCDP, actual)
}

func TestMsgCloseCDP_ValidateBasic_Valid(t *testing.T) {
	actual := testMsgCloseCDP.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgCloseCDP_ValidateBasic_InvalidSigner(t *testing.T) {
	msg := MsgCloseCDP{
		Signer:    nil,
		Timestamp: "",
	}
	actual := msg.ValidateBasic()
	assert.Equal(t, sdk.ErrInvalidAddress(msg.Signer.String()), actual)
	assert.Error(t, actual)
}

func TestMsgCloseCDP_ValidateBasic_InvalidTimestamp(t *testing.T) {
	msg := MsgCloseCDP{
		Signer:    TestOwner,
		Timestamp: "    ",
	}
	actual := msg.ValidateBasic()
	assert.Equal(t, sdk.ErrUnknownRequest("cdp's timestamp can't be empty"), actual)
	assert.Error(t, actual)
}

func TestMsgCloseCDP_GetSignBytes(t *testing.T) {
	actual := testMsgCloseCDP.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(testMsgCloseCDP))
	assert.Equal(t, expected, actual)
}

func TestMsgCloseCDP_GetSigners(t *testing.T) {
	actual := testMsgCloseCDP.GetSigners()
	expected := []sdk.AccAddress{testMsgCloseCDP.Signer}
	assert.Equal(t, expected, actual)
}
