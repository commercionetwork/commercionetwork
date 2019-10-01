package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgOpenCDP struct {
	Request CDPRequest `json:"cdp_request"`
}

func NewMsgOpenCDP(request CDPRequest) MsgOpenCDP {
	return MsgOpenCDP{
		Request: request,
	}
}

// Route Implements Msg.
func (msg MsgOpenCDP) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgOpenCDP) Type() string { return MsgTypeOpenCDP }

func (msg MsgOpenCDP) ValidateBasic() sdk.Error {
	if msg.Request.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Request.Signer.String())
	}
	if msg.Request.DepositedAmount.Empty() || msg.Request.DepositedAmount.IsAnyNegative() {
		return sdk.ErrInvalidCoins("Deposited amount cannot be empty or negative")
	}
	if len(strings.TrimSpace(msg.Request.Timestamp)) == 0 {
		return sdk.ErrUnknownRequest("Cdp request's timestamp can't be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgOpenCDP) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgOpenCDP) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Request.Signer}
}

type MsgCloseCDP struct {
	Signer    sdk.AccAddress `json:"sender"`
	Timestamp string         `json:"timestamp"`
}

func NewMsgCloseCDP(signer sdk.AccAddress, timestamp string) MsgCloseCDP {
	return MsgCloseCDP{
		Signer:    signer,
		Timestamp: timestamp,
	}
}

// Route Implements Msg.
func (msg MsgCloseCDP) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCloseCDP) Type() string { return MsgTypeCloseCDP }

func (msg MsgCloseCDP) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if len(strings.TrimSpace(msg.Timestamp)) == 0 {
		return sdk.ErrUnknownRequest("Cdp's timestamp can't be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCloseCDP) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgCloseCDP) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
