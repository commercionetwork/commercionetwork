package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

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
