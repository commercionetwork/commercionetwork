package types

import (
	"errors"
	"strings"
)

// MetadataSchema contains the information about an
// officially supported document metadata schema
type MetadataSchema struct {
	Type      string `json:"type"`
	SchemaURI string `json:"schema_uri"`
	Version   string `json:"version"`
}

// Equals returns true iff other contains the same data as this metadata schema
func (m MetadataSchema) Equals(other MetadataSchema) bool {
	return m.Type == other.Type &&
		m.SchemaURI == other.SchemaURI &&
		m.Version == other.Version
}

// Validate allows to validate the content of the schema, returning
// an error if something is not valid
func (m MetadataSchema) Validate() error {
	if len(strings.TrimSpace(m.Type)) == 0 {
		return errors.New("type cannot be empty")
	}
	if len(strings.TrimSpace(m.SchemaURI)) == 0 {
		return errors.New("uri cannot be empty")
	}
	if len(strings.TrimSpace(m.Version)) == 0 {
		return errors.New("version cannot be empty")
	}

	return nil
}
