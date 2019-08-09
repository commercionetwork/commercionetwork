package types

import (
	"github.com/commercionetwork/commercionetwork/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

//TEST VARS
var addr = "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"
var sender, _ = sdk.AccAddressFromBech32(addr)
var recipient, _ = sdk.AccAddressFromBech32(addr)
var validUuid = "6a2f41a3-c54c-fce8-32d2-0324e1c32e22"
var validChecksum = types.DocumentChecksum{
	Value:     "48656c6c6f20476f7068657221234567",
	Algorithm: MD5,
}
var validMetadataSchema = types.DocumentMetadataSchema{
	Uri:     "http://www.contentUri.com",
	Version: "test",
}
var validDocumentMetadata = types.DocumentMetadata{
	ContentUri: "http://www.contentUri.com",
	Schema:     validMetadataSchema,
	Proof:      "proof",
}

var invalidChecksum = types.DocumentChecksum{
	Value:     "",
	Algorithm: "",
}
var invalidMetadataSchema = types.DocumentMetadataSchema{
	Uri:     "",
	Version: "",
}
var invalidDocumentMetadata = types.DocumentMetadata{
	ContentUri: "",
	Schema:     invalidMetadataSchema,
	Proof:      "",
}

var validMsg = MsgShareDocument{
	Sender:     sender,
	Recipient:  recipient,
	Uuid:       "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	ContentUri: "http://www.contentUri.com",
	Metadata:   validDocumentMetadata,
	Checksum:   validChecksum,
}

var invalidMsg = MsgShareDocument{
	Sender:     sender,
	Recipient:  recipient,
	Uuid:       validUuid,
	ContentUri: "http://www.contentUri.com",
	Metadata: types.DocumentMetadata{
		ContentUri: "http://www.contentUri.com",
		Schema: types.DocumentMetadataSchema{
			Uri:     "http://www.contentUri.com",
			Version: "test",
		},
		Proof: "proof",
	},
	Checksum: types.DocumentChecksum{
		Value:     "testValue",
		Algorithm: SHA256,
	},
}

func TestMsgShareDocument_Route(t *testing.T) {
	actual := validMsg.Route()
	expected := ModuleName

	assert.Equal(t, expected, actual)
}

func TestMsgShareDocument_Type(t *testing.T) {
	actual := validMsg.Type()
	expected := MsgType

	assert.Equal(t, expected, actual)
}

func TestValidateChecksum_valid(t *testing.T) {
	actual := validateChecksum(validChecksum)
	assert.Nil(t, actual)
}

func TestValidateChecksum_invalid(t *testing.T) {
	actual := validateChecksum(invalidChecksum)
	assert.NotNil(t, actual)
}

func TestValidateDocMetadata_valid(t *testing.T) {
	actual := validateDocMetadata(validDocumentMetadata)
	assert.Nil(t, actual)
}

func TestValidateDocMetadata_invalid(t *testing.T) {
	actual := validateDocMetadata(invalidDocumentMetadata)
	assert.NotNil(t, actual)
}

func TestValidateUuid_valid(t *testing.T) {
	actual := validateUuid(validUuid)
	assert.True(t, actual)
}

func TestValidateUuid_invalid(t *testing.T) {
	actual := validateUuid("ebkfkd")
	assert.False(t, actual)
}

func TestMsgShareDocument_ValidateBasic_valid(t *testing.T) {
	actual := validMsg.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgShareDocument_ValidateBasic_invalid(t *testing.T) {
	actual := invalidMsg.ValidateBasic()
	assert.NotNil(t, actual)
}

func TestMsgShareDocument_GetSignBytes(t *testing.T) {
	actual := validMsg.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(validMsg))
	assert.Equal(t, expected, actual)
}

func TestMsgShareDocument_GetSigners(t *testing.T) {
	actual := validMsg.GetSigners()
	expected := validMsg.Sender

	assert.Equal(t, expected, actual[0])
}
