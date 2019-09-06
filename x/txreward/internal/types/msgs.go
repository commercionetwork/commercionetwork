package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgIncrementBlockRewardsPool struct {
	Funder Funder   `json:"funder"`
	Amount sdk.Coin `json:"amount"`
}

func NewMsgIncrementBlockRewardsPool(funder Funder, amount sdk.Coin) MsgIncrementBlockRewardsPool {
	return MsgIncrementBlockRewardsPool{
		Funder: funder,
		Amount: amount,
	}
}

func (msg MsgIncrementBlockRewardsPool) Route() string { return ModuleName }

func (msg MsgIncrementBlockRewardsPool) Type() string { return MsgTypeIncrementBlockRewardsPool }

func (msg MsgIncrementBlockRewardsPool) ValidateBasic() sdk.Error {
	if msg.Funder.Address.Empty() {
		return sdk.ErrInvalidAddress(msg.Funder.Address.String())
	}
	if msg.Amount.Amount.IsZero() {
		return sdk.ErrUnknownRequest("You can't transfer a null amount")
	}
	if msg.Amount.Denom != DefaultBondDenom {
		return sdk.ErrUnknownRequest(fmt.Sprintf("You can't transfer others than %s", DefaultBondDenom))
	}

	return nil
}

func (msg MsgIncrementBlockRewardsPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgIncrementBlockRewardsPool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Funder.Address}
}
