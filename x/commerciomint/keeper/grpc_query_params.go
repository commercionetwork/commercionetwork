package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Params(c context.Context, req *types.QueryGetParams) (*types.QueryGetParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryGetParamsResponse{
		ConversionRate: params.ConversionRate.String(),
		FreezePeriod:   params.FreezePeriod.String(),
	}, nil
}

func (k Keeper) ConversionRate(c context.Context, req *types.QueryConversionRate) (*types.QueryConversionRateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	rate := k.GetConversionRate(ctx)
	return &types.QueryConversionRateResponse{Rate: rate}, nil
}

func (k Keeper) FreezePeriod(c context.Context, req *types.QueryFreezePeriod) (*types.QueryFreezePeriodResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	freezePeriod := k.GetFreezePeriod(ctx)
	return &types.QueryFreezePeriodResponse{FreezePeriod: freezePeriod.String()}, nil
}
