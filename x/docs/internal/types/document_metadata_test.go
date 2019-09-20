package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// -------------------------
// --- MetadataSchema validation
// -------------------------

func TestValidateDocMetadata_valid(t *testing.T) {
	validDocumentMetadata := DocumentMetadata{
		ContentUri: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			Uri:     "http://www.contentUri.com",
			Version: "test",
		},
		Proof: "proof",
	}

	actual := validDocumentMetadata.Validate()
	assert.Nil(t, actual)
}

func TestValidateDocMetadata_emptyContentUri(t *testing.T) {
	invalidDocumentMetadata := DocumentMetadata{
		ContentUri: "",
		Schema: &DocumentMetadataSchema{
			Uri:     "http://www.contentUri.com",
			Version: "test",
		},
		Proof: "proof",
	}

	actual := invalidDocumentMetadata.Validate()
	assert.NotNil(t, actual)
}

func TestValidateDocMetadata_emptySchemaUri(t *testing.T) {
	invalidDocumentMetadata := DocumentMetadata{
		ContentUri: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			Uri:     "",
			Version: "test",
		},
		Proof: "proof",
	}

	actual := invalidDocumentMetadata.Validate()
	assert.NotNil(t, actual)
}

func TestValidateDocMetadata_emptySchemaVersion(t *testing.T) {
	invalidDocumentMetadata := DocumentMetadata{
		ContentUri: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			Uri:     "http://www.contentUri.com",
			Version: "",
		},
		Proof: "proof",
	}

	actual := invalidDocumentMetadata.Validate()
	assert.NotNil(t, actual)
}

func TestValidateDocMetadata_emptyProof(t *testing.T) {
	invalidDocumentMetadata := DocumentMetadata{
		ContentUri: "http://www.contentUri.com",
		Schema: &DocumentMetadataSchema{
			Uri:     "http://www.contentUri.com",
			Version: "test",
		},
		Proof: "",
	}

	actual := invalidDocumentMetadata.Validate()
	assert.NotNil(t, actual)
}
