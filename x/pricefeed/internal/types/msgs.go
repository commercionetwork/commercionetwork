package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgSetPrice struct {
	Signer    sdk.AccAddress `json:"sender"`
	Price     sdk.Dec        `json:"price"`
	TokenName string         `json:"token_name"`
}

func NewMsgSetPrice(signer sdk.AccAddress, price sdk.Dec, tokenName string) MsgSetPrice {
	return MsgSetPrice{
		Signer:    signer,
		Price:     price,
		TokenName: tokenName,
	}
}

// Route Implements Msg.
func (msg MsgSetPrice) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSetPrice) Type() string { return MsgTypeSetPrice }

func (msg MsgSetPrice) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if msg.Price.IsZero() || msg.Price.IsNegative() {
		return sdk.ErrUnknownRequest("Token's price cannot be zero or negative")
	}
	if msg.TokenName != "ucommercio" {
		return sdk.ErrUnknownRequest("Cannot set price for this token")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetPrice) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSetPrice) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
