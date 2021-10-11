package v2_2_0

import (
	"time"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "commerciokyc"
)

type InviteStatus uint8

const (
	InviteStatusPending InviteStatus = iota
	InviteStatusRewarded
	InviteStatusInvalid
)

type Membership struct {
	Owner          sdk.AccAddress `json:"owner"`
	TspAddress     sdk.AccAddress `json:"tsp_address"`
	MembershipType string         `json:"membership_type"`
	ExpiryAt       time.Time      `json:"expiry_at"` // Time at which the membership expired

}

type Memberships []Membership

type Invite struct {
	Sender           sdk.AccAddress `json:"sender"`            // User that has sent the invitation
	SenderMembership string         `json:"sender_membership"` // Membership of Sender when the invite was created
	User             sdk.AccAddress `json:"user"`              // Invited user
	Status           InviteStatus   `json:"status"`            // Tells if the invite is pending, rewarded or invalid
}

type Invites []Invite

type GenesisState struct {
	LiquidityPoolAmount     sdk.Coins        `json:"liquidity_pool_amount"`     // Liquidity pool from which to get the rewards
	Invites                 Invites          `json:"invites"`                   // List of invites
	TrustedServiceProviders ctypes.Addresses `json:"trusted_service_providers"` // List of trusted service providers
	Memberships             Memberships      `json:"memberships"`               // List of all the existing memberships
}
