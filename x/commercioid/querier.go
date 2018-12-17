package commercioid

/**
 * This is the place to define which queries against application state users will be able to make.
 * Our nameservice module will expose two queries:
 *
 * • resolve: This takes a name and returns the value that is stored by the nameservice. This is similar to a DNS query.
 * • whois: This takes a name and returns the price, value, and owner of the name.
 *   Used for figuring out how much names cost when you want to buy them.
 */

import (
	"commercio-network/types"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryResolveDid  = "resolve"
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

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, identityResult)
	if err2 != nil {
		panic("Could not marshal result to JSON")
	}

	return bz, nil
}

// Identity represents a Did -> Did Document lookup
// If your application needs some custom response types (Identity in this example), define them in this file.
type IdentityResult struct {
	Did          types.Did `json:"did"`
	DdoReference string    `json:"ddoReference"`
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
		panic("Could not marshal result to JSON")
	}

	return bz, nil
}

type ConnectionsResult struct {
	Did         types.Did   `json:"did"`
	Connections []types.Did `json:"connections"`
}
