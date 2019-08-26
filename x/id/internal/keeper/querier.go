package keeper

/**
 * This is the place to define which queries against application state users will be able to make.
 * Our module will expose two queries:
 * • resolve: This takes a Did and returns the associated Did Document reference
 * • connections: This takes a Did and returns the list of all the connections associated with it
 */

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	idtypes "github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

const (
	QueryResolveDid = "identities"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryResolveDid:
			return queryResolveIdentity(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown %s query endpoint", idtypes.ModuleName))
		}
	}
}

func queryResolveIdentity(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	address, err2 := sdk.AccAddressFromBech32(path[0])
	if err2 != nil {
		return nil, sdk.ErrInvalidAddress(path[0])
	}

	identityResult := IdentityResult{}
	identityResult.Did = address
	identityResult.DdoReference = keeper.GetDidDocumentUriByDid(ctx, address)

	if identityResult.DdoReference == "" {
		return nil, sdk.ErrUnknownRequest("No Did Document reference associated to the given address")
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.Cdc, identityResult)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

// Identity represents a Did -> Did Document lookup
type IdentityResult struct {
	Did          sdk.AccAddress `json:"did"`
	DdoReference string         `json:"did_document_uri"`
}
