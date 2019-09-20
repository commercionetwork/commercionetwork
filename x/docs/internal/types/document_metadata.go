package types

import (
	"errors"
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
// TODO: Test this
func (metadata DocumentMetadata) Equals(other DocumentMetadata) bool {
	if metadata.Schema == nil && other.Schema == nil {
		return true
	}

	if metadata.Schema == nil || other.Schema == nil {
		return false
	}

	return metadata.ContentUri == other.ContentUri &&
		metadata.Proof == other.Proof &&
		metadata.Schema.Equals(*other.Schema)
}

// Validate tries to validate all the data contained inside the given
// DocumentMetadata and returns an error if something is wrong
// TODO: Test this
func (metadata DocumentMetadata) Validate() error {
	if len(metadata.ContentUri) == 0 {
		return errors.New("metadataSchema content URI can't be empty")
	}

	if (metadata.Schema == nil) && len(metadata.SchemaType) == 0 {
		return errors.New("either schema or schema_type must be defined")
	}

	if metadata.Schema != nil {
		if len(metadata.Schema.Uri) == 0 {
			return errors.New("schema URI can't be empty")
		}
		if len(metadata.Schema.Version) == 0 {
			return errors.New("schema version can't be empty")
		}
	}

	if len(metadata.Proof) == 0 {
		return errors.New("computation proof can't be empty")
	}
	return nil
}
