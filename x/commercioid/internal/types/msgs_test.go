package types

/*
import (
	"commercio-network/types"
	"commercio-network/x/commercioid"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

// ----------------------------------
// --- SetIdentity
// ----------------------------------

func TestMsgSetIdentity_Route(t *testing.T) {
	key := "commercioid"

	actual := commercioid.msgSetId.Route()

	assert.Equal(t, key, actual)
}

func TestMsgSetIdentity_Type(t *testing.T) {
	ttype := "set_identity"

	actual := commercioid.msgSetId.Type()

	assert.Equal(t, ttype, actual)
}

func TestMsgSetIdentity_ValidateBasic_AllFieldsCorrect(t *testing.T) {

	actual := commercioid.msgSetId.ValidateBasic()

	assert.Nil(t, actual)
}

func TestMsgSetIdentity_ValidateBasic_InvalidAddress(t *testing.T) {
	invalidMsg := types2.MsgSetIdentity{
		DID:          commercioid.ownerIdentity,
		DDOReference: commercioid.identityRef,
		Owner:        sdk.AccAddress{},
	}

	actual := invalidMsg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgSetIdentity_ValidateBasic_InvalidDID(t *testing.T) {
	invalidMsg := types2.MsgSetIdentity{
		DID:          types.Did(""),
		DDOReference: commercioid.identityRef,
		Owner:        sdk.AccAddress{},
	}

	actual := invalidMsg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgSetIdentity_GetSignBytes(t *testing.T) {
	expected := sdk.MustSortJSON(commercioid.input.cdc.MustMarshalJSON(commercioid.msgSetId))

	actual := commercioid.msgSetId.GetSignBytes()

	assert.Equal(t, expected, actual)
}

func TestNewMsgSetIdentity_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{commercioid.msgSetId.Owner}

	actual := commercioid.msgSetId.GetSigners()

	assert.Equal(t, expected, actual)
}

// ----------------------------------
// --- CreateConnection
// ----------------------------------

func TestMsgCreateConnection_Route(t *testing.T) {
	key := "commercioid"

	actual := commercioid.msgCreateConn.Route()

	assert.Equal(t, key, actual)
}

func TestMsgCreateConnection_Type(t *testing.T) {
	ttype := "create_connection"

	actual := commercioid.msgCreateConn.Type()

	assert.Equal(t, ttype, actual)
}

func TestMsgCreateConnection_ValidateBasic_AllFieldsCorrect(t *testing.T) {

	actual := commercioid.msgCreateConn.ValidateBasic()

	assert.Nil(t, actual)
}

func TestMsgCreateConnection_ValidateBasic_InvalidSignerAddress(t *testing.T) {
	invMsg := types2.MsgCreateConnection{
		FirstUser:  commercioid.ownerIdentity,
		SecondUser: commercioid.recipient,
		Signer:     sdk.AccAddress{},
	}

	actual := invMsg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgCreateConnection_ValidateBasic_InvalidUser(t *testing.T) {
	invMsg := types2.MsgCreateConnection{
		FirstUser:  types.Did(""),
		SecondUser: commercioid.recipient,
		Signer:     commercioid.owner,
	}

	actual := invMsg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgCreateConnection_GetSignBytes(t *testing.T) {
	expected := sdk.MustSortJSON(types2.msgCdc.MustMarshalJSON(commercioid.msgSetId))

	actual := commercioid.msgSetId.GetSignBytes()

	assert.Equal(t, expected, actual)
}

func TestMsgCreateConnection_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{commercioid.msgSetId.Owner}

	actual := commercioid.msgSetId.GetSigners()

	assert.Equal(t, expected, actual)
}

*/
