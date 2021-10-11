package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func (k Keeper) getInviteStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.InviteStorePrefix + user.String())
}

// InviteUser allows to set a given user as being invited by the given invite sender.
func (k Keeper) Invite(ctx sdk.Context, recipient, sender sdk.AccAddress) error {
	store := ctx.KVStore(k.storeKey)
	inviteKey := k.getInviteStoreKey(recipient)

	if store.Has(inviteKey) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("%s has already been invited", recipient))
	}

	inviterMembership, err := k.GetMembership(ctx, sender)
	if err != nil {
		return err
	}
	// Build and save the invite
	accreditation := types.NewInvite(sender, recipient, inviterMembership.MembershipType)
	store.Set(inviteKey, k.cdc.MustMarshalBinaryBare(&accreditation))

	// TODO emits events
	/*
		ctx.EventManager().EmitEvent(sdk.NewEvent(
			eventInvite,
			sdk.NewAttribute("recipient", recipient.String()),
			sdk.NewAttribute("sender", sender.String()),
			sdk.NewAttribute("sender_membership_type", inviterMembership.MembershipType), // Maybe
		))
	*/

	return nil
}

// GetInvite allows to get the invitation related to a user
func (k Keeper) GetInvite(ctx sdk.Context, user sdk.AccAddress) (invite types.Invite, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := k.getInviteStoreKey(user)

	if store.Has(key) {
		k.cdc.MustUnmarshalBinaryBare(store.Get(key), &invite)
		return invite, true
	}

	return types.Invite{}, false
}

// InvitesIterator returns an Iterator which iterates over all the invites.
func (k Keeper) InvitesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte(types.InviteStorePrefix))
}

// GetInvites returns all the invites ever made
func (k Keeper) GetInvites(ctx sdk.Context) []*types.Invite {
	iterator := k.InvitesIterator(ctx)

	invites := []*types.Invite{}
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var invite types.Invite
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &invite)
		invites = append(invites, &invite)
	}

	return invites
}

// SaveInvite allows to save the given invite inside the store
func (k Keeper) SaveInvite(ctx sdk.Context, invite types.Invite) {
	store := ctx.KVStore(k.storeKey)
	store.Set(k.getInviteStoreKey(sdk.AccAddress(invite.User)), k.cdc.MustMarshalBinaryBare(&invite))
}
