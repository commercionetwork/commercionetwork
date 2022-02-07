package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// Test vars

var msgShareDocumentSchema = MsgShareDocument(Document{
	Sender:     sender.String(),
	Recipients: append([]string{}, recipient.String()),
	UUID:       "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	Metadata: &DocumentMetadata{
		ContentURI: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			URI:     "http://www.contentUri.com",
			Version: "test",
		},
	},
	ContentURI: "http://www.contentUri.com",
	Checksum: &DocumentChecksum{
		Value:     "48656c6c6f20476f7068657221234567",
		Algorithm: "md5",
	},
})

// ----------------------
// --- MsgShareDocument
// ----------------------

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
				Metadata: &DocumentMetadata{
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
				Sender:     sender.String(),
				Recipients: append([]string{}, recipient.String()),
			},
			msgShareDocumentSchema,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, *NewMsgShareDocument(tt.document))
		})
	}
}
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
		wantErr bool
	}{
		{
			"MsgShareDocument with valid schema",
			msgShareDocumentSchema,
			true,
		},
		{
			"MsgShareDocument with no schema",
			MsgShareDocument(Document{
				UUID:       "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				ContentURI: "http://www.contentUri.com",
				Metadata: &DocumentMetadata{
					ContentURI: "http://www.contentUri.com",
				},
				Checksum: &DocumentChecksum{
					Value:     "48656c6c6f20476f7068657221234567",
					Algorithm: "md5",
				},
				Sender:     sender.String(),
				Recipients: append([]string{}, recipient.String()),
			}),
			true,
		},
		{
			"MsgShareDocument with no schema type",
			MsgShareDocument(Document{
				UUID:       "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				ContentURI: "http://www.contentUri.com",
				Metadata: &DocumentMetadata{
					ContentURI: "http://www.contentUri.com",
				},
				Checksum: &DocumentChecksum{
					Value:     "48656c6c6f20476f7068657221234567",
					Algorithm: "md5",
				},
				Sender:     sender.String(),
				Recipients: append([]string{}, recipient.String()),
			}),
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sdr.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgShareDocument_GetSignBytes(t *testing.T) {
	actual := msgShareDocumentSchema.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msgShareDocumentSchema))
	require.Equal(t, expected, actual)
}

func TestMsgShareDocument_GetSigners(t *testing.T) {
	actual := msgShareDocumentSchema.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgShareDocumentSchema.Sender, actual[0].String())
}

func TestMsgShareDocument_UnmarshalJson_Schema(t *testing.T) {
	json := `{"sender":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0",
			"recipients":["cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"], 
			"UUID":"6a2f41a3-c54c-fce8-32d2-0324e1c32e22", 
			"metadata":{
				"contentURI":"http://www.contentUri.com", 
				"schema":{
					"URI":"http://www.contentUri.com", 
					"version":"test"}},
			"contentURI":"http://www.contentUri.com",
			"checksum":
				{"value":"48656c6c6f20476f7068657221234567",
				"algorithm":"md5"}
			}`
	var msg MsgShareDocument
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	require.Equal(t, "http://www.contentUri.com", msg.Metadata.Schema.URI)
	require.Equal(t, "test", msg.Metadata.Schema.Version)
}

// -----------------------------
// --- MsgSendDocumentReceipt
// -----------------------------

var msgDocumentReceipt = MsgSendDocumentReceipt{
	UUID:         "cfbb5b51-6ac0-43b0-8e09-022236285e31",
	Sender:       sender.String(),
	Recipient:    recipient.String(),
	TxHash:       "txHash",
	DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
	Proof:        "proof",
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
				Sender:       sender.String(),
				Recipient:    recipient.String(),
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
			actual := NewMsgSendDocumentReceipt(tt.document.UUID, tt.document.Sender, tt.document.Recipient, tt.document.TxHash, tt.document.DocumentUUID, tt.document.Proof)
			require.Equal(t, tt.want, *actual)
		})
	}
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
		wantErr bool
	}{
		{
			"valid SendDocumentReceipt",
			msgDocumentReceipt,
			false,
		},
		{
			"empty UUID",
			MsgSendDocumentReceipt{
				Sender:       sender.String(),
				Recipient:    recipient.String(),
				TxHash:       "txHash",
				DocumentUUID: "",
				Proof:        "proof",
			},
			true,
		},
		{
			"invalid UUID",
			MsgSendDocumentReceipt{
				Sender:       sender.String(),
				Recipient:    recipient.String(),
				TxHash:       "txHash",
				DocumentUUID: "6a2f41a3",
				Proof:        "proof",
			},
			true,
		},
		{
			"empty sender",
			MsgSendDocumentReceipt{
				UUID:         "cfbb5b51-6ac0-43b0-8e09-022236285e31",
				Recipient:    recipient.String(),
				TxHash:       "txHash",
				DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				Proof:        "proof",
			},
			true,
		},
		{
			"empty recipient",
			MsgSendDocumentReceipt{
				UUID:         "cfbb5b51-6ac0-43b0-8e09-022236285e31",
				Sender:       sender.String(),
				TxHash:       "txHash",
				DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				Proof:        "proof",
			},
			true,
		},
		{
			"invalid sender",
			MsgSendDocumentReceipt{
				UUID:         "cfbb5b51-6ac0-43b0-8e09-022236285e31",
				Recipient:    recipient.String() + "$",
				TxHash:       "txHash",
				DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				Proof:        "proof",
			},
			true,
		},
		{
			"invalid recipient",
			MsgSendDocumentReceipt{
				UUID:         "cfbb5b51-6ac0-43b0-8e09-022236285e31",
				Sender:       sender.String() + "$",
				TxHash:       "txHash",
				DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				Proof:        "proof",
			},
			true,
		},
		{
			"empty TxHash",
			MsgSendDocumentReceipt{
				UUID:         "cfbb5b51-6ac0-43b0-8e09-022236285e31",
				Sender:       sender.String(),
				Recipient:    recipient.String(),
				DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
				Proof:        "proof",
			},
			true,
		},
		{
			"invalid document UUID",
			MsgSendDocumentReceipt{
				UUID:      "cfbb5b51-6ac0-43b0-8e09-022236285e31",
				Sender:    sender.String(),
				Recipient: recipient.String(),
				TxHash:    "txHash",
				Proof:     "proof",
			},
			true,
		},
		{
			"empty proof",
			MsgSendDocumentReceipt{
				UUID:      "cfbb5b51-6ac0-43b0-8e09-022236285e31",
				Sender:    sender.String(),
				Recipient: recipient.String(),
				TxHash:    "txHash",
			},
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sdr.ValidateBasic()
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestMsgDocumentReceipt_GetSignBytes(t *testing.T) {
	actual := msgDocumentReceipt.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msgDocumentReceipt))
	require.Equal(t, expected, actual)
}

func TestMsgDocumentReceipt_GetSigners(t *testing.T) {
	actual := msgDocumentReceipt.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgDocumentReceipt.Sender, actual[0].String())
}
