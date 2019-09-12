package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/government/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryGovernmentAddress:
			return queryGetReceivedDocuments(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryGetReceivedDocuments(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {
	address := keeper.GetGovernmentAddress(ctx)

	bz, err := codec.MarshalJSONIndent(keeper.cdc, address)
	if err != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}
