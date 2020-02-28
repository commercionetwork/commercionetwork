package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgSetTumblerAddress is used to modify the current government tumbler address.
type MsgSetTumblerAddress struct {
	Government sdk.AccAddress `json:"government"`
	NewTumbler sdk.AccAddress `json:"new_tumbler"`
}

func NewMsgSetTumblerAddress(government sdk.AccAddress, newTumbler sdk.AccAddress) MsgSetTumblerAddress {
	return MsgSetTumblerAddress{
		Government: government,
		NewTumbler: newTumbler,
	}
}

// Route Implements Msg.
func (msg MsgSetTumblerAddress) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSetTumblerAddress) Type() string { return MsgTypeSetTumblerAddress }

// ValidateBasic Implements Msg.
func (msg MsgSetTumblerAddress) ValidateBasic() error {
	if msg.Government.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("invalid government address: %s", msg.Government))
	}
	if msg.NewTumbler.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, fmt.Sprintf("invalid new tumbler address: %s", msg.NewTumbler))
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetTumblerAddress) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSetTumblerAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Government}
}
