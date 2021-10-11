package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// ----------------------------------
// --- MsgScheduleUpgrade
// ----------------------------------
var _ sdk.Msg = &MsgScheduleUpgrade{}

func NewMsgScheduleUpgrade(proposer string, plan upgradeTypes.Plan) *MsgScheduleUpgrade {
	return &MsgScheduleUpgrade{
		Proposer: proposer,
		Plan: &plan,
	}
}

const ScheduleUpgradeConst = "ScheduleUpgrade"

func (msg *MsgScheduleUpgrade) Route() string {
	return RouterKey
}

func (msg *MsgScheduleUpgrade) Type() string {
	return ScheduleUpgradeConst
}

func (msg *MsgScheduleUpgrade) GetSigners() []sdk.AccAddress {
	proposer, err := sdk.AccAddressFromBech32(msg.Proposer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{proposer}
}

func (msg *MsgScheduleUpgrade) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgScheduleUpgrade) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Proposer)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid proposer address (%s)", err)
	}

	if err := msg.Plan.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// ----------------------------------
// --- MsgDeleteUpgrade
// ----------------------------------

var _ sdk.Msg = &MsgDeleteUpgrade{}

func NewMsgDeleteUpgrade(proposer string) *MsgDeleteUpgrade {
  return &MsgDeleteUpgrade{
	  Proposer: proposer,
	}
}

const DeleteUpgradeConst = "DeleteUpgrade"

func (msg *MsgDeleteUpgrade) Route() string {
  return RouterKey
}

func (msg *MsgDeleteUpgrade) Type() string {
  return DeleteUpgradeConst
}

func (msg *MsgDeleteUpgrade) GetSigners() []sdk.AccAddress {
  proposer, err := sdk.AccAddressFromBech32(msg.Proposer)
  if err != nil {
    panic(err)
  }
  return []sdk.AccAddress{proposer}
}

func (msg *MsgDeleteUpgrade) GetSignBytes() []byte {
  bz := ModuleCdc.MustMarshalJSON(msg)
  return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteUpgrade) ValidateBasic() error {
  _, err := sdk.AccAddressFromBech32(msg.Proposer)
  	if err != nil {
  		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid proposer address (%s)", err)
  	}
  return nil
}


