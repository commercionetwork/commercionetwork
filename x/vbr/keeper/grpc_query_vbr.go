package keeper

import (
	"context"

    "github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)
func (k Keeper) GetBlockRewardsPoolFunds(goCtx context.Context,  req *types.QueryGetBlockRewardsPoolFundsRequest) (*types.QueryGetBlockRewardsPoolFundsResponse, error) {
	if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

	ctx := sdk.UnwrapSDKContext(goCtx)
	funds := k.GetTotalRewardPool(ctx)

	/*fundsBz, err2 := codec.MarshalJSONIndent(keeper.cdc, funds)
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}*/
	return &types.QueryGetBlockRewardsPoolFundsResponse{Funds: funds}, nil
}

func (k Keeper) GetRewardRate(goCtx context.Context,  req *types.QueryGetRewardRateRequest) (*types.QueryGetRewardRateResponse, error) {
	if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

	ctx := sdk.UnwrapSDKContext(goCtx)
	rate := k.GetRewardRateKeeper(ctx)

	return &types.QueryGetRewardRateResponse{RewardRate: rate}, nil
}

func (k Keeper) GetAutomaticWithdraw(goCtx context.Context,  req *types.QueryGetAutomaticWithdrawRequest) (*types.QueryGetAutomaticWithdrawResponse, error) {
	if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

	ctx := sdk.UnwrapSDKContext(goCtx)
	autoW := k.GetAutomaticWithdrawKeeper(ctx)

	return &types.QueryGetAutomaticWithdrawResponse{AutoW: autoW}, nil
}
