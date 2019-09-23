package types

import (
	"errors"
	"strings"
)

// DocumentMetadataSchema represents the information about the schema that should be used in order to
// validate the metadata associated with a document.
type DocumentMetadataSchema struct {
	Uri     string `json:"uri"`
	Version string `json:"version"`
}

func (metaSchema DocumentMetadataSchema) Equals(metSchema2 DocumentMetadataSchema) bool {
	return metaSchema.Uri == metSchema2.Uri &&
		metaSchema.Version == metSchema2.Version
}

// DocumentMetadata represents the information about the metadata associated to a document
type DocumentMetadata struct {
	ContentUri string                  `json:"content_uri"`
	SchemaType string                  `json:"schema_type"` // Optional - Either this or schema must be defined
	Schema     *DocumentMetadataSchema `json:"schema"`      // Optional - Either this or schema_type must be defined
	Proof      string                  `json:"proof"`
}

// Equals returns true iff this metadata and other contain the same data
func (metadata DocumentMetadata) Equals(other DocumentMetadata) bool {
	if metadata.ContentUri != other.ContentUri || metadata.Proof != other.Proof {
		return false
	}

	if metadata.Schema != nil && other.Schema != nil {
		return metadata.Schema.Equals(*other.Schema)
	}

	return metadata.Schema == other.Schema
}

// Validate tries to validate all the data contained inside the given
// DocumentMetadata and returns an error if something is wrong
func (metadata DocumentMetadata) Validate() error {
	if len(strings.TrimSpace(metadata.ContentUri)) == 0 {
		return errors.New("metadata.content_uri can't be empty")
	}

	if (metadata.Schema == nil) && len(strings.TrimSpace(metadata.SchemaType)) == 0 {
		return errors.New("either metadata.schema or metadata.schema_type must be defined")
	}

	if metadata.Schema != nil {
		if len(strings.TrimSpace(metadata.Schema.Uri)) == 0 {
			return errors.New("metadata.schema.uri can't be empty")
		}
		if len(strings.TrimSpace(metadata.Schema.Version)) == 0 {
			return errors.New("metadata.schema.version can't be empty")
		}
	}

	if len(strings.TrimSpace(metadata.Proof)) == 0 {
		return errors.New("metadata.proof can't be empty")
	}
	return nil
}
