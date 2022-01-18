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

	got, err := time.Parse(ComplaintW3CTime, timestamp)
	require.NoError(t, err)

	require.Equal(t, now.Year(), got.Year())
	require.Equal(t, now.Month(), got.Month())
	require.Equal(t, now.Day(), got.Day())
	require.Equal(t, now.Hour(), got.Hour())
	require.Equal(t, now.Minute(), got.Minute())
	require.Equal(t, now.Second(), got.Second())
}
