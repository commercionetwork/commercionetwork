package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var data = DocumentEncryptionData{
	Keys:          []DocumentEncryptionKey{{Recipient: recipient, Value: "6F7468657276616C7565"}},
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
			data,
			data,
			true,
		},
		{
			"different key length",
			data,
			DocumentEncryptionData{Keys: []DocumentEncryptionKey{}, EncryptedData: data.EncryptedData},
			false,
		},
		{
			"different keys",
			data,
			DocumentEncryptionData{
				Keys:          []DocumentEncryptionKey{{Recipient: sender, Value: data.Keys[0].Value}},
				EncryptedData: data.EncryptedData,
			},
			false,
		},
		{
			"different data length",
			data,
			DocumentEncryptionData{Keys: data.Keys, EncryptedData: []string{"content"}},
			false,
		},
		{
			"different encrypted data",
			data,
			DocumentEncryptionData{Keys: data.Keys, EncryptedData: []string{"metadata.schema.uri", "content_uri", "content", "metadata.content_uri"}},
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
			data,
			nil,
		},
		{
			"empty keys",
			DocumentEncryptionData{Keys: nil, EncryptedData: data.EncryptedData},
			errors.New("encryption data keys cannot be empty"),
		},
		{
			"invalid keys (invalid address)",
			DocumentEncryptionData{
				Keys: []DocumentEncryptionKey{
					{Recipient: nil, Value: ""},
				},
				EncryptedData: data.EncryptedData,
			},
			errors.New("invalid address "),
		},
		{
			"invalid encryption data",
			DocumentEncryptionData{Keys: data.Keys, EncryptedData: []string{"invalid.data"}},
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
