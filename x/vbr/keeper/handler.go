package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgSetRewardRate:
			return handleMsgSetRewardRate(ctx, keeper, msg)
		case types.MsgSetAutomaticWithdraw:
			return handleMsgSetAutomaticWithdraw(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgSetRewardRate(ctx sdk.Context, keeper Keeper, msg types.MsgSetRewardRate) (*sdk.Result, error) {
	gov := keeper.govKeeper.GetGovernmentAddress(ctx)
	if !(gov.Equals(msg.Government)) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("%s cannot set reward rate", msg.Government))
	}
	err := keeper.SetRewardRate(ctx, msg.RewardRate)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}
	return &sdk.Result{Log: fmt.Sprintf("Reward rate changed successfully to %s", msg.RewardRate)}, nil
}

func handleMsgSetAutomaticWithdraw(ctx sdk.Context, keeper Keeper, msg types.MsgSetAutomaticWithdraw) (*sdk.Result, error) {
	gov := keeper.govKeeper.GetGovernmentAddress(ctx)
	if !(gov.Equals(msg.Government)) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("%s cannot set reward rate", msg.Government))
	}
	err := keeper.SetAutomaticWithdraw(ctx, msg.AutomaticWithdraw)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}
	return &sdk.Result{Log: fmt.Sprintf("Automatic withdraw changed successfully to %v", msg.AutomaticWithdraw)}, nil
}
