package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// --------------
// --- Invite
// --------------

// Invite represents an invitation that a user has made towards another user
type Invite struct {
	Sender   sdk.AccAddress `json:"sender"`   // User that has sent the invitation
	User     sdk.AccAddress `json:"user"`     // Invited user
	Rewarded bool           `json:"rewarded"` // Tells if the invitee has already been rewarded
}

// NewInvite creates a new invite object representing an invitation from the sender to the specified user
func NewInvite(sender, user sdk.AccAddress) Invite {
	return Invite{
		Sender:   sender,
		User:     user,
		Rewarded: false,
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
		invite.Rewarded == other.Rewarded
}

// --------------
// --- Invites
// --------------

type Invites []Invite
