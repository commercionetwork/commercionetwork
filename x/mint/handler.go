package mint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDepositToken:
			return handleMsgOpenCdp(ctx, keeper, msg)
		case MsgWithdrawToken:
			return handleMsgCloseCdp(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgOpenCdp(ctx sdk.Context, keeper Keeper, msg MsgDepositToken) sdk.Result {
	err := keeper.OpenCdp(ctx, CdpRequest(msg))
	if err != nil {
		return sdk.ResultFromError(err)
	}

	return sdk.Result{Log: "Cdp opened successfully"}
}

func handleMsgCloseCdp(ctx sdk.Context, keeper Keeper, msg MsgWithdrawToken) sdk.Result {
	err := keeper.CloseCdp(ctx, msg.Signer, msg.Timestamp)
	if err != nil {
		return sdk.ResultFromError(err)
	}

	return sdk.Result{Log: "Cdp closed successfully"}
}
