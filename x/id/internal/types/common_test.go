package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestValidateHex(t *testing.T) {
	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{
			"empty string",
			"",
			false,
		},
		{
			"just spaces",
			"    ",
			false,
		},
		{
			"random sequence of characters",
			"dasfasdfdf897987",
			false,
		},
		{
			"a well-formed hex string",
			"6369616f6369616f63",
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.valid, types.ValidateHex(tt.input))
		})
	}
}
