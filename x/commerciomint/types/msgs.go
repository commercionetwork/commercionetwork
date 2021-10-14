package types

import (
	"time"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
	uuid "github.com/satori/go.uuid"
)

var _ sdk.Msg = &MsgMintCCC{}

// MsgMintCCC

func NewMsgMintCCC(position Position) *MsgMintCCC {
	var depositAmount []*sdk.Coin
	coin := sdk.NewInt64Coin(CreditsDenom, position.Collateral)
	depositAmount = append(depositAmount, &coin)

	return &MsgMintCCC{
		Depositor:     position.Owner,
		DepositAmount: depositAmount,
		ID:            position.ID,
	}
}

func (msg *MsgMintCCC) Route() string {
	return ModuleName
}

func (msg *MsgMintCCC) Type() string {
	return MsgTypeMintCCC
}

func (msg *MsgMintCCC) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMintCCC) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMintCCC) ValidateBasic() error {
	if sdk.AccAddress(msg.Depositor).Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Depositor)
	}

	if msg.ID == "" {
		return errors.Wrap(errors.ErrInvalidRequest, "missing position ID")
	}

	coins := sdk.NewCoins()
	for _, coin := range msg.DepositAmount {
		coins = append(coins, *coin)
	}
	sdk.NewCoins()
	if !ValidateDeposit(coins) {
		return errors.Wrap(errors.ErrInvalidCoins, coins.String())
	}

	return nil
}

// MsgBurnCCC

// TODO REVIEW MESSAGES CREATOR
func NewMsgBurnCCC(signer sdk.AccAddress, id string, amount sdk.Coin) *MsgBurnCCC {
	return &MsgBurnCCC{
		Signer: signer.String(),
		Amount: &amount,
		ID:     id,
	}
}

func (msg *MsgBurnCCC) Route() string {
	return ModuleName
}

func (msg *MsgBurnCCC) Type() string {
	return MsgTypeBurnCCC
}

func (msg *MsgBurnCCC) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBurnCCC) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBurnCCC) ValidateBasic() error {
	if sdk.AccAddress(msg.Signer).Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Signer)
	}

	if msg.Amount.IsZero() || msg.Amount.IsNegative() || msg.Amount.Denom != CreditsDenom {
		return errors.Wrap(errors.ErrInvalidRequest, "invalid amount")
	}

	if _, err := uuid.FromString(msg.ID); err != nil {
		return errors.Wrap(errors.ErrInvalidRequest, "id must be a well-defined UUID")
	}
	return nil
}

// NewMsgSetCCCConversionRate

// TODO REVIEW MESSAGES CREATOR
func NewMsgSetCCCConversionRate(signer sdk.AccAddress, rate types.DecProto) *MsgSetCCCConversionRate {
	return &MsgSetCCCConversionRate{
		Signer: signer.String(),
		Rate:   &rate,
	}
}

func (msg *MsgSetCCCConversionRate) Route() string {
	return ModuleName
}

func (msg *MsgSetCCCConversionRate) Type() string {
	return MsgTypeSetCCCConversionRate
}

func (msg *MsgSetCCCConversionRate) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetCCCConversionRate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// TODO remove duplicate validation
func (msg *MsgSetCCCConversionRate) ValidateBasic() error {
	if sdk.AccAddress(msg.Signer).Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Signer)
	}
	return ValidateConversionRate(msg.Rate.Dec)
}

// NewMsgSetCCCFreezePeriod

// TODO REVIEW MESSAGES CREATOR
func NewMsgSetCCCFreezePeriod(signer sdk.AccAddress, freezePeriod string) *MsgSetCCCFreezePeriod {
	return &MsgSetCCCFreezePeriod{
		Signer:       signer.String(),
		FreezePeriod: freezePeriod,
	}
}

func (msg *MsgSetCCCFreezePeriod) Route() string {
	return ModuleName
}

func (msg *MsgSetCCCFreezePeriod) Type() string {
	return MsgTypeSetCCCFreezePeriod
}

func (msg *MsgSetCCCFreezePeriod) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetCCCFreezePeriod) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetCCCFreezePeriod) ValidateBasic() error {
	if sdk.AccAddress(msg.Signer).Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Signer)
	}
	// TODO move all control into method ValidateFreezePeriod
	freezePeriod, err := time.ParseDuration(msg.FreezePeriod)
	if err != nil {
		return errors.Wrap(errors.ErrInvalidRequest, err.Error())
	}
	return ValidateFreezePeriod(freezePeriod)
}
