package mint

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgDepositToken:
			return handleMsgDepositToken(ctx, keeper, msg)
		case MsgWithdrawToken:
			return handleMsgWithdrawToken(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgDepositToken(ctx sdk.Context, keeper Keeper, msg MsgDepositToken) sdk.Result {
	credits, err := keeper.OpenCDP(ctx, msg.Signer, msg.Tokens)
	if err != nil {
		return sdk.ResultFromError(err)
	}

	goodRes := fmt.Sprintf("Token deposit successful, credits' transfered amount: %s", credits.AmountOf(DefaultCreditsDenom))
	return sdk.Result{Log: goodRes}
}

func handleMsgWithdrawToken(ctx sdk.Context, keeper Keeper, msg MsgWithdrawToken) sdk.Result {
	token, err := keeper.CloseCDP(ctx, msg.Signer, msg.Amount)
	if err != nil {
		return sdk.ResultFromError(err)
	}

	goodRes := fmt.Sprintf("Token withdrawal successful, tokens' withdrawed amount: %s", token.AmountOf(DefaultBondDenom))
	return sdk.Result{Log: goodRes}
}
