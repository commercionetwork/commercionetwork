package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/encapsulated/bank/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier returns a new sdk.Keeper instance.
func NewQuerier(q sdk.Querier, k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryBlockedAccounts:
			return queryBlockedAccounts(ctx, req, k)

		default:
			return q(ctx, path, req)
		}
	}
}

// queryBlockedAccounts fetch an account's balance for the supplied height.
// Height and account address are passed as first and second path components respectively.
func queryBlockedAccounts(ctx sdk.Context, _ abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, k.GetBlockedAddresses(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
