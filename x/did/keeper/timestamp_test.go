package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_Timestamp(t *testing.T) {
	_, ctx := setupKeeper(t)

	now := time.Now().UTC()

	ctx = ctx.WithBlockTime(now)

	timestamp := obtainTimestamp(ctx)

	got, err := readTimestamp(timestamp)
	require.NoError(t, err)

	areTimingsEqualForW3CComplaintFormat(t, now, got)
}

// func Test_obtainTimestamp(t *testing.T) {

// 	type args struct {
// 		ctx sdk.Context
// 	}
// 	tests := []struct {
// 		name    string
// 		want    string
// 		wantErr bool
// 	}{
// 		{"OK 1", "2021-03-18T09:30:10Z", false},
// 		{"OK 2", "2020-12-20T19:17:47Z", false},
// 		{"OK 3", "2021-02-28T09:30:10Z", false},
// 		{"OK NOW", "NOW", false},
// 		{"BAD", "2021-02-29T23:61:47Z", true},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, ctx := setupKeeper(t)

// 			var want time.Time
// 			var err error
// 			if tt.want == "NOW" {
// 				want = time.Now().UTC()
// 			} else {
// 				want, err = readTimestamp(tt.want)
// 			}
// 			// if tt.wantErr {
// 			// 	require.Error(t, err)
// 			// }
// 			t.Log(err)

// 			ctx = ctx.WithBlockTime(want)

// 			got := obtainTimestamp(ctx)

// 			actual, err := readTimestamp(got)
// 			// if tt.wantErr {
// 			// 	require.Error(t, err)
// 			// }

// 			areTimingsEqualForW3CComplaintFormat(t, want, actual)
// 		})
// 	}
// }

func areTimingsEqualForW3CComplaintFormat(t *testing.T, wanted time.Time, actual time.Time) {
	require.Equal(t, wanted.Year(), actual.Year())
	require.Equal(t, wanted.Month(), actual.Month())
	require.Equal(t, wanted.Day(), actual.Day())
	require.Equal(t, wanted.Hour(), actual.Hour())
	require.Equal(t, wanted.Minute(), actual.Minute())
	require.Equal(t, wanted.Second(), actual.Second())
}
