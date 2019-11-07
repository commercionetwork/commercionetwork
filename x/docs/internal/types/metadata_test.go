package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadataSchema_Equals(t *testing.T) {
	tests := []struct {
		name  string
		us    MetadataSchema
		other MetadataSchema
		want  bool
	}{
		{
			"two equal MetadataSchemas",
			MetadataSchema{
				"type",
				"schemaUri",
				"version",
			},
			MetadataSchema{
				"type",
				"schemaUri",
				"version",
			},
			true,
		},
		{
			"two non-equal MetadataSchemas",
			MetadataSchema{
				"type",
				"schemaUri",
				"",
			},
			MetadataSchema{
				"type",
				"schemaUri",
				"version",
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			value := tt.us.Equals(tt.other)
			assert.Equal(t, tt.want, value)
		})
	}
}

func TestMetadataSchema_Validate(t *testing.T) {
	tests := []struct {
		name    string
		us      MetadataSchema
		wantErr string
	}{
		{
			"a valid MetadataSchema",
			MetadataSchema{
				"type",
				"schemaUri",
				"version",
			},
			"",
		},
		{
			"missing type",
			MetadataSchema{
				"",
				"schemaUri",
				"version",
			},
			"type cannot be empty",
		},
		{
			"missing uri",
			MetadataSchema{
				"type",
				"",
				"version",
			},
			"uri cannot be empty",
		},
		{
			"missing version",
			MetadataSchema{
				"type",
				"schemaUri",
				"",
			},
			"version cannot be empty",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.us.Validate()
			if tt.wantErr != "" {
				assert.EqualError(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMetadataSchemes_Contains(t *testing.T) {
	tests := []struct {
		name            string
		metadataSchemes MetadataSchemes
		metadataSchema  MetadataSchema
		want            bool
	}{
		{
			"does not contain this MetadataSchema",
			MetadataSchemes{},
			MetadataSchema{
				"type",
				"schemaUri",
				"version",
			},
			false,
		},
		{
			"does contain this MetadataSchema",
			MetadataSchemes{
				MetadataSchema{
					"type",
					"schemaUri",
					"version",
				},
			},
			MetadataSchema{
				"type",
				"schemaUri",
				"version",
			},
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			val := tt.metadataSchemes.Contains(tt.metadataSchema)
			assert.Equal(t, tt.want, val)
		})
	}
}

func TestMetadataSchemes_IsTypeSupported(t *testing.T) {
	tests := []struct {
		name            string
		metadataSchemes MetadataSchemes
		metadataType    string
		want            bool
	}{
		{
			"specified type is contained in a MetadataScheme inside MetadataSchemes",
			MetadataSchemes{
				MetadataSchema{
					"type",
					"schemaUri",
					"version",
				},
			},
			"type",
			true,
		},
		{
			"specified type is not contained in a MetadataScheme inside MetadataSchemes",
			MetadataSchemes{
				MetadataSchema{
					"type",
					"schemaUri",
					"version",
				},
			},
			"anotherType",
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			val := tt.metadataSchemes.IsTypeSupported(tt.metadataType)
			assert.Equal(t, tt.want, val)
		})
	}
}

func TestMetadataSchemes_AppendIfMissing(t *testing.T) {
	tests := []struct {
		name           string
		us             MetadataSchemes
		newData        MetadataSchema
		want           MetadataSchemes
		alreadyPresent bool
	}{
		{
			"adding a new element",
			MetadataSchemes{},
			MetadataSchema{
				"type",
				"schemaUri",
				"version",
			},
			MetadataSchemes{
				{
					"type",
					"schemaUri",
					"version",
				},
			},
			true,
		},
		{
			"adding an already present element",
			MetadataSchemes{
				{
					"type",
					"schemaUri",
					"version",
				},
			},
			MetadataSchema{
				"type",
				"schemaUri",
				"version",
			},
			MetadataSchemes{
				{
					"type",
					"schemaUri",
					"version",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			val, present := tt.us.AppendIfMissing(tt.newData)

			assert.Equal(t, tt.alreadyPresent, present)
			assert.Equal(t, tt.want, val)
		})
	}
}
