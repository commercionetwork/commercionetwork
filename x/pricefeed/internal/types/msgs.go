package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

type MsgSetPrice struct {
	Signer    sdk.AccAddress `json:"sender"`
	Price     sdk.Int        `json:"price"`
	TokenName string         `json:"token_name"`
	Expiry    sdk.Int        `json:"expiry"`
}

func NewMsgSetPrice(signer sdk.AccAddress, price sdk.Int, tokenName string, expiry sdk.Int) MsgSetPrice {
	return MsgSetPrice{
		Signer:    signer,
		Price:     price,
		TokenName: tokenName,
		Expiry:    expiry,
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
	//
	if len(strings.TrimSpace(msg.TokenName)) == 0 {
		return sdk.ErrUnknownRequest("Cannot set price for unnamed token")
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

type MsgAddOracle struct {
	Signer sdk.AccAddress
	Oracle sdk.AccAddress
}

func NewMsgAddOracle(signer sdk.AccAddress, oracle sdk.AccAddress) MsgAddOracle {
	return MsgAddOracle{
		Signer: signer,
		Oracle: oracle,
	}
}

// Route Implements Msg.
func (msg MsgAddOracle) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgAddOracle) Type() string { return MsgTypeAddOracle }

func (msg MsgAddOracle) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if msg.Oracle.Empty() {
		return sdk.ErrInvalidAddress(msg.Oracle.String())
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgAddOracle) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgAddOracle) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
