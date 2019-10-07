package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgOpenCdp struct {
	Request CdpRequest `json:"cdp_request"`
}

func NewMsgOpenCdp(request CdpRequest) MsgOpenCdp {
	return MsgOpenCdp{
		Request: request,
	}
}

// Route Implements Msg.
func (msg MsgOpenCdp) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgOpenCdp) Type() string { return MsgTypeOpenCdp }

func (msg MsgOpenCdp) ValidateBasic() sdk.Error {
	if msg.Request.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Request.Signer.String())
	}
	if msg.Request.DepositedAmount.Empty() || msg.Request.DepositedAmount.IsAnyNegative() {
		return sdk.ErrInvalidCoins(msg.Request.DepositedAmount.String())
	}
	if len(strings.TrimSpace(msg.Request.Timestamp)) == 0 {
		return sdk.ErrUnknownRequest("cdp request's timestamp can't be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgOpenCdp) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgOpenCdp) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Request.Signer}
}

///////////////////
///MsgCloseCdp////
/////////////////
type MsgCloseCdp struct {
	Signer    sdk.AccAddress `json:"signer"`
	Timestamp string         `json:"timestamp"`
}

func NewMsgCloseCdp(signer sdk.AccAddress, timestamp string) MsgCloseCdp {
	return MsgCloseCdp{
		Signer:    signer,
		Timestamp: timestamp,
	}
}

// Route Implements Msg.
func (msg MsgCloseCdp) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCloseCdp) Type() string { return MsgTypeCloseCdp }

func (msg MsgCloseCdp) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if len(strings.TrimSpace(msg.Timestamp)) == 0 {
		return sdk.ErrUnknownRequest("cdp's timestamp can't be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCloseCdp) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgCloseCdp) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
