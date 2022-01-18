package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/did/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryResolveDid:
			return queryResolveIdentity(ctx, path[1:], k, legacyQuerierCdc)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryResolveIdentity(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {

	identity, err := k.GetLastIdentityOfAddress(ctx, path[0])
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownAddress, err.Error())
	}

	response := types.QueryResolveIdentityResponse{
		Identity: identity,
	}

	bz, err2 := codec.MarshalJSONIndent(legacyQuerierCdc, response)
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}
