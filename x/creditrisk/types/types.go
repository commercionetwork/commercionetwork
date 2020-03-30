package types

import (
	"github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName   = "creditrisk"
	RouterKey    = ModuleName
	QuerierRoute = ModuleName
	StoreKey     = ModuleName

	QueryPool = "pool"
)

type GenesisState struct {
	Pool types.Coins `json:"pool"`
}
