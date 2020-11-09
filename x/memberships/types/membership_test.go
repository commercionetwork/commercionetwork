package types_test

import (
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/memberships/types"
)

func TestMembership_Equals(t *testing.T) {
	owner, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	tsp, _ := sdk.AccAddressFromBech32("cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee")
	height := int64(10)
	membership := types.NewMembership(types.MembershipTypeBronze, owner, tsp, height)

	tests := []struct {
		name          string
		first         types.Membership
		second        types.Membership
		shouldBeEqual bool
	}{
		{
			name:          "Different type returns false",
			first:         membership,
			second:        types.NewMembership(types.MembershipTypeSilver, membership.Owner, membership.TspAddress, membership.ExpiryAt),
			shouldBeEqual: false,
		},
		{
			name:          "Different type returns false",
			first:         membership,
			second:        types.NewMembership(types.MembershipTypeGold, membership.Owner, membership.TspAddress, membership.ExpiryAt),
			shouldBeEqual: false,
		},
		{
			name:          "Different type returns false",
			first:         membership,
			second:        types.NewMembership(types.MembershipTypeBlack, membership.Owner, membership.TspAddress, membership.ExpiryAt),
			shouldBeEqual: false,
		},
		{
			name:          "Different owner returns false",
			first:         membership,
			second:        types.NewMembership(membership.MembershipType, sdk.AccAddress{}, membership.TspAddress, membership.ExpiryAt),
			shouldBeEqual: false,
		},
		{
			name:          "Different tsp returns false",
			first:         membership,
			second:        types.NewMembership(membership.MembershipType, membership.Owner, sdk.AccAddress{}, membership.ExpiryAt),
			shouldBeEqual: false,
		},
		{
			name:          "Different expiry at returns false",
			first:         membership,
			second:        types.NewMembership(membership.MembershipType, membership.Owner, membership.TspAddress, int64(11)),
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

func TestMembership_IsMembershipTypeValid(t *testing.T) {
	tests := []struct {
		membershipType string
		shouldBeValid  bool
	}{
		{membershipType: types.MembershipTypeGreen, shouldBeValid: true},
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

func TestMemberships_AppendIfMissing(t *testing.T) {
	owner, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	tsp, _ := sdk.AccAddressFromBech32("cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee")
	height := int64(10)

	membership1 := types.NewMembership(types.MembershipTypeBronze, owner, tsp, height)
	membership2 := types.NewMembership(types.MembershipTypeSilver, owner, tsp, height)

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

func TestMembership_ValidateBasic(t *testing.T) {
	owner, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	tsp, _ := sdk.AccAddressFromBech32("cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee")
	height := int64(10)

	tests := []struct {
		name    string
		invite  types.Membership
		wantErr bool
	}{
		{
			"A valid bronze membership",
			types.Membership{
				Owner:          owner,
				TspAddress:     tsp,
				MembershipType: types.MembershipTypeBronze,
				ExpiryAt:       height,
			},
			false,
		},
		{
			"A invalid membership",
			types.Membership{
				Owner:          owner,
				TspAddress:     tsp,
				MembershipType: "nvm",
				ExpiryAt:       height,
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
