package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) InviteUser(goCtx context.Context, msg *types.MsgInviteUser) (*types.MsgInviteUserResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.accountKeeper.GetAccount(ctx, sdk.AccAddress(msg.Recipient)) != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, "cannot invite existing user")
	}

	// Verify that the user that is inviting has already a membership
	if _, err := k.GetMembership(ctx, sdk.AccAddress(msg.Sender)); err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, "Cannot send an invitation without having a membership")
	}

	// Try inviting the user
	if err := k.Invite(ctx, sdk.AccAddress(msg.Recipient), sdk.AccAddress(msg.Sender)); err != nil {
		return nil, err
	}
	//ctypes.EmitCommonEvents(ctx, msg.Sender)

	return &types.MsgInviteUserResponse{
		Status: "1",
	}, nil

}
