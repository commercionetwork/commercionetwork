package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case types.MsgIncrementBlockRewardsPool:
			return handleMsgIncrementBlockRewardsPool(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, msg.Type())
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)
		}
	}
}

func handleMsgIncrementBlockRewardsPool(ctx sdk.Context, k Keeper, msg types.MsgIncrementBlockRewardsPool) (*sdk.Result, error) {
	// Subtract the coins from the account
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, msg.Funder, types.ModuleName, msg.Amount); err != nil {
		return nil, err
	}

	// Set the total rewards pool
	k.SetTotalRewardPool(ctx, k.GetTotalRewardPool(ctx).Add(sdk.NewDecCoinsFromCoins(msg.Amount...)...))

	return &sdk.Result{}, nil
}
