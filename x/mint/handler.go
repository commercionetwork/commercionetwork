package mint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDepositToken:
			return handleMsgOpenCDP(ctx, keeper, msg)
		case MsgWithdrawToken:
			return handleMsgCloseCDP(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgOpenCDP(ctx sdk.Context, keeper Keeper, msg MsgDepositToken) sdk.Result {
	err := keeper.OpenCDP(ctx, msg.Request)
	if err != nil {
		return sdk.ResultFromError(err)
	}

	return sdk.Result{Log: "CDP opened successfully"}
}

func handleMsgCloseCDP(ctx sdk.Context, keeper Keeper, msg MsgWithdrawToken) sdk.Result {
	err := keeper.CloseCDP(ctx, msg.Signer, msg.Timestamp)
	if err != nil {
		return sdk.ResultFromError(err)
	}

	return sdk.Result{Log: "CDP closed successfully"}
}
