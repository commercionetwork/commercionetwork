package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// --------------------------
// --- MsgSetAccrediter
// --------------------------

type MsgSetAccrediter struct {
	User       sdk.AccAddress `json:"user"`
	Accrediter sdk.AccAddress `json:"accrediter"`
	Signer     sdk.AccAddress `json:"signer"`
}

func NewMsgSetAccrediter(user, accrediter, signer sdk.AccAddress) MsgSetAccrediter {
	return MsgSetAccrediter{
		User:       user,
		Accrediter: accrediter,
		Signer:     signer,
	}
}

// RouterKey Implements Msg.
func (msg MsgSetAccrediter) Route() string { return ModuleName }

// Type Implements Msg.
func (msg MsgSetAccrediter) Type() string { return MsgTypeSetAccrediter }

// ValidateBasic Implements Msg.
func (msg MsgSetAccrediter) ValidateBasic() sdk.Error {
	if msg.User.Empty() {
		return sdk.ErrInvalidAddress(msg.User.String())
	}
	if msg.Accrediter.Empty() {
		return sdk.ErrInvalidAddress(msg.Accrediter.String())
	}
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetAccrediter) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSetAccrediter) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
