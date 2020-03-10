package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryResolveDid:
			return queryResolveIdentity(ctx, path[1:], keeper)
		case types.QueryResolvePowerUpRequest:
			return queryResolvePowerUpRequest(ctx, path[1:], keeper)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

// ------------------
// --- Identities
// ------------------

func queryResolveIdentity(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err error) {
	address, err2 := sdk.AccAddressFromBech32(path[0])
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, path[0])
	}

	var response ResolveIdentityResponse
	response.Owner = address

	didDocument, err := keeper.GetDidDocumentByOwner(ctx, address)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownAddress, err.Error())
	}

	response.DidDocument = &didDocument

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, response)
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
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
func queryResolvePowerUpRequest(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err error) {

	// Get the request
	request, err := keeper.GetPowerUpRequestByID(ctx, path[0])
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
	}

	bz, sErr := codec.MarshalJSONIndent(keeper.cdc, &request)
	if sErr != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("could not marshal result to JSON: %s", sErr.Error()))
	}

	return bz, nil
}
