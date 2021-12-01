package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgBuyMembership{}

// ------------------
// MsgBuyMembership
// ------------------

// TODO change function passing parameters
//func NewMsgBuyMembership(membership Membership) *MsgBuyMembership {
func NewMsgBuyMembership(membershipType string, buyer sdk.AccAddress, tsp sdk.AccAddress) *MsgBuyMembership {

	return &MsgBuyMembership{
		MembershipType: membershipType,
		Buyer:          buyer.String(),
		Tsp:            tsp.String(),
	}
}

func (msg *MsgBuyMembership) Route() string {
	return ModuleName
}

func (msg *MsgBuyMembership) Type() string {
	return MsgTypeBuyMembership
}

func (msg *MsgBuyMembership) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Tsp)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBuyMembership) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBuyMembership) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid buyer address: %s", msg.Buyer))
	}

	_, err = sdk.AccAddressFromBech32(msg.Tsp)
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid tsp address: %s", msg.Tsp))
	}

	membershipType := strings.TrimSpace(msg.MembershipType)
	if !IsMembershipTypeValid(membershipType) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid membership type: %s", msg.MembershipType))
	}

	return nil
}

// ------------------
// MsgInviteUser
// ------------------

func NewMsgInviteUser(sender string, recipient string) *MsgInviteUser {
	return &MsgInviteUser{
		Sender:    sender,
		Recipient: recipient,
	}
}

func (msg *MsgInviteUser) Route() string {
	return ModuleName
}

func (msg *MsgInviteUser) Type() string {
	return MsgTypeInviteUser
}

func (msg *MsgInviteUser) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgInviteUser) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic Implements Msg Validate Basic
func (msg *MsgInviteUser) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid recipient address: %s (%s)", msg.Recipient, err))
	}
	_, err = sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid sender address: %s (%s)", msg.Sender, err))
	}
	return nil
}

// ------------------
// MsgAddTsp
// ------------------

func NewMsgAddTsp(tsp string, government string) *MsgAddTsp {
	return &MsgAddTsp{
		Tsp:        tsp,
		Government: government,
	}
}

func (msg *MsgAddTsp) Route() string {
	return ModuleName
}

func (msg *MsgAddTsp) Type() string {
	return MsgTypeAddTsp
}

func (msg *MsgAddTsp) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Government)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddTsp) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddTsp) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Tsp)
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid TSP address: %s", msg.Tsp))
	}

	_, err = sdk.AccAddressFromBech32(msg.Government)
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid government address: %s", msg.Government))
	}
	return nil
}

// ------------------
// MsgDepositIntoLiquidityPool
// ------------------

func NewMsgDepositIntoLiquidityPool(amount sdk.Coins, depositor string) *MsgDepositIntoLiquidityPool {
	return &MsgDepositIntoLiquidityPool{
		Amount:    amount,
		Depositor: depositor,
	}
}

func (msg *MsgDepositIntoLiquidityPool) Route() string {
	return ModuleName
}

func (msg *MsgDepositIntoLiquidityPool) Type() string {
	return MsgTypesDepositIntoLiquidityPool
}

func (msg *MsgDepositIntoLiquidityPool) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDepositIntoLiquidityPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDepositIntoLiquidityPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid depositor address: %s", msg.Depositor))
	}

	if msg.Amount.Empty() || !msg.Amount.IsValid() {
		return sdkErr.Wrap(sdkErr.ErrInvalidCoins, fmt.Sprintf("Invalid deposit amount: %s", msg.Amount))
	}
	return nil
}

// ------------------
// MsgRemoveTsp
// ------------------

func NewMsgRemoveTsp(tsp string, government string) *MsgRemoveTsp {
	return &MsgRemoveTsp{
		Tsp:        tsp,
		Government: government,
	}
}

func (msg *MsgRemoveTsp) Route() string {
	return ModuleName
}

func (msg *MsgRemoveTsp) Type() string {
	return MsgTypeRemoveTsp
}

func (msg *MsgRemoveTsp) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Government)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveTsp) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveTsp) ValidateBasic() error {
	if msg.Tsp == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid TSP address: %s", msg.Tsp))
	}
	if msg.Government == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid government address: %s", msg.Government))
	}
	return nil
}

// ------------------
// MsgRemoveMembership
// ------------------

func NewMsgRemoveMembership(government string, subscriber string) *MsgRemoveMembership {
	return &MsgRemoveMembership{
		Government: government,
		Subscriber: subscriber,
	}
}

func (msg *MsgRemoveMembership) Route() string {
	return ModuleName
}

func (msg *MsgRemoveMembership) Type() string {
	return MsgTypeRemoveMembership
}

func (msg *MsgRemoveMembership) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Government)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveMembership) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemoveMembership) ValidateBasic() error {
	if msg.Subscriber == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid subscriber address: %s", msg.Subscriber))
	}

	if msg.Government == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid government address: %s", msg.Government))
	}
	return nil
}

// ------------------
// MsgSetMembership
// ------------------

func NewMsgSetMembership(subscriber string, government string, newMembership string) *MsgSetMembership {

	return &MsgSetMembership{
		Subscriber:    subscriber,
		Government:    government,
		NewMembership: newMembership,
	}
}

func (msg *MsgSetMembership) Route() string {
	return ModuleName
}

func (msg *MsgSetMembership) Type() string {
	return MsgTypeSetMembership
}

func (msg *MsgSetMembership) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Government)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetMembership) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetMembership) ValidateBasic() error {
	if msg.Subscriber == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid subscriber address: %s", msg.Subscriber))
	}

	if msg.Government == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("Invalid government address: %s", msg.Government))
	}

	if msg.NewMembership == "" {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized, "new membership must not be empty")
	}

	return nil
}
