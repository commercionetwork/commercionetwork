package types

import (
	"github.com/commercionetwork/commercionetwork/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testAddress = "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"
var testOwner, _ = sdk.AccAddressFromBech32(testAddress)
var testOwnerIdentity = types.Did("newReader")
var testIdentityRef = "ddo-reference"
var testReference = "testReference"
var testMetadata = "testMetadata"
var testRecipient = types.Did("recipient")

var msgSetId = MsgSetIdentity{
	Did:          testOwnerIdentity,
	DDOReference: testIdentityRef,
	Owner:        testOwner,
}

var msgCreateConn = MsgCreateConnection{
	FirstUser:  testOwnerIdentity,
	SecondUser: testRecipient,
	Signer:     testOwner,
}

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
		Did:          testOwnerIdentity,
		DDOReference: testIdentityRef,
		Owner:        sdk.AccAddress{},
	}

	actual := invalidMsg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgSetIdentity_ValidateBasic_InvalidDID(t *testing.T) {
	invalidMsg := MsgSetIdentity{
		Did:          types.Did(""),
		DDOReference: testIdentityRef,
		Owner:        sdk.AccAddress{},
	}

	actual := invalidMsg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgSetIdentity_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercioid/SetIdentity","value":{"ddo_reference":"ddo-reference","did":"newReader","owner":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"}}`

	actual := msgSetId.GetSignBytes()

	assert.Equal(t, expected, string(actual))
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
		FirstUser:  testOwnerIdentity,
		SecondUser: testRecipient,
		Signer:     sdk.AccAddress{},
	}

	actual := invMsg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgCreateConnection_ValidateBasic_InvalidUser(t *testing.T) {
	invMsg := MsgCreateConnection{
		FirstUser:  types.Did(""),
		SecondUser: testRecipient,
		Signer:     testOwner,
	}

	actual := invMsg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgCreateConnection_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercioid/CreateConnection","value":{"first_user":"newReader","second_user":"recipient","signer":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"}}`

	actual := msgCreateConn.GetSignBytes()

	assert.Equal(t, expected, string(actual))
}

func TestMsgCreateConnection_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgSetId.Owner}

	actual := msgSetId.GetSigners()

	assert.Equal(t, expected, actual)
}
