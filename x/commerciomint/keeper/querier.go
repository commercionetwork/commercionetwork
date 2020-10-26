package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryGetEtps:
			return queryGetEtps(ctx, path[1:], keeper)
		case types.QueryConversionRate:
			return queryConversionRate(ctx, keeper)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryGetEtps(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	ownerAddr, _ := sdk.AccAddressFromBech32(path[0])
	cdps := keeper.GetAllPositionsOwnedBy(ctx, ownerAddr)
	cdpsBz, err := codec.MarshalJSONIndent(keeper.cdc, cdps)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "could not marshal result to JSON")
	}

	return cdpsBz, nil
}

func queryConversionRate(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	return codec.MarshalJSONIndent(keeper.cdc, keeper.GetConversionRate(ctx))
}
