// You can edit this code!
// Click here and start typing.
package types

import "testing"

func Test_hasDuplicate(t *testing.T) {

	tests := []struct {
		name  string
		slice []*Service
		want  bool
	}{
		{},
		{"sliceEmpty", []*Service{}, false},
		{"sliceNoDup", []*Service{{ID: "A"}}, false},
		{"sliceNoDup1", []*Service{{ID: "A"}, {ID: "B"}, {ID: "C"}}, false},
		{"sliceDup", []*Service{{ID: "A"}, {ID: "B"}, {ID: "A"}}, true},
		{"sliceDup1", []*Service{{ID: "A"}, {ID: "B"}, {ID: "A"}, {ID: ""}}, true},
		{"sliceDup2", []*Service{{ID: ""}, {ID: "B"}, {ID: "A"}, {ID: ""}}, true},
		{"sliceDup3", []*Service{{ID: "A"}, {ID: "A"}, {ID: "A"}, {ID: ""}}, true},
		{"sliceDup4", []*Service{{ID: "A"}, {ID: "A"}}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := ServiceSlice(tt.slice).hasDuplicate(); got != tt.want {
				t.Errorf("hasDuplicate() = %v, want %v", got, tt.want)
			}
		})
	}
}
