package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ----------------------------------
// --- MsgSetGovAddress
// ----------------------------------
var _ sdk.Msg = &MsgSetGovAddress{}

func NewMsgSetGovAddress() *MsgSetGovAddress {
	return &MsgSetGovAddress{}
}

const SetGovAddressConst = "SetGovAddress"

func (msg *MsgSetGovAddress) Route() string {
	return RouterKey
}

func (msg *MsgSetGovAddress) Type() string {
	return SetGovAddressConst
}

func (msg *MsgSetGovAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgSetGovAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetGovAddress) ValidateBasic() error {
	return nil
}

// ----------------------------------
// --- MsgFixSupply
// ----------------------------------
var _ sdk.Msg = &MsgSetGovAddress{}

func NewMsgFixSupplys(signer sdk.AccAddress, amount sdk.Coin) *MsgFixSupply {
	return &MsgFixSupply{
		Sender: signer.String(),
		Amount: &amount,
	}
}

const FixSupplyConst = "FixSupply"

func (msg *MsgFixSupply) Route() string {
	return RouterKey
}

func (msg *MsgFixSupply) Type() string {
	return FixSupplyConst
}

func (msg *MsgFixSupply) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgFixSupply) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgFixSupply) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Sender)
	}

	coin := msg.Amount

	if !coin.IsValid() {
		return errors.Wrap(errors.ErrInvalidCoins, msg.Amount.String())
	}

	if !coin.IsPositive() {
		return errors.Wrap(errors.ErrInvalidCoins, msg.Amount.String())
	}

	return nil
}
