package v1_3_3

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "tbr"
)

// ---------------
// --- Genesis
// ---------------
// v1.3.3 vbr genesis state
type GenesisState struct {
	PoolAmount       sdk.DecCoins `json:"pool_amount"`
	YearlyPoolAmount sdk.DecCoins `json:"yearly_pool_amount"`
	YearNumber       int64        `json:"year_number"`
}
