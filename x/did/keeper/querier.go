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
func NewQuerier(keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryResolveDid:
			return queryResolveIdentity(ctx, path[1:], keeper, legacyQuerierCdc)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

// ------------------
// --- Identities
// ------------------

func queryResolveIdentity(ctx sdk.Context, path []string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {

	var response ResolveIdentityResponse
	response.Owner = path[0]

	didDocument, err := keeper.GetDidDocumentOfAddress(ctx, path[0])
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownAddress, err.Error())
	}

	response.DidDocument = &didDocument

	bz, err2 := codec.MarshalJSONIndent(legacyQuerierCdc, response)
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

type ResolveIdentityResponse struct {
	Owner       string             `json:"owner" swaggertype:"string" example:"did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf"`
	DidDocument *types.DidDocument `json:"did_document"`
}
