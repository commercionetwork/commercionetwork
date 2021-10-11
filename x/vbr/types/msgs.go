package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// -------------------------
// --- MsgIncrementBlockRewardsPool
// -------------------------

var _ sdk.Msg = &MsgIncrementBlockRewardsPool{}

func NewMsgIncrementBlockRewardsPool(funder string, amount []sdk.Coin) *MsgIncrementBlockRewardsPool {
	return &MsgIncrementBlockRewardsPool{
		Funder: funder,
		Amount: amount,
	}
}

func (msg *MsgIncrementBlockRewardsPool) Route() string {
	return RouterKey
}

func (msg *MsgIncrementBlockRewardsPool) Type() string {
	return MsgTypeIncrementBlockRewardsPool
}

func (msg *MsgIncrementBlockRewardsPool) GetSigners() []sdk.AccAddress {
	funder, err := sdk.AccAddressFromBech32(msg.Funder)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{funder}
}

func (msg *MsgIncrementBlockRewardsPool) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgIncrementBlockRewardsPool) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Funder)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid funder address (%s)", err)
	}

	var am sdk.Coins = msg.Amount

	if am.IsZero() || am.IsAnyNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "You can't transfer a null or negative amount")
	}

	return nil
}

// -------------------------
// --- MsgSetRewardRate
// -------------------------

var _ sdk.Msg = &MsgSetRewardRate{}

func NewMsgSetRewardRate(government string, rewardRate sdk.DecProto) *MsgSetRewardRate {
	return &MsgSetRewardRate{
		Government: government,
		RewardRate: &rewardRate,
	}
}

func (msg *MsgSetRewardRate) Route() string {
	return RouterKey
}

func (msg *MsgSetRewardRate) Type() string {
	return MsgTypeSetRewardRate
}

func (msg *MsgSetRewardRate) GetSigners() []sdk.AccAddress {
	gov, err := sdk.AccAddressFromBech32(msg.Government)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{gov}
}

func (msg *MsgSetRewardRate) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetRewardRate) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Government)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid government address (%s)", err)
	}

	if err := ValidateRewardRate(msg.RewardRate.Dec); err != nil {
		return err
	}

	return nil
}

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

var _ sdk.Msg = &MsgSetAutomaticWithdraw{}

func NewMsgSetAutomaticWithdraw(government string, automaticWithdraw bool) *MsgSetAutomaticWithdraw {
	return &MsgSetAutomaticWithdraw{
		Government: government,
		AutomaticWithdraw: automaticWithdraw,
	}
}

func (msg *MsgSetAutomaticWithdraw) Route() string {
	return RouterKey
}

func (msg *MsgSetAutomaticWithdraw) Type() string {
	return MsgTypeSetAutomaticWithdraw
}

func (msg *MsgSetAutomaticWithdraw) GetSigners() []sdk.AccAddress {
	gov, err := sdk.AccAddressFromBech32(msg.Government)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{gov}
}

func (msg *MsgSetAutomaticWithdraw) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetAutomaticWithdraw) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Government)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid government address (%s)", err)
	}
	return nil
}
