package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
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
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

type InviteResponse struct {
	Recipient sdk.AccAddress `json:"recipient"`
	Sender    sdk.AccAddress `json:"sender"`
}

func queryGetInvites(ctx sdk.Context, path []string, keeper Keeper) ([]byte, sdk.Error) {

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

	if invites == nil {
		invites = make([]types.Invite, 0)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.Cdc, invites)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetSigners(ctx sdk.Context, _ []string, keeper Keeper) (res []byte, err sdk.Error) {
	signers := keeper.GetTrustedServiceProviders(ctx)
	if signers == nil {
		signers = make([]sdk.AccAddress, 0)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.Cdc, signers)
	if err2 != nil {
		return nil, sdk.ErrUnknownRequest("Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetPoolFunds(ctx sdk.Context, _ []string, keeper Keeper) (res []byte, err sdk.Error) {
	value := keeper.GetPoolFunds(ctx)
	if value == nil {
		value = make([]sdk.Coin, 0)
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.Cdc, value)
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
