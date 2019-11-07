package keeper_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_InviteUser_NoInvite(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)

	err := k.InviteUser(ctx, testUser, TestUser2)
	assert.Nil(t, err)

	var invite types.Invite
	accreditationBz := store.Get([]byte(types.InviteStorePrefix + testUser.String()))
	k.Cdc.MustUnmarshalBinaryBare(accreditationBz, &invite)

	assert.Equal(t, testUser, invite.User)
	assert.Equal(t, TestUser2, invite.Sender)
	assert.False(t, invite.Rewarded)
}

func TestKeeper_InviteUser_ExistentInvite(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	existingAccredit := types.Invite{User: testUser, Sender: TestUser2, Rewarded: false}

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.InviteStorePrefix+testUser.String()), k.Cdc.MustMarshalBinaryBare(existingAccredit))

	err := k.InviteUser(ctx, testUser, TestUser2)
	assert.NotNil(t, err)
}

func TestKeeper_GetInvite_NoInvite(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	_, found := k.GetInvite(ctx, testUser)
	assert.False(t, found)
}

func TestKeeper_GetInvite_ExistingInvite(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	expected := types.Invite{User: testUser, Sender: TestUser2, Rewarded: false}

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.InviteStorePrefix+testUser.String()), k.Cdc.MustMarshalBinaryBare(expected))

	stored, found := k.GetInvite(ctx, testUser)
	assert.True(t, found)
	assert.Equal(t, expected, stored)
}

func TestKeeper_GetInvites_EmptyList(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	invites := k.GetInvites(ctx)
	assert.Empty(t, invites)
}

func TestKeeper_GetInvites_NonEmptyList(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	inv1 := types.Invite{Sender: TestUser2, User: testUser, Rewarded: false}
	inv2 := types.Invite{Sender: TestUser2, User: TestUser2, Rewarded: false}

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.InviteStorePrefix+testUser.String()), k.Cdc.MustMarshalBinaryBare(inv1))
	store.Set([]byte(types.InviteStorePrefix+TestUser2.String()), k.Cdc.MustMarshalBinaryBare(inv2))

	invites := k.GetInvites(ctx)
	assert.Equal(t, 2, len(invites))
	assert.Contains(t, invites, inv1)
	assert.Contains(t, invites, inv2)
}
