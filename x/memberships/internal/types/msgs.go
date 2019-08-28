package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgAssignMembership struct {
	Signer         sdk.AccAddress `json:"signer"`          // Represents the user that must sign the transaction. Should be a trusted user
	User           sdk.AccAddress `json:"user"`            // Represents the user that will receive the membership token
	MembershipType string         `json:"membership_type"` // Represents the type of the membership to be given
}

func NewMsgAssignMembership(signer sdk.AccAddress, user sdk.AccAddress, membershipType string) MsgAssignMembership {
	return MsgAssignMembership{
		Signer:         signer,
		User:           user,
		MembershipType: membershipType,
	}
}

// Route Implements Msg.
func (msg MsgAssignMembership) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgAssignMembership) Type() string { return "assign_membership" }

// ValidateBasic Implements Msg.
func (msg MsgAssignMembership) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if msg.User.Empty() {
		return sdk.ErrInvalidAddress(msg.User.String())
	}
	if len(msg.MembershipType) == 0 {
		return sdk.ErrUnknownRequest("Did Document reference cannot be empty")
	}
	if !IsMembershipTypeValid(msg.MembershipType) {
		return sdk.ErrUnknownRequest("Invalid membership type")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgAssignMembership) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgAssignMembership) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
