// DONTCOVER
// nolint
package v2_1_2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "vbr"
)

// ---------------
// --- Genesis
// ---------------
// v2.1.2 vbr genesis state
type GenesisState struct {
	PoolAmount       sdk.DecCoins `json:"pool_amount"`
	YearlyPoolAmount sdk.DecCoins `json:"yearly_pool_amount"`
	YearNumber       int64        `json:"year_number"`
}
