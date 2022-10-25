package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var msgShareDocumentSchema = MsgShareDocument(ValidDocument)

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
				Sender:         ValidDocument.Sender,
				Recipients:     ValidDocument.Recipients,
				UUID:           ValidDocument.UUID,
				Metadata:       ValidDocument.Metadata,
				ContentURI:     ValidDocument.ContentURI,
				Checksum:       ValidDocument.Checksum,
				EncryptionData: ValidDocument.EncryptionData,
				DoSign:         ValidDocument.DoSign,
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
		message MsgShareDocument
		wantErr bool
	}{
		{
			name:    "ok",
			message: MsgShareDocument(ValidDocument),
		},
		{
			name:    "invalid message",
			message: MsgShareDocument(InvalidDocument),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.ValidateBasic()
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

	defer func() { recover() }()
	invalidMsg := msgShareDocumentSchema
	invalidMsg.Sender = ""
	invalidMsg.GetSigners()
	defer func() {
		t.Error("should have panicked")
	}()
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

var msgDocumentReceipt = MsgSendDocumentReceipt(ValidDocumentReceiptRecipient1)

func TestNewMsgSendDocumentReceipt(t *testing.T) {
	tests := []struct {
		name     string
		document DocumentReceipt
		want     MsgSendDocumentReceipt
	}{
		{
			"document receipt creation",
			DocumentReceipt{
				UUID:         ValidDocumentReceiptRecipient1.UUID,
				Sender:       ValidDocumentReceiptRecipient1.Sender,
				Recipient:    ValidDocumentReceiptRecipient1.Recipient,
				TxHash:       ValidDocumentReceiptRecipient1.TxHash,
				DocumentUUID: ValidDocumentReceiptRecipient1.DocumentUUID,
				Proof:        ValidDocumentReceiptRecipient1.Proof,
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

func TestMsgDocumentReceipt_GetSignBytes(t *testing.T) {
	actual := msgDocumentReceipt.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msgDocumentReceipt))
	require.Equal(t, expected, actual)
}

func TestMsgDocumentReceipt_GetSigners(t *testing.T) {
	actual := msgDocumentReceipt.GetSigners()
	require.Equal(t, 1, len(actual))
	require.Equal(t, msgDocumentReceipt.Sender, actual[0].String())

	defer func() { recover() }()
	invalidMsg := msgDocumentReceipt
	invalidMsg.Sender = ""
	invalidMsg.GetSigners()
	defer func() {
		t.Error("should have panicked")
	}()
}

func TestMsgSendDocumentReceipt_ValidateBasic(t *testing.T) {

	tests := []struct {
		name    string
		message MsgSendDocumentReceipt
		wantErr bool
	}{
		{
			name:    "ok",
			message: MsgSendDocumentReceipt(ValidDocumentReceiptRecipient1),
		},
		{
			name:    "invalid message",
			message: MsgSendDocumentReceipt(InvalidDocumentReceipt),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.message.ValidateBasic(); (err != nil) != tt.wantErr {
				t.Errorf("MsgSendDocumentReceipt.ValidateBasic() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
