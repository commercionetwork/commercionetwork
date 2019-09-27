package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryGetCurrentPrices:
			return queryGetCurrentPrices(ctx, path[1:], keeper)
		case types.QueryGetCurrentPrice:
			return queryGetCurrentPrice(ctx, path[1:], keeper)
		case types.QueryGetOracles:
			return queryGetOracles(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryGetCurrentPrices(ctx sdk.Context, _ []string, keeper Keeper) ([]byte, sdk.Error) {
	prices := keeper.GetCurrentPrices(ctx)
	if prices == nil {
		prices = make(types.CurrentPrices, 0)
	}

	pricesBz, err := codec.MarshalJSONIndent(keeper.cdc, prices)
	if err != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return pricesBz, nil
}

func queryGetCurrentPrice(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	asset := path[0]

	price, found := keeper.GetCurrentPrice(ctx, asset)
	if !found {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Could not find price for asset %s", asset))
	}

	priceBz, err := codec.MarshalJSONIndent(keeper.cdc, price)
	if err != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return priceBz, nil
}

func queryGetOracles(ctx sdk.Context, _ []string, keeper Keeper) ([]byte, sdk.Error) {
	oracles := keeper.GetOracles(ctx)
	if oracles == nil {
		oracles = make([]sdk.AccAddress, 0)
	}

	oraclesBz, err := codec.MarshalJSONIndent(keeper.cdc, oracles)
	if err != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return oraclesBz, nil
}
