package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentReceiptsIDs_AppendIfMissing(t *testing.T) {
	tests := []struct {
		name           string
		us             DocumentReceiptsIDs
		newData        string
		want           DocumentReceiptsIDs
		alreadyPresent bool
	}{
		{
			"adding a new element",
			DocumentReceiptsIDs{},
			"newElement",
			DocumentReceiptsIDs{
				"newElement",
			},
			true,
		},
		{
			"adding an existing element",
			DocumentReceiptsIDs{
				"newElement",
			},
			"newElement",
			DocumentReceiptsIDs{
				"newElement",
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

func TestDocumentReceiptsIDs_Empty(t *testing.T) {
	tests := []struct {
		name     string
		receipts DocumentReceiptsIDs
		empty    bool
	}{
		{
			"really empty",
			DocumentReceiptsIDs{},
			true,
		},
		{
			"not empty",
			DocumentReceiptsIDs{
				"a thing",
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.empty, tt.receipts.IsEmpty())
		})
	}
}
