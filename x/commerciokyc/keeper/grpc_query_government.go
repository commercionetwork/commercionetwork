package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func (k Keeper) Tsps(c context.Context, req *types.QueryTspsRequest) (*types.QueryTspsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var tsps []string
	tsps = k.GetTrustedServiceProviders(ctx).Addresses
	return &types.QueryTspsResponse{Tsps: tsps}, nil
}

func (k Keeper) Funds(c context.Context, req *types.QueryFundsRequest) (*types.QueryFundsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	// TODO check this Get functions
	liquidityPoolAmount := k.GetLiquidityPoolAmountCoins(ctx)
	return &types.QueryFundsResponse{Funds: liquidityPoolAmount}, nil
}
