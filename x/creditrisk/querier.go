package creditrisk

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"

	types "github.com/commercionetwork/commercionetwork/x/creditrisk/types"
)

func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryPool:
			return queryPoolFunds(ctx, keeper)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryPoolFunds(ctx sdk.Context, keeper Keeper) ([]byte, error) {
	return keeper.cdc.MarshalJSON(keeper.GetPoolFunds(ctx))
}
