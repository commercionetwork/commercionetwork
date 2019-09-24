package types

import (
	"errors"
	"strings"
)

// MetadataSchema contains the information about an
// officially supported document metadata schema
type MetadataSchema struct {
	Type      string `json:"type"`
	SchemaUri string `json:"schema_uri"`
	Version   string `json:"version"`
}

// Equals returns true iff other contains the same data as this metadata schema
func (m MetadataSchema) Equals(other MetadataSchema) bool {
	return m.Type == other.Type &&
		m.SchemaUri == other.SchemaUri &&
		m.Version == other.Version
}

// Validate allows to validate the content of the schema, returning
// an error if something is not valid
func (m MetadataSchema) Validate() error {
	if len(strings.TrimSpace(m.Type)) == 0 {
		return errors.New("type cannot be empty")
	}
	if len(strings.TrimSpace(m.SchemaUri)) == 0 {
		return errors.New("uri cannot be empty")
	}
	if len(strings.TrimSpace(m.Version)) == 0 {
		return errors.New("version cannot be empty")
	}

	return nil
}

// MetadataSchemes represents a list of MetadataSchema
type MetadataSchemes []MetadataSchema

// Contains returns true iff the specified metadata is present inside this list
func (metadataSchemes MetadataSchemes) Contains(metadata MetadataSchema) bool {
	for _, m := range metadataSchemes {
		if m.Equals(metadata) {
			return true
		}
	}
	return false
}

// IsTypeSupported allows to tell if there is one metadata
// scheme having the given type inside this list
func (metadataSchemes MetadataSchemes) IsTypeSupported(metadataType string) bool {
	for _, m := range metadataSchemes {
		if m.Type == metadataType {
			return true
		}
	}
	return false
}

// AppendIfMissing allows to add to this list of schemes the given schema, if it isn't already present
func (metadataSchemes MetadataSchemes) AppendIfMissing(schema MetadataSchema) MetadataSchemes {
	if metadataSchemes.Contains(schema) {
		return metadataSchemes
	} else {
		return append(metadataSchemes, schema)
	}
}
