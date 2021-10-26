package keeper

import (
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier returns a new sdk.Keeper instance.
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryGetPoolFunds:
			return queryGetPoolFunds(ctx, req, k, legacyQuerierCdc)
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
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}

func queryGetPoolFunds(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {

	poolFunds := k.GetPoolFunds(ctx)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, poolFunds)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
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
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

func queryGetSigners(ctx sdk.Context, _ []string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	signers := keeper.GetTrustedServiceProviders(ctx)
	if signers == nil {
		signers = make([]sdk.AccAddress, 0)
	}

	bz, err2 := codec.MarshalJSONIndent(legacyQuerierCdc, signers)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

// queryGetMembership allows to retrieve the current membership of a user having a specified address
func queryGetMembership(ctx sdk.Context, path []string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	address, err2 := sdk.AccAddressFromBech32(path[0])
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, path[0])
	}
	// Search the membership
	membership, err := keeper.GetMembership(ctx, address)
	if err != nil {
		return nil, err
	}

	bz, err2 := codec.MarshalJSONIndent(legacyQuerierCdc, membership)
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil
}

// queryGetMemberships allows to retrieve all the current membership
func queryGetMemberships(ctx sdk.Context, _ []string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	// Extract all memberships
	var memberships []*types.Membership
	memberships = keeper.GetMemberships(ctx)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, memberships)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil

}

// queryGetTspMemberships allows to retrieve all the current membership
func queryGetTspMemberships(ctx sdk.Context, path []string, keeper Keeper, legacyQuerierCdc *codec.LegacyAmino) (res []byte, err error) {
	tsp, err2 := sdk.AccAddressFromBech32(path[0])
	if err2 != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, path[0])
	}

	if !keeper.IsTrustedServiceProvider(ctx, tsp) {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Requested address is not a valid tsp")
	}

	// Search the membership
	memberships := keeper.GetTspMemberships(ctx, tsp)

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, memberships)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Could not marshal result to JSON")
	}

	return bz, nil

}
