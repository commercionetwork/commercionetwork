package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryResolveDid:
			return queryResolveIdentity(ctx, path[1:], keeper)
		case types.QueryResolveDepositRequest:
			return queryResolveDepositRequest(ctx, path[1:], keeper)
		case types.QueryResolvePowerUpRequest:
			return queryResolvePowerUpRequest(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

// ------------------
// --- Identities
// ------------------

func queryResolveIdentity(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	address, err2 := sdk.AccAddressFromBech32(path[0])
	if err2 != nil {
		return nil, sdk.ErrInvalidAddress(path[0])
	}

	var response ResolveIdentityResponse
	response.Owner = address

	didDocument, err := keeper.GetDidDocumentByOwner(ctx, address)
	if err != nil {
		return nil, sdk.ErrUnknownAddress(err.Error())
	}

	response.DidDocument = &didDocument

	res, sErr := codec.MarshalJSONIndent(keeper.cdc, response)
	if sErr != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return res, nil
}

type ResolveIdentityResponse struct {
	Owner       sdk.AccAddress     `json:"owner"`
	DidDocument *types.DidDocument `json:"did_document"`
}

// -------------------
// --- Pairwise Did
//--------------------

func queryResolveDepositRequest(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {

	// Get the request
	request, err := keeper.GetDidDepositRequestByProof(ctx, path[0])
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, &request)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

func queryResolvePowerUpRequest(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {

	// Get the request
	request, err := keeper.GetPowerUpRequestByProof(ctx, path[0])
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}

	bz, sErr := codec.MarshalJSONIndent(keeper.cdc, &request)
	if sErr != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}
