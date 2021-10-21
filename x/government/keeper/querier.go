package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/government/types"
)

type QueryGovernmentResponse struct {
	GovernmentAddress string `json:"government_address"`
}

type QueryTumblerResponse struct {
	TumblerAddress string `json:"tumbler_address"`
}

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryGovernmentAddress:
			return queryGetGovernmentAddress(ctx, keeper, legacyQuerierCdc)
		case types.QueryTumblerAddress:
			return queryGetTumblerAddress(ctx, keeper, legacyQuerierCdc)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryGetGovernmentAddress(ctx sdk.Context, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	address := keeper.GetGovernmentAddress(ctx)

	r := QueryGovernmentResponse{
		GovernmentAddress: address.String(),
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, r)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetTumblerAddress(ctx sdk.Context, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	address := keeper.GetTumblerAddress(ctx)

	r := QueryTumblerResponse{
		TumblerAddress: address.String(),
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, r)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}
