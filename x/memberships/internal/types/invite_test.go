package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/stretchr/testify/require"

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
			name:          "Only rewarded invite is not empty",
			invite:        types.Invite{Rewarded: true},
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
	invite := types.NewInvite(sender, user)

	tests := []struct {
		name          string
		first         types.Invite
		second        types.Invite
		shouldBeEqual bool
	}{
		{
			name:          "Different sender returns false",
			first:         invite,
			second:        types.NewInvite(user, user),
			shouldBeEqual: false,
		},
		{
			name:          "Different user returns false",
			first:         invite,
			second:        types.NewInvite(sender, sender),
			shouldBeEqual: false,
		},
		{
			name:          "Different rewarded returns false",
			first:         invite,
			second:        types.Invite{User: user, Sender: sender, Rewarded: true},
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
