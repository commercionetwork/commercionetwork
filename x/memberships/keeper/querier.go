package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/memberships/types"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case types.QueryGetInvites:
			return queryGetInvites(ctx, path[1:], keeper)
		case types.QueryGetTrustedServiceProviders:
			return queryGetSigners(ctx, path[1:], keeper)
		case types.QueryGetPoolFunds:
			return queryGetPoolFunds(ctx, path[1:], keeper)
		case types.QueryGetMembership:
			return queryResolveMembership(ctx, path[1:], keeper)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

type InviteResponse struct {
	Recipient sdk.AccAddress `json:"recipient"`
	Sender    sdk.AccAddress `json:"sender"`
}

func queryGetInvites(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {

	// Get the list of invites
	var invites []types.Invite
	if len(path) == 0 {

		// Get all the invites
		invites = keeper.GetInvites(ctx)

	} else {

		// A user has been specified, get only his invites
		address, _ := sdk.AccAddressFromBech32(path[0])
		if invite, found := keeper.GetInvite(ctx, address); found {
			invites = []types.Invite{invite}
		}
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.Cdc, invites)
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetSigners(ctx sdk.Context, _ []string, keeper Keeper) (res []byte, err error) {
	signers := keeper.GetTrustedServiceProviders(ctx)
	if signers == nil {
		signers = make([]sdk.AccAddress, 0)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.Cdc, signers)
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetPoolFunds(ctx sdk.Context, _ []string, keeper Keeper) (res []byte, err error) {
	value := keeper.GetPoolFunds(ctx)
	if value == nil {
		value = make([]sdk.Coin, 0)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.Cdc, value)
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

// MembershipResult represents the data returned when a search is performed to know the membership of a given user
type MembershipResult struct {
	User           sdk.AccAddress `json:"user"`
	MembershipType string         `json:"membership_type"`
}

// queryResolveMembership allows to retrieve the current membership of a user having a specified address
func queryResolveMembership(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err error) {
	address, err2 := sdk.AccAddressFromBech32(path[0])
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, path[0])
	}

	// Create the result type
	result := MembershipResult{
		User: address,
	}

	// Search the membership
	membership, err := keeper.GetMembership(ctx, address)
	if err != nil {
		return nil, err
	}

	result.MembershipType = membership.MembershipType
	bz, err2 := codec.MarshalJSONIndent(keeper.Cdc, result)
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}
