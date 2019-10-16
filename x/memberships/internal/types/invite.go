package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Invite struct {
	Sender   sdk.AccAddress `json:"sender"`   // User that has sent the invitation
	User     sdk.AccAddress `json:"user"`     // Invited user
	Rewarded bool           `json:"rewarded"` // Tells if the invitee has already been rewarded
}
