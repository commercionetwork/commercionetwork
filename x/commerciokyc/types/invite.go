package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// --------------
// --- Invite
// --------------

type InviteStatus uint8

const (
	InviteStatusPending InviteStatus = iota
	InviteStatusRewarded
	InviteStatusInvalid
)

// Invite represents an invitation that a user has made towards another user
type Invite struct {
	Sender           sdk.AccAddress `json:"sender"`            // User that has sent the invitation
	SenderMembership string         `json:"sender_membership"` // Membership of Sender when the invite was created
	User             sdk.AccAddress `json:"user"`              // Invited user
	Status           InviteStatus   `json:"status"`            // Tells if the invite is pending, rewarded or invalid
}

// NewInvite creates a new invite object representing an invitation from the sender to the specified user.
// By default, NewInvite returns a Pending invite.
func NewInvite(sender, user sdk.AccAddress, senderMembership string) Invite {
	return Invite{
		Sender:           sender,
		User:             user,
		SenderMembership: senderMembership,
		Status:           InviteStatusPending,
	}
}

// Empty returns true of the given invite is empty
func (invite Invite) Empty() bool {
	return Invite{}.Equals(invite)
}

// Equals returns true iff invite contains the same data of the other invite
func (invite Invite) Equals(other Invite) bool {
	return invite.Sender.Equals(other.Sender) &&
		invite.User.Equals(other.User) &&
		invite.SenderMembership == other.SenderMembership &&
		invite.Status == other.Status
}

// --------------
// --- Invites
// --------------

// Invites represents a slice of Invite objects
type Invites []Invite

// Equals returns true iff this slice contains the same data of the
// other one and in the same order
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
func (invite Invite) ValidateBasic() error {
	switch invite.Status {
	case InviteStatusPending, InviteStatusRewarded, InviteStatusInvalid:
		return nil
	default:
		return fmt.Errorf("invite has invalid status: %d", invite.Status)
	}
}
