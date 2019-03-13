package commercioauth

/**
 * This is the place to define which queries against application state users will be able to make.
 * Our commercioauth module will expose:
 *
 * • account: This takes an address and returns the details of the account associated with that address
 */

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryAccount = "account"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryAccount:
			return queryAddress(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown commerciodocs query endpoint")
		}
	}
}

// ----------------------------------
// --- Read account details
// ----------------------------------

func queryAddress(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	address := path[0]

	var result interface{}

	if address == "list" {
		// List all the accounts
		result = keeper.ListAccounts(ctx)
	} else {
		// Read the account details
		result, err = keeper.GetAccount(ctx, address)
		if err != nil {
			panic(err)
		}
	}

	// Serialize the account details as a JSON object
	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, result)
	if err2 != nil {
		panic("Could not marshal result to JSON")
	}

	return bz, nil
}
