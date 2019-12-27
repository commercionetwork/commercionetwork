package keeper_test

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
				err := k.AssignMembership(ctx, test.user, test.existingMembership)
				assert.NoError(t, err)
			}

			err := k.AssignMembership(ctx, test.user, test.membershipType)
			assert.Equal(t, test.error, err)
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
			_ = k.AssignMembership(ctx, m.Owner, m.MembershipType)
		}

		_ = k.RemoveMembership(ctx, test.membership.Owner)
		set, _ := k.GetMembershipsSet(ctx)
		assert.True(t, test.expected.Equals(set))
	}
}

func TestKeeper_GetMembership(t *testing.T) {
	tests := []struct {
		name                   string
		existingMembershipType string
		user                   sdk.AccAddress
		expectedError          sdk.Error
		expectedMembership     types.Membership
	}{
		{
			name: "Non existing membership is returned properly",
			user: testUser,
			expectedError: sdk.ErrUnknownRequest(
				fmt.Sprintf("membership not found for user \"%s\"", testUser.String()),
			),
		},
		{
			name:                   "Existing membership is returned properly",
			existingMembershipType: types.MembershipTypeBronze,
			user:                   testUser,
			expectedError:          nil,
			expectedMembership: types.Membership{
				Owner:          testUser,
				MembershipType: types.MembershipTypeBronze,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			_ = k.AssignMembership(ctx, test.user, test.existingMembershipType)

			foundMembership, err := k.GetMembership(ctx, testUser)
			if test.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
			assert.Equal(t, test.expectedMembership, foundMembership)
		})
	}
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
				err := k.AssignMembership(ctx, m.Owner, m.MembershipType)
				assert.NoError(t, err)
			}

			set, _ := k.GetMembershipsSet(ctx)
			for _, m := range test.storedMemberships {
				assert.Contains(t, set, m)
			}
		})
	}
}

func TestKeeper_SetStableCreditsDenom(t *testing.T) {
	tests := []struct {
		name  string
		denom string
	}{
		{
			name:  "Empty credits denom is saved properly",
			denom: "",
		},
		{
			name:  "Non empty credits denom is saved properly",
			denom: "uatom",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			k.SetStableCreditsDenom(ctx, test.denom)

			store := ctx.KVStore(k.StoreKey)
			assert.Equal(t, test.denom, string(store.Get([]byte(types.StableCreditsStoreKey))))
		})
	}
}

func TestKeeper_GetStableCreditsDenom(t *testing.T) {
	tests := []struct {
		name  string
		denom string
	}{
		{
			name:  "Empty credits denom is returned properly",
			denom: "",
		},
		{
			name:  "Non empty credits denom is returned properly",
			denom: "uatom",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			store := ctx.KVStore(k.StoreKey)
			store.Set([]byte(types.StableCreditsStoreKey), []byte(test.denom))

			assert.Equal(t, test.denom, k.GetStableCreditsDenom(ctx))
		})
	}
}
