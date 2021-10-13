package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func (k Keeper) Tsps(c context.Context, req *types.QueryTspsRequest) (*types.QueryTspsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	var tsps []string
	for _, tsp := range k.GetTrustedServiceProviders(ctx) {
		tsps = append(tsps, tsp.String())
	}
	return &types.QueryTspsResponse{Tsps: tsps}, nil
}

func (k Keeper) Funds(c context.Context, req *types.QueryFundsRequest) (*types.QueryFundsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	liquidityPoolAmount := k.GetPoolFunds(ctx)
	return &types.QueryFundsResponse{Funds: liquidityPoolAmount}, nil
}
