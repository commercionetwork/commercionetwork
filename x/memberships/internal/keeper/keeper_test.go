package keeper_test

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/nft"
	"github.com/cosmos/modules/incubator/nft/exported"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_AssignMembership(t *testing.T) {
	tests := []struct {
		name               string
		existingMembership string
		membershipType     string
		user               sdk.AccAddress
		error              sdk.Error
	}{
		{
			name:           "Invalid membership type returns error",
			membershipType: "grn",
			user:           testUser,
			error:          sdk.ErrUnknownRequest("Invalid membership type: grn"),
		},
		{
			name:           "Non existing membership is properly saved",
			user:           testUser,
			membershipType: types.MembershipTypeBronze,
		},
		{
			name:               "Existing membership is replaced",
			user:               testUser,
			existingMembership: types.MembershipTypeBronze,
			membershipType:     types.MembershipTypeGold,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			if len(test.existingMembership) != 0 {
				_, err := k.AssignMembership(ctx, test.user, test.existingMembership)
				assert.NoError(t, err)
			}

			tokenURI, err := k.AssignMembership(ctx, test.user, test.membershipType)
			assert.Equal(t, test.error, err)

			if test.error == nil {
				expectedURI := fmt.Sprintf("membership:%s:membership-%s", test.membershipType, test.user)
				assert.Equal(t, expectedURI, tokenURI)
			}
		})
	}
}

func TestKeeper_RemoveMembership(t *testing.T) {
	tests := []struct {
		name        string
		memberships types.Memberships
		membership  types.Membership
		expected    types.Memberships
	}{
		{
			name:        "Non existing membership works properly",
			memberships: types.Memberships{},
			membership:  types.NewMembership(types.MembershipTypeBronze, testUser),
			expected:    types.Memberships{},
		},
		{
			name: "Existing membership is removed properly",
			memberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser),
				types.NewMembership(types.MembershipTypeGold, testUser2),
			},
			membership: types.NewMembership(types.MembershipTypeBronze, testUser),
			expected: types.Memberships{
				types.NewMembership(types.MembershipTypeGold, testUser2),
			},
		},
	}

	for _, test := range tests {
		ctx, _, _, k := SetupTestInput()

		for _, m := range test.memberships {
			_, _ = k.AssignMembership(ctx, m.Owner, m.MembershipType)
		}

		_ = k.RemoveMembership(ctx, test.membership.Owner)
		assert.True(t, test.expected.Equals(k.GetMembershipsSet(ctx)))
	}
}

func TestKeeper_GetMembership(t *testing.T) {
	tests := []struct {
		name                   string
		existingMembershipType string
		user                   sdk.AccAddress
		expectedFound          bool
		expectedMembership     exported.NFT
	}{
		{
			name:          "Non existing membership is returned properly",
			user:          testUser,
			expectedFound: false,
		},
		{
			name:                   "Existing membership is returned properly",
			existingMembershipType: types.MembershipTypeBronze,
			user:                   testUser,
			expectedFound:          true,
			expectedMembership: &nft.BaseNFT{
				ID:       fmt.Sprintf("membership-%s", testUser),
				Owner:    testUser,
				TokenURI: fmt.Sprintf("membership:%s:%s", types.MembershipTypeBronze, fmt.Sprintf("membership-%s", testUser)),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			_, _ = k.AssignMembership(ctx, test.user, test.existingMembershipType)

			foundMembership, found := k.GetMembership(ctx, testUser)
			assert.Equal(t, test.expectedFound, found)
			assert.Equal(t, test.expectedMembership, foundMembership)
		})
	}
}

func TestKeeper_GetMembershipType(t *testing.T) {
	_, _, _, k := SetupTestInput()

	id := "123"
	membershipType := "black"
	membership := nft.NewBaseNFT(id, testUser, fmt.Sprintf("membership:%s:%s", membershipType, id))

	assert.Equal(t, membershipType, k.GetMembershipType(&membership))
}

func TestKeeper_GetMembershipsSet(t *testing.T) {
	tests := []struct {
		name              string
		storedMemberships types.Memberships
	}{
		{
			name:              "Empty set is returned properly",
			storedMemberships: types.Memberships{},
		},
		{
			name: "Existing set is returned properly",
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser),
				types.NewMembership(types.MembershipTypeGold, testUser2),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, m := range test.storedMemberships {
				_, err := k.AssignMembership(ctx, m.Owner, m.MembershipType)
				assert.NoError(t, err)
			}

			set := k.GetMembershipsSet(ctx)
			for _, m := range test.storedMemberships {
				assert.Contains(t, set, m)
			}
		})
	}
}
