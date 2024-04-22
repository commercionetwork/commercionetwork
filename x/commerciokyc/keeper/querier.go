package keeper

import (
	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	errorsmod "cosmossdk.io/errors"
)

// NewQuerier returns a new sdk.Keeper instance.
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGetPoolFunds:
			return queryGetPoolFunds(ctx, req, k, legacyQuerierCdc)
		case types.QueryGetInvite:
			return queryGetInvite(ctx, path[1:], k, legacyQuerierCdc)
		case types.QueryGetInvites:
			return queryGetInvites(ctx, path[1:], k, legacyQuerierCdc)
		case types.QueryGetTrustedServiceProviders:
			return queryGetSigners(ctx, path[1:], k, legacyQuerierCdc)
		case types.QueryGetMembership:
			return queryGetMembership(ctx, path[1:], k, legacyQuerierCdc)
		case types.QueryGetMemberships:
			return queryGetMemberships(ctx, path[1:], k, legacyQuerierCdc)
		case types.QueryGetTspMemberships:
			return queryGetTspMemberships(ctx, path[1:], k, legacyQuerierCdc)
		default:
			return nil, errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}

func queryGetPoolFunds(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {

	poolFunds := k.GetPoolFunds(ctx)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, poolFunds)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}
	return bz, nil
}

func queryGetInvites(ctx sdk.Context, _ []string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {

	// Get the list of invites
	var invites []*types.Invite

	// Get all the invites
	invites = keeper.GetInvites(ctx)

	bz, err2 := codec.MarshalJSONIndent(legacyQuerierCdc, invites)
	if err2 != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

// queryGetMembership allows to retrieve the current membership of a user having a specified address
func queryGetInvite(ctx sdk.Context, path []string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, path[0])
	}
	// Search the membership
	invite, found := keeper.GetInvite(ctx, address)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "Could not find invitation")
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, invite)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

// queryGetSigners allows to retrieve the all current trust service providers
func queryGetSigners(ctx sdk.Context, _ []string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	signers := keeper.GetTrustedServiceProviders(ctx)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, signers)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

// queryGetMembership allows to retrieve the current membership of a user having a specified address
func queryGetMembership(ctx sdk.Context, path []string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	address, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, path[0])
	}
	// Search the membership
	membership, err := keeper.GetMembership(ctx, address)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, membership)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

// queryGetMemberships allows to retrieve all the current membership
func queryGetMemberships(ctx sdk.Context, _ []string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	// Extract all memberships
	var memberships []*types.Membership
	memberships = keeper.GetMemberships(ctx)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, memberships)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil

}

// queryGetTspMemberships allows to retrieve all the current membership
func queryGetTspMemberships(ctx sdk.Context, path []string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	tsp, err := sdk.AccAddressFromBech32(path[0])
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrInvalidAddress, path[0])
	}

	if !keeper.IsTrustedServiceProvider(ctx, tsp) {
		return nil, errorsmod.Wrap(sdkerrors.ErrNotFound, "Requested address is not a valid tsp")
	}

	// Search the membership
	memberships := keeper.GetTspMemberships(ctx, tsp)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, memberships)
	if err != nil {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil

}
