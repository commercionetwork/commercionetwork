package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) IncrementBlockRewardsPool(goCtx context.Context, msg *types.MsgIncrementBlockRewardsPool) (*types.MsgIncrementBlockRewardsPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	//return &types.MsgIncrementBlockRewardsPoolResponse{}, nil
	funderAddr, e := sdk.AccAddressFromBech32(msg.Funder)
	if e != nil {
		return nil, e
	}
	// Subtract the coins from the account
	if err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, funderAddr, types.ModuleName, msg.Amount); err != nil {
		return nil, err
	}

	// Set the total rewards pool
	//k.SetTotalRewardPool(ctx, k.GetTotalRewardPool(ctx).Add(sdk.NewDecCoinsFromCoins(msg.Amount...)...))

	//return &sdk.Result{Events: ctx.EventManager().Events(), Log: "Block reward pool successfully increased"}, nil
	return &types.MsgIncrementBlockRewardsPoolResponse{}, nil
}

func (k msgServer) SetRewardRate(goCtx context.Context, msg *types.MsgSetRewardRate) (*types.MsgSetRewardRateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgSetRewardRateResponse{}, nil
}

func (k msgServer) SetAutomaticWithdraw(goCtx context.Context, msg *types.MsgSetAutomaticWithdraw) (*types.MsgSetAutomaticWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgSetAutomaticWithdrawResponse{}, nil
}
