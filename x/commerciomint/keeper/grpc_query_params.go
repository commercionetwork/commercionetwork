package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) Params(c context.Context, req *types.QueryGetParams) (*types.QueryGetParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := k.GetParams(ctx)

	return &types.QueryGetParamsResponse{ConversionRate: params.ConversionRate.String(), FreezePeriod: params.FreezePeriod.String()}, nil
}
