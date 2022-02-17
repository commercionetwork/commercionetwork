package keeper

import (
	"context"
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	eventDistributeRewardFail = "distribute_reward_fail"
)

// BuyMembership handle message MsgBuyMembership
func (k msgServer) BuyMembership(goCtx context.Context, msg *types.MsgBuyMembership) (*types.MsgBuyMembershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify invite exists
	msgBuyer, _ := sdk.AccAddressFromBech32(msg.Buyer)
	invite, found := k.GetInvite(ctx, msgBuyer)
	if !found {
		return &types.MsgBuyMembershipResponse{}, sdkErr.Wrap(sdkErr.ErrUnauthorized, "Cannot buy a membership without being invited")
	}

	// Verify invite status
	inviteStatus := types.InviteStatus(invite.Status)
	if inviteStatus == types.InviteStatusInvalid {
		return &types.MsgBuyMembershipResponse{}, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("invite for account %s has been marked as invalid previously, cannot continue", msg.Buyer))
	}

	// Forbidden black membership buying
	if msg.MembershipType == types.MembershipTypeBlack {
		return &types.MsgBuyMembershipResponse{}, sdkErr.Wrap(sdkErr.ErrUnauthorized, "cannot buy black membership")
	}

	membershipPrice := membershipCosts[msg.MembershipType] * 1000000 // Always multiply by one million
	membershipCost := sdk.NewCoins(sdk.NewInt64Coin(types.CreditsDenom, membershipPrice))

	govAddr := k.GovKeeper.GetGovernmentAddress(ctx)
	// TODO Not send coins but control if account has enough
	msgTsp, _ := sdk.AccAddressFromBech32(msg.Tsp)
	if err := k.bankKeeper.SendCoins(ctx, msgTsp, govAddr, membershipCost); err != nil {
		return &types.MsgBuyMembershipResponse{}, err
	}

	expirationAt := k.ComputeExpiryHeight(ctx.BlockTime())

	err := k.AssignMembership(
		ctx,
		msgBuyer,
		msg.MembershipType,
		msgTsp,
		expirationAt,
	)

	// If AssignMembership fail return coins to tsp
	// TODO: Resolve nested error and potential no return funds to tsp
	if err != nil {
		if errRet := k.bankKeeper.SendCoins(ctx, govAddr, msgTsp, membershipCost); errRet != nil {
			return &types.MsgBuyMembershipResponse{}, errRet
		}
		return &types.MsgBuyMembershipResponse{}, err
	}

	// Give the reward to the invitee
	// Emits events if error occours. No transaction error
	if err := k.DistributeReward(ctx, invite); err != nil {
		// Emits events
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			eventDistributeRewardFail,
			sdk.NewAttribute("error", err.Error()),
		))
	}

	return &types.MsgBuyMembershipResponse{
		ExpiryAt: expirationAt.String(),
	}, nil
}

// RemoveMembership allows to handle a MsgRemoveMembership message.
// It checks that whoever sent the message is actually the government and remove membership
func (k msgServer) RemoveMembership(goCtx context.Context, msg *types.MsgRemoveMembership) (*types.MsgRemoveMembershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	govAddr := k.GovKeeper.GetGovernmentAddress(ctx)
	if !govAddr.Equals(msg.GetSigners()[0]) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownAddress,
			fmt.Sprintf("%s is government address and %s is not a government address", govAddr.String(), msg.Government),
		)
	}
	subscriber, _ := sdk.AccAddressFromBech32(msg.Subscriber)
	err := k.DeleteMembership(ctx, subscriber)
	// TODO emits events
	//ctypes.EmitCommonEvents(ctx, msg.Government)

	return &types.MsgRemoveMembershipResponse{
		Subscriber: msg.Subscriber,
	}, err
}

// SetMembership handles MsgSetMembership messages.
// It checks that whoever sent the message is actually the government, assigns the membership and then
// distribute the reward to the inviter.
// If the user isn't invited already, an invite will be created.
func (k msgServer) SetMembership(goCtx context.Context, msg *types.MsgSetMembership) (*types.MsgSetMembershipResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	govAddr := k.GovKeeper.GetGovernmentAddress(ctx)
	if !govAddr.Equals(msg.GetSigners()[0]) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownAddress,
			fmt.Sprintf("%s is government address and %s is not a government address", govAddr.String(), msg.Government),
		)
	}

	if !types.IsMembershipTypeValid(msg.NewMembership) {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("invalid membership type: %s", msg.NewMembership))
	}
	subscriber, _ := sdk.AccAddressFromBech32(msg.Subscriber)

	invite, err := k.governmentInvitesUser(ctx, subscriber)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnauthorized, fmt.Sprintf("government could not invite user: %s", err.Error()))
	}

	expiredAt := k.ComputeExpiryHeight(ctx.BlockTime())

	err = k.AssignMembership(ctx, subscriber, msg.NewMembership, govAddr, expiredAt)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest,
			fmt.Sprintf("could not assign membership to user %s: %s", msg.Subscriber, err.Error()),
		)
	}

	err = k.DistributeReward(ctx, invite)
	if err != nil {
		return nil, sdkErr.Wrap(sdkErr.ErrUnknownRequest,
			fmt.Sprintf("could not distribute membership reward to user %s: %s", invite.Sender, err.Error()),
		)
	}
	// TODO emits events
	// ctypes.EmitCommonEvents(ctx, msg.Government)
	return &types.MsgSetMembershipResponse{}, err

}

// governmentInvitesUser makes government invite an user if it isn't already invited and validated.
// This function is used when there's the need to assign an arbitrary membership to a given user.
func (k msgServer) governmentInvitesUser(ctx sdk.Context, user sdk.AccAddress) (types.Invite, error) {
	govAddr := k.GovKeeper.GetGovernmentAddress(ctx)

	// check the user has already been invited
	// if there's an invite, save a credential for it,
	// this way invited, but non-verified users will be able to receive a membership
	invite, found := k.GetInvite(ctx, user)
	if found {
		return invite, nil
	}
	_ = govAddr
	// otherwise, create an invite from the government
	// TODO create invite
	err := k.SetInvite(ctx, user, govAddr)
	if err != nil {
		return types.Invite{}, err
	}

	// get the invite again, mark it as rewarded, and return it
	invite, found = k.GetInvite(ctx, user)
	if !found {
		return types.Invite{}, fmt.Errorf("invite from government created correctly, but invite lookup failed")
	}
	invite.Status = uint64(types.InviteStatusRewarded)
	k.SaveInvite(ctx, invite)

	return invite, nil
}
