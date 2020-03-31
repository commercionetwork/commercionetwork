package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocumentMetadata_Equals(t *testing.T) {
	tests := []struct {
		name   string
		us     DocumentMetadata
		them   DocumentMetadata
		equals bool
	}{
		{
			"two equal DocumentMetadata",
			DocumentMetadata{
				ContentURI: "http://example.com/metadata",
				Schema: &DocumentMetadataSchema{
					URI:     "https://example.com/metadata/schema",
					Version: "1.0.0",
				},
			},
			DocumentMetadata{
				ContentURI: "http://example.com/metadata",
				Schema: &DocumentMetadataSchema{
					URI:     "https://example.com/metadata/schema",
					Version: "1.0.0",
				},
			},
			true,
		},
		{
			"difference in contentURI",
			DocumentMetadata{ContentURI: "http://example.com/metadata"},
			DocumentMetadata{ContentURI: "http://example.com/metadat"},
			false,
		},
		{
			"difference in schema",
			DocumentMetadata{ContentURI: "http://example.com/metadata"},
			DocumentMetadata{
				ContentURI: "http://example.com/metadata",
				SchemaType: "",
				Schema: &DocumentMetadataSchema{
					URI:     "https://example.com/metadata/schema",
					Version: "1.0.0",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.equals, tt.us.Equals(tt.them))
		})
	}
}

// ---------------
// --- Validate
// ---------------

func TestDocumentMetadata_Validate(t *testing.T) {
	tests := []struct {
		name    string
		dm      DocumentMetadata
		wantErr error
	}{
		{
			"a valid DocumentMetadata",
			DocumentMetadata{
				ContentURI: "http://www.contentUri.com",
				Schema: &DocumentMetadataSchema{
					URI:     "http://www.contentUri.com",
					Version: "test",
				},
			},
			nil,
		},
		{
			"empty content uri",
			DocumentMetadata{
				ContentURI: "   ",
				Schema: &DocumentMetadataSchema{
					URI:     "http://www.contentUri.com",
					Version: "test",
				},
			},
			errors.New("metadata.content_uri can't be empty"),
		},
		{
			"empty metadata info",
			DocumentMetadata{
				ContentURI: "https://example.com/metadata",
				Schema:     nil,
				SchemaType: "",
			},
			errors.New("either metadata.schema or metadata.schema_type must be defined"),
		},
		{
			"empty schema uri",
			DocumentMetadata{
				ContentURI: "http://www.contentUri.com",
				Schema: &DocumentMetadataSchema{
					URI:     "",
					Version: "test",
				},
			},
			errors.New("metadata.schema.uri can't be empty"),
		},
		{
			"empty schema version",
			DocumentMetadata{
				ContentURI: "http://www.contentUri.com",
				Schema: &DocumentMetadataSchema{
					URI:     "http://www.contentUri.com",
					Version: "",
				},
			},
			errors.New("metadata.schema.version can't be empty"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr != nil {
				require.EqualError(t, tt.dm.Validate(), tt.wantErr.Error())
			} else {
				require.NoError(t, tt.dm.Validate())
			}
		})
	}
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
	require.Equal(t, expected, metadata)
}
