package types

import (
	"fmt"
	"strings"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

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

// NewMsgInviteUser create MsgInviteUser
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
func (msg MsgInviteUser) ValidateBasic() error {
	if msg.Recipient.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid recipient address: %s", msg.Recipient))
	}
	if msg.Sender.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid sender address: %s", msg.Sender))
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
func (msg MsgDepositIntoLiquidityPool) ValidateBasic() error {
	if msg.Depositor.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid depositor address: %s", msg.Depositor))
	}
	if msg.Amount.Empty() || !msg.Amount.IsValid() {
		return sdkErr.Wrap(sdkErr.ErrInvalidCoins, fmt.Sprintf("Invalid deposit amount: %s", msg.Amount))
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
func (msg MsgAddTsp) ValidateBasic() error {
	if msg.Tsp.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid TSP address: %s", msg.Tsp))
	}
	if msg.Government.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid government address: %s", msg.Government))
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

// --------------------------------
// --- MsgRemoveTsp
// --------------------------------

// MsgRemoveTsp should be used when wanting to remove a address
// from list of  Trusted Service Providers (TSPs).
type MsgRemoveTsp struct {
	Tsp        sdk.AccAddress `json:"tsp"`
	Government sdk.AccAddress `json:"government"`
}

// NewMsgRemoveTsp return MsgRemoveTsp
func NewMsgRemoveTsp(tsp, government sdk.AccAddress) MsgRemoveTsp {
	return MsgRemoveTsp{Tsp: tsp, Government: government}
}

// Route Implements Msg.
func (msg MsgRemoveTsp) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRemoveTsp) Type() string { return MsgTypeRemoveTsp }

// ValidateBasic Implements Msg.
func (msg MsgRemoveTsp) ValidateBasic() error {
	if msg.Tsp.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid TSP address: %s", msg.Tsp))
	}
	if msg.Government.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid government address: %s", msg.Government))
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRemoveTsp) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgRemoveTsp) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Government}
}

// --------------------------------
// --- MsgBuyMembership
// --------------------------------

// MsgBuyMembership allows a user to buy a membership.
// In order to be able to perform such an action, the following requirements
// should be met:
// 1. The user has been invited from a member already having a membership
// 2. The user has been verified from a TSP
// 3. The user has enough stable credits in his wallet
type MsgBuyMembership struct {
	MembershipType string         `json:"membership_type"` // Membership type to be bought
	Buyer          sdk.AccAddress `json:"buyer"`           // Buyer address
	Tsp            sdk.AccAddress `json:"tsp"`             // Buyer address
}

func NewMsgBuyMembership(membershipType string, buyer sdk.AccAddress, tsp sdk.AccAddress) MsgBuyMembership {
	return MsgBuyMembership{
		MembershipType: membershipType,
		Buyer:          buyer,
		Tsp:            tsp,
	}
}

// Route Implements Msg.
func (msg MsgBuyMembership) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgBuyMembership) Type() string { return MsgTypeBuyMembership }

// ValidateBasic Implements Msg.
func (msg MsgBuyMembership) ValidateBasic() error {
	if msg.Buyer.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid buyer address: %s", msg.Buyer))
	}

	if msg.Tsp.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid tsp address: %s", msg.Tsp))
	}

	membershipType := strings.TrimSpace(msg.MembershipType)
	if !IsMembershipTypeValid(membershipType) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid membership type: %s", msg.MembershipType))
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBuyMembership) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgBuyMembership) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Tsp}
}

// --------------------------------
// --- MsgRemoveMembership
// --------------------------------

type MsgRemoveMembership struct {
	GovernmentAddress sdk.AccAddress `json:"government_address"` // Buyer address
	Subscriber        sdk.AccAddress `json:"subscriber"`         // Buyer address
}

func NewMsgRemoveMembership(govAddr sdk.AccAddress, subscriber sdk.AccAddress) MsgRemoveMembership {
	return MsgRemoveMembership{
		GovernmentAddress: govAddr,
		Subscriber:        subscriber,
	}
}

// Route Implements Msg.
func (msg MsgRemoveMembership) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgRemoveMembership) Type() string { return MsgTypeRemoveMembership }

// ValidateBasic Implements Msg.
func (msg MsgRemoveMembership) ValidateBasic() error {
	if msg.Subscriber.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid subscriber address: %s", msg.Subscriber))
	}

	if msg.GovernmentAddress.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid government address: %s", msg.GovernmentAddress))
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgRemoveMembership) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgRemoveMembership) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.GovernmentAddress}
}

// ---------------------------
// --- MsgSetMembership
// ---------------------------

// MsgSetMembership allows government to assign a membership to Subscriber,
// which has been already invited by another black membership user.
type MsgSetMembership struct {
	GovernmentAddress sdk.AccAddress `json:"government_address"`
	Subscriber        sdk.AccAddress `json:"subscriber"`
	NewMembership     string         `json:"new_membership"`
}

// NewMsgSetMembership create MsgSetMembership
func NewMsgSetMembership(subscriber sdk.AccAddress, govAddr sdk.AccAddress, newMembership string) MsgSetMembership {
	return MsgSetMembership{
		Subscriber:        subscriber,
		GovernmentAddress: govAddr,
		NewMembership:     newMembership,
	}
}

// Route Implements Msg.
func (msg MsgSetMembership) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSetMembership) Type() string { return MsgTypeSetMembership }

// ValidateBasic Implements Msg.
func (msg MsgSetMembership) ValidateBasic() error {
	if msg.Subscriber.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid subscriber address: %s", msg.Subscriber))
	}

	if msg.GovernmentAddress.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid government address: %s", msg.GovernmentAddress))
	}

	if msg.NewMembership == "" {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized, "new membership must not be empty")
	}

	/*membershipType := strings.TrimSpace(msg.NewMembership)
	if !IsMembershipTypeValid(membershipType) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid membership type: %s", msg.NewMembership))
	}*/

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetMembership) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSetMembership) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.GovernmentAddress}
}
