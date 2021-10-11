package types

import (
	"errors"
	"strings"
)

func (metaSchema DocumentMetadataSchema) Equals(metSchema2 DocumentMetadataSchema) bool {
	return metaSchema.URI == metSchema2.URI &&
		metaSchema.Version == metSchema2.Version
}

// Equals returns true iff this metadata and other contain the same data
func (metadata DocumentMetadata) Equals(other DocumentMetadata) bool {
	if metadata.ContentURI != other.ContentURI {
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
	if len(strings.TrimSpace(metadata.ContentURI)) == 0 {
		return errors.New("metadata.content_uri can't be empty")
	}

	if metadata.Schema == nil {
		return errors.New("metadata.schema must be defined")
	}

	if metadata.Schema != nil {
		if len(strings.TrimSpace(metadata.Schema.URI)) == 0 {
			return errors.New("metadata.schema.uri can't be empty")
		}
		if len(strings.TrimSpace(metadata.Schema.Version)) == 0 {
			return errors.New("metadata.schema.version can't be empty")
		}
	}
	return nil
}
