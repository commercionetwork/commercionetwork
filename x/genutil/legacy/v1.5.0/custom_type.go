package v1_5_0

// DONTCOVER
// nolint

// PSA: this file only serves as a bridge between Cosmos SDK we used before v0.38 (we used master, commit 92ea174ea6e6),
// and should NEVER be used to instantiate new accounts or stuff like that.

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	//v134memberships "github.com/commercionetwork/commercionetwork/x/memberships/legacy/v1.3.4"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
)

type GenesisState struct {
	LiquidityPoolAmount     sdk.Coins         `json:"liquidity_pool_amount"`     // Liquidity pool from which to get the rewards
	TrustedServiceProviders ctypes.Addresses  `json:"trusted_service_providers"` // List of trusted service providers
	StableCreditsDenom      string            `json:"stable_credits_denom"`      // Stable credits denom used during membership buying
	Memberships             types.Memberships `json:"memberships"`               // List of all the existing memberships
}
