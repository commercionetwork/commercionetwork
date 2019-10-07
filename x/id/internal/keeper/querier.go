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

func queryResolveIdentity(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	address, err2 := sdk.AccAddressFromBech32(path[0])
	if err2 != nil {
		return nil, sdk.ErrInvalidAddress(path[0])
	}

	var response ResolveIdentityResponse
	response.Owner = address

	// Get the Did Document
	if didDocument, found := keeper.GetDidDocumentByOwner(ctx, address); found {
		response.DidDocument = &didDocument
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, response)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

type ResolveIdentityResponse struct {
	Owner       sdk.AccAddress     `json:"owner"`
	DidDocument *types.DidDocument `json:"did_document"`
}

// -------------------
// --- Pairwise Did
//--------------------

func queryResolveDepositRequest(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {

	// Get the request
	request, found := keeper.GetDidDepositRequestByProof(ctx, path[0])
	if !found {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Deposit request with proof %s not found", path[0]))
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, &request)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

func queryResolvePowerUpRequest(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {

	// Get the request
	request, found := keeper.GetPowerUpRequestByProof(ctx, path[0])
	if !found {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Poer up request with proof %s not found", path[0]))
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, &request)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}
