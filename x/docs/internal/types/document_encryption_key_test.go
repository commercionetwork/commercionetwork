package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var key = DocumentEncryptionKey{
	Recipient: recipient,
	Value:     "76616C7565",
	Encoding:  "hex",
}

// --------------
// --- Equals
// --------------

func TestDocumentEncryptionKey_Equals(t *testing.T) {
	assert.True(t, key.Equals(key))
}

func TestDocumentEncryptionKey_Equals_DifferentRecipient(t *testing.T) {
	var other = DocumentEncryptionKey{
		Recipient: sender,
		Value:     key.Value,
		Encoding:  key.Encoding,
	}
	assert.False(t, key.Equals(other))
}

func TestDocumentEncryptionKey_Equals_DifferentValue(t *testing.T) {
	var other = DocumentEncryptionKey{
		Recipient: key.Recipient,
		Value:     "6F7468657276616C7565",
		Encoding:  key.Encoding,
	}
	assert.False(t, key.Equals(other))
}

func TestDocumentEncryptionKey_Equals_DifferentEncoding(t *testing.T) {
	var other = DocumentEncryptionKey{
		Recipient: key.Recipient,
		Value:     key.Value,
		Encoding:  key.Encoding + "difference",
	}
	assert.False(t, key.Equals(other))
}

// ---------------
// --- Validate
// ---------------

func TestDocumentEncryptionKey_Validate(t *testing.T) {
	assert.Nil(t, key.Validate())
}

func TestDocumentEncryptionKey_Validate_EmptyValue(t *testing.T) {
	key := DocumentEncryptionKey{Recipient: key.Recipient, Value: "   ", Encoding: key.Encoding}
	err := key.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "encryption key value cannot be empty", err.Error())
}

func TestDocumentEncryptionKey_Validate_EmptyEncoding(t *testing.T) {
	key := DocumentEncryptionKey{Recipient: key.Recipient, Value: key.Encoding, Encoding: "   "}
	err := key.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "encryption key encoding cannot be empty", err.Error())
}

func TestDocumentEncryptionKey_Validate_InvalidEncoding(t *testing.T) {
	key := DocumentEncryptionKey{Recipient: key.Recipient, Value: key.Encoding, Encoding: "faslgjfsdlgj"}
	err := key.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "encryption key encoding method unknown", err.Error())
}

func TestDocumentEncryptionKey_Validate_InvalidHex(t *testing.T) {
	key := DocumentEncryptionKey{Recipient: key.Recipient, Value: "^&*(^*(&*", Encoding: "hex"}
	err := key.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "hex")
}

func TestDocumentEncryptionKey_Validate_InvalidBase64(t *testing.T) {
	key := DocumentEncryptionKey{Recipient: key.Recipient, Value: "^&*(^*(&*", Encoding: "base64"}
	err := key.Validate()

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "base64")
}
