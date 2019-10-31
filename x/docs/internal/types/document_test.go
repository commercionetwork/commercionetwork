package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocument_Equals_NilValues(t *testing.T) {
	document := Document{
		UUID: "uuid",
		Metadata: DocumentMetadata{
			ContentURI: "document_metadata_content_uri",
			SchemaType: "document_metadata_schema_type",
		},
		ContentURI:     "",
		Checksum:       nil,
		EncryptionData: nil,
	}
	assert.True(t, document.Equals(document))
}

func Test_validateUUID(t *testing.T) {
	tests := []struct {
		name    string
		UUID    string
		badUUID bool
	}{
		{
			"empty string",
			"",
			true,
		},
		{
			"a well-formed UUID",
			"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
			false,
		},
		{
			"a seemingly well-formed UUID, with the last character removed",
			"6ba7b810-9dad-11d1-80b4-00c04fd430c",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := validateUUID(tt.UUID)
			if tt.badUUID {
				assert.False(t, val, "got true")
			} else {
				assert.True(t, val, "got false")
			}
		})
	}
}
