package commercioid

/**
 * This is the place to define which queries against application state users will be able to make.
 * Our commercioid module will expose two queries:
 *
 * • resolve: This takes a Did and returns the associated Did Document reference
 * • connections: This takes a Did and returns the list of all the connections associated with it
 */

import (
	"commercio-network/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryResolveDid  = "identities"
	QueryConnections = "connections"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryResolveDid:
			return queryResolveIdentity(ctx, path[1:], keeper)
		case QueryConnections:
			return queryGetConnections(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown commercioid query endpoint")
		}
	}
}

// ----------------------------------
// --- Resolve identity
// ----------------------------------

// nolint: unparam
func queryResolveIdentity(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	did := types.Did(path[0])

	identityResult := IdentityResult{}
	identityResult.Did = did
	identityResult.DdoReference = keeper.GetDdoReferenceByDid(ctx, did)

	if identityResult.DdoReference == "" {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("No ddo reference related to given did"))
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, identityResult)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Could not marshal result to JSON"))
	}

	return bz, nil
}

// Identity represents a Did -> Did Document lookup
type IdentityResult struct {
	Did          types.Did `json:"did"`
	DdoReference string    `json:"ddo_reference"`
}

// ----------------------------------
// --- Get connections
// ----------------------------------

func queryGetConnections(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	did := types.Did(path[0])

	connectionsResult := ConnectionsResult{}
	connectionsResult.Did = did
	connectionsResult.Connections = keeper.GetConnections(ctx, did)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, connectionsResult)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Could not marshal result to JSON"))
	}

	return bz, nil
}

type ConnectionsResult struct {
	Did         types.Did   `json:"did"`
	Connections []types.Did `json:"connections"`
}
