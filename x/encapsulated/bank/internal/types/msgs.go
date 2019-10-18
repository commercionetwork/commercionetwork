package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ---------------------------
// --- MsgBlockAccountSend
// ---------------------------

type MsgBlockAccountSend struct {
	Address sdk.AccAddress `json:"address"`
	Signer  sdk.AccAddress `json:"signer"`
}

// Route Implements Msg.
func (msg MsgBlockAccountSend) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgBlockAccountSend) Type() string { return MsgTypeBlockAccountSend }

// ValidateBasic Implements Msg.
func (msg MsgBlockAccountSend) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBlockAccountSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgBlockAccountSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// ---------------------------
// --- MsgUnlockAccountSend
// ---------------------------

type MsgUnlockAccountSend struct {
	Address sdk.AccAddress `json:"address"`
	Signer  sdk.AccAddress `json:"signer"`
}

// Route Implements Msg.
func (msg MsgUnlockAccountSend) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgUnlockAccountSend) Type() string { return MsgTypeUnlockAccountSend }

// ValidateBasic Implements Msg.
func (msg MsgUnlockAccountSend) ValidateBasic() sdk.Error {
	if msg.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Address.String())
	}
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgUnlockAccountSend) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgUnlockAccountSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
