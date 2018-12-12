package commercioid

/**
 * This is the place to define which queries against application state users will be able to make.
 * Our nameservice module will expose two queries:
 *
 * • resolve: This takes a name and returns the value that is stored by the nameservice. This is similar to a DNS query.
 * • whois: This takes a name and returns the price, value, and owner of the name.
 *   Used for figuring out how much names cost when you want to buy them.
 */

import (
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the governance Querier
const (
	QueryRead = "resolve"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryRead:
			return queryRead(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("Unknown commercioid query endpoint")
		}
	}
}

// nolint: unparam
func queryRead(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) (res []byte, err sdk.Error) {
	did := path[0]

	whois := Identity{}

	whois.Value = keeper.GetIdentity(ctx, did)
	whois.Owner = keeper.GetOwner(ctx, did)

	bz, err2 := codec.MarshalJSONIndent(keeper.cdc, whois)
	if err2 != nil {
		panic("Could not marshal result to JSON")
	}

	return bz, nil
}

// Identity represents a did -> DDO lookup
// If your application needs some custom response types (Identity in this example), define them in this file.
type Identity struct {
	Value string         `json:"value"`
	Owner sdk.AccAddress `json:"owner"`
}
