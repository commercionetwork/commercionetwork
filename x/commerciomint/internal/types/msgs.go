package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// -----------------
// --- MsgOpenCdp
// -----------------

type MsgOpenCdp struct {
	Depositor       sdk.AccAddress `json:"depositor"`
	DepositedAmount sdk.Coins      `json:"deposit_amount"`
}

func NewMsgOpenCdp(depositAmount sdk.Coins, depositor sdk.AccAddress) MsgOpenCdp {
	return MsgOpenCdp{
		DepositedAmount: depositAmount,
		Depositor:       depositor,
	}
}

// Route Implements Msg.
func (msg MsgOpenCdp) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgOpenCdp) Type() string { return MsgTypeOpenCdp }

func (msg MsgOpenCdp) ValidateBasic() sdk.Error {
	if msg.Depositor.Empty() {
		return sdk.ErrInvalidAddress(msg.Depositor.String())
	}
	if msg.DepositedAmount.Empty() || msg.DepositedAmount.IsAnyNegative() {
		return sdk.ErrInvalidCoins(msg.DepositedAmount.String())
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgOpenCdp) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgOpenCdp) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Depositor}
}

///////////////////
///MsgCloseCdp////
/////////////////

type MsgCloseCdp struct {
	Signer    sdk.AccAddress `json:"signer"`
	Timestamp int64          `json:"cdp_timestamp"` // Block height at which the CDP has been created
}

func NewMsgCloseCdp(signer sdk.AccAddress, timestamp int64) MsgCloseCdp {
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
	if msg.Timestamp == 0 {
		return sdk.ErrUnknownRequest("CDP timestamp is invalid")
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
