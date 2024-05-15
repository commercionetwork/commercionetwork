package keeper
/*
import (
//	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"

//	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

//	"github.com/commercionetwork/commercionetwork/x/did/types"
)

// NewQuerier is the module level router for state queries
// func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
// 	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
// 		switch path[0] {
// 		case types.QueryResolveIdentity:
// 			return queryGetLastIdentityOfAddress(ctx, path[1:], k, legacyQuerierCdc)
// 		case types.QueryResolveIdentityHistory:
// 			return queryGetIdentityHistoryOfAddress(ctx, path[1:], k, legacyQuerierCdc)
// 		default:
// 			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
// 		}
// 	}
// }

func queryGetLastIdentityOfAddress(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	identity, err := k.GetLastIdentityOfAddress(ctx, path[0])
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownAddress, err.Error())
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, identity)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetIdentityHistoryOfAddress(ctx sdk.Context, path []string, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	identities := k.GetIdentityHistoryOfAddress(ctx, path[0])

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, identities)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}
*/