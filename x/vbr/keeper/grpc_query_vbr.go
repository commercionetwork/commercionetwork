package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetBlockRewardsPoolFunds(goCtx context.Context, req *types.QueryGetBlockRewardsPoolFundsRequest) (*types.QueryGetBlockRewardsPoolFundsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	funds := k.GetTotalRewardPool(ctx)
	return &types.QueryGetBlockRewardsPoolFundsResponse{Funds: funds}, nil
}

func (k Keeper) GetParams(goCtx context.Context, req *types.QueryGetParamsRequest) (*types.QueryGetParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParamSet(ctx)

	return &types.QueryGetParamsResponse{Params: params}, nil
}
