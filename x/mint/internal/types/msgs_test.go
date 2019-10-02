package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var testMsgOpenCDP = MsgOpenCDP{Request: TestCdpRequest}

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
	assert.Equal(t, sdk.ErrUnknownRequest("Cdp request's timestamp can't be empty"), err)
}
