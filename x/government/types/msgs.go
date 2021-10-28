package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ----------------------------------
// --- MsgSetGovAddress
// ----------------------------------
var _ sdk.Msg = &MsgSetGovAddress{}

func NewMsgSetGovAddress() *MsgSetGovAddress {
	return &MsgSetGovAddress{}
}

const SetGovAddressConst = "SetGovAddress"

func (msg *MsgSetGovAddress) Route() string {
	return RouterKey
}

func (msg *MsgSetGovAddress) Type() string {
	return SetGovAddressConst
}

func (msg *MsgSetGovAddress) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgSetGovAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetGovAddress) ValidateBasic() error {
	return nil
}
