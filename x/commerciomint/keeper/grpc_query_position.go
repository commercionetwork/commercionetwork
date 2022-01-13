package keeper

import (
	"context"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) EtpsByOwner(c context.Context, req *types.QueryEtpRequestByOwner) (*types.QueryEtpsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var positions []*types.Position

	owner, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInvalidAddress, fmt.Sprintf("Error while converting address: %s", err.Error()))
	}

	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	//etpStore := prefix.NewStore(store, []byte(types.EtpStorePrefix))
	etpStore := prefix.NewStore(store, getEtpByOwnerIdsStoreKey(owner))

	//positions := k.GetAllPositionsOwnedBy(ctx, owner)
	pageRes, err := query.Paginate(etpStore, req.Pagination, func(key []byte, value []byte) error {
		/*positions = k.GetAllPositionsOwnedBy(ctx, owner)
		position, ok := k.GetPositionById(ctx, string(value))
		if !ok {
			return status.Error(codes.NotFound, fmt.Sprintf("Position with id %s not found!", string(value)))
		}
		positions = append(positions, &position)*/
		var position types.Position
		e := k.cdc.UnmarshalBinaryBare(value, &position)
		if e != nil {
			return e
		}
		positions = append(positions, &position)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Pagination error: %s", err.Error()))
	}
	return &types.QueryEtpsResponse{Positions: positions, Pagination: pageRes}, nil
}

func (k Keeper) Etps(c context.Context, req *types.QueryEtpsRequest) (*types.QueryEtpsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var positions []*types.Position
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	etpStore := prefix.NewStore(store, []byte(types.EtpStorePrefix))
	pageRes, err := query.Paginate(etpStore, req.Pagination, func(key []byte, value []byte) error {
		positions = k.GetAllPositions(ctx)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryEtpsResponse{Positions: positions, Pagination: pageRes}, nil
}

func (k Keeper) Etp(c context.Context, req *types.QueryEtpRequest) (*types.QueryEtpResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var position types.Position
	ctx := sdk.UnwrapSDKContext(c)

	position, ok := k.GetPositionById(ctx, req.ID)
	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Position with id %s not found!", req.ID))
	}

	//positions := k.GetAllPositions(ctx)

	return &types.QueryEtpResponse{Position: &position}, nil
}

/*
	Etp(ctx context.Context, in *QueryEtpRequest, opts ...grpc.CallOption) (*QueryEtpResponse, error)
*/
