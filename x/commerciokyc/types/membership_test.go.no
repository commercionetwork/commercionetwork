package types_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMembership_Equals(t *testing.T) {
	owner, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	tsp, _ := sdk.AccAddressFromBech32("cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee")
	expiration := time.Now().Add(60)
	membership := types.NewMembership(types.MembershipTypeBronze, owner, tsp, expiration)

	tests := []struct {
		name          string
		first         types.Membership
		second        types.Membership
		shouldBeEqual bool
	}{
		{
			name:          "Different type returns false",
			first:         membership,
			second:        types.NewMembership(types.MembershipTypeSilver, owner, tsp, expiration),
			shouldBeEqual: false,
		},
		{
			name:          "Different type returns false",
			first:         membership,
			second:        types.NewMembership(types.MembershipTypeGold, owner, tsp, expiration),
			shouldBeEqual: false,
		},
		{
			name:          "Different type returns false",
			first:         membership,
			second:        types.NewMembership(types.MembershipTypeBlack, owner, tsp, expiration),
			shouldBeEqual: false,
		},
		{
			name:          "Different owner returns false",
			first:         membership,
			second:        types.NewMembership(membership.MembershipType, sdk.AccAddress{}, tsp, expiration),
			shouldBeEqual: false,
		},
		{
			name:          "Different tsp returns false",
			first:         membership,
			second:        types.NewMembership(membership.MembershipType, owner, sdk.AccAddress{}, expiration),
			shouldBeEqual: false,
		},
		{
			name:          "Different expiry at returns false",
			first:         membership,
			second:        types.NewMembership(membership.MembershipType, owner, tsp, time.Now().Add(90)),
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

func TestMembership_Validate(t *testing.T) {
	timestamp, err := time.Parse(time.RFC3339, "2019-03-23T06:35:22Z")
	require.NoError(t, err)

	tests := []struct {
		name       string
		membership types.Membership
		wantErr    bool
	}{
		{
			name: "ok",
			membership: types.Membership{
				Owner:          government.String(),
				TspAddress:     tsp.String(),
				MembershipType: types.MembershipTypeBlack,
				ExpiryAt:       &timestamp,
			},
			wantErr: false,
		},
		{
			name: "invalid membership",
			membership: types.Membership{
				Owner:          government.String(),
				TspAddress:     tsp.String(),
				MembershipType: "teal",
				ExpiryAt:       &timestamp,
			},
			wantErr: true,
		},
		{
			name: "invalid owner",
			membership: types.Membership{
				Owner:          "",
				TspAddress:     tsp.String(),
				MembershipType: types.MembershipTypeBlack,
				ExpiryAt:       &timestamp,
			},
			wantErr: true,
		},
		{
			name: "invalid tsp",
			membership: types.Membership{
				Owner:          government.String(),
				TspAddress:     "",
				MembershipType: types.MembershipTypeBlack,
				ExpiryAt:       &timestamp,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.membership
			if err := m.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Membership.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})

	}
}
