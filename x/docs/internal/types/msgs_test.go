package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// Test vars
var sender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var recipient, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var msgShareDocumentSchema = MsgShareDocument{
	Sender:     sender,
	Recipient:  recipient,
	Uuid:       "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	ContentUri: "http://www.contentUri.com",
	Metadata: DocumentMetadata{
		ContentUri: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			Uri:     "http://www.contentUri.com",
			Version: "test",
		},
		Proof: "proof",
	},
	Checksum: DocumentChecksum{
		Value:     "48656c6c6f20476f7068657221234567",
		Algorithm: "md5",
	},
}
var msgShareDocumentSchemaType = MsgShareDocument{
	Sender:     sender,
	Recipient:  recipient,
	Uuid:       "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	ContentUri: "http://www.contentUri.com",
	Metadata: DocumentMetadata{
		ContentUri: "http://www.contentUri.com",
		SchemaType: "uni-sincro",
		Proof:      "proof",
	},
	Checksum: DocumentChecksum{
		Value:     "48656c6c6f20476f7068657221234567",
		Algorithm: "md5",
	},
}

// ----------------------
// --- Msg methods
// ----------------------

func TestMsgShareDocument_Route(t *testing.T) {
	actual := msgShareDocumentSchema.Route()
	expected := ModuleName

	assert.Equal(t, expected, actual)
}

func TestMsgShareDocument_Type(t *testing.T) {
	actual := msgShareDocumentSchema.Type()
	expected := MsgTypeShareDocument

	assert.Equal(t, expected, actual)
}

func TestMsgShareDocument_ValidateBasic_Schema_valid(t *testing.T) {
	actual := msgShareDocumentSchema.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgShareDocument_ValidateBasic_SchemaType_valid(t *testing.T) {
	actual := msgShareDocumentSchemaType.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgShareDocument_ValidateBasic_invalid(t *testing.T) {
	invalidMsg := MsgShareDocument{
		Sender:     sender,
		Recipient:  recipient,
		Uuid:       "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
		ContentUri: "http://www.contentUri.com",
		Metadata: DocumentMetadata{
			ContentUri: "http://www.contentUri.com",
			Schema: &DocumentMetadataSchema{
				Uri:     "http://www.contentUri.com",
				Version: "test",
			},
			Proof: "proof",
		},
		Checksum: DocumentChecksum{
			Value:     "testValue",
			Algorithm: "sha-256",
		},
	}

	actual := invalidMsg.ValidateBasic()
	assert.NotNil(t, actual)
}

func TestMsgShareDocument_GetSignBytes(t *testing.T) {
	actual := msgShareDocumentSchema.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgShareDocumentSchema))
	assert.Equal(t, expected, actual)
}

func TestMsgShareDocument_GetSigners(t *testing.T) {
	actual := msgShareDocumentSchema.GetSigners()
	expected := msgShareDocumentSchema.Sender

	assert.Equal(t, expected, actual[0])
}

func TestMsgShareDocument_UnmarshalJson_Schema(t *testing.T) {
	json := `{"type":"commercio/MsgShareDocument","value":{"sender":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","recipient":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","uuid":"6a2f41a3-c54c-fce8-32d2-0324e1c32e22","content_uri":"http://www.contentUri.com","metadata":{"content_uri":"http://www.contentUri.com","schema":{"uri":"http://www.contentUri.com","version":"test"},"proof":"proof"},"checksum":{"value":"48656c6c6f20476f7068657221234567","algorithm":"md5"}}}`

	var msg MsgShareDocument
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	assert.Equal(t, "http://www.contentUri.com", msg.Metadata.Schema.Uri)
	assert.Equal(t, "test", msg.Metadata.Schema.Version)
}

func TestMsgShareDocument_UnmarshalJson_SchemaType(t *testing.T) {
	json := `{"type":"commercio/MsgShareDocument","value":{"sender":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","recipient":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","uuid":"6a2f41a3-c54c-fce8-32d2-0324e1c32e22","content_uri":"http://www.contentUri.com","metadata":{"content_uri":"http://www.contentUri.com","schema_type":"uni-sincro","proof":"proof"},"checksum":{"value":"48656c6c6f20476f7068657221234567","algorithm":"md5"}}}`

	var msg MsgShareDocument
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	assert.Equal(t, "uni-sincro", msg.Metadata.SchemaType)
}

// ----------------------
// --- UUID Validation
// ----------------------

func TestValidateUuid_valid(t *testing.T) {
	actual := validateUuid("6a2f41a3-c54c-fce8-32d2-0324e1c32e22")
	assert.True(t, actual)
}

func TestValidateUuid_empty(t *testing.T) {
	actual := validateUuid("")
	assert.False(t, actual)
}

func TestValidateUuid_invalid(t *testing.T) {
	actual := validateUuid("ebkfkd")
	assert.False(t, actual)
}

// -------------------------
// --- MetadataSchema validation
// -------------------------

func TestValidateDocMetadata_valid(t *testing.T) {
	validDocumentMetadata := DocumentMetadata{
		ContentUri: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			Uri:     "http://www.contentUri.com",
			Version: "test",
		},
		Proof: "proof",
	}

	actual := validateDocMetadata(validDocumentMetadata)
	assert.Nil(t, actual)
}

func TestValidateDocMetadata_emptyContentUri(t *testing.T) {
	invalidDocumentMetadata := DocumentMetadata{
		ContentUri: "",
		Schema: &DocumentMetadataSchema{
			Uri:     "http://www.contentUri.com",
			Version: "test",
		},
		Proof: "proof",
	}

	actual := validateDocMetadata(invalidDocumentMetadata)
	assert.NotNil(t, actual)
}

func TestValidateDocMetadata_emptySchemaUri(t *testing.T) {
	invalidDocumentMetadata := DocumentMetadata{
		ContentUri: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			Uri:     "",
			Version: "test",
		},
		Proof: "proof",
	}

	actual := validateDocMetadata(invalidDocumentMetadata)
	assert.NotNil(t, actual)
}

func TestValidateDocMetadata_emptySchemaVersion(t *testing.T) {
	invalidDocumentMetadata := DocumentMetadata{
		ContentUri: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			Uri:     "http://www.contentUri.com",
			Version: "",
		},
		Proof: "proof",
	}

	actual := validateDocMetadata(invalidDocumentMetadata)
	assert.NotNil(t, actual)
}

func TestValidateDocMetadata_emptyProof(t *testing.T) {
	invalidDocumentMetadata := DocumentMetadata{
		ContentUri: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			Uri:     "http://www.contentUri.com",
			Version: "test",
		},
		Proof: "",
	}

	actual := validateDocMetadata(invalidDocumentMetadata)
	assert.NotNil(t, actual)
}

// --------------------------
// --- Checksum validation
// --------------------------

func TestValidateChecksum_validChecksum(t *testing.T) {
	var checksumList = map[string]string{
		"md5":     "0cc175b9c0f1b6a831c399e269772661",
		"sha-1":   "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8",
		"sha-224": "abd37534c7d9a2efb9465de931cd7055ffdb8879563ae98078d6d6d5",
		"sha-256": "ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb",
		"sha-384": "54a59b9f22b0b80880d8427e548b7c23abd873486e1f035dce9cd697e85175033caa88e6d57bc35efae0b5afd3145f31",
		"sha-512": "1f40fc92da241694750979ee6cf582f2d5d7d28e18335de05abc54d0560e0f5302860c652bf08d560252aa5e74210546f369fbbbce8c12cfc7957b2652fe9a75",
	}

	for key, value := range checksumList {
		checksum := DocumentChecksum{Algorithm: key, Value: value}
		actual := validateChecksum(checksum)
		assert.Nil(t, actual)
	}
}

func TestValidateChecksum_emptyValue(t *testing.T) {
	invalidChecksum := DocumentChecksum{
		Value:     "",
		Algorithm: "md5",
	}

	actual := validateChecksum(invalidChecksum)
	assert.NotNil(t, actual)
}

func TestValidateChecksum_emptyAlgorithm(t *testing.T) {
	invalidChecksum := DocumentChecksum{
		Value:     "48656c6c6f20476f7068657221234567",
		Algorithm: "",
	}

	actual := validateChecksum(invalidChecksum)
	assert.NotNil(t, actual)
}

func TestValidateChecksum_invalidHexValue(t *testing.T) {
	invalidChecksum := DocumentChecksum{
		Value:     "qr54g7srg5674fsg4sfg",
		Algorithm: "md5",
	}

	actual := validateChecksum(invalidChecksum)
	assert.NotNil(t, actual)
}

func TestValidateChecksum_invalidAlgorithmType(t *testing.T) {
	invalidChecksum := DocumentChecksum{
		Value:     "48656c6c6f20476f7068657221234567",
		Algorithm: "md6",
	}

	actual := validateChecksum(invalidChecksum)
	assert.NotNil(t, actual)
}

func TestValidateChecksum_invalidChecksumLengths(t *testing.T) {
	var checksumList = map[string]string{
		"md5":     "0cc175bc0f1b6a831c399e269772661",
		"sha-1":   "86f7e437faa5a7fce15dddcb9eaeaea377667b8",
		"sha-224": "abd37534c7d9a2efb946de931cd7055ffdb8879563ae98078d6d6d5",
		"sha-256": "ca978112ca1bbdcafac21b39a23dc4da786eff8147c4e72b9807785afee48bb",
		"sha-384": "54a59b9f22b0b80880d427e548b7c23abd873486e1f035dce9cd697e85175033caa88e6d57bc35efae0b5afd3145f31",
		"sha-512": "1f40fc92da24169475099ee6cf582f2d5d7d28e18335de05abc54d0560e0f5302860c652bf08d560252aa5e74210546f369fbbbce8c12cfc7957b2652fe9a75",
	}

	for key, value := range checksumList {
		checksum := DocumentChecksum{Algorithm: key, Value: value}
		actual := validateChecksum(checksum)
		assert.NotNil(t, actual)
	}
}

// ----------------------------------
// --- DocumentReceipt tests
// ----------------------------------
var msgDocumentReceipt = MsgSendDocumentReceipt{
	Sender:       sender,
	Recipient:    recipient,
	TxHash:       "txHash",
	DocumentUuid: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	Proof:        "proof",
}

func TestMsgDocumentReceipt_Route(t *testing.T) {
	actual := msgDocumentReceipt.Route()
	expected := ModuleName

	assert.Equal(t, expected, actual)
}

func TestMsgDocumentReceipt_Type(t *testing.T) {
	actual := msgDocumentReceipt.Type()
	expected := MsgTypeSendDocumentReceipt

	assert.Equal(t, expected, actual)
}

func TestMsgDocumentReceipt_ValidateBasic_valid(t *testing.T) {
	actual := msgDocumentReceipt.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgDocumentReceipt_ValidateBasic_invalid(t *testing.T) {
	var msgDocReceipt = MsgSendDocumentReceipt{
		Sender:       sender,
		Recipient:    recipient,
		TxHash:       "txHash",
		DocumentUuid: "123456789",
		Proof:        "proof",
	}
	actual := msgDocReceipt.ValidateBasic()

	assert.NotNil(t, actual)
}

func TestMsgDocumentReceipt_GetSignBytes(t *testing.T) {
	actual := msgDocumentReceipt.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgDocumentReceipt))
	assert.Equal(t, expected, actual)
}

func TestMsgDocumentReceipt_GetSigners(t *testing.T) {
	actual := msgDocumentReceipt.GetSigners()
	expected := msgDocumentReceipt.Sender

	assert.Equal(t, expected, actual[0])
}
