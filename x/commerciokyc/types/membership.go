package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MembershipTypeGreen  = "green"
	MembershipTypeBronze = "bronze"
	MembershipTypeSilver = "silver"
	MembershipTypeGold   = "gold"
	MembershipTypeBlack  = "black"
)

// -------------------
// --- Membership
// -------------------

// Membership contains the data of a membership associated to a specific user
type Membership struct {
	Owner          sdk.AccAddress `json:"owner"`
	TspAddress     sdk.AccAddress `json:"tsp_address"`
	MembershipType string         `json:"membership_type"`
	ExpiryAt       int64          `json:"expiry_at"` // Block height at which the membership expired

}

// NewMembership returns a new memberships containing the given data
func NewMembership(membershipType string, owner sdk.AccAddress, tsp sdk.AccAddress, expiryAt int64) Membership {
	return Membership{
		Owner:          owner,
		TspAddress:     tsp,
		MembershipType: membershipType,
		ExpiryAt:       expiryAt,
	}
}

// Equals returns true iff m and other contain the same data
func (m Membership) Equals(other Membership) bool {
	return m.Owner.Equals(other.Owner) &&
		m.TspAddress.Equals(other.TspAddress) &&
		m.ExpiryAt == other.ExpiryAt &&
		m.MembershipType == other.MembershipType
}

// IsMembershipTypeValid returns true iff the given membership type if valid
func IsMembershipTypeValid(membershipType string) bool {
	return membershipType == MembershipTypeGreen ||
		membershipType == MembershipTypeBronze ||
		membershipType == MembershipTypeSilver ||
		membershipType == MembershipTypeGold ||
		membershipType == MembershipTypeBlack
}

// -------------------
// --- Memberships
// -------------------

// Memberships represents a slice of Membership objects
type Memberships []Membership

// AppendIfMissing appends the other membership to the given slice, returning the result of the appending
func (slice Memberships) AppendIfMissing(other Membership) (Memberships, bool) {
	for _, membership := range slice {
		if membership.Equals(other) {
			return slice, false
		}
	}
	return append(slice, other), true
}

// Equals returns true if this slice and the other one contain the same memberships
// in the same exact order
func (slice Memberships) Equals(other Memberships) bool {
	if len(slice) != len(other) {
		return false
	}

	for index, m := range slice {
		if !m.Equals(other[index]) {
			return false
		}
	}

	return true
}

// ValidateBasic returns error if Membership type is not valid
func (m Membership) ValidateBasic() error {
	if !IsMembershipTypeValid(m.MembershipType) {
		return fmt.Errorf("membership has invalid type: %s", m.MembershipType)
	}
	return nil
}
