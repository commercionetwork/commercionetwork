package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgIncrementBlockRewardsPool struct {
	Funder sdk.AccAddress `json:"funder"`
	Amount sdk.Coin       `json:"amount"`
}

func NewMsgIncrementBlockRewardsPool(funder sdk.AccAddress, amount sdk.Coin) MsgIncrementBlockRewardsPool {
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
	if msg.Amount.Amount.IsZero() {
		return sdk.ErrUnknownRequest("You can't transfer a null amount")
	}

	return nil
}

func (msg MsgIncrementBlockRewardsPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIncrementBlockRewardsPool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Funder}
}
