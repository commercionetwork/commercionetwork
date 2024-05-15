package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	errors "cosmossdk.io/errors"
)

// GetSigners implements "github.com/cosmos/cosmos-sdk/types".Msg
func (msg *MsgSetIdentity) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.DidDocument.ID)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes implements "github.com/cosmos/cosmos-sdk/types".Msg
func (msg *MsgSetIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Route implements "github.com/cosmos/cosmos-sdk/types".Msg
func (msg *MsgSetIdentity) Route() string {
	return RouterKey
}

// Type implements "github.com/cosmos/cosmos-sdk/types".Msg
func (msg *MsgSetIdentity) Type() string {
	return MsgTypeSetIdentity
}

// ValidateBasic implements "github.com/cosmos/cosmos-sdk/types".Msg
func (msg *MsgSetIdentity) ValidateBasic() error {

	if msg == nil {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "msgSetDidDocument not defined")
	}

	if err := msg.DidDocument.Validate(); err != nil {
		return err
	}

	return nil
}
