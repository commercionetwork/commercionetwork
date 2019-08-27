package types

const (
	MembershipTypeGreen  = "green"
	MembershipTypeBronze = "bronze"
	MembershipTypeSilver = "silver"
	MembershipTypeGold   = "gold"
	MembershipTypeBlack  = "black"
)

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
func CanUpgrade(currentMembershipType string, newMembershipType string) bool {
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
