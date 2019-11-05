package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentMetadata_Equals(t *testing.T) {
	metadata := DocumentMetadata{
		ContentURI: "http://example.com/metadata",
		Schema: &DocumentMetadataSchema{
			URI:     "https://example.com/metadata/schema",
			Version: "1.0.0",
		},
	}
	assert.True(t, metadata.Equals(metadata))
}

func TestDocumentMetadata_Equals_DifferentContents(t *testing.T) {
	metadata := DocumentMetadata{ContentURI: "http://example.com/metadata"}

	other := DocumentMetadata{ContentURI: "https://example.com"}
	assert.False(t, metadata.Equals(other))
}

func TestDocumentMetadata_Equals_DifferentSchema(t *testing.T) {
	metadata := DocumentMetadata{ContentURI: "http://example.com/metadata"}

	other := DocumentMetadata{
		ContentURI: metadata.ContentURI,
		SchemaType: metadata.SchemaType,
		Schema: &DocumentMetadataSchema{
			URI:     "https://example.com/metadata/schema",
			Version: "1.0.0",
		},
	}
	assert.False(t, metadata.Equals(other))
}

// ---------------
// --- Validate
// ---------------

func TestDocumentMetadata_Validate(t *testing.T) {
	validDocumentMetadata := DocumentMetadata{
		ContentURI: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			URI:     "http://www.contentUri.com",
			Version: "test",
		},
	}

	actual := validDocumentMetadata.Validate()
	assert.Nil(t, actual)
}

func TestDocumentMetadata_Validate_EmptyContentUri(t *testing.T) {
	invalidDocumentMetadata := DocumentMetadata{
		ContentURI: "   ",
		Schema: &DocumentMetadataSchema{
			URI:     "http://www.contentUri.com",
			Version: "test",
		},
	}

	err := invalidDocumentMetadata.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "metadata.content_uri can't be empty", err.Error())
}

func TestDocumentMetadata_Validate_EmptyMetadataInfo(t *testing.T) {
	invalidDocumentMetadata := DocumentMetadata{
		ContentURI: "https://example.com/metadata",
		Schema:     nil,
		SchemaType: "",
	}

	err := invalidDocumentMetadata.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "either metadata.schema or metadata.schema_type must be defined", err.Error())
}

func TestDocumentMetadata_Validate_EmptySchemaUri(t *testing.T) {
	invalidDocumentMetadata := DocumentMetadata{
		ContentURI: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			URI:     "",
			Version: "test",
		},
	}

	err := invalidDocumentMetadata.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "metadata.schema.uri can't be empty", err.Error())
}

func TestDocumentMetadata_Validate_EmptySchemaVersion(t *testing.T) {
	invalidDocumentMetadata := DocumentMetadata{
		ContentURI: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			URI:     "http://www.contentUri.com",
			Version: "",
		},
	}

	err := invalidDocumentMetadata.Validate()
	assert.NotNil(t, err)
	assert.Equal(t, "metadata.schema.version can't be empty", err.Error())
}

func TestDocumentMetadata_JSONUnmarshal(t *testing.T) {
	json := `{"content_uri":"http://www.contentUri.com","schema":{"uri":"http://www.contentUri.com","version":"1.0.0"}}`

	var metadata DocumentMetadata
	ModuleCdc.MustUnmarshalJSON([]byte(json), &metadata)

	expected := DocumentMetadata{
		ContentURI: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			URI:     "http://www.contentUri.com",
			Version: "1.0.0",
		},
	}
	assert.Equal(t, expected, metadata)
}
