package id

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/government"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for type messages and is essentially a sub-router that directs
// messages coming into this module to the proper handler.
func NewHandler(keeper Keeper, govKeeper government.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSetIdentity:
			return handleMsgSetIdentity(ctx, keeper, msg)
		case MsgRequestDidDeposit:
			return handleMsgRequestDidDeposit(ctx, keeper, msg)
		case MsgChangeDidDepositRequestStatus:
			return handleMsgChangeDidDepositRequestStatus(ctx, keeper, govKeeper, msg)
		case MsgRequestDidPowerup:
			return handleMsgRequestDidPowerup(ctx, keeper, msg)
		case MsgChangeDidPowerupRequestStatus:
			return handleMsgChangeDidPowerupRequestStatus(ctx, keeper, govKeeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgSetIdentity(ctx sdk.Context, keeper Keeper, msg MsgSetIdentity) sdk.Result {
	keeper.SaveIdentity(ctx, msg.Owner, msg.DidDocument)
	return sdk.Result{}
}

func handleMsgRequestDidDeposit(ctx sdk.Context, keeper Keeper, msg MsgRequestDidDeposit) sdk.Result {

}

func handleMsgChangeDidDepositRequestStatus(ctx sdk.Context, keeper Keeper, govKeeper government.Keeper,
	msg MsgChangeDidDepositRequestStatus) sdk.Result {

}

func handleMsgRequestDidPowerup(ctx sdk.Context, keeper Keeper, msg MsgRequestDidPowerup) sdk.Result {

}

func handleMsgChangeDidPowerupRequestStatus(ctx sdk.Context, keeper Keeper, govKeeper government.Keeper,
	msg MsgChangeDidPowerupRequestStatus) sdk.Result {

}
