package upgrade

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/upgrade/keeper"
	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
)

func NewHandler(keeper keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
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

func handleScheduleUpgrade(ctx sdk.Context, keeper keeper.Keeper, msg types.MsgScheduleUpgrade) (*sdk.Result, error) {
	err := keeper.ScheduleUpgrade(ctx, msg.Proposer, msg.Plan)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}

	return &sdk.Result{Log: "Upgrade scheduled successfully"}, nil
}

func handleDeleteUpgrade(ctx sdk.Context, k keeper.Keeper, msg types.MsgDeleteUpgrade) (*sdk.Result, error) {
	err := k.DeleteUpgrade(ctx, msg.Proposer)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}

	return &sdk.Result{Log: "Upgrade deleted successfully"}, nil
}
