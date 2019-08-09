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
	"github.com/commercionetwork/commercionetwork/types"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryDocument           = "document"
	QuerySentDocuments      = "sent"
	QueryReceivedDocuments  = "received"
	QuerySharedDocsWithUser = "shared-with"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryReceivedDocuments:
			return queryGetReceivedDocuments(ctx, path[1:], keeper)
		case QuerySentDocuments:
			return queryGetSentDocuments(ctx, path[1:], keeper)
		/*
			case QueryDocument:
				return queryGetDocument(ctx, path[1:], keeper)
			case QuerySharedDocsWithUser:
				return queryGetSharedDocumentsWithUser(ctx, path[1:], keeper)
		*/
		default:
			return nil, sdk.ErrUnknownRequest("Unknown commerciodocs query endpoint")
		}
	}
}

// ----------------------------------
// --- Get Received documents
// ----------------------------------

func queryGetReceivedDocuments(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	addr := path[0]

	address, _ := sdk.AccAddressFromBech32(addr)

	var receivedResult []types.Document
	receivedResult = keeper.GetUserReceivedDocuments(ctx, address)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, receivedResult)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetSentDocuments(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	addr := path[0]

	address, _ := sdk.AccAddressFromBech32(addr)

	var receivedResult []types.Document
	receivedResult = keeper.GetUserSentDocuments(ctx, address)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, receivedResult)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}
