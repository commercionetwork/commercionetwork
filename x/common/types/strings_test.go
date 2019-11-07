package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/stretchr/testify/assert"
)

func TestStrings_AppendIfMissing(t *testing.T) {
	tests := []struct {
		name             string
		strings          types.Strings
		string           string
		shouldBeAppended bool
	}{
		{
			name:             "Existing string is not appended",
			strings:          types.Strings{"first", "second"},
			string:           "first",
			shouldBeAppended: false,
		},
		{
			name:             "New string is appended into existing list",
			strings:          types.Strings{"first", "second"},
			string:           "third",
			shouldBeAppended: true,
		},
		{
			name:             "New string is appended into empty list",
			strings:          types.Strings{},
			string:           "first",
			shouldBeAppended: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, appended := test.strings.AppendIfMissing(test.string)
			assert.Equal(t, test.shouldBeAppended, appended)
			assert.Contains(t, result, test.string)
		})
	}
}

func TestStrings_Contains(t *testing.T) {
	assert.False(t, types.Strings{}.Contains("first"))
	assert.False(t, types.Strings{"first"}.Contains("seceond"))
	assert.True(t, types.Strings{"first", "second"}.Contains("first"))
	assert.True(t, types.Strings{"first", "second"}.Contains("second"))
}

func TestStrings_Equals(t *testing.T) {
	assert.False(t, types.Strings{}.Equals(types.Strings{"first"}))
	assert.False(t, types.Strings{"first"}.Equals(types.Strings{""}))
	assert.False(t, types.Strings{"first"}.Equals(types.Strings{"second"}))
	assert.True(t, types.Strings{"first", "second"}.Equals(types.Strings{"first", "second"}))
	assert.False(t, types.Strings{"first", "second"}.Equals(types.Strings{"second", "first"}))
}
