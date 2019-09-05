package accreditation

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler is essentially a sub-router that directs messages coming into this module to the proper handler.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetAccrediter:
			return handleSetAccrediter(ctx, keeper, msg)
		case MsgDistributeReward:
			return handleDistributeReward(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleSetAccrediter(ctx sdk.Context, keeper Keeper, msg MsgSetAccrediter) sdk.Result {

	// Check the signer
	if !keeper.IsTrustworthySigner(ctx, msg.Signer) {
		errMsg := fmt.Sprintf("The signer %s is not trustworthy", msg.Signer.String())
		return sdk.ErrUnknownRequest(errMsg).Result()
	}

	// Check the accrediter
	if accrediter, found := keeper.GetAccrediter(ctx, msg.User); found {
		errMsg := fmt.Sprintf("User %s already has an accrediter (%s)", msg.User.String(), accrediter.String())
		return sdk.ErrUnknownRequest(errMsg).Result()
	}

	// If everything passes the checks, set the accrediter
	keeper.SetAccrediter(ctx, msg.Accrediter, msg.User)

	return sdk.Result{}
}

func handleDistributeReward(ctx sdk.Context, keeper Keeper, msg MsgDistributeReward) sdk.Result {
	// TODO

	// 1. Check that the pair user - accrediter has not status accreditated
	// 2. Distribute the reward

	return sdk.Result{}
}
