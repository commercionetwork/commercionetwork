package keeper

import (
	"context"

    "github.com/commercionetwork/commercionetwork/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k Keeper) CurrentUpgrade(goCtx context.Context,  req *types.QueryCurrentUpgradeRequest) (*types.QueryCurrentUpgradeResponse, error) {
	if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

	ctx := sdk.UnwrapSDKContext(goCtx)
	
	plan, has := k.GetUpgradePlan(ctx)
	if !has {
		return nil, nil
	}

	res, err := types.ModuleCdc.MarshalBinaryBare(&plan) //codec.MarshalJSONIndent(k.cdc, plan)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result")
	}

	return &types.QueryCurrentUpgradeResponse{CurrentUpgrade: res }, nil
}
