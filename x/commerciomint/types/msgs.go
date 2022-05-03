package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
	uuid "github.com/satori/go.uuid"
)

var _ sdk.Msg = &MsgMintCCC{}

// -------------------------
// MsgMintCCC
// -------------------------

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
	_, err := sdk.AccAddressFromBech32(msg.Depositor)
	if err != nil {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Depositor)
	}

	if msg.ID == "" {
		return errors.Wrap(errors.ErrInvalidRequest, "missing position ID")
	}

	coins := sdk.NewCoins()
	for _, coin := range msg.DepositAmount {
		coins = append(coins, *coin)
	}
	if !ValidateDeposit(coins) {
		return errors.Wrap(errors.ErrInvalidCoins, coins.String())
	}

	return nil
}

// -------------------------
// MsgBurnCCC
// -------------------------

var _ sdk.Msg = &MsgBurnCCC{}

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
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
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

// -------------------------
// --- MsgSetParams
// -------------------------

var _ sdk.Msg = &MsgSetParams{}

func NewMsgSetParams(government string, conversionRate sdk.Dec, freezePeriod time.Duration) *MsgSetParams {
	params := Params{
		ConversionRate: conversionRate,
		FreezePeriod:   freezePeriod,
	}

	return &MsgSetParams{
		Signer: government,
		Params: &params,
	}
}

func (msg *MsgSetParams) Route() string {
	return RouterKey
}

func (msg *MsgSetParams) Type() string {
	return MsgTypeSetParams
}

func (msg *MsgSetParams) GetSigners() []sdk.AccAddress {
	gov, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{gov}
}

func (msg *MsgSetParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return errors.Wrapf(errors.ErrInvalidAddress, "invalid government address (%s)", err)
	}

	err = msg.Params.Validate()
	if err != nil {
		return errors.Wrapf(errors.ErrUnknownRequest, "invalid params (%s)", err)
	}

	return nil
}
