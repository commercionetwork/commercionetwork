package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var validDocumentEncryptionKey = DocumentEncryptionKey{
	Recipient: recipient.String(),
	Value:     "76616C7565",
}

// --------------
// --- Equals
// --------------

func TestDocumentEncryptionKey_Equals(t *testing.T) {
	tests := []struct {
		name  string
		us    DocumentEncryptionKey
		them  DocumentEncryptionKey
		equal bool
	}{
		{
			"two equal keys",
			validDocumentEncryptionKey,
			validDocumentEncryptionKey,
			true,
		},
		{
			"different recipient",
			validDocumentEncryptionKey,
			DocumentEncryptionKey{
				Recipient: sender.String(),
				Value:     validDocumentEncryptionKey.Value,
			},
			false,
		},
		{
			"different value",
			validDocumentEncryptionKey,
			DocumentEncryptionKey{
				Recipient: validDocumentEncryptionKey.Recipient,
				Value:     "6F7468657276616C7565",
			},
			false,
		},
		{
			"different encoding",
			validDocumentEncryptionKey,
			DocumentEncryptionKey{
				Recipient: validDocumentEncryptionKey.Recipient,
				Value:     validDocumentEncryptionKey.Value + "difference",
			},
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

func TestDocumentEncryptionKey_Validate(t *testing.T) {
	tests := []struct {
		name    string
		ek      DocumentEncryptionKey
		wantErr error
	}{
		{
			"a valid key",
			validDocumentEncryptionKey,
			nil,
		},
		{
			"empty value",
			DocumentEncryptionKey{Recipient: validDocumentEncryptionKey.Recipient, Value: "   "},
			errors.New("encryption key value cannot be empty"),
		},
		{
			"invalid hex",
			DocumentEncryptionKey{Recipient: validDocumentEncryptionKey.Recipient, Value: "^&*(^*(&*"},
			errors.New("invalid encryption key value (must be hex)"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr != nil {
				require.EqualError(t, tt.ek.Validate(), tt.wantErr.Error())
			} else {
				require.NoError(t, tt.ek.Validate())
			}
		})
	}
}
