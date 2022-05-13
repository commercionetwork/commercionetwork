package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Memberships(c context.Context, req *types.QueryMembershipsRequest) (*types.QueryMembershipsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	memberships := []*types.Membership{}
	store := ctx.KVStore(k.storeKey)
	membershipsStore := prefix.NewStore(store, []byte(types.MembershipsStorageKey))

	pageRes, err := query.Paginate(
		membershipsStore,
		req.Pagination,
		func(key []byte, value []byte) error {
			membership := types.Membership{}
			k.cdc.MustUnmarshal(value, &membership)
			if IsValidMembership(ctx, *membership.ExpiryAt, membership.MembershipType) {
				memberships = append(memberships, &membership)
			}
			return nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryMembershipsResponse{Memberships: memberships, Pagination: pageRes}, nil
}

func (k Keeper) Membership(c context.Context, req *types.QueryMembershipRequest) (*types.QueryMembershipResponse, error) {
	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, req.Address)
	}
	// Search  membership
	ctx := sdk.UnwrapSDKContext(c)
	membership, err := k.GetMembership(ctx, address)
	if err != nil {
		return nil, err
	}
	return &types.QueryMembershipResponse{Membership: &membership}, nil
}
