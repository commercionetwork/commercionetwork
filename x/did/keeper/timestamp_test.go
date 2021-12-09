package keeper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_getTimestamp(t *testing.T) {
	_, ctx := setupKeeper(t)

	now := time.Now().UTC()

	ctx = ctx.WithBlockTime(now)

	timestamp := obtainTimestamp(ctx)

	tt, err := readTimestamp(timestamp)
	require.NoError(t, err)

	require.Equal(t, now.Year(), tt.Year())
	require.Equal(t, now.Month(), tt.Month())
	require.Equal(t, now.Day(), tt.Day())
	require.Equal(t, now.Hour(), tt.Hour())
	require.Equal(t, now.Minute(), tt.Minute())
	require.Equal(t, now.Second(), tt.Second())
}
