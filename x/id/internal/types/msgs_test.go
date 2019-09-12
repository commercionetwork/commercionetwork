package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var TestOwnerAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestDidDocumentUri = "https://test.example.com/did-document#1"
var TestConnectionAddress, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")

var msgSetId = MsgSetIdentity{
	Owner:          TestOwnerAddress,
	DidDocumentUri: TestDidDocumentUri,
}

// ----------------------------------
// --- SetIdentity
// ----------------------------------

func TestMsgSetIdentity_Route(t *testing.T) {
	actual := msgSetId.Route()
	assert.Equal(t, ModuleName, actual)
}

func TestMsgSetIdentity_Type(t *testing.T) {
	actual := msgSetId.Type()
	assert.Equal(t, MsgTypeSetIdentity, actual)
}

func TestMsgSetIdentity_ValidateBasic_AllFieldsCorrect(t *testing.T) {
	actual := msgSetId.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgSetIdentity_ValidateBasic_InvalidAddress(t *testing.T) {
	invalidMsg := MsgSetIdentity{
		DidDocumentUri: TestDidDocumentUri,
		Owner:          sdk.AccAddress{},
	}

	actual := invalidMsg.ValidateBasic()
	assert.Error(t, actual)
}

func TestMsgSetIdentity_ValidateBasic_InvalidDidDocumentUri(t *testing.T) {
	invalidMsg := MsgSetIdentity{
		DidDocumentUri: "",
		Owner:          TestConnectionAddress,
	}

	actual := invalidMsg.ValidateBasic()
	assert.Error(t, actual)
}

func TestMsgSetIdentity_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgSetIdentity","value":{"did_document_uri":"https://test.example.com/did-document#1","owner":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"}}`

	actual := msgSetId.GetSignBytes()
	assert.Equal(t, expected, string(actual))
}

func TestNewMsgSetIdentity_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgSetId.Owner}
	actual := msgSetId.GetSigners()
	assert.Equal(t, expected, actual)
}
