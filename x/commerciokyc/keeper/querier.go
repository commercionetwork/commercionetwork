package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"

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
			return queryGetMembership(ctx, path[1:], keeper)
		case types.QueryGetMemberships:
			return queryGetMemberships(ctx, path[1:], keeper)
		case types.QueryGetTspMemberships:
			return queryGetTspMemberships(ctx, path[1:], keeper)
		default:
			return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Unknown %s query endpoint", types.ModuleName))
		}
	}
}

func queryGetInvites(ctx sdk.Context, _ []string, keeper Keeper) ([]byte, error) {

	// Get the list of invites
	var invites []types.Invite

	// Get all the invites
	invites = keeper.GetInvites(ctx)

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

// queryGetMembership allows to retrieve the current membership of a user having a specified address
func queryGetMembership(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err error) {
	address, err2 := sdk.AccAddressFromBech32(path[0])
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, path[0])
	}
	// Search the membership
	membership, err := keeper.GetMembership(ctx, address)
	if err != nil {
		return nil, err
	}

	bz, err2 := codec.MarshalJSONIndent(keeper.Cdc, membership)
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

// queryGetMemberships allows to retrieve all the current membership
func queryGetMemberships(ctx sdk.Context, _ []string, keeper Keeper) (res []byte, err error) {
	// Extract all memberships
	var memberships types.Memberships
	memberships = keeper.GetMemberships(ctx)

	bz, err := codec.MarshalJSONIndent(keeper.Cdc, memberships)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil

}

// queryGetTspMemberships allows to retrieve all the current membership
func queryGetTspMemberships(ctx sdk.Context, path []string, keeper Keeper) (res []byte, err error) {
	tsp, err2 := sdk.AccAddressFromBech32(path[0])
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, path[0])
	}

	if !keeper.IsTrustedServiceProvider(ctx, tsp) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, "Requested address is not a valid tsp")
	}

	// Search the membership
	memberships := keeper.GetTspMemberships(ctx, tsp)

	bz, err := codec.MarshalJSONIndent(keeper.Cdc, memberships)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil

}
