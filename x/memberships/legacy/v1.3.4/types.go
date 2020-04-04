package v1_3_4

import (
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "accreditations"
)

type Invites []Invite

// Invite represents an invitation that a user has made towards another user
type Invite struct {
	Sender   sdk.AccAddress `json:"sender"`   // User that has sent the invitation
	User     sdk.AccAddress `json:"user"`     // Invited user
	Rewarded bool           `json:"rewarded"` // Tells if the invitee has already been rewarded
}

type GenesisState struct {
	LiquidityPoolAmount     sdk.Coins         `json:"liquidity_pool_amount"`     // Liquidity pool from which to get the rewards
	Invites                 Invites           `json:"invites"`                   // List of invites
	TrustedServiceProviders ctypes.Addresses  `json:"trusted_service_providers"` // List of trusted service providers
	StableCreditsDenom      string            `json:"stable_credits_denom"`      // Stable credits denom used during membership buying
	Memberships             types.Memberships `json:"memberships"`               // List of all the existing memberships
}
