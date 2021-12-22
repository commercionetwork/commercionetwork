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

func (k Keeper) GetVbrParams(goCtx context.Context, req *types.QueryGetVbrParamsRequest) (*types.QueryGetVbrParamsResponse, error){
	if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)
	
	return &types.QueryGetVbrParamsResponse{VbrParams: params}, nil
}