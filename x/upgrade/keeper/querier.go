package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/upgrade/types"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryCurrent:
			return queryCurrentUpgrade(ctx, k, legacyQuerierCdc)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint: %s", types.ModuleName, path[0]))
		}
	}
}

func queryCurrentUpgrade(ctx sdk.Context, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	plan, has := k.GetUpgradePlan(ctx)
	if !has {
		return nil, nil
	}

	res, err := codec.MarshalJSONIndent(legacyQuerierCdc, plan)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return res, nil
}
