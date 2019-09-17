package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgDepositToken struct {
	Signer sdk.AccAddress `json:"sender"`
	Amount sdk.Coins      `json:"amount"`
}

func NewMsgDepositToken(signer sdk.AccAddress, amount sdk.Coins) MsgDepositToken {
	return MsgDepositToken{
		Signer: signer,
		Amount: amount,
	}
}

// Route Implements Msg.
func (msg MsgDepositToken) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgDepositToken) Type() string { return MsgTypeDepositToken }

func (msg MsgDepositToken) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if msg.Amount.Empty() || msg.Amount.IsAnyNegative() {
		return sdk.ErrUnknownRequest("Token's amount cannot be empty or negative")
	}

	for _, ele := range msg.Amount {
		if ele.Denom != "ucommercio" {
			return sdk.ErrUnknownRequest("Only commercio tokens can be deposited")
		}
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgDepositToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgDepositToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type MsgWithdrawToken struct {
	Signer sdk.AccAddress `json:"sender"`
	Amount sdk.Coins      `json:"amount"`
}

func NewMsgWithdrawToken(signer sdk.AccAddress, amount sdk.Coins) MsgDepositToken {
	return MsgDepositToken{
		Signer: signer,
		Amount: amount,
	}
}

// Route Implements Msg.
func (msg MsgWithdrawToken) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgWithdrawToken) Type() string { return MsgTypeWithdrawToken }

func (msg MsgWithdrawToken) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if msg.Amount.Empty() || msg.Amount.IsAnyNegative() {
		return sdk.ErrUnknownRequest("You can't withdraw an empty or negative amount")
	}

	for _, ele := range msg.Amount {
		if ele.Denom != "ucommercio" {
			return sdk.ErrUnknownRequest("Only commercio tokens can be withdrawed")
		}
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgWithdrawToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgWithdrawToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
