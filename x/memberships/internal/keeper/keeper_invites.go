package keeper

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) getInviteStoreKey(user sdk.AccAddress) []byte {
	return []byte(types.InviteStorePrefix + user.String())
}

// InviteUser allows to set a given user as being invited by the given invite sender.
func (k Keeper) InviteUser(ctx sdk.Context, recipient, sender sdk.AccAddress) error {
	store := ctx.KVStore(k.StoreKey)
	inviteKey := k.getInviteStoreKey(recipient)

	if store.Has(inviteKey) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("%s has already been invited", recipient))
	}

	// Build and save the invite
	accreditation := types.NewInvite(sender, recipient)
	store.Set(inviteKey, k.Cdc.MustMarshalBinaryBare(&accreditation))
	return nil
}

// GetInvite allows to get the invitation related to a user
func (k Keeper) GetInvite(ctx sdk.Context, user sdk.AccAddress) (invite types.Invite, found bool) {
	store := ctx.KVStore(k.StoreKey)
	key := k.getInviteStoreKey(user)

	if store.Has(key) {
		k.Cdc.MustUnmarshalBinaryBare(store.Get(key), &invite)
		return invite, true
	}

	return types.Invite{}, false
}

// GetInvites returns all the invites ever made
func (k Keeper) GetInvites(ctx sdk.Context) types.Invites {
	store := ctx.KVStore(k.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.InviteStorePrefix))

	invites := types.Invites{}
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var invite types.Invite
		k.Cdc.MustUnmarshalBinaryBare(iterator.Value(), &invite)
		invites = append(invites, invite)
	}

	return invites
}

// SaveInvite allows to save the given invite inside the store
func (k Keeper) SaveInvite(ctx sdk.Context, invite types.Invite) {
	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getInviteStoreKey(invite.User), k.Cdc.MustMarshalBinaryBare(&invite))
}
