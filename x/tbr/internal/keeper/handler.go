package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/tbr/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgIncrementBlockRewardsPool:
			return handleMsgIncrementBlockRewardsPool(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgIncrementBlockRewardsPool(ctx sdk.Context, keeper Keeper, msg types.MsgIncrementBlockRewardsPool) sdk.Result {
	if err := keeper.IncrementBlockRewardsPool(ctx, msg.Funder, msg.Amount); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}
