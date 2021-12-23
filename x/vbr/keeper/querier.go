package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryBlockRewardsPoolFunds:
			return queryGetBlockRewardsPoolFunds(ctx, path[1:], k, legacyQuerierCdc)
		case types.QueryVbrParams:
			return queryVbrParams(ctx, k, legacyQuerierCdc)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint: %s", types.ModuleName, path[0]))
		}
	}
}

func queryGetBlockRewardsPoolFunds(ctx sdk.Context, _ []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	funds := k.GetTotalRewardPool(ctx)

	fundsBz, err2 := codec.MarshalJSONIndent(legacyQuerierCdc, funds)
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return fundsBz, nil
}

func queryVbrParams(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error){
	params := k.GetParams(ctx)

	paramsBz, err := codec.MarshalJSONIndent(legacyQuerierCdc, params)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}
	
	return paramsBz, nil
}