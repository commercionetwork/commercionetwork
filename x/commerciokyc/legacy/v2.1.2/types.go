package v2_1_2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	commKycTypes "github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
)

const (
	ModuleName = "accreditations"
)

type GenesisState struct {
	LiquidityPoolAmount     sdk.Coins            `json:"liquidity_pool_amount"`     // Liquidity pool from which to get the rewards
	Invites                 commKycTypes.Invites `json:"invites"`                   // List of invites
	TrustedServiceProviders ctypes.Addresses     `json:"trusted_service_providers"` // List of trusted service providers
	StableCreditsDenom      string               `json:"stable_credits_denom"`      // Stable credits denom used during membership buying
	Memberships             Memberships          `json:"memberships"`               // List of all the existing memberships
}

type Memberships []Membership

type Membership struct {
	Owner          sdk.AccAddress `json:"owner"`
	MembershipType string         `json:"membership_type"`
}
