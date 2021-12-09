package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// XML Datetime normalized to UTC 00:00:00 and without sub-second decimal precision
	ComplaintW3CTime = time.RFC3339
)

func obtainTimestamp(ctx sdk.Context) string {
	// Following the W3C Decentralized Identifiers (DIDs) v1.0 for

	return ctx.BlockTime().Format(ComplaintW3CTime)
}

func readTimestamp(t string) (time.Time, error) {
	return time.Parse(ComplaintW3CTime, t)
}
