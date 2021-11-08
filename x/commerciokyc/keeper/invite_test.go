package keeper_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func TestKeeper_InviteUser(t *testing.T) {
	tests := []struct {
		name              string
		inviterMembership string
		existingInvite    types.Invite
		invite            types.Invite
		expected          types.Invite
		error             string
	}{
		{
			name:           "Existing invitation returns error",
			existingInvite: types.NewInvite(testInviteSender, testUser, "bronze"),
			invite:         types.NewInvite(testUser2, testUser, "bronze"),
			expected:       types.NewInvite(testInviteSender, testUser, "bronze"),
			error:          fmt.Sprintf("%s has already been invited: unknown request", testUser),
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
			store.Set([]byte(types.InviteStorePrefix+testUser.String()), k.Cdc.MustMarshalBinaryBare(&test.existingInvite))
		}

		test_invite_User, _ := sdk.AccAddressFromBech32(test.invite.User)
		test_invite_Sender, _ := sdk.AccAddressFromBech32(test.invite.Sender)

		if test.inviterMembership != "" {
			err := k.AssignMembership(ctx, test_invite_Sender, test.inviterMembership, test_invite_User, testExpiration)
			require.NoError(t, err)
		}

		err := k.Invite(ctx, test_invite_User, test_invite_Sender)
		if test.error != "" {
			require.Equal(t, test.error, err.Error())
		} else {
			require.NoError(t, err)
		}

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
			store.Set([]byte(types.InviteStorePrefix+test.storedInvite.User), k.Cdc.MustMarshalBinaryBare(&test.storedInvite))
		}

		actual, found := k.GetInvite(ctx, test.user)
		require.Equal(t, test.expected, actual)
		require.Equal(t, test.shouldBeFound, found)
	}
}

func TestInvites_Equals(t *testing.T) {
	user, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	sender, _ := sdk.AccAddressFromBech32("cosmos1007jzaanx5kmqnn3akgype2jseawfj80dne9t6")
	invite := types.NewInvite(sender, user, "bronze")

	tests := []struct {
		name          string
		first         types.Invite
		second        types.Invite
		shouldBeEqual bool
	}{
		{
			name:          "Different sender returns false",
			first:         invite,
			second:        types.NewInvite(user, user, "bronze"),
			shouldBeEqual: false,
		},
		{
			name:          "Different user returns false",
			first:         invite,
			second:        types.NewInvite(sender, sender, "bronze"),
			shouldBeEqual: false,
		},
		{
			name:          "Different memebership returns false",
			first:         invite,
			second:        types.Invite{User: user.String(), Sender: sender.String(), SenderMembership: types.MembershipTypeGold},
			shouldBeEqual: false,
		},
		{
			name:          "Different rewarded returns false",
			first:         invite,
			second:        types.Invite{User: user.String(), Sender: sender.String(), Status: uint64(types.InviteStatusRewarded)},
			shouldBeEqual: false,
		},
		{
			name:          "Same data returns true",
			first:         invite,
			second:        invite,
			shouldBeEqual: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.shouldBeEqual, test.first.Equals(test.second))
		})
	}
}

func TestKeeper_GetInvites(t *testing.T) {
	storedList := []*types.Invite{}
	expectedList := []*types.Invite{}
	storedInvite := types.NewInvite(testInviteSender, testUser2, "bronze")
	storedInvite2 := types.NewInvite(testUser2, testUser, "bronze")
	storedList = append(storedList, &storedInvite, &storedInvite2)
	expectedList = append(expectedList, &storedInvite, &storedInvite2)

	tests := []struct {
		name     string
		stored   []*types.Invite
		expected []*types.Invite
	}{
		{
			name:     "Empty list is returned properly",
			stored:   []*types.Invite([]*types.Invite{}),
			expected: []*types.Invite([]*types.Invite{}),
		},
		{
			name:     "Existing list is returned properly",
			stored:   storedList,
			expected: expectedList,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			store := ctx.KVStore(k.StoreKey)

			for _, invite := range test.stored {
				store.Set([]byte(types.InviteStorePrefix+invite.User), k.Cdc.MustMarshalBinaryBare(invite))
			}

			actual := k.GetInvites(ctx)
			require.Equal(t, test.expected, actual)
		})
	}
}
