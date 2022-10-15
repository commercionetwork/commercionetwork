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

func (k Keeper) Invites(c context.Context, req *types.QueryInvitesRequest) (*types.QueryInvitesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	invites := []*types.Invite{}
	store := ctx.KVStore(k.storeKey)
	invitesStore := prefix.NewStore(store, []byte(types.InviteStorePrefix))

	pageRes, err := query.Paginate(
		invitesStore,
		req.Pagination,
		func(key []byte, value []byte) error {
			invite := types.Invite{}
			k.cdc.MustUnmarshal(value, &invite)
			invites = append(invites, &invite)
			return nil
		},
	)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryInvitesResponse{Invites: invites, Pagination: pageRes}, nil
}

func (k Keeper) Invite(c context.Context, req *types.QueryInviteRequest) (*types.QueryInviteResponse, error) {
	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, req.Address)
	}
	// Search the invite
	ctx := sdk.UnwrapSDKContext(c)
	invite, found := k.GetInvite(ctx, address)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Could not find invitation")
	}
	return &types.QueryInviteResponse{Invite: &invite}, nil
}
