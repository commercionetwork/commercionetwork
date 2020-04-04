package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// -----------------
// --- MsgOpenCdp
// -----------------

type MsgOpenCdp struct {
	Owner   sdk.AccAddress `json:"owner"`
	Deposit sdk.Coin       `json:"deposit"`
}

func NewMsgOpenCdp(owner sdk.AccAddress, deposit sdk.Coin) MsgOpenCdp {
	return MsgOpenCdp{
		Deposit: deposit,
		Owner:   owner,
	}
}

// Route Implements Msg.
func (msg MsgOpenCdp) Route() string                { return RouterKey }
func (msg MsgOpenCdp) Type() string                 { return MsgTypeOpenCdp }
func (msg MsgOpenCdp) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }
func (msg MsgOpenCdp) GetSignBytes() []byte         { return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)) }
func (msg MsgOpenCdp) ValidateBasic() error {
	if msg.Owner.Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Owner.String())
	}
	if !msg.Deposit.IsValid() || !msg.Deposit.IsPositive() {
		return errors.Wrap(errors.ErrInvalidCoins, msg.Deposit.String())
	}
	return nil
}

type MsgCloseCdp struct {
	Signer  sdk.AccAddress `json:"signer"`
	Created int64          `json:"cdp_created_at"` // Block height at which the CDP has been created
}

func NewMsgCloseCdp(signer sdk.AccAddress, timestamp int64) MsgCloseCdp {
	return MsgCloseCdp{
		Signer:  signer,
		Created: timestamp,
	}
}

// Route Implements Msg.
func (msg MsgCloseCdp) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCloseCdp) Type() string { return MsgTypeCloseCdp }

func (msg MsgCloseCdp) ValidateBasic() error {
	if msg.Signer.Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Signer.String())
	}
	if msg.Created == 0 {
		return errors.Wrap(errors.ErrInvalidCoins, "CDP timestamp is invalid")
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

// -------------------
// --- MsgSetCdpCollateralRate
// -------------------

type MsgSetCdpCollateralRate struct {
	Signer            sdk.AccAddress `json:"signer"`
	CdpCollateralRate sdk.Dec        `json:"cdp_collateral_rate"`
}

func NewMsgSetCdpCollateralRate(signer sdk.AccAddress, cdpCollateralRate sdk.Dec) MsgSetCdpCollateralRate {
	return MsgSetCdpCollateralRate{Signer: signer, CdpCollateralRate: cdpCollateralRate}
}

func (MsgSetCdpCollateralRate) Route() string                    { return RouterKey }
func (MsgSetCdpCollateralRate) Type() string                     { return MsgTypeSetCdpCollateralRate }
func (msg MsgSetCdpCollateralRate) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Signer} }
func (msg MsgSetCdpCollateralRate) ValidateBasic() error {
	if msg.Signer.Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Signer.String())
	}
	return ValidateCollateralRate(msg.CdpCollateralRate)
}

func (msg MsgSetCdpCollateralRate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func ValidateCollateralRate(rate sdk.Dec) error {
	if rate.IsNil() {
		return fmt.Errorf("cdp collateral rate must be not nil")
	}
	if !rate.IsPositive() {
		return fmt.Errorf("cdp collateral rate must be positive: %s", rate)
	}
	return nil
}
