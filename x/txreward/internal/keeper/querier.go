package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"

	tbrTypes "github.com/commercionetwork/commercionetwork/x/txreward/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

//Querier unused
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case tbrTypes.QueryBlockRewardsPoolFunds:
			return queryGetBlockRewardsPoolFunds(ctx, path[1:], keeper)
		case tbrTypes.QueryBlockRewardsPoolFunders:
			return queryGetBlockRewardsPoolFunders(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown %s query endpoint", tbrTypes.ModuleName))
		}
	}
}

func queryGetBlockRewardsPoolFunds(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	funds := keeper.getBrPool(ctx)

	fundsBz, err2 := codec.MarshalJSONIndent(keeper.Cdc, funds)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return fundsBz, nil
}

func queryGetBlockRewardsPoolFunders(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	funders := keeper.getFunders(ctx)

	fundersBz, err2 := codec.MarshalJSONIndent(keeper.Cdc, funders)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return fundersBz, nil
}
