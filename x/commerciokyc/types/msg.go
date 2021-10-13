package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgBuyMembership{}

// MsgBuyMembership

func NewMsgBuyMembership(membership Membership) *MsgBuyMembership {

	return &MsgBuyMembership{
		MembershipType: membership.MembershipType,
		Buyer:          membership.Owner,
		Tsp:            membership.TspAddress,
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

	return nil
}

// MsgInviteUser

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

func (msg *MsgInviteUser) ValidateBasic() error {

	return nil
}

// MsgAddTsp

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

	return nil
}

// MsgDepositIntoLiquidityPool

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
	return MsgTypeAddTsp
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

	return nil
}

// MsgRemoveTsp

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
	return MsgTypeAddTsp
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

	return nil
}

// MsgRemoveMembership

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
	return MsgTypeAddTsp
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

	return nil
}

// MsgSetMembership

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
	return MsgTypeAddTsp
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

	return nil
}
