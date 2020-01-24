package types_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMembership_Equals(t *testing.T) {
	owner, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	membership := types.NewMembership(types.MembershipTypeBronze, owner)

	tests := []struct {
		name          string
		first         types.Membership
		second        types.Membership
		shouldBeEqual bool
	}{
		{
			name:          "Different type returns false",
			first:         membership,
			second:        types.NewMembership(types.MembershipTypeSilver, membership.Owner),
			shouldBeEqual: false,
		},
		{
			name:          "Different type returns false",
			first:         membership,
			second:        types.NewMembership(types.MembershipTypeGold, membership.Owner),
			shouldBeEqual: false,
		},
		{
			name:          "Different type returns false",
			first:         membership,
			second:        types.NewMembership(types.MembershipTypeBlack, membership.Owner),
			shouldBeEqual: false,
		},
		{
			name:          "Different owner returns false",
			first:         membership,
			second:        types.NewMembership(membership.MembershipType, sdk.AccAddress{}),
			shouldBeEqual: false,
		},
		{
			name:          "Same data returns true",
			first:         membership,
			second:        membership,
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

func TestIsMembershipTypeValid(t *testing.T) {
	tests := []struct {
		membershipType string
		shouldBeValid  bool
	}{
		{membershipType: types.MembershipTypeBronze, shouldBeValid: true},
		{membershipType: types.MembershipTypeSilver, shouldBeValid: true},
		{membershipType: types.MembershipTypeGold, shouldBeValid: true},
		{membershipType: types.MembershipTypeBlack, shouldBeValid: true},
		{membershipType: strings.ToUpper(types.MembershipTypeBronze), shouldBeValid: false},
	}

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%s is valid", test.membershipType), func(t *testing.T) {
			require.Equal(t, test.shouldBeValid, types.IsMembershipTypeValid(test.membershipType))
		})
	}
}

func TestCanUpgrade(t *testing.T) {
	tests := []struct {
		first              string
		second             string
		shouldBeUpgradable bool
	}{
		{first: types.MembershipTypeBronze, second: types.MembershipTypeBronze, shouldBeUpgradable: false},
		{first: types.MembershipTypeBronze, second: types.MembershipTypeSilver, shouldBeUpgradable: true},
		{first: types.MembershipTypeBronze, second: types.MembershipTypeGold, shouldBeUpgradable: true},
		{first: types.MembershipTypeBronze, second: types.MembershipTypeBlack, shouldBeUpgradable: true},

		{first: types.MembershipTypeSilver, second: types.MembershipTypeBronze, shouldBeUpgradable: false},
		{first: types.MembershipTypeSilver, second: types.MembershipTypeSilver, shouldBeUpgradable: false},
		{first: types.MembershipTypeSilver, second: types.MembershipTypeGold, shouldBeUpgradable: true},
		{first: types.MembershipTypeSilver, second: types.MembershipTypeBlack, shouldBeUpgradable: true},

		{first: types.MembershipTypeGold, second: types.MembershipTypeBronze, shouldBeUpgradable: false},
		{first: types.MembershipTypeGold, second: types.MembershipTypeSilver, shouldBeUpgradable: false},
		{first: types.MembershipTypeGold, second: types.MembershipTypeGold, shouldBeUpgradable: false},
		{first: types.MembershipTypeGold, second: types.MembershipTypeBlack, shouldBeUpgradable: true},

		{first: types.MembershipTypeBlack, second: types.MembershipTypeBronze, shouldBeUpgradable: false},
		{first: types.MembershipTypeBlack, second: types.MembershipTypeSilver, shouldBeUpgradable: false},
		{first: types.MembershipTypeBlack, second: types.MembershipTypeGold, shouldBeUpgradable: false},
		{first: types.MembershipTypeBlack, second: types.MembershipTypeBlack, shouldBeUpgradable: false},
	}

	for _, test := range tests {
		test := test
		t.Run(fmt.Sprintf("%s can upgrade to %s", test.first, test.second), func(t *testing.T) {
			require.Equal(t, test.shouldBeUpgradable, types.CanUpgrade(test.first, test.second))
		})
	}
}

func TestMemberships_AppendIfMissing(t *testing.T) {
	owner, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	membership1 := types.NewMembership(types.MembershipTypeBronze, owner)
	membership2 := types.NewMembership(types.MembershipTypeSilver, owner)

	tests := []struct {
		name             string
		memberships      types.Memberships
		membership       types.Membership
		shouldBeAppended bool
	}{
		{
			name:             "Membership is appended to empty slice",
			memberships:      types.Memberships{},
			membership:       membership1,
			shouldBeAppended: true,
		},
		{
			name:             "Membership is appended to existing list",
			memberships:      types.Memberships{membership1},
			membership:       membership2,
			shouldBeAppended: true,
		},
		{
			name:             "Membership is not appended when existing",
			memberships:      types.Memberships{membership1, membership2},
			membership:       membership1,
			shouldBeAppended: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, appended := test.memberships.AppendIfMissing(test.membership)
			require.Equal(t, test.shouldBeAppended, appended)
			require.Contains(t, result, test.membership)
		})
	}
}
