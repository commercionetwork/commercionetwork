// DONTCOVER
// nolint
package v1_3_4

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "vbr"
)

// ---------------
// --- Genesis
// ---------------
// v1.3.4 vbr genesis state
type GenesisState struct {
	PoolAmount       sdk.DecCoins `json:"pool_amount"`
	YearlyPoolAmount sdk.DecCoins `json:"yearly_pool_amount"`
	YearNumber       int64        `json:"year_number"`
}
