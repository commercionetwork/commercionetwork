package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var TestOwnerAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var msgSetId = MsgSetIdentity{
	Owner: TestOwnerAddress,
	DidDocument: DidDocument{
		Uri:         "https://test.example.com/did-document#1",
		ContentHash: "ebd7cfd95e67f57b4f5bcb3cf4554741a9b5ea666052106625a93a6f8881dc43",
	},
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
		Owner:       sdk.AccAddress{},
		DidDocument: msgSetId.DidDocument,
	}

	actual := invalidMsg.ValidateBasic()
	assert.Error(t, actual)
}

func TestMsgSetIdentity_ValidateBasic_InvalidDidDocumentUri(t *testing.T) {
	invalidMsg := MsgSetIdentity{
		Owner: TestOwnerAddress,
		DidDocument: DidDocument{
			Uri:         "",
			ContentHash: "ebd7cfd95e67f57b4f5bcb3cf4554741a9b5ea666052106625a93a6f8881dc43",
		},
	}

	actual := invalidMsg.ValidateBasic()
	assert.Error(t, actual)
}

func TestMsgSetIdentity_ValidateBasic_InvalidDidDocumentContentHash(t *testing.T) {
	invalidMsg := MsgSetIdentity{
		Owner: TestOwnerAddress,
		DidDocument: DidDocument{
			Uri:         msgSetId.DidDocument.Uri,
			ContentHash: "ebd7cfd95e67f57b4f5bcb3cf4554741a9b5ea6662106625a93a6f8881dc43",
		},
	}

	actual := invalidMsg.ValidateBasic()
	assert.Error(t, actual)
}

func TestMsgSetIdentity_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgSetIdentity","value":{"did_document":{"content_hash":"ebd7cfd95e67f57b4f5bcb3cf4554741a9b5ea666052106625a93a6f8881dc43","uri":"https://test.example.com/did-document#1"},"owner":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"}}`

	actual := msgSetId.GetSignBytes()
	assert.Equal(t, expected, string(actual))
}

func TestNewMsgSetIdentity_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgSetId.Owner}
	actual := msgSetId.GetSigners()
	assert.Equal(t, expected, actual)
}
