package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var key = DocumentEncryptionKey{
	Recipient: recipient,
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
			key,
			key,
			true,
		},
		{
			"different recipient",
			key,
			DocumentEncryptionKey{
				Recipient: sender,
				Value:     key.Value,
			},
			false,
		},
		{
			"different value",
			key,
			DocumentEncryptionKey{
				Recipient: key.Recipient,
				Value:     "6F7468657276616C7565",
			},
			false,
		},
		{
			"different encoding",
			key,
			DocumentEncryptionKey{
				Recipient: key.Recipient,
				Value:     key.Value + "difference",
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.equal, tt.us.Equals(tt.them))
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
			key,
			nil,
		},
		{
			"empty value",
			DocumentEncryptionKey{Recipient: key.Recipient, Value: "   "},
			errors.New("encryption key value cannot be empty"),
		},
		{
			"innvalid hex",
			DocumentEncryptionKey{Recipient: key.Recipient, Value: "^&*(^*(&*"},
			errors.New("invalid encryption key value (must be hex)"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr != nil {
				assert.Error(t, tt.wantErr, tt.ek.Validate())
			} else {
				assert.NoError(t, tt.ek.Validate())
			}
		})
	}
}
