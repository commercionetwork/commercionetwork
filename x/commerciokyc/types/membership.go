package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMembership returns a new memberships containing the given data
// TODO: fix conversion
func NewMembership(membershipType string, owner sdk.AccAddress, tsp sdk.AccAddress, expiryAt time.Time) Membership {
	return Membership{
		Owner:          owner.String(),
		TspAddress:     tsp.String(),
		MembershipType: membershipType,
		ExpiryAt:       &expiryAt, //TODO CONVERSION
	}
}

// Equals returns true iff m and other contain the same data
func (m Membership) Equals(other Membership) bool {
	return m.Owner == other.Owner &&
		m.TspAddress == other.TspAddress &&
		m.ExpiryAt.Equal(*other.ExpiryAt) &&
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
