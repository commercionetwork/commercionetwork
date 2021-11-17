package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/store/prefix"
)

func (k Keeper) EtpByOwner(c context.Context, req *types.QueryEtpRequestByOwner) (*types.QueryEtpsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var positions  []*types.Position
	ctx := sdk.UnwrapSDKContext(c)
	
	store := ctx.KVStore(k.storeKey)
	owner, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, err
	}
	etpStore := prefix.NewStore(store, getEtpByOwnerIdsStoreKey(owner))

	//positions := k.GetAllPositionsOwnedBy(ctx, owner)
	pageRes, err := query.Paginate(etpStore, req.Pagination, func(key []byte, value []byte) error {
		positions = k.GetAllPositionsOwnedBy(ctx, owner)

		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryEtpsResponse{Positions: positions, Pagination: pageRes}, nil
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
