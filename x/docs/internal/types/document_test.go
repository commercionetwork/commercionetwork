package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocument_Equals_NilValues(t *testing.T) {
	document := Document{
		Uuid: "uuid",
		Metadata: DocumentMetadata{
			ContentUri: "document_metadata_content_uri",
			SchemaType: "document_metadata_schema_type",
		},
		ContentUri:     "",
		Checksum:       nil,
		EncryptionData: nil,
	}
	assert.True(t, document.Equals(document))
}
