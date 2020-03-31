package types

import (
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
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
	require.Equal(t, QuerierRoute, actual)
}

func TestMsgShareDocument_Type(t *testing.T) {
	actual := msgShareDocumentSchema.Type()
	require.Equal(t, MsgTypeShareDocument, actual)
}

func TestMsgShareDocument_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		sdr     MsgShareDocument
		haveErr error
	}{
		{
			"MsgShareDocument with valid schema",
			msgShareDocumentSchema,
			nil,
		},
		{
			"MsgShareDocument with no schema",
			MsgShareDocument(Document{
				UUID:       "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				ContentURI: "http://www.contentUri.com",
				Metadata: DocumentMetadata{
					ContentURI: "http://www.contentUri.com",
				},
				Checksum: &DocumentChecksum{
					Value:     "48656c6c6f20476f7068657221234567",
					Algorithm: "md5",
				},
				Sender:     sender,
				Recipients: types.Addresses{recipient},
			}),
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "either metadata.schema or metadata.schema_type must be defined"),
		},
		{
			"MsgShareDocument with valid schema type",
			msgShareDocumentSchemaType,
			nil,
		},
		{
			"MsgShareDocument with no schema type",
			MsgShareDocument(Document{
				UUID:       "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				ContentURI: "http://www.contentUri.com",
				Metadata: DocumentMetadata{
					ContentURI: "http://www.contentUri.com",
				},
				Checksum: &DocumentChecksum{
					Value:     "48656c6c6f20476f7068657221234567",
					Algorithm: "md5",
				},
				Sender:     sender,
				Recipients: types.Addresses{recipient},
			}),
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "either metadata.schema or metadata.schema_type must be defined"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sdr.ValidateBasic()
			if tt.haveErr != nil {
				require.EqualError(t, err, tt.haveErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgShareDocument_GetSignBytes(t *testing.T) {
	actual := msgShareDocumentSchema.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgShareDocumentSchema))
	require.Equal(t, expected, actual)
}

func TestMsgShareDocument_GetSigners(t *testing.T) {
	actual := msgShareDocumentSchema.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgShareDocumentSchema.Sender, actual[0])
}

func TestMsgShareDocument_UnmarshalJson_Schema(t *testing.T) {
	json := `{"type":"commercio/MsgShareDocument","value":{"sender":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","recipients":["cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"], "uuid":"6a2f41a3-c54c-fce8-32d2-0324e1c32e22","content_uri":"http://www.contentUri.com","metadata":{"content_uri":"http://www.contentUri.com","schema":{"uri":"http://www.contentUri.com","version":"test"},"proof":"proof"},"checksum":{"value":"48656c6c6f20476f7068657221234567","algorithm":"md5"}}}`

	var msg MsgShareDocument
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	require.Equal(t, "http://www.contentUri.com", msg.Metadata.Schema.URI)
	require.Equal(t, "test", msg.Metadata.Schema.Version)
}

func TestMsgShareDocument_UnmarshalJson_SchemaType(t *testing.T) {
	json := `{"type":"commercio/MsgShareDocument","value":{"sender":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","recipients":["cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"],"uuid":"6a2f41a3-c54c-fce8-32d2-0324e1c32e22","content_uri":"http://www.contentUri.com","metadata":{"content_uri":"http://www.contentUri.com","schema_type":"uni-sincro","proof":"proof"},"checksum":{"value":"48656c6c6f20476f7068657221234567","algorithm":"md5"}}}`

	var msg MsgShareDocument
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	require.Equal(t, "uni-sincro", msg.Metadata.SchemaType)
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
	require.Equal(t, QuerierRoute, actual)
}

func TestMsgDocumentReceipt_Type(t *testing.T) {
	actual := msgDocumentReceipt.Type()
	require.Equal(t, MsgTypeSendDocumentReceipt, actual)
}

func TestMsgSendDocumentReceipt_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		sdr     MsgSendDocumentReceipt
		haveErr error
	}{
		{
			"valid SendDocumentReceipt",
			msgDocumentReceipt,
			nil,
		},
		{
			"invalid UUID",
			MsgSendDocumentReceipt{
				Sender:       sender,
				Recipient:    recipient,
				TxHash:       "txHash",
				DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				Proof:        "proof",
			},
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid uuid: "),
		},
		{
			"empty sender",
			MsgSendDocumentReceipt{
				UUID:         "cfbb5b51-6ac0-43b0-8e09-022236285e31",
				Recipient:    recipient,
				TxHash:       "txHash",
				DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				Proof:        "proof",
			},
			sdkErr.Wrap(sdkErr.ErrInvalidAddress, ""),
		},
		{
			"empty recipient",
			MsgSendDocumentReceipt{
				UUID:         "cfbb5b51-6ac0-43b0-8e09-022236285e31",
				Sender:       sender,
				TxHash:       "txHash",
				DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				Proof:        "proof",
			},
			sdkErr.Wrap(sdkErr.ErrInvalidAddress, ""),
		},
		{
			"empty TxHash",
			MsgSendDocumentReceipt{
				UUID:         "cfbb5b51-6ac0-43b0-8e09-022236285e31",
				Sender:       sender,
				Recipient:    recipient,
				DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				Proof:        "proof",
			},
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Send Document's Transaction Hash can't be empty"),
		},
		{
			"invalid document UUID",
			MsgSendDocumentReceipt{
				UUID:      "cfbb5b51-6ac0-43b0-8e09-022236285e31",
				Sender:    sender,
				Recipient: recipient,
				TxHash:    "txHash",
				Proof:     "proof",
			},
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid document UUID"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sdr.ValidateBasic()
			if tt.haveErr != nil {
				require.EqualError(t, err, tt.haveErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgDocumentReceipt_GetSignBytes(t *testing.T) {
	actual := msgDocumentReceipt.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgDocumentReceipt))
	require.Equal(t, expected, actual)
}

func TestMsgDocumentReceipt_GetSigners(t *testing.T) {
	actual := msgDocumentReceipt.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgDocumentReceipt.Sender, actual[0])
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
	require.Equal(t, QuerierRoute, actual)
}

func Test_MsgAddSupportedMetadataSchema_Type(t *testing.T) {
	actual := msgAddSupportedMetadataSchema.Type()
	require.Equal(t, MsgTypeAddSupportedMetadataSchema, actual)
}

func Test_MsgAddSupportedMetadataSchema_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		sdr     MsgAddSupportedMetadataSchema
		haveErr error
	}{
		{
			"a valid AddSuppoertedMetadataSchema message",
			msgAddSupportedMetadataSchema,
			nil,
		},
		{
			"missing signer",
			MsgAddSupportedMetadataSchema{
				Schema: MetadataSchema{
					Type:      "schema",
					SchemaURI: "https://example.com/schema",
					Version:   "1.0.0",
				},
			},
			sdkErr.Wrap(sdkErr.ErrInvalidAddress, ""),
		},
		{
			"invalid schema",
			MsgAddSupportedMetadataSchema{
				Signer: recipient,
				Schema: MetadataSchema{
					Type:      "schema-2",
					SchemaURI: "",
					Version:   "",
				},
			},
			sdkErr.Wrap(sdkErr.ErrUnknownRequest, "uri cannot be empty"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sdr.ValidateBasic()
			if tt.haveErr != nil {
				require.EqualError(t, err, tt.haveErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_MsgAddSupportedMetadataSchema_GetSignBytes(t *testing.T) {
	actual := msgAddSupportedMetadataSchema.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgAddSupportedMetadataSchema))
	require.Equal(t, expected, actual)
}

func Test_MsgAddSupportedMetadataSchema_GetSigners(t *testing.T) {
	actual := msgAddSupportedMetadataSchema.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgAddSupportedMetadataSchema.Signer, actual[0])
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
	require.Equal(t, QuerierRoute, actual)
}

func Test_MsgAddTrustedMetadataSchemaProposer_Type(t *testing.T) {
	actual := msgAddTrustedMetadataSchemaProposer.Type()
	require.Equal(t, MsgTypeAddTrustedMetadataSchemaProposer, actual)
}

func Test_MsgAddTrustedMetadataSchemaProposer_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		sdr     MsgAddTrustedMetadataSchemaProposer
		haveErr error
	}{
		{
			"a valid AddSuppoertedMetadataSchema message",
			msgAddTrustedMetadataSchemaProposer,
			nil,
		},
		{
			"missing proposer",
			MsgAddTrustedMetadataSchemaProposer{
				Signer: recipient,
			},
			sdkErr.Wrap(sdkErr.ErrInvalidAddress, ""),
		},
		{
			"missing signer",
			MsgAddTrustedMetadataSchemaProposer{
				Proposer: sender,
			},
			sdkErr.Wrap(sdkErr.ErrInvalidAddress, ""),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sdr.ValidateBasic()
			if tt.haveErr != nil {
				require.EqualError(t, err, tt.haveErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_MsgAddTrustedMetadataSchemaProposer_ValidateBasic_valid(t *testing.T) {
	err := msgAddTrustedMetadataSchemaProposer.ValidateBasic()
	require.Nil(t, err)
}

func Test_MsgAddTrustedMetadataSchemaProposer_ValidateBasic_invalid(t *testing.T) {
	var msgDocReceipt = MsgAddTrustedMetadataSchemaProposer{
		Proposer: nil,
		Signer:   recipient,
	}
	err := msgDocReceipt.ValidateBasic()
	require.NotNil(t, err)
}

func Test_MsgAddTrustedMetadataSchemaProposer_GetSignBytes(t *testing.T) {
	actual := msgAddTrustedMetadataSchemaProposer.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgAddTrustedMetadataSchemaProposer))
	require.Equal(t, expected, actual)
}

func Test_MsgAddTrustedMetadataSchemaProposer_GetSigners(t *testing.T) {
	actual := msgAddTrustedMetadataSchemaProposer.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgAddTrustedMetadataSchemaProposer.Signer, actual[0])
}

func TestNewMsgShareDocument(t *testing.T) {
	tests := []struct {
		name     string
		document Document
		want     MsgShareDocument
	}{
		{
			"document creation",
			Document{
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
			},
			msgShareDocumentSchema,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, NewMsgShareDocument(tt.document))
		})
	}
}

func TestNewMsgSendDocumentReceipt(t *testing.T) {
	tests := []struct {
		name     string
		document DocumentReceipt
		want     MsgSendDocumentReceipt
	}{
		{
			"document receipt creation",
			DocumentReceipt{
				UUID:         "cfbb5b51-6ac0-43b0-8e09-022236285e31",
				Sender:       sender,
				Recipient:    recipient,
				TxHash:       "txHash",
				DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				Proof:        "proof",
			},
			msgDocumentReceipt,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, NewMsgSendDocumentReceipt(tt.document))
		})
	}
}
