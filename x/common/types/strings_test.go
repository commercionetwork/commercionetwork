package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/stretchr/testify/require"
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
			require.Equal(t, test.shouldBeAppended, appended)
			require.Contains(t, result, test.string)
		})
	}
}

func TestStrings_Contains(t *testing.T) {
	require.False(t, types.Strings{}.Contains("first"))
	require.False(t, types.Strings{"first"}.Contains("seceond"))
	require.True(t, types.Strings{"first", "second"}.Contains("first"))
	require.True(t, types.Strings{"first", "second"}.Contains("second"))
}

func TestStrings_Equals(t *testing.T) {
	require.False(t, types.Strings{}.Equals(types.Strings{"first"}))
	require.False(t, types.Strings{"first"}.Equals(types.Strings{""}))
	require.False(t, types.Strings{"first"}.Equals(types.Strings{"second"}))
	require.True(t, types.Strings{"first", "second"}.Equals(types.Strings{"first", "second"}))
	require.False(t, types.Strings{"first", "second"}.Equals(types.Strings{"second", "first"}))
	require.False(t, types.Strings{"first", "second"}.Equals(types.Strings{"first"}))
	require.False(t, types.Strings{"first"}.Equals(types.Strings{"first", "second"}))
}
