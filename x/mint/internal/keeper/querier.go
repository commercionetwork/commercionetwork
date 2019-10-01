package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryGetCDP:
			return queryGetCDP(ctx, path[1:], keeper)
		case types.QueryGetCDPs:
			return queryGetCDPs(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryGetCDP(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	ownerAddr, _ := sdk.AccAddressFromBech32(path[0])
	timestamp := path[1]
	cdp := keeper.GetCDP(ctx, ownerAddr, timestamp)
	if cdp == nil {
		return nil, sdk.ErrUnknownRequest("couldn't find any cdp associated with the given address and timestamp")
	}

	cdpBz, err := codec.MarshalJSONIndent(keeper.cdc, &cdp)
	if err != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return cdpBz, nil
}

func queryGetCDPs(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	ownerAddr, _ := sdk.AccAddressFromBech32(path[0])
	cdps := keeper.GetCDPs(ctx, ownerAddr)
	if cdps == nil {
		cdps = make(types.CDPs, 0)
	}

	cdpsBz, err := codec.MarshalJSONIndent(keeper.cdc, cdps)
	if err != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return cdpsBz, nil
}
