package keeper

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
)

const (
	newUpgradeScheduled = "new_scheduled_upgrade"
	delUpgradeScheduled = "del_scheduled_upgrade"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgScheduleUpgrade:
			return handleScheduleUpgrade(ctx, keeper, msg)
		case types.MsgDeleteUpgrade:
			return handleDeleteUpgrade(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)
		}
	}
}

func handleScheduleUpgrade(ctx sdk.Context, keeper Keeper, msg types.MsgScheduleUpgrade) (*sdk.Result, error) {
	err := keeper.ScheduleUpgrade(ctx, msg.Proposer, msg.Plan)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		newUpgradeScheduled,
		sdk.NewAttribute("name", msg.Plan.Name),
		sdk.NewAttribute("time", msg.Plan.Time.String()),
		sdk.NewAttribute("height", strconv.FormatInt(msg.Plan.Height, 10)),
		sdk.NewAttribute("info", msg.Plan.Info)))

	ctypes.EmitCommonEvents(ctx, msg.Proposer)
	return &sdk.Result{Events: ctx.EventManager().Events(), Log: "Upgrade scheduled successfully"}, nil
}

func handleDeleteUpgrade(ctx sdk.Context, k Keeper, msg types.MsgDeleteUpgrade) (*sdk.Result, error) {
	err := k.DeleteUpgrade(ctx, msg.Proposer)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		delUpgradeScheduled,
		sdk.NewAttribute("action", "deleted")))

	ctypes.EmitCommonEvents(ctx, msg.Proposer)
	return &sdk.Result{Events: ctx.EventManager().Events(), Log: "Upgrade deleted successfully"}, nil
}
