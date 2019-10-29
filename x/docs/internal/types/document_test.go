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
