package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func (msg MsgIncrementBlockRewardsPool) ValidateBasic() sdk.Error {
	if msg.Funder.Empty() {
		return sdk.ErrInvalidAddress(msg.Funder.String())
	}
	if msg.Amount.IsZero() || msg.Amount.IsAnyNegative() {
		return sdk.ErrUnknownRequest("You can't transfer a null or negative amount")
	}

	return nil
}

func (msg MsgIncrementBlockRewardsPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIncrementBlockRewardsPool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Funder}
}
