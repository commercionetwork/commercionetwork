package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_InviteUser_NoInvite(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.storeKey)

	err := k.InviteUser(ctx, TestUser, TestUser2)
	assert.Nil(t, err)

	var invite types.Invite
	accreditationBz := store.Get(k.getInviteStoreKey(TestUser))
	cdc.MustUnmarshalBinaryBare(accreditationBz, &invite)

	assert.Equal(t, TestUser, invite.User)
	assert.Equal(t, TestUser2, invite.Sender)
	assert.False(t, invite.Rewarded)
}

func TestKeeper_InviteUser_ExistentInvite(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	existingAccredit := types.Invite{User: TestUser, Sender: TestUser2, Rewarded: false}

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getInviteStoreKey(TestUser), cdc.MustMarshalBinaryBare(existingAccredit))

	err := k.InviteUser(ctx, TestUser, TestUser2)
	assert.NotNil(t, err)
}

func TestKeeper_GetInvite_NoInvite(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()

	_, found := k.GetInvite(ctx, TestUser)
	assert.False(t, found)
}

func TestKeeper_GetInvite_ExistingInvite(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	expected := types.Invite{User: TestUser, Sender: TestUser2, Rewarded: false}

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getInviteStoreKey(TestUser), cdc.MustMarshalBinaryBare(expected))

	stored, found := k.GetInvite(ctx, TestUser)
	assert.True(t, found)
	assert.Equal(t, expected, stored)
}

func TestKeeper_GetInvites_EmptyList(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()

	invites := k.GetInvites(ctx)
	assert.Empty(t, invites)
}

func TestKeeper_GetInvites_NonEmptyList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	inv1 := types.Invite{Sender: TestUser2, User: TestUser, Rewarded: false}
	inv2 := types.Invite{Sender: TestUser2, User: TestUser2, Rewarded: false}

	store := ctx.KVStore(k.storeKey)
	store.Set(k.getInviteStoreKey(TestUser), cdc.MustMarshalBinaryBare(inv1))
	store.Set(k.getInviteStoreKey(TestUser2), cdc.MustMarshalBinaryBare(inv2))

	invites := k.GetInvites(ctx)
	assert.Equal(t, 2, len(invites))
	assert.Contains(t, invites, inv1)
	assert.Contains(t, invites, inv2)
}
