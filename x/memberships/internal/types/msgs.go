package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgBuyMembership allows a user to buy a membership.
// In order to be able to perform such an action, the following requirements
// should be met:
// 1. The user has been invited from a member already having a membership
// 2. The user has been verified from a TSP
// 3. The user has enough stable credits in his wallet
type MsgBuyMembership struct {
	MembershipType string         `json:"membership_type"` // Membership type to be bought
	Buyer          sdk.AccAddress `json:"buyer"`           // Buyer address
}

func NewMsgBuyMembership(membershipType string, buyer sdk.AccAddress) MsgBuyMembership {
	return MsgBuyMembership{
		MembershipType: membershipType,
		Buyer:          buyer,
	}
}

// Route Implements Msg.
func (msg MsgBuyMembership) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgBuyMembership) Type() string { return MsgTypeBuyMembership }

// ValidateBasic Implements Msg.
func (msg MsgBuyMembership) ValidateBasic() sdk.Error {
	if msg.Buyer.Empty() {
		return sdk.ErrInvalidAddress(msg.Buyer.String())
	}

	membershipType := strings.TrimSpace(msg.MembershipType)
	if len(membershipType) == 0 {
		return sdk.ErrUnknownRequest("Did Document reference cannot be empty")
	}
	if !IsMembershipTypeValid(membershipType) {
		return sdk.ErrUnknownRequest("Invalid membership type")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBuyMembership) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgBuyMembership) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Buyer}
}
