package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func (k Keeper) Invites(c context.Context, req *types.QueryInvitesRequest) (*types.QueryInvitesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	invites := k.GetInvites(ctx)
	return &types.QueryInvitesResponse{Invites: invites}, nil
}
