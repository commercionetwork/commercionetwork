package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)
/*
const (
	newUpgradeScheduled = "new_scheduled_upgrade"
	delUpgradeScheduled = "del_scheduled_upgrade"
)
*/
func (k msgServer) ScheduleUpgrade(goCtx context.Context, msg *types.MsgScheduleUpgrade) (*types.MsgScheduleUpgradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	proposer, e := sdk.AccAddressFromBech32(msg.Proposer)
	if e != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid proposer address (%s)", e)
	}

	err := k.ScheduleUpgradeGov(ctx, proposer, *msg.Plan)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	/*ctx.EventManager().EmitEvent(sdk.NewEvent(
		newUpgradeScheduled,
		sdk.NewAttribute("name", msg.Plan.Name),
		sdk.NewAttribute("time", msg.Plan.Time.String()),
		sdk.NewAttribute("height", strconv.FormatInt(msg.Plan.Height, 10)),
		sdk.NewAttribute("info", msg.Plan.Info)))

	ctypes.EmitCommonEvents(ctx, proposer)*/
	return &types.MsgScheduleUpgradeResponse{PlanName: msg.Plan.Name}, nil
}

func (k msgServer) DeleteUpgrade(goCtx context.Context,  msg *types.MsgDeleteUpgrade) (*types.MsgDeleteUpgradeResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	proposer, e := sdk.AccAddressFromBech32(msg.Proposer)
	if e != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid proposer address (%s)", e)
	}
	err := k.DeleteUpgradeGov(ctx, proposer)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}
	/*ctx.EventManager().EmitEvent(sdk.NewEvent(
		delUpgradeScheduled,
		sdk.NewAttribute("action", "deleted")))

	ctypes.EmitCommonEvents(ctx, msg.Proposer)*/
	return &types.MsgDeleteUpgradeResponse{}, nil
}
