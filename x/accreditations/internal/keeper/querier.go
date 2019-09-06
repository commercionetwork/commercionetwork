package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryGetAccrediter:
			return queryGetAccrediter(ctx, path[1:], keeper)
		case types.QueryGetSigners:
			return queryGetSigners(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryGetAccrediter(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	addr := path[0]
	address, _ := sdk.AccAddressFromBech32(addr)

	accrediter := keeper.GetAccrediter(ctx, address)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, accrediter)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetSigners(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	signers := keeper.GetTrustworthySigners(ctx)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, signers)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}
