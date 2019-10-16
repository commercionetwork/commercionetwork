package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (keeper Keeper) getInviteStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.InviteStorePrefix + user.String())
}

// InviteUser allows to set a given user as being invited by the given invite sender.
func (keeper Keeper) InviteUser(ctx sdk.Context, recipient, sender sdk.AccAddress) sdk.Error {
	store := ctx.KVStore(keeper.StoreKey)
	inviteKey := keeper.getInviteStoreKey(recipient)

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
	accreditationBz := keeper.cdc.MustMarshalBinaryBare(&accreditation)
	store.Set(inviteKey, accreditationBz)
	return nil
}

// GetInvite allows to get the invitation related to a user
func (keeper Keeper) GetInvite(ctx sdk.Context, user sdk.AccAddress) (invite types.Invite, found bool) {
	store := ctx.KVStore(keeper.StoreKey)
	key := keeper.getInviteStoreKey(user)

	if store.Has(key) {
		keeper.cdc.MustUnmarshalBinaryBare(store.Get(key), &invite)
		return invite, true
	}

	return types.Invite{}, false
}

// GetInvites returns all the invites ever made
func (keeper Keeper) GetInvites(ctx sdk.Context) (invites []types.Invite) {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.InviteStorePrefix))

	for ; iterator.Valid(); iterator.Next() {
		var invite types.Invite
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &invite)
		invites = append(invites, invite)
	}

	return
}

// SaveInvite allows to save the given invite inside the store
func (keeper Keeper) SaveInvite(ctx sdk.Context, invite types.Invite) {
	store := ctx.KVStore(keeper.StoreKey)
	store.Set(keeper.getInviteStoreKey(invite.User), keeper.cdc.MustMarshalBinaryBare(&invite))
}
