package keeper

import (
	"fmt"
	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryCurrent:
			return queryCurrentUpgrade(ctx, keeper)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryCurrentUpgrade(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	plan, has := keeper.GetUpgradePlan(ctx)
	if !has {
		return nil, nil
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, plan)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return res, nil
}
