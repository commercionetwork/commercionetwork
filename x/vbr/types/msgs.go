package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

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

func (msg MsgIncrementBlockRewardsPool) Route() string { return ModuleName }

func (msg MsgIncrementBlockRewardsPool) Type() string { return MsgTypeIncrementBlockRewardsPool }

func (msg MsgIncrementBlockRewardsPool) ValidateBasic() error {
	if msg.Funder.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, msg.Funder.String())
	}
	if msg.Amount.IsZero() || msg.Amount.IsAnyNegative() {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "You can't transfer a null or negative amount")
	}

	return nil
}

func (msg MsgIncrementBlockRewardsPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIncrementBlockRewardsPool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Funder}
}
