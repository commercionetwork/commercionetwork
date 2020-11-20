package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// -------------------------
// --- MsgIncrementBlockRewardsPool
// -------------------------
// MsgIncrementBlockRewardsPool - struct for increment the vbr pool
type MsgIncrementBlockRewardsPool struct {
	Funder sdk.AccAddress `json:"funder"`
	Amount sdk.Coins      `json:"amount"`
}

func NewMsgIncrementBlockRewardsPool(funder sdk.AccAddress, amount sdk.Coins) MsgIncrementBlockRewardsPool {
	return MsgIncrementBlockRewardsPool{
		Funder: funder,
		Amount: amount,
	}
}

// Route implements the sdk.Msg interface.
func (msg MsgIncrementBlockRewardsPool) Route() string { return ModuleName }

func (msg MsgIncrementBlockRewardsPool) Type() string { return MsgTypeIncrementBlockRewardsPool }

func (msg MsgIncrementBlockRewardsPool) ValidateBasic() error {
	if msg.Funder.Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Funder.String())
	}
	if msg.Amount.IsZero() || msg.Amount.IsAnyNegative() {
		return errors.Wrap(errors.ErrUnknownRequest, "You can't transfer a null or negative amount")
	}

	return nil
}

func (msg MsgIncrementBlockRewardsPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIncrementBlockRewardsPool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Funder}
}

// -------------------------
// --- MsgSetRewardRate
// -------------------------

// MsgSetRewardRate is messagge structure
type MsgSetRewardRate struct {
	Government sdk.AccAddress `json:"government"`
	RewardRate sdk.Dec        `json:"reward_rate"`
}

// NewMsgSetRewardRate return new MsgSetRewardRate
func NewMsgSetRewardRate(governement sdk.AccAddress, rewardRate sdk.Dec) MsgSetRewardRate {
	return MsgSetRewardRate{
		Government: governement,
		RewardRate: rewardRate,
	}
}

// Route Implements Msg.
func (msg MsgSetRewardRate) Route() string { return RouterKey }

// Type returns human-readable string for the message.
func (msg MsgSetRewardRate) Type() string { return MsgTypeSetRewardRate }

// ValidateBasic does a basic validation
func (msg MsgSetRewardRate) ValidateBasic() error {
	if msg.Government.Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Government.String())
	}
	if err := ValidateRewardRate(msg.RewardRate); err != nil {
		return err
	}
	return nil
}

// GetSignBytes returns the canonical byte representation of the Msg.
func (msg MsgSetRewardRate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners returns the addrs of signers that must sign.
func (msg MsgSetRewardRate) GetSigners() []sdk.AccAddress { return []sdk.AccAddress{msg.Government} }

// ValidateRewardRate validate reward rete.
func ValidateRewardRate(rate sdk.Dec) error {
	if rate.IsNil() {
		return fmt.Errorf("reward rate must be not nil")
	}
	if !rate.IsPositive() {
		return fmt.Errorf("reward rate must be positive: %s", rate)
	}
	return nil
}

// -------------------------
// --- MsgSetAutomaticWithdraw
// -------------------------

// MsgSetAutomaticWithdraw is messagge structure
type MsgSetAutomaticWithdraw struct {
	Government        sdk.AccAddress `json:"government"`
	AutomaticWithdraw bool           `json:"automatic_withdraw"`
}

// NewMsgSetAutomaticWithdraw return new MsgSetAutomaticWithdraw
func NewMsgSetAutomaticWithdraw(governement sdk.AccAddress, aWith bool) MsgSetAutomaticWithdraw {
	return MsgSetAutomaticWithdraw{
		Government:        governement,
		AutomaticWithdraw: aWith,
	}
}

// Route Implements Msg.
func (msg MsgSetAutomaticWithdraw) Route() string { return RouterKey }

// Type returns human-readable string for the message.
func (msg MsgSetAutomaticWithdraw) Type() string { return MsgTypeSetAutomaticWithdraw }

// ValidateBasic does a basic validation
func (msg MsgSetAutomaticWithdraw) ValidateBasic() error {
	if msg.Government.Empty() {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Government.String())
	}

	return nil
}

// GetSignBytes returns the canonical byte representation of the Msg.
func (msg MsgSetAutomaticWithdraw) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners returns the addrs of signers that must sign.
func (msg MsgSetAutomaticWithdraw) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Government}
}
