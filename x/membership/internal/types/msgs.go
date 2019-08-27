package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgAssignMembership struct {
	Signer         sdk.AccAddress // Represents the user that must sign the transaction. Should be a trusted user
	User           sdk.AccAddress // Represents the user that will receive the membership token
	MembershipType string         // Represents the type of the membership to be given
}

func NewMsgAssignMembership(signer sdk.AccAddress, user sdk.AccAddress, membershipType string) MsgAssignMembership {
	return MsgAssignMembership{
		Signer:         signer,
		User:           user,
		MembershipType: membershipType,
	}
}

// Route Implements Msg.
func (msg MsgAssignMembership) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgAssignMembership) Type() string { return MsgTypeAssignMembership }

// ValidateBasic Implements Msg.
func (msg MsgAssignMembership) ValidateBasic() sdk.Error {
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
