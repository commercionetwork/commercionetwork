package keeper

import (
	"fmt"
	"strconv"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryGetCdp:
			return queryGetCdp(ctx, path[1:], keeper)
		case types.QueryGetCdps:
			return queryGetCdps(ctx, path[1:], keeper)
		case types.QueryCollateralRate:
			return queryCollateralRate(ctx, keeper)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryGetCdp(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	ownerAddr, _ := sdk.AccAddressFromBech32(path[0])
	timestamp, err := strconv.ParseInt(path[1], 10, 64)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("timestamp not valid: %s", path[1]))
	}
	cdp, found := keeper.GetPosition(ctx, ownerAddr, timestamp)
	if !found {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "couldn't find any cdp associated with the given address and timestamp")
	}

	cdpBz, err := codec.MarshalJSONIndent(keeper.cdc, &cdp)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return cdpBz, nil
}

func queryGetCdps(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	ownerAddr, _ := sdk.AccAddressFromBech32(path[0])
	cdps := keeper.GetAllPositionsOwnedBy(ctx, ownerAddr)
	cdpsBz, err := codec.MarshalJSONIndent(keeper.cdc, cdps)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return cdpsBz, nil
}

func queryCollateralRate(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	return codec.MarshalJSONIndent(keeper.cdc, keeper.GetCollateralRate(ctx))
}
