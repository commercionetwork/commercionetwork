package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/pricefeed/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryGetCurrentPrices:
			return queryGetCurrentPrices(ctx, path[1:], keeper)
		case types.QueryGetCurrentPrice:
			return queryGetCurrentPrice(ctx, path[1:], keeper)
		case types.QueryGetOracles:
			return queryGetOracles(ctx, path[1:], keeper)
		case types.QueryGetBlacklistedDenoms:
			return queryGetBlacklistedDenoms(ctx, path[1:], keeper)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryGetCurrentPrices(ctx sdk.Context, _ []string, keeper Keeper) ([]byte, error) {
	prices := keeper.GetCurrentPrices(ctx)
	if prices == nil {
		prices = make(types.Prices, 0)
	}

	pricesBz, err := codec.MarshalJSONIndent(keeper.cdc, prices)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return pricesBz, nil
}

func queryGetCurrentPrice(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	asset := path[0]

	price, found := keeper.GetCurrentPrice(ctx, asset)
	if !found {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Could not find price for asset %s", asset))
	}

	priceBz, err := codec.MarshalJSONIndent(keeper.cdc, price)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return priceBz, nil
}

func queryGetOracles(ctx sdk.Context, _ []string, keeper Keeper) ([]byte, error) {
	oracles := keeper.GetOracles(ctx)
	if oracles == nil {
		oracles = make([]sdk.AccAddress, 0)
	}

	oraclesBz, err := codec.MarshalJSONIndent(keeper.cdc, oracles)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return oraclesBz, nil
}

func queryGetBlacklistedDenoms(ctx sdk.Context, _ []string, keeper Keeper) ([]byte, error) {
	denomBlacklist := keeper.DenomBlacklist(ctx)

	denomBlacklistBz, err := codec.MarshalJSONIndent(keeper.cdc, denomBlacklist)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return denomBlacklistBz, nil
}
