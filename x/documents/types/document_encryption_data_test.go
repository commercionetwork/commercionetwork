package types

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var validDocumentEncryptionData = DocumentEncryptionData{
	Keys:          []*DocumentEncryptionKey{{Recipient: recipient.String(), Value: "6F7468657276616C7565"}},
	EncryptedData: []string{"content", "content_uri", "metadata.content_uri", "metadata.schema.uri"},
}

// --------------
// --- Equals
// --------------

func TestDocumentEncryptionData_Equals(t *testing.T) {
	tests := []struct {
		name  string
		us    DocumentEncryptionData
		them  DocumentEncryptionData
		equal bool
	}{
		{
			"two equal encryption data",
			validDocumentEncryptionData,
			validDocumentEncryptionData,
			true,
		},
		{
			"different key length",
			validDocumentEncryptionData,
			DocumentEncryptionData{Keys: []*DocumentEncryptionKey{}, EncryptedData: validDocumentEncryptionData.EncryptedData},
			false,
		},
		{
			"different keys",
			validDocumentEncryptionData,
			DocumentEncryptionData{
				Keys:          []*DocumentEncryptionKey{{Recipient: sender.String(), Value: validDocumentEncryptionData.Keys[0].Value}},
				EncryptedData: validDocumentEncryptionData.EncryptedData,
			},
			false,
		},
		{
			"different data length",
			validDocumentEncryptionData,
			DocumentEncryptionData{Keys: validDocumentEncryptionData.Keys, EncryptedData: []string{"content"}},
			false,
		},
		{
			"different encrypted data",
			validDocumentEncryptionData,
			DocumentEncryptionData{Keys: validDocumentEncryptionData.Keys, EncryptedData: []string{"metadata.schema.uri", "content_uri", "content", "metadata.content_uri"}},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.equal, tt.us.Equals(tt.them))
		})
	}
}

// ---------------
// --- Validate
// ---------------

func TestDocumentEncryptionData_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ed      DocumentEncryptionData
		wantErr error
	}{
		{
			"valid DocumentEncryptionData",
			validDocumentEncryptionData,
			nil,
		},
		{
			"empty keys",
			DocumentEncryptionData{Keys: nil, EncryptedData: validDocumentEncryptionData.EncryptedData},
			errors.New("encryption data keys cannot be empty"),
		},
		{
			"invalid keys (invalid address)",
			DocumentEncryptionData{
				Keys: []*DocumentEncryptionKey{
					{Recipient: "", Value: ""},
				},
				EncryptedData: validDocumentEncryptionData.EncryptedData,
			},
			errors.New("invalid address "),
		},
		{
			"invalid encryption data",
			DocumentEncryptionData{Keys: validDocumentEncryptionData.Keys, EncryptedData: []string{"invalid.data"}},
			errors.New("encrypted data not supported"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr != nil {
				require.EqualError(t, tt.ed.Validate(), tt.wantErr.Error())
			} else {
				require.NoError(t, tt.ed.Validate())
			}
		})
	}
}

// ---------------
// --- ContainsRecipient
// ---------------

func TestDocumentEncryptionData_ContainsRecipient(t *testing.T) {

	type args struct {
		recipient sdk.AccAddress
	}
	tests := []struct {
		name string
		data DocumentEncryptionData
		args args
		want bool
	}{
		{
			name: "contains",
			data: validDocumentEncryptionData,
			args: args{
				recipient: recipient,
			},
			want: true,
		},
		{
			name: "not contains",
			data: validDocumentEncryptionData,
			args: args{
				recipient: sender,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := validDocumentEncryptionData.ContainsRecipient(tt.args.recipient); got != tt.want {
				t.Errorf("DocumentEncryptionData.ContainsRecipient() = %v, want %v", got, tt.want)
			}
		})
	}
}
