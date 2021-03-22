package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/government/types"
)

// NewHandler is essentially a sub-router that directs messages coming into this module to the proper handler.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgSetTumblerAddress:
			return handleMsgSetTumblerAddress(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgSetTumblerAddress handles MsgSetTumblerAddress messages
func handleMsgSetTumblerAddress(ctx sdk.Context, keeper Keeper, msg types.MsgSetTumblerAddress) (*sdk.Result, error) {
	if !keeper.GetGovernmentAddress(ctx).Equals(msg.GetSigners()[0]) {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, "government address specified not allowed")
	}

	err := keeper.SetTumblerAddress(ctx, msg.NewTumbler)

	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}

	ctypes.EmitCommonEvents(ctx, msg.GetSigners()[0])
	return &sdk.Result{Events: ctx.EventManager().Events(), Log: "Tumbler address configured"}, nil
}
