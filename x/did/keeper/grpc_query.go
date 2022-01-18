package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/did/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Identity(c context.Context, req *types.QueryResolveIdentityRequest) (*types.QueryResolveIdentityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	identity, err := k.GetLastIdentityOfAddress(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &types.QueryResolveIdentityResponse{Identity: identity}, nil
}
