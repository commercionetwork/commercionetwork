package keeper

import (
	"context"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) InviteUser(goCtx context.Context, msg *types.MsgInviteUser) (*types.MsgInviteUserResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	msgRecipient, _ := sdk.AccAddressFromBech32(msg.Recipient)
	if k.accountKeeper.GetAccount(ctx, msgRecipient) != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, "cannot invite existing user")
	}

	// Verify that the user that is inviting has already a membership
	msgSender, _ := sdk.AccAddressFromBech32(msg.Sender)
	if _, err := k.GetMembership(ctx, msgSender); err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("Cannot send an invitation without having a membership: %s", err.Error()))
	}

	// Try inviting the user
	if err := k.Invite(ctx, msgRecipient, msgSender); err != nil {
		return nil, err
	}
	//ctypes.EmitCommonEvents(ctx, msg.Sender)

	return &types.MsgInviteUserResponse{
		Status: "1",
	}, nil

}
