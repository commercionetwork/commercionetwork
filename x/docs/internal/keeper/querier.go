package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryReceivedDocuments:
			return queryGetReceivedDocuments(ctx, path[1:], keeper)
		case types.QuerySentDocuments:
			return queryGetSentDocuments(ctx, path[1:], keeper)
		case types.QueryReceivedReceipts:
			return queryGetReceivedDocsReceipts(ctx, path[1:], keeper)
		case types.QuerySentReceipts:
			return queryGetSentDocsReceipts(ctx, path[1:], keeper)
		case types.QuerySupportedMetadataSchemes:
			return querySupportedMetadataSchemes(ctx, path[1:], keeper)
		case types.QueryTrustedMetadataProposers:
			return queryTrustedMetadataProposers(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

// ----------------------------------
// --- Documents
// ----------------------------------

func queryGetReceivedDocuments(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	addr := path[0]
	address, _ := sdk.AccAddressFromBech32(addr)

	receivedResult, err := keeper.GetUserReceivedDocuments(ctx, address)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}

	if receivedResult == nil {
		receivedResult = make([]types.Document, 0)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, receivedResult)
	if err != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetSentDocuments(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	addr := path[0]
	address, _ := sdk.AccAddressFromBech32(addr)

	receivedResult, err := keeper.GetUserSentDocuments(ctx, address)
	if err != nil {
		return nil, sdk.ErrUnknownRequest(err.Error())
	}

	if receivedResult == nil {
		receivedResult = make([]types.Document, 0)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, receivedResult)
	if err != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

// ----------------------------------
// --- Documents receipts
// ----------------------------------

func queryGetReceivedDocsReceipts(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	addr := path[0]
	address, _ := sdk.AccAddressFromBech32(addr)

	var uuid string
	if len(path) == 2 {
		uuid = path[1]
	}

	var receipts []types.DocumentReceipt

	//If user wants all his receipts
	if uuid == "" {
		receipts = keeper.GetUserReceivedReceipts(ctx, address)
	} else {
		receipts = keeper.GetUserReceivedReceiptsForDocument(ctx, address, uuid)
	}

	if receipts == nil {
		receipts = make([]types.DocumentReceipt, 0)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, &receipts)

	if err != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetSentDocsReceipts(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	addr := path[0]
	address, err := sdk.AccAddressFromBech32(addr)

	if err != nil {
		return nil, sdk.ErrInvalidAddress(addr)
	}

	receipts := keeper.GetUserSentReceipts(ctx, address)
	if receipts == nil {
		receipts = make([]types.DocumentReceipt, 0)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, &receipts)

	if err != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

// ----------------------------------
// --- Document metadata schemes
// ----------------------------------

func querySupportedMetadataSchemes(ctx sdk.Context, _ []string, keeper Keeper) ([]byte, sdk.Error) {
	schemes := keeper.GetSupportedMetadataSchemes(ctx)
	if schemes == nil {
		schemes = make([]types.MetadataSchema, 0)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, &schemes)

	if err != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

// -----------------------------------------
// --- Document metadata schemes proposers
// -----------------------------------------

func queryTrustedMetadataProposers(ctx sdk.Context, _ []string, keeper Keeper) ([]byte, sdk.Error) {
	proposers := keeper.GetTrustedSchemaProposers(ctx)
	if proposers == nil {
		proposers = make([]sdk.AccAddress, 0)
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, &proposers)

	if err != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}
