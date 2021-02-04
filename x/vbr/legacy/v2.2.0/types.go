// DONTCOVER
// nolint
package v2_2_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "vbr"
)

// ---------------
// --- Genesis
// ---------------
// v2.2.0 vbr genesis state
type GenesisState struct {
	PoolAmount        sdk.DecCoins `json:"pool_amount"`
	RewardRate        sdk.Dec      `json:"reward_rate"`
	AutomaticWithdraw bool         `json:"automatic_withdraw"`
}
