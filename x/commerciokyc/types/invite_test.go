package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestInvite_Empty(t *testing.T) {
	address, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")

	tests := []struct {
		name          string
		invite        types.Invite
		expectedEmpty bool
	}{
		{
			name:          "Empty invite is returned empty",
			invite:        types.Invite{},
			expectedEmpty: true,
		},
		{
			name:          "Only sender invite is not empty",
			invite:        types.Invite{Sender: address},
			expectedEmpty: false,
		},
		{
			name:          "Only user invite is not empty",
			invite:        types.Invite{User: address},
			expectedEmpty: false,
		},
		{
			name:          "Only membership invite is not empty",
			invite:        types.Invite{SenderMembership: types.MembershipTypeGold},
			expectedEmpty: false,
		},
		{
			name:          "Only rewarded invite is not empty",
			invite:        types.Invite{Status: types.InviteStatusRewarded},
			expectedEmpty: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.expectedEmpty, test.invite.Empty())
		})
	}
}

func TestInvite_Equals(t *testing.T) {
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
			second:        types.Invite{User: user, Sender: sender, SenderMembership: types.MembershipTypeGold},
			shouldBeEqual: false,
		},
		{
			name:          "Different rewarded returns false",
			first:         invite,
			second:        types.Invite{User: user, Sender: sender, Status: types.InviteStatusRewarded},
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

/*func TestInvites_Equals(t *testing.T) {
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
			second:        types.Invite{User: user, Sender: sender, SenderMembership: types.MembershipTypeGold},
			shouldBeEqual: false,
		},
		{
			name:          "Different rewarded returns false",
			first:         invite,
			second:        types.Invite{User: user, Sender: sender, Status: types.InviteStatusRewarded},
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
}*/

func TestInvite_ValidateBasic(t *testing.T) {
	tests := []struct {
		name    string
		invite  types.Invite
		wantErr bool
	}{
		{
			"A valid rewarded invite",
			types.Invite{
				Sender: sender,
				User:   user,
				Status: types.InviteStatusRewarded,
			},
			false,
		},
		{
			"A valid pending invite",
			types.Invite{
				Sender: sender,
				User:   user,
				Status: types.InviteStatusPending,
			},
			false,
		},
		{
			"A valid invalid invite",
			types.Invite{
				Sender: sender,
				User:   user,
				Status: types.InviteStatusInvalid,
			},
			false,
		},
		{
			"An invite with invalid status",
			types.Invite{
				Sender: sender,
				User:   user,
				Status: 42,
			},
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			res := tt.invite.ValidateBasic()
			if tt.wantErr {
				require.Error(t, res)
			} else {
				require.NoError(t, res)
			}
		})
	}
}
