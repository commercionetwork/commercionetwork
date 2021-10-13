package keeper

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

		store := ctx.KVStore(k.storeKey)

		if !test.existingInvite.Empty() {
			store.Set([]byte(types.InviteStorePrefix+testUser.String()), k.cdc.MustMarshalBinaryBare(&test.existingInvite))
		}

		if test.inviterMembership != "" {
			membership := types.Membership{
				Owner:          test.invite.Sender,
				TspAddress:     testTsp.String(),
				MembershipType: test.inviterMembership,
				ExpiryAt:       &testExpiration,
			}
			err := k.AssignMembership(ctx, membership)
			require.NoError(t, err)
		}
		test_invite_User, _ := sdk.AccAddressFromBech32(test.invite.User)
		test_invite_Sender, _ := sdk.AccAddressFromBech32(test.invite.Sender)

		err := k.Invite(ctx, test_invite_User, test_invite_Sender)
		if test.error != "" {
			require.Equal(t, test.error, err.Error())
		} else {
			require.NoError(t, err)
		}

		var invite types.Invite
		accreditationBz := store.Get([]byte(types.InviteStorePrefix + testUser.String()))
		k.cdc.MustUnmarshalBinaryBare(accreditationBz, &invite)
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
		store := ctx.KVStore(k.storeKey)

		if !test.storedInvite.Empty() {
			store.Set([]byte(types.InviteStorePrefix+test.storedInvite.User), k.cdc.MustMarshalBinaryBare(&test.storedInvite))
		}

		actual, found := k.GetInvite(ctx, test.user)
		require.Equal(t, test.expected, actual)
		require.Equal(t, test.shouldBeFound, found)
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
			store := ctx.KVStore(k.storeKey)

			for _, invite := range test.stored {
				store.Set([]byte(types.InviteStorePrefix+invite.User), k.cdc.MustMarshalBinaryBare(invite))
			}

			actual := k.GetInvites(ctx)
			require.Equal(t, test.expected, actual)
		})
	}
}
