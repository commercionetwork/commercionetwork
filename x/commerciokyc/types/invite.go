package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type InviteStatus uint64

const (
	InviteStatusPending InviteStatus = iota
	InviteStatusRewarded
	InviteStatusInvalid
)

// NewMembership returns a new memberships containing the given data
// TODO: fix conversion
func NewInvite(sender, user sdk.AccAddress, senderMembership string) Invite {
	return Invite{
		Sender:           sender.String(),
		User:             user.String(),
		SenderMembership: senderMembership,
		Status:           uint64(InviteStatusPending), // TODO fix conversion
	}
}

// Empty returns true of the given invite is empty
func (invite Invite) Empty() bool {
	return Invite{}.Equals(invite)
}

// Equals returns true iff invite contains the same data of the other invite
func (invite Invite) Equals(other Invite) bool {
	return invite.Sender == other.Sender &&
		invite.User == other.User &&
		invite.SenderMembership == other.SenderMembership &&
		invite.Status == other.Status
}

type Invites []Invite

// Equals returns true iff this slice contains the same data of the
// other one and in the same order
// TODO evalute remove Equals: only for tests
func (slice Invites) Equals(other Invites) bool {
	if len(slice) != len(other) {
		return false
	}

	for index, invite := range slice {
		if !invite.Equals(other[index]) {
			return false
		}
	}

	return true
}

// ValidateBasic returns error if Invite status is not Pending, Reward or Invalid
// TODO validate basic not used
func (invite Invite) ValidateBasic() error {
	switch invite.Status {
	case uint64(InviteStatusPending), uint64(InviteStatusRewarded), uint64(InviteStatusInvalid):
		return nil
	default:
		return fmt.Errorf("invite has invalid status: %d", invite.Status)
	}
}
