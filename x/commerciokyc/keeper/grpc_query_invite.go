package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func (k Keeper) Invites(c context.Context, req *types.QueryInvitesRequest) (*types.QueryInvitesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	invites := k.GetInvites(ctx)
	return &types.QueryInvitesResponse{Invites: invites}, nil
}

func (k Keeper) Invite(c context.Context, req *types.QueryInviteRequest) (*types.QueryInviteResponse, error) {
	address, err2 := sdk.AccAddressFromBech32(req.Address)
	if err2 != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrInvalidAddress, req.Address)
	}
	// Search the invite
	ctx := sdk.UnwrapSDKContext(c)
	invite, found := k.GetInvite(ctx, address)
	if !found {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Could not find invitation")
	}
	return &types.QueryInviteResponse{Invite: &invite}, nil
}
