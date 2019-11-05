package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) getInviteStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.InviteStorePrefix + user.String())
}

// InviteUser allows to set a given user as being invited by the given invite sender.
func (k Keeper) InviteUser(ctx sdk.Context, recipient, sender sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(k.storeKey)
	inviteKey := k.getInviteStoreKey(recipient)

	if store.Has(inviteKey) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("%s has already been invited", recipient.String()))
	}

	// Build the accreditation
	accreditation := types.Invite{
		Sender:   sender,
		User:     recipient,
		Rewarded: false,
	}

	// Save the accreditation
	accreditationBz := k.cdc.MustMarshalBinaryBare(&accreditation)
	store.Set(inviteKey, accreditationBz)
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

// GetInvites returns all the invites ever made
func (k Keeper) GetInvites(ctx sdk.Context) (invites []types.Invite) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.InviteStorePrefix))

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var invite types.Invite
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &invite)
		invites = append(invites, invite)
	}

	return
}

// SaveInvite allows to save the given invite inside the store
func (k Keeper) SaveInvite(ctx sdk.Context, invite types.Invite) {
	store := ctx.KVStore(k.storeKey)
	store.Set(k.getInviteStoreKey(invite.User), k.cdc.MustMarshalBinaryBare(&invite))
}
