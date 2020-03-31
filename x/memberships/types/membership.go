package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
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
	MembershipType string         `json:"membership_type"`
}

// NewMembership returns a new memberships containing the given data
func NewMembership(membershipType string, owner sdk.AccAddress) Membership {
	return Membership{
		Owner:          owner,
		MembershipType: membershipType,
	}
}

// Equals returns true iff m and other contain the same data
func (m Membership) Equals(other Membership) bool {
	return m.Owner.Equals(other.Owner) &&
		m.MembershipType == other.MembershipType
}

// IsMembershipTypeValid returns true iff the given membership type if valid
func IsMembershipTypeValid(membershipType string) bool {
	return membershipType == MembershipTypeBronze ||
		membershipType == MembershipTypeSilver ||
		membershipType == MembershipTypeGold ||
		membershipType == MembershipTypeBlack
}

// CanUpgrade returns true iff the currentMembershipType is a less important than the newMembership type and thus a
// user having a membership of the first type can upgrade to a one of the second type.
func CanUpgrade(currentMembershipType string, newMembershipType string) bool {
	if !IsMembershipTypeValid(currentMembershipType) || !IsMembershipTypeValid(newMembershipType) {
		return false
	}

	if currentMembershipType == newMembershipType {
		return false
	}

	if currentMembershipType == MembershipTypeBronze {
		return true
	}

	if currentMembershipType == MembershipTypeSilver {
		return newMembershipType != MembershipTypeBronze
	}

	if currentMembershipType == MembershipTypeGold {
		return newMembershipType == MembershipTypeBlack
	}

	return false
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
