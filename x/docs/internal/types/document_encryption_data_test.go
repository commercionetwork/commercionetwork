package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var data = DocumentEncryptionData{
	Keys: []DocumentEncryptionKey{
		{Recipient: recipient, Value: "6F7468657276616C7565", Encoding: "hex"},
		{Recipient: sender, Value: "b3RoZXIgdmFsdWU=", Encoding: "base64"},
	},
	EncryptedData: []string{"content", "content_uri", "metadata.content_uri", "metadata.schema.uri"},
}

// --------------
// --- Equals
// --------------

func TestDocumentEncryptionData_Equals(t *testing.T) {
	assert.True(t, data.Equals(data))
}

func TestDocumentEncryptionData_Equals_DifferentKeysLength(t *testing.T) {
	other := DocumentEncryptionData{Keys: []DocumentEncryptionKey{data.Keys[0]}, EncryptedData: data.EncryptedData}
	assert.False(t, data.Equals(other))
}

func TestDocumentEncryptionData_Equals_DifferentKeys(t *testing.T) {
	other := DocumentEncryptionData{Keys: []DocumentEncryptionKey{data.Keys[1], data.Keys[0]}, EncryptedData: data.EncryptedData}
	assert.False(t, data.Equals(other))
}

func TestDocumentEncryptionData_Equals_DifferentEncryptedDataLength(t *testing.T) {
	other := DocumentEncryptionData{Keys: data.Keys, EncryptedData: []string{"content"}}
	assert.False(t, data.Equals(other))
}

func TestDocumentEncryptionData_Equals_DifferentEncryptedData(t *testing.T) {
	other := DocumentEncryptionData{Keys: data.Keys, EncryptedData: []string{"metadata.schema.uri", "content_uri", "content", "metadata.content_uri"}}
	assert.False(t, data.Equals(other))
}

// ---------------
// --- Validate
// ---------------

func TestDocumentEncryptionData_Validate(t *testing.T) {
	assert.Nil(t, data.Validate())
}

func TestDocumentEncryptionData_Validate_EmptyKeys(t *testing.T) {
	invalid := DocumentEncryptionData{Keys: nil, EncryptedData: data.EncryptedData}
	err := invalid.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "encryption data keys cannot be empty", err.Error())
}

func TestDocumentEncryptionData_Validate_InvalidKey(t *testing.T) {
	invalid := DocumentEncryptionData{
		Keys: []DocumentEncryptionKey{
			{Recipient: nil, Value: "", Encoding: ""},
		},
		EncryptedData: data.EncryptedData,
	}
	err := invalid.Validate()

	assert.NotNil(t, err)
}

func TestDocumentEncryptionData_Validate_InvalidEncryptedData(t *testing.T) {
	invalid := DocumentEncryptionData{Keys: data.Keys, EncryptedData: []string{"invalid.data"}}
	err := invalid.Validate()

	assert.NotNil(t, err)
	assert.Equal(t, "encrypted data not supported", err.Error())
}
