package types

import (
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
