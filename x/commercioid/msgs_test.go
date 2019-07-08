package commercioid

import (
	"commercio-network/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

// ----------------------------------
// --- SetIdentity
// ----------------------------------

func TestMsgSetIdentity_Route(t *testing.T) {
	key := "commercioid"

	actual := msgSetId.Route()

	assert.Equal(t, key, actual)
}

func TestMsgSetIdentity_Type(t *testing.T) {
	ttype := "set_identity"

	actual := msgSetId.Type()

	assert.Equal(t, ttype, actual)
}

func TestMsgSetIdentity_ValidateBasic_AllFieldsCorrect(t *testing.T) {

	actual := msgSetId.ValidateBasic()

	assert.Nil(t, actual)
}

func TestMsgSetIdentity_ValidateBasic_InvalidAddress(t *testing.T) {
	invalidMsg := MsgSetIdentity{
		DID:          ownerIdentity,
		DDOReference: identityRef,
		Owner:        sdk.AccAddress{},
	}

	actual := invalidMsg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgSetIdentity_ValidateBasic_InvalidDID(t *testing.T) {
	invalidMsg := MsgSetIdentity{
		DID:          types.Did(""),
		DDOReference: identityRef,
		Owner:        sdk.AccAddress{},
	}

	actual := invalidMsg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgSetIdentity_GetSignBytes(t *testing.T) {
	expected := sdk.MustSortJSON(input.cdc.MustMarshalJSON(msgSetId))

	actual := msgSetId.GetSignBytes()

	assert.Equal(t, expected, actual)
}

func TestNewMsgSetIdentity_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgSetId.Owner}

	actual := msgSetId.GetSigners()

	assert.Equal(t, expected, actual)
}

// ----------------------------------
// --- CreateConnection
// ----------------------------------

func TestMsgCreateConnection_Route(t *testing.T) {
	key := "commercioid"

	actual := msgCreateConn.Route()

	assert.Equal(t, key, actual)
}

func TestMsgCreateConnection_Type(t *testing.T) {
	ttype := "create_connection"

	actual := msgCreateConn.Type()

	assert.Equal(t, ttype, actual)
}

func TestMsgCreateConnection_ValidateBasic_AllFieldsCorrect(t *testing.T) {

	actual := msgCreateConn.ValidateBasic()

	assert.Nil(t, actual)
}

func TestMsgCreateConnection_ValidateBasic_InvalidSignerAddress(t *testing.T) {
	invMsg := MsgCreateConnection{
		FirstUser:  ownerIdentity,
		SecondUser: recipient,
		Signer:     sdk.AccAddress{},
	}

	actual := invMsg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgCreateConnection_ValidateBasic_InvalidUser(t *testing.T) {
	invMsg := MsgCreateConnection{
		FirstUser:  types.Did(""),
		SecondUser: recipient,
		Signer:     owner,
	}

	actual := invMsg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgCreateConnection_GetSignBytes(t *testing.T) {
	expected := sdk.MustSortJSON(msgCdc.MustMarshalJSON(msgSetId))

	actual := msgSetId.GetSignBytes()

	assert.Equal(t, expected, actual)
}

func TestMsgCreateConnection_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgSetId.Owner}

	actual := msgSetId.GetSigners()

	assert.Equal(t, expected, actual)
}
