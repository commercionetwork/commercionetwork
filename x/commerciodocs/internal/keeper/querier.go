package keeper

/**
 * This is the place to define which queries against application state users will be able to make.
 * Our commerciodocd module will expose:
 *
 * • TestMetadata: This takes a document TestReference and retrieve the associated TestMetadata.
 * • readers: This takes a document TestReference and return the list of all the users that are authorized to access it
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
	QueryMetadata = "TestMetadata"
	QuerySharing  = "readers"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryMetadata:
			return queryGetMetadata(ctx, path[1:], keeper)
		case QuerySharing:
			return queryGetAuthorized(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown commerciodocs query endpoint")
		}
	}
}

// ----------------------------------
// --- Get Metadata
// ----------------------------------

func queryGetMetadata(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	documentReference := path[0]

	identityResult := MetadataResult{}
	identityResult.Document = documentReference
	identityResult.Metadata = keeper.GetMetadata(ctx, documentReference)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, identityResult)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

// Metadata represents a TestReference -> TestMetadata lookup
type MetadataResult struct {
	Document string `json:"document_reference"`
	Metadata string `json:"metadata_reference"`
}

// ----------------------------------
// --- Get connections
// ----------------------------------

func queryGetAuthorized(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	documentReference := path[0]

	connectionsResult := AuthorizedResult{}
	connectionsResult.Document = documentReference
	connectionsResult.Readers = keeper.GetAuthorizedReaders(ctx, documentReference)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, connectionsResult)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

type AuthorizedResult struct {
	Document string      `json:"document_reference"`
	Readers  []types.Did `json:"authorized_readers"`
}
