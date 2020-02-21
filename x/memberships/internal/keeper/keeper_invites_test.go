package keeper_test

import (
	"fmt"
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestKeeper_InviteUser(t *testing.T) {
	tests := []struct {
		name              string
		inviterMembership string
		existingInvite    types.Invite
		invite            types.Invite
		expected          types.Invite
		error             error
	}{
		{
			name:           "Existing invitation returns error",
			existingInvite: types.NewInvite(testInviteSender, testUser, "bronze"),
			invite:         types.NewInvite(testUser2, testUser, "bronze"),
			expected:       types.NewInvite(testInviteSender, testUser, "bronze"),
			error:          sdk.ErrUnknownRequest(fmt.Sprintf("%s has already been invited", testUser)),
		},
		{
			name:              "New invite works properly",
			inviterMembership: "gold",
			invite:            types.NewInvite(testInviteSender, testUser, "gold"),
			expected:          types.NewInvite(testInviteSender, testUser, "gold"),
		},
	}

	for _, test := range tests {
		ctx, _, _, k := SetupTestInput()

		store := ctx.KVStore(k.StoreKey)

		if !test.existingInvite.Empty() {
			store.Set([]byte(types.InviteStorePrefix+testUser.String()), k.Cdc.MustMarshalBinaryBare(test.existingInvite))
		}

		if test.inviterMembership != "" {
			err := k.AssignMembership(ctx, test.invite.Sender, test.inviterMembership)
			require.NoError(t, err)
		}
		err := k.InviteUser(ctx, test.invite.User, test.invite.Sender)
		require.Equal(t, test.error, err)

		var invite types.Invite
		accreditationBz := store.Get([]byte(types.InviteStorePrefix + testUser.String()))
		k.Cdc.MustUnmarshalBinaryBare(accreditationBz, &invite)
		require.Equal(t, test.expected, invite)
	}
}

func TestKeeper_GetInvite(t *testing.T) {
	tests := []struct {
		name          string
		user          sdk.AccAddress
		storedInvite  types.Invite
		expected      types.Invite
		shouldBeFound bool
	}{
		{
			name:          "Non existing invite is handled properly",
			user:          testUser,
			expected:      types.Invite{},
			shouldBeFound: false,
		},
		{
			name:          "Existing invite is handled properly",
			user:          testUser,
			storedInvite:  types.NewInvite(testInviteSender, testUser, "bronze"),
			expected:      types.NewInvite(testInviteSender, testUser, "bronze"),
			shouldBeFound: true,
		},
		{
			name:          "Existing invite for different user returns empty",
			user:          testUser,
			storedInvite:  types.NewInvite(testInviteSender, testUser2, "bronze"),
			expected:      types.Invite{},
			shouldBeFound: false,
		},
	}

	for _, test := range tests {
		ctx, _, _, k := SetupTestInput()
		store := ctx.KVStore(k.StoreKey)

		if !test.storedInvite.Empty() {
			store.Set([]byte(types.InviteStorePrefix+test.storedInvite.User.String()), k.Cdc.MustMarshalBinaryBare(&test.storedInvite))
		}

		actual, found := k.GetInvite(ctx, test.user)
		require.Equal(t, test.expected, actual)
		require.Equal(t, test.shouldBeFound, found)
	}
}

func TestKeeper_GetInvites_EmptyList(t *testing.T) {
	tests := []struct {
		name     string
		stored   types.Invites
		expected types.Invites
	}{
		{
			name:     "Empty list is returned properly",
			stored:   types.Invites{},
			expected: types.Invites{},
		},
		{
			name: "Existing list is returned properly",
			stored: types.Invites{
				types.NewInvite(testInviteSender, testUser2, "bronze"),
				types.NewInvite(testUser2, testUser, "bronze"),
			},
			expected: types.Invites{
				types.NewInvite(testInviteSender, testUser2, "bronze"),
				types.NewInvite(testUser2, testUser, "bronze"),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			store := ctx.KVStore(k.StoreKey)

			for _, invite := range test.stored {
				store.Set([]byte(types.InviteStorePrefix+invite.User.String()), k.Cdc.MustMarshalBinaryBare(invite))
			}

			actual := k.GetInvites(ctx)
			require.Equal(t, test.expected, actual)
		})
	}
}
