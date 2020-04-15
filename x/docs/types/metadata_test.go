package types

import (
	"testing"

	"github.com/stretchr/testify/require"
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
			require.Equal(t, tt.want, value)
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
				require.EqualError(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
