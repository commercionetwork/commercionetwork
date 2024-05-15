package keeper

import (
	"context"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	errors "cosmossdk.io/errors"
)

// InviteUser handle message MsgInviteUser
func (k msgServer) InviteUser(goCtx context.Context, msg *types.MsgInviteUser) (*types.MsgInviteUserResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify that the user that is invited is not present on the chain
	msgRecipient, _ := sdk.AccAddressFromBech32(msg.Recipient)
	if k.accountKeeper.GetAccount(ctx, msgRecipient) != nil {
		return nil, errors.Wrap(sdkErr.ErrUnauthorized, "cannot invite existing user")
	}

	msgSender, _ := sdk.AccAddressFromBech32(msg.Sender)

	// Try inviting the user
	if err := k.SetInvite(ctx, msgRecipient, msgSender); err != nil {
		return nil, err
	}
	ctypes.EmitCommonEvents(ctx, msg.Sender)
	return &types.MsgInviteUserResponse{
		Status: "1",
	}, nil

}
