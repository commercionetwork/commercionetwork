package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgOpenCdp:
			return handleMsgOpenCdp(ctx, keeper, msg)
		case types.MsgCloseCdp:
			return handleMsgCloseCdp(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgOpenCdp(ctx sdk.Context, keeper Keeper, msg types.MsgOpenCdp) (*sdk.Result, error) {
	err := keeper.OpenCdp(ctx, msg.Depositor, msg.DepositedAmount)
	if err != nil {
		return sdk.ResultFromError(err)
	}

	return sdk.Result{Log: "Cdp opened successfully"}
}

func handleMsgCloseCdp(ctx sdk.Context, keeper Keeper, msg types.MsgCloseCdp) (*sdk.Result, error) {
	err := keeper.CloseCdp(ctx, msg.Signer, msg.Timestamp)
	if err != nil {
		return sdk.ResultFromError(err)
	}

	return sdk.Result{Log: "Cdp closed successfully"}
}
