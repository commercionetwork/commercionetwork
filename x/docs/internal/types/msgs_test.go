package types

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// Test vars
var sender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var recipient, _ = sdk.AccAddressFromBech32("cosmos1v0yk4hs2nry020ufmu9yhpm39s4scdhhtecvtr")
var msgShareDocumentSchema = MsgShareDocument(Document{
	UUID:       "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	ContentURI: "http://www.contentUri.com",
	Metadata: DocumentMetadata{
		ContentURI: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			URI:     "http://www.contentUri.com",
			Version: "test",
		},
	},
	Checksum: &DocumentChecksum{
		Value:     "48656c6c6f20476f7068657221234567",
		Algorithm: "md5",
	},
	Sender:     sender,
	Recipients: types.Addresses{recipient},
})

var msgShareDocumentSchemaType = MsgShareDocument(Document{
	UUID:       "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	ContentURI: "http://www.contentUri.com",
	Metadata: DocumentMetadata{
		ContentURI: "http://www.contentUri.com",
		SchemaType: "uni-sincro",
	},
	Checksum: &DocumentChecksum{
		Value:     "48656c6c6f20476f7068657221234567",
		Algorithm: "md5",
	},
	Sender:     sender,
	Recipients: types.Addresses{recipient},
})

// ----------------------
// --- MsgShareDocument
// ----------------------

func TestMsgShareDocument_Route(t *testing.T) {
	actual := msgShareDocumentSchema.Route()
	assert.Equal(t, QuerierRoute, actual)
}

func TestMsgShareDocument_Type(t *testing.T) {
	actual := msgShareDocumentSchema.Type()
	assert.Equal(t, MsgTypeShareDocument, actual)
}

func TestMsgShareDocument_ValidateBasic_Schema_valid(t *testing.T) {
	err := msgShareDocumentSchema.ValidateBasic()
	assert.Nil(t, err)
}

func TestMsgShareDocument_ValidateBasic_SchemaType_valid(t *testing.T) {
	err := msgShareDocumentSchemaType.ValidateBasic()
	assert.Nil(t, err)
}

func TestMsgShareDocument_GetSignBytes(t *testing.T) {
	actual := msgShareDocumentSchema.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgShareDocumentSchema))
	assert.Equal(t, expected, actual)
}

func TestMsgShareDocument_GetSigners(t *testing.T) {
	actual := msgShareDocumentSchema.GetSigners()
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, msgShareDocumentSchema.Sender, actual[0])
}

func TestMsgShareDocument_UnmarshalJson_Schema(t *testing.T) {
	json := `{"type":"commercio/MsgShareDocument","value":{"sender":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","recipients":["cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"], "uuid":"6a2f41a3-c54c-fce8-32d2-0324e1c32e22","content_uri":"http://www.contentUri.com","metadata":{"content_uri":"http://www.contentUri.com","schema":{"uri":"http://www.contentUri.com","version":"test"},"proof":"proof"},"checksum":{"value":"48656c6c6f20476f7068657221234567","algorithm":"md5"}}}`

	var msg MsgShareDocument
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	assert.Equal(t, "http://www.contentUri.com", msg.Metadata.Schema.URI)
	assert.Equal(t, "test", msg.Metadata.Schema.Version)
}

func TestMsgShareDocument_UnmarshalJson_SchemaType(t *testing.T) {
	json := `{"type":"commercio/MsgShareDocument","value":{"sender":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","recipients":["cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"],"uuid":"6a2f41a3-c54c-fce8-32d2-0324e1c32e22","content_uri":"http://www.contentUri.com","metadata":{"content_uri":"http://www.contentUri.com","schema_type":"uni-sincro","proof":"proof"},"checksum":{"value":"48656c6c6f20476f7068657221234567","algorithm":"md5"}}}`

	var msg MsgShareDocument
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	assert.Equal(t, "uni-sincro", msg.Metadata.SchemaType)
}

// ----------------------
// --- UUID Validation
// ----------------------

func TestValidateUuid_valid(t *testing.T) {
	actual := validateUUID("6a2f41a3-c54c-fce8-32d2-0324e1c32e22")
	assert.True(t, actual)
}

func TestValidateUuid_empty(t *testing.T) {
	actual := validateUUID("")
	assert.False(t, actual)
}

func TestValidateUuid_invalid(t *testing.T) {
	actual := validateUUID("ebkfkd")
	assert.False(t, actual)
}

// -----------------------------
// --- MsgSendDocumentReceipt
// -----------------------------

var msgDocumentReceipt = MsgSendDocumentReceipt{
	UUID:         "cfbb5b51-6ac0-43b0-8e09-022236285e31",
	Sender:       sender,
	Recipient:    recipient,
	TxHash:       "txHash",
	DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	Proof:        "proof",
}

func TestMsgDocumentReceipt_Route(t *testing.T) {
	actual := msgDocumentReceipt.Route()
	assert.Equal(t, QuerierRoute, actual)
}

func TestMsgDocumentReceipt_Type(t *testing.T) {
	actual := msgDocumentReceipt.Type()
	assert.Equal(t, MsgTypeSendDocumentReceipt, actual)
}

func TestMsgDocumentReceipt_ValidateBasic_valid(t *testing.T) {
	err := msgDocumentReceipt.ValidateBasic()
	assert.Nil(t, err)
}

func TestMsgDocumentReceipt_ValidateBasic_invalid(t *testing.T) {
	var msgDocReceipt = MsgSendDocumentReceipt{
		Sender:       sender,
		Recipient:    recipient,
		TxHash:       "txHash",
		DocumentUUID: "123456789",
		Proof:        "proof",
	}
	err := msgDocReceipt.ValidateBasic()
	assert.NotNil(t, err)
}

func TestMsgDocumentReceipt_GetSignBytes(t *testing.T) {
	actual := msgDocumentReceipt.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgDocumentReceipt))
	assert.Equal(t, expected, actual)
}

func TestMsgDocumentReceipt_GetSigners(t *testing.T) {
	actual := msgDocumentReceipt.GetSigners()
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, msgDocumentReceipt.Sender, actual[0])
}

// ------------------------------------
// --- MsgAddSupportedMetadataSchema
// ------------------------------------

var msgAddSupportedMetadataSchema = MsgAddSupportedMetadataSchema{
	Signer: sender,
	Schema: MetadataSchema{
		Type:      "schema",
		SchemaURI: "https://example.com/schema",
		Version:   "1.0.0",
	},
}

func Test_MsgAddSupportedMetadataSchema_Route(t *testing.T) {
	actual := msgAddSupportedMetadataSchema.Route()
	assert.Equal(t, QuerierRoute, actual)
}

func Test_MsgAddSupportedMetadataSchema_Type(t *testing.T) {
	actual := msgAddSupportedMetadataSchema.Type()
	assert.Equal(t, MsgTypeAddSupportedMetadataSchema, actual)
}

func Test_MsgAddSupportedMetadataSchema_ValidateBasic_valid(t *testing.T) {
	err := msgAddSupportedMetadataSchema.ValidateBasic()
	assert.Nil(t, err)
}

func Test_MsgAddSupportedMetadataSchema_ValidateBasic_invalid(t *testing.T) {
	var msgDocReceipt = MsgAddSupportedMetadataSchema{
		Signer: recipient,
		Schema: MetadataSchema{
			Type:      "schema-2",
			SchemaURI: "",
			Version:   "",
		},
	}
	err := msgDocReceipt.ValidateBasic()
	assert.NotNil(t, err)
}

func Test_MsgAddSupportedMetadataSchema_GetSignBytes(t *testing.T) {
	actual := msgAddSupportedMetadataSchema.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgAddSupportedMetadataSchema))
	assert.Equal(t, expected, actual)
}

func Test_MsgAddSupportedMetadataSchema_GetSigners(t *testing.T) {
	actual := msgAddSupportedMetadataSchema.GetSigners()
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, msgAddSupportedMetadataSchema.Signer, actual[0])
}

// -----------------------------------------
// --- MsgAddTrustedMetadataSchemaProposer
// -----------------------------------------

var msgAddTrustedMetadataSchemaProposer = MsgAddTrustedMetadataSchemaProposer{
	Proposer: sender,
	Signer:   recipient,
}

func Test_MsgAddTrustedMetadataSchemaProposer_Route(t *testing.T) {
	actual := msgAddTrustedMetadataSchemaProposer.Route()
	assert.Equal(t, QuerierRoute, actual)
}

func Test_MsgAddTrustedMetadataSchemaProposer_Type(t *testing.T) {
	actual := msgAddTrustedMetadataSchemaProposer.Type()
	assert.Equal(t, MsgTypeAddTrustedMetadataSchemaProposer, actual)
}

func Test_MsgAddTrustedMetadataSchemaProposer_ValidateBasic_valid(t *testing.T) {
	err := msgAddTrustedMetadataSchemaProposer.ValidateBasic()
	assert.Nil(t, err)
}

func Test_MsgAddTrustedMetadataSchemaProposer_ValidateBasic_invalid(t *testing.T) {
	var msgDocReceipt = MsgAddTrustedMetadataSchemaProposer{
		Proposer: nil,
		Signer:   recipient,
	}
	err := msgDocReceipt.ValidateBasic()
	assert.NotNil(t, err)
}

func Test_MsgAddTrustedMetadataSchemaProposer_GetSignBytes(t *testing.T) {
	actual := msgAddTrustedMetadataSchemaProposer.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgAddTrustedMetadataSchemaProposer))
	assert.Equal(t, expected, actual)
}

func Test_MsgAddTrustedMetadataSchemaProposer_GetSigners(t *testing.T) {
	actual := msgAddTrustedMetadataSchemaProposer.GetSigners()
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, msgAddTrustedMetadataSchemaProposer.Signer, actual[0])
}
