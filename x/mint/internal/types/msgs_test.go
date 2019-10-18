package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var testMsgOpenCdp = MsgOpenCdp(TestCdpRequest)
var testMsgCloseCdp = MsgCloseCdp{
	Signer:    TestOwner,
	Timestamp: TestTimestamp,
}

func TestMsgOpenCdp_Route(t *testing.T) {
	actual := testMsgOpenCdp.Route()
	assert.Equal(t, RouterKey, actual)
}

func TestMsgOpenCdp_Type(t *testing.T) {
	actual := testMsgOpenCdp.Type()
	assert.Equal(t, MsgTypeOpenCdp, actual)
}

func TestMsgOpenCdp_ValidateBasic_Valid(t *testing.T) {
	actual := testMsgOpenCdp.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgOpenCdp_ValidateBasic_InvalidOwnerAddr(t *testing.T) {
	invalidMsg := MsgOpenCdp(CdpRequest{
		Signer:          nil,
		DepositedAmount: nil,
		Timestamp:       "",
	})
	err := invalidMsg.ValidateBasic()
	assert.Error(t, err)
	assert.Equal(t, sdk.ErrInvalidAddress(invalidMsg.Signer.String()), err)
}

func TestMsgOpenCdp_ValidateBasic_InvalidDepositedAmount(t *testing.T) {
	invalidMsg := MsgOpenCdp(CdpRequest{
		Signer:          TestOwner,
		DepositedAmount: nil,
		Timestamp:       "",
	})
	err := invalidMsg.ValidateBasic()
	assert.Error(t, err)
	assert.Equal(t, sdk.ErrInvalidCoins(invalidMsg.DepositedAmount.String()), err)
}

func TestMsgOpenCdp_ValidateBasic_InvalidTimestamp(t *testing.T) {
	invalidMsg := MsgOpenCdp(CdpRequest{
		Signer:          TestOwner,
		DepositedAmount: TestDepositedAmount,
		Timestamp:       "  ",
	})
	err := invalidMsg.ValidateBasic()
	assert.Error(t, err)
	assert.Equal(t, sdk.ErrUnknownRequest("cdp request's timestamp can't be empty"), err)
}

func TestMsgOpenCdp_GetSignBytes(t *testing.T) {
	actual := testMsgOpenCdp.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(testMsgOpenCdp))
	assert.Equal(t, expected, actual)
}

func TestMsgOpenCdp_GetSigners(t *testing.T) {
	actual := testMsgOpenCdp.GetSigners()
	expected := []sdk.AccAddress{testMsgOpenCdp.Signer}
	assert.Equal(t, expected, actual)
}

///////////////////
///MsgCloseCdp////
/////////////////
func TestMsgCloseCdp_Route(t *testing.T) {
	actual := testMsgCloseCdp.Route()
	assert.Equal(t, RouterKey, actual)
}

func TestMsgCloseCdp_Type(t *testing.T) {
	actual := testMsgCloseCdp.Type()
	assert.Equal(t, MsgTypeCloseCdp, actual)
}

func TestMsgCloseCdp_ValidateBasic_Valid(t *testing.T) {
	actual := testMsgCloseCdp.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgCloseCdp_ValidateBasic_InvalidSigner(t *testing.T) {
	msg := MsgCloseCdp{
		Signer:    nil,
		Timestamp: "",
	}
	actual := msg.ValidateBasic()
	assert.Equal(t, sdk.ErrInvalidAddress(msg.Signer.String()), actual)
	assert.Error(t, actual)
}

func TestMsgCloseCdp_ValidateBasic_InvalidTimestamp(t *testing.T) {
	msg := MsgCloseCdp{
		Signer:    TestOwner,
		Timestamp: "    ",
	}
	actual := msg.ValidateBasic()
	assert.Equal(t, sdk.ErrUnknownRequest("cdp's timestamp can't be empty"), actual)
	assert.Error(t, actual)
}

func TestMsgCloseCdp_GetSignBytes(t *testing.T) {
	actual := testMsgCloseCdp.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(testMsgCloseCdp))
	assert.Equal(t, expected, actual)
}

func TestMsgCloseCdp_GetSigners(t *testing.T) {
	actual := testMsgCloseCdp.GetSigners()
	expected := []sdk.AccAddress{testMsgCloseCdp.Signer}
	assert.Equal(t, expected, actual)
}
