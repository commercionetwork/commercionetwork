package types

import "errors"

type MetadataSchema struct {
	Type      string `json:"type"`
	SchemaUri string `json:"schema_uri"`
	Version   string `json:"version"`
}

func (m MetadataSchema) Equals(other MetadataSchema) bool {
	return m.Type == other.Type &&
		m.SchemaUri == other.SchemaUri &&
		m.Version == other.Version
}

func (m MetadataSchema) Validate() error {
	if len(m.Type) == 0 {
		return errors.New("type cannot be empty")
	}
	if len(m.SchemaUri) == 0 {
		return errors.New("uri cannot be empty")
	}
	if len(m.Version) == 0 {
		return errors.New("version cannot be empty")
	}

	return nil
}

type MetadataSchemes []MetadataSchema

func (metadataSchemes MetadataSchemes) Contains(metadata MetadataSchema) bool {
	for _, m := range metadataSchemes {
		if m.Equals(metadata) {
			return true
		}
	}
	return false
}

func (metadataSchemes MetadataSchemes) IsTypeSupported(metadataType string) bool {
	for _, m := range metadataSchemes {
		if m.Type == metadataType {
			return true
		}
	}
	return false
}

func (metadataSchemes MetadataSchemes) AppendIfMissing(metadata MetadataSchema) MetadataSchemes {
	if metadataSchemes.Contains(metadata) {
		return metadataSchemes
	} else {
		return append(metadataSchemes, metadata)
	}
}
