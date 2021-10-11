package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func (k Keeper) Memberships(c context.Context, req *types.QueryMembershipsRequest) (*types.QueryMembershipsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	memberships := k.GetMemberships(ctx)
	return &types.QueryMembershipsResponse{Memberships: memberships}, nil
}

func (k Keeper) Membership(c context.Context, req *types.QueryMembershipRequest) (*types.QueryMembershipResponse, error) {
	address, err2 := sdk.AccAddressFromBech32(req.Address)
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, req.Address)
	}
	// Search the membership
	ctx := sdk.UnwrapSDKContext(c)
	membership, err := k.GetMembership(ctx, address)
	if err != nil {
		return nil, err
	}
	return &types.QueryMembershipResponse{Membership: &membership}, nil
}
