package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	uuid "github.com/satori/go.uuid"
)

// -----------------
// --- MsgMintCCC
// -----------------

type MsgMintCCC struct {
	Owner   sdk.AccAddress `json:"depositor"`
	Credits sdk.Coins      `json:"deposit_amount"`
	ID      string         `json:"id"`
}

func NewMsgMintCCC(owner sdk.AccAddress, deposit sdk.Coins, id string) MsgMintCCC {
	return MsgMintCCC{
		Credits: deposit,
		Owner:   owner,
		ID:      id,
	}
}

// Route Implements Msg.
func (msg MsgMintCCC) Route() string                { return RouterKey }
func (msg MsgMintCCC) Type() string                 { return MsgTypeMintCCC }
func (msg MsgMintCCC) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Owner} }
func (msg MsgMintCCC) GetSignBytes() []byte         { return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg)) }
func (msg MsgMintCCC) ValidateBasic() error {
	if msg.Owner.Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Owner.String())
	}

	if msg.ID == "" {
		return errors.Wrap(errors.ErrInvalidRequest, "missing position ID")
	}

	if !ValidateDeposit(msg.Credits) {
		return errors.Wrap(errors.ErrInvalidCoins, msg.Credits.String())
	}

	return nil
}

type MsgBurnCCC struct {
	Signer sdk.AccAddress `json:"signer"`
	Amount sdk.Coin       `json:"amount"`
	ID     string         `json:"id"`
}

func NewMsgBurnCCC(signer sdk.AccAddress, id string, amount sdk.Coin) MsgBurnCCC {
	return MsgBurnCCC{
		Signer: signer,
		ID:     id,
		Amount: amount,
	}
}

// Route Implements Msg.
func (msg MsgBurnCCC) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgBurnCCC) Type() string { return MsgTypeBurnCCC }

func (msg MsgBurnCCC) ValidateBasic() error {
	if msg.Signer.Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Signer.String())
	}

	if msg.Amount.IsZero() || msg.Amount.IsNegative() || msg.Amount.Denom != CreditsDenom {
		return errors.Wrap(errors.ErrInvalidRequest, "invalid amount")
	}

	if _, err := uuid.FromString(msg.ID); err != nil {
		return errors.Wrap(errors.ErrInvalidRequest, "id must be a well-defined UUID")
	}

	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgBurnCCC) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgBurnCCC) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

// -------------------
// --- MsgSetCCCConversionRate
// -------------------

type MsgSetCCCConversionRate struct {
	Signer sdk.AccAddress `json:"signer"`
	Rate   sdk.Dec        `json:"rate"`
}

func NewMsgSetCCCConversionRate(signer sdk.AccAddress, rate sdk.Dec) MsgSetCCCConversionRate {
	return MsgSetCCCConversionRate{Signer: signer, Rate: rate}
}

func (MsgSetCCCConversionRate) Route() string                    { return RouterKey }
func (MsgSetCCCConversionRate) Type() string                     { return MsgTypeSetCCCConversionRate }
func (msg MsgSetCCCConversionRate) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Signer} }
func (msg MsgSetCCCConversionRate) ValidateBasic() error {
	if msg.Signer.Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Signer.String())
	}

	return ValidateConversionRate(msg.Rate)
}

func (msg MsgSetCCCConversionRate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func ValidateConversionRate(rate sdk.Dec) error {
	if rate.IsZero() {
		return fmt.Errorf("conversion rate cannot be zero")
	}
	if rate.IsNegative() {
		return fmt.Errorf("conversion rate must be positive")
	}
	return nil
}
