package keeper

import (
	"fmt"

	doctypes "github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case doctypes.QueryReceivedDocuments:
			return queryGetReceivedDocuments(ctx, path[1:], keeper)
		case doctypes.QuerySentDocuments:
			return queryGetSentDocuments(ctx, path[1:], keeper)
		case doctypes.QueryReceipts:
			return queryGetReceivedDocsReceipts(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown %s query endpoint", doctypes.ModuleName))
		}
	}
}

func queryGetReceivedDocuments(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	addr := path[0]
	address, _ := sdk.AccAddressFromBech32(addr)

	receivedResult := keeper.GetUserReceivedDocuments(ctx, address)
	if receivedResult == nil {
		receivedResult = make([]doctypes.Document, 0)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, receivedResult)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetSentDocuments(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	addr := path[0]
	address, _ := sdk.AccAddressFromBech32(addr)

	receivedResult := keeper.GetUserSentDocuments(ctx, address)
	if receivedResult == nil {
		receivedResult = make([]doctypes.Document, 0)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, receivedResult)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetReceivedDocsReceipts(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	addr := path[0]
	address, _ := sdk.AccAddressFromBech32(addr)

	var uuid string
	if len(path) == 2 {
		uuid = path[1]
	}

	var receipts []doctypes.DocumentReceipt
	var receipt doctypes.DocumentReceipt
	var bz []byte
	var err2 error

	//If user wants all his receipts
	if uuid == "" {
		receipts = keeper.GetUserReceivedReceipts(ctx, address)

		if receipts == nil {
			receipts = make([]doctypes.DocumentReceipt, 0)
		}
		bz, err2 = codec.MarshalJSONIndent(keeper.cdc, &receipts)

	} else {
		receipt = keeper.GetReceiptByDocumentUuid(ctx, address, uuid)
		bz, err2 = codec.MarshalJSONIndent(keeper.cdc, &receipt)
	}

	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}
