package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryGetMembership:
			return queryResolveMembership(ctx, path[1:], keeper)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

// queryResolveMembership allows to retrieve the current membership of a user having a specified address
func queryResolveMembership(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err sdk.Error) {
	address, err2 := sdk.AccAddressFromBech32(path[0])
	if err2 != nil {
		return nil, sdk.ErrInvalidAddress(path[0])
	}

	// Create the result type
	result := MembershipResult{
		User: address,
	}

	// Search the membership
	if membership, found := keeper.GetMembership(ctx, address); found {
		result.MembershipType = keeper.GetMembershipType(membership)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.Cdc, result)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

// MembershipResult represents the data returned when a search is performed to know the membership of a given user
type MembershipResult struct {
	User           sdk.AccAddress `json:"user"`
	MembershipType string         `json:"membership_type"`
}
