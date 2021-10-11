package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) EtpByOwner(c context.Context, req *types.QueryEtpRequestByOwner) (*types.QueryEtpsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	positions := k.GetAllPositionsOwnedBy(ctx, sdk.AccAddress(req.Owner))
	return &types.QueryEtpsResponse{Positions: positions}, nil
}

func (k Keeper) Etps(c context.Context, req *types.QueryEtpsRequest) (*types.QueryEtpsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	positions := k.GetAllPositions(ctx)
	return &types.QueryEtpsResponse{Positions: positions}, nil
}

func (k Keeper) Etp(c context.Context, req *types.QueryEtpRequest) (*types.QueryEtpResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	_ = ctx
	//positions := k.GetAllPositions(ctx)
	return &types.QueryEtpResponse{}, nil
}

/*
	Etp(ctx context.Context, in *QueryEtpRequest, opts ...grpc.CallOption) (*QueryEtpResponse, error)
*/
