package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// --------------------------
// --- MsgInviteUser
// --------------------------

// MsgInviteUser allows to properly invite a user.
// Te invitation system should be a one-invite-only system, where invites
// consecutive to the first one should be discarded.
type MsgInviteUser struct {
	Recipient sdk.AccAddress `json:"recipient"`
	Sender    sdk.AccAddress `json:"sender"`
}

func NewMsgInviteUser(sender, recipient sdk.AccAddress) MsgInviteUser {
	return MsgInviteUser{
		Recipient: recipient,
		Sender:    sender,
	}
}

// Route Implements Msg.
func (msg MsgInviteUser) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgInviteUser) Type() string { return MsgTypeInviteUser }

// ValidateBasic Implements Msg.
func (msg MsgInviteUser) ValidateBasic() sdk.Error {
	if msg.Recipient.Empty() {
		return sdk.ErrInvalidAddress(msg.Recipient.String())
	}
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress(msg.Sender.String())
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgInviteUser) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgInviteUser) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// --------------------------
// --- MsgSetUserVerified
// --------------------------

// MsgSetUserVerified is used to set a specific user as properly verified.
// Note that the verifier address should identify a Trusted Service Provider account.
type MsgSetUserVerified struct {
	User     sdk.AccAddress `json:"user"`
	Verifier sdk.AccAddress `json:"verifier"`
}

func NewMsgSetUserVerified(user, verifier sdk.AccAddress) MsgSetUserVerified {
	return MsgSetUserVerified{
		User:     user,
		Verifier: verifier,
	}
}

// Route Implements Msg.
func (msg MsgSetUserVerified) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSetUserVerified) Type() string { return MsgTypeSetUserVerified }

// ValidateBasic Implements Msg.
func (msg MsgSetUserVerified) ValidateBasic() sdk.Error {
	if msg.User.Empty() {
		return sdk.ErrInvalidAddress(msg.User.String())
	}
	if msg.Verifier.Empty() {
		return sdk.ErrInvalidAddress(msg.Verifier.String())
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetUserVerified) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSetUserVerified) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Verifier}
}

// --------------------------------
// --- MsgDepositIntoLiquidityPool
// --------------------------------

// MsgDepositIntoLiquidityPool should be used when wanting to deposit a specific
// amount into the liquidity pool which contains the total amount of rewards to
// be distributed upon an accreditation
type MsgDepositIntoLiquidityPool struct {
	Depositor sdk.AccAddress `json:"depositor"`
	Amount    sdk.Coins      `json:"amount"`
}

func NewMsgDepositIntoLiquidityPool(amount sdk.Coins, depositor sdk.AccAddress) MsgDepositIntoLiquidityPool {
	return MsgDepositIntoLiquidityPool{
		Depositor: depositor,
		Amount:    amount,
	}
}

// Route Implements Msg.
func (msg MsgDepositIntoLiquidityPool) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDepositIntoLiquidityPool) Type() string { return MsgTypesDepositIntoLiquidityPool }

// ValidateBasic Implements Msg.
func (msg MsgDepositIntoLiquidityPool) ValidateBasic() sdk.Error {
	if msg.Depositor.Empty() {
		return sdk.ErrInvalidAddress(msg.Depositor.String())
	}
	if msg.Amount.Empty() {
		return sdk.ErrUnknownRequest("amount cannot be empty")
	}
	if msg.Amount.IsAnyNegative() {
		return sdk.ErrUnknownRequest("amount cannot be negative")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDepositIntoLiquidityPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgDepositIntoLiquidityPool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Depositor}
}

// --------------------------------
// --- MsgAddTsp
// --------------------------------

// MsgAddTsp should be used when wanting to add a new address
// as a Trusted Service Provider (TSP).
// TSPs will be able to sign rewarding-giving transactions
// so should be only a handful of very trusted accounts.
type MsgAddTsp struct {
	Tsp        sdk.AccAddress `json:"tsp"`
	Government sdk.AccAddress `json:"government"`
}

func NewMsgAddTsp(tsp, government sdk.AccAddress) MsgAddTsp {
	return MsgAddTsp{Tsp: tsp, Government: government}
}

// Route Implements Msg.
func (msg MsgAddTsp) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgAddTsp) Type() string { return MsgTypeAddTsp }

// ValidateBasic Implements Msg.
func (msg MsgAddTsp) ValidateBasic() sdk.Error {
	if msg.Tsp.Empty() {
		return sdk.ErrInvalidAddress(msg.Tsp.String())
	}
	if msg.Government.Empty() {
		return sdk.ErrInvalidAddress(msg.Government.String())
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgAddTsp) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgAddTsp) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Government}
}

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
