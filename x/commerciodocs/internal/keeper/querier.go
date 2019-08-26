package keeper

import (
	"fmt"
	"github.com/commercionetwork/commercionetwork/types"
	comtypes "github.com/commercionetwork/commercionetwork/x/commerciodocs/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case comtypes.QueryReceivedDocuments:
			return queryGetReceivedDocuments(ctx, path[1:], keeper)
		case comtypes.QuerySentDocuments:
			return queryGetSentDocuments(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown %s query endpoint", comtypes.ModuleName))
		}
	}
}

func queryGetReceivedDocuments(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	addr := path[0]
	address, _ := sdk.AccAddressFromBech32(addr)

	receivedResult := keeper.GetUserReceivedDocuments(ctx, address)
	if receivedResult == nil {
		receivedResult = make([]types.Document, 0)
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
		receivedResult = make([]types.Document, 0)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, receivedResult)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}
