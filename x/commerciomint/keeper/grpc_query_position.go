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

func (k Keeper) EtpsByOwner(c context.Context, req *types.QueryEtpsByOwnerRequest) (*types.QueryEtpsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	owner, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInvalidAddress, fmt.Sprintf("could not convert address: %s", err.Error()))
	}

	ctx := sdk.UnwrapSDKContext(c)

	positions := []*types.Position{}

	store := ctx.KVStore(k.storeKey)
	etpsByOwnerStore := prefix.NewStore(store, getEtpsByOwnerStoreKey(owner))

	pageRes, err := query.Paginate(etpsByOwnerStore, req.Pagination, func(key []byte, value []byte) error {
		var position types.Position
		e := k.cdc.Unmarshal(value, &position)
		if e != nil {
			return e
		}
		positions = append(positions, &position)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("pagination error: %s", err.Error()))
	}
	return &types.QueryEtpsResponse{Positions: positions, Pagination: pageRes}, nil
}

func (k Keeper) Etps(c context.Context, req *types.QueryEtpsRequest) (*types.QueryEtpsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	positions := []*types.Position{}
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	etpStore := prefix.NewStore(store, []byte(types.EtpStorePrefix))
	pageRes, err := query.Paginate(etpStore, req.Pagination, func(key []byte, value []byte) error {
		var position types.Position
		e := k.cdc.Unmarshal(value, &position)
		if e != nil {
			return e
		}
		positions = append(positions, &position)
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
	ctx := sdk.UnwrapSDKContext(c)

	position, ok := k.GetPositionById(ctx, req.ID)
	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("could not find position with id %s", req.ID))
	}

	return &types.QueryEtpResponse{Position: &position}, nil
}
