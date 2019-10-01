package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	MembershipTypeGreen  = "green"
	MembershipTypeBronze = "bronze"
	MembershipTypeSilver = "silver"
	MembershipTypeGold   = "gold"
	MembershipTypeBlack  = "black"
)

// Membership contains the data of a membership associated to a specific user
type Membership struct {
	Owner          sdk.AccAddress `json:"owner"`
	MembershipType string         `json:"membership_type"`
}

// IsMembershipTypeValid returns true iff the given membership type if valid
func IsMembershipTypeValid(membershipType string) bool {
	return membershipType == MembershipTypeGreen ||
		membershipType == MembershipTypeBronze ||
		membershipType == MembershipTypeSilver ||
		membershipType == MembershipTypeGold ||
		membershipType == MembershipTypeBlack
}

// CanUpgrade returns true iff the currentMembershipType is a less important than the newMembership type and thus a
// user having a membership of the first type can upgrade to a one of the second type.
// TODO: Test this
func CanUpgrade(currentMembershipType string, newMembershipType string) bool {
	if currentMembershipType == newMembershipType {
		return false
	}

	if currentMembershipType == MembershipTypeGreen {
		return true
	}

	if currentMembershipType == MembershipTypeBronze {
		return newMembershipType != MembershipTypeGreen
	}

	if currentMembershipType == MembershipTypeSilver {
		return newMembershipType != MembershipTypeGreen && newMembershipType != MembershipTypeSilver
	}

	if currentMembershipType == MembershipTypeGold {
		return newMembershipType == MembershipTypeBlack
	}

	return false
}
