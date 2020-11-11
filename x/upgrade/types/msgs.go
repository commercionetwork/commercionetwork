package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
)

var _ sdk.Msg = &MsgScheduleUpgrade{}

type MsgScheduleUpgrade struct {
	Proposer sdk.AccAddress `json:"proposer"`
	Plan     upgrade.Plan   `json:"plan"`
}

const ScheduleUpgradeConst = "ScheduleUpgrade"

func (msg MsgScheduleUpgrade) Route() string {
	return RouterKey
}

func (msg MsgScheduleUpgrade) Type() string {
	return ScheduleUpgradeConst
}

func (msg MsgScheduleUpgrade) ValidateBasic() error {
	if msg.Proposer.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, msg.Proposer.String())
	}

	if err := msg.Plan.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

func (msg MsgScheduleUpgrade) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgScheduleUpgrade) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proposer}
}

var _ sdk.Msg = &MsgDeleteUpgrade{}

type MsgDeleteUpgrade struct {
	Proposer sdk.AccAddress `json:"proposer"`
}

const DeleteUpgradeConst = "DeleteUpgrade"

func (msg MsgDeleteUpgrade) Route() string {
	return RouterKey
}

func (msg MsgDeleteUpgrade) Type() string {
	return DeleteUpgradeConst
}

func (msg MsgDeleteUpgrade) ValidateBasic() error {
	if msg.Proposer.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, msg.Proposer.String())
	}

	return nil
}

func (msg MsgDeleteUpgrade) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgDeleteUpgrade) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Proposer}
}
