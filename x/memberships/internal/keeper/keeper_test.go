package keeper_test

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
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
				require.NoError(t, err)
			}

			err := k.AssignMembership(ctx, test.user, test.membershipType)
			require.Equal(t, test.error, err)
		})
	}
}

func TestKeeper_RemoveMembership(t *testing.T) {
	tests := []struct {
		name       string
		membership types.Membership
		mustError  bool
	}{
		{
			name:       "Non existing membership throws an error",
			membership: types.NewMembership(types.MembershipTypeBronze, testUser),
			mustError:  true,
		},
		{
			name:       "Existing membership is removed properly",
			membership: types.NewMembership(types.MembershipTypeBronze, testUser),
			mustError:  false,
		},
	}

	for _, test := range tests {
		ctx, _, _, k := SetupTestInput()

		// if the test should not throw an error when removing, we must add
		// a membership first
		if !test.mustError {
			_ = k.AssignMembership(ctx, test.membership.Owner, test.membership.MembershipType)
		}

		err := k.RemoveMembership(ctx, test.membership.Owner)
		if !test.mustError {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}

		_, err = k.GetMembership(ctx, test.membership.Owner)
		require.Error(t, err)
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
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			require.Equal(t, test.expectedMembership, foundMembership)
		})
	}
}

func TestKeeper_MembershipIterator(t *testing.T) {
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
				require.NoError(t, err)
			}
			i := k.MembershipIterator(ctx)
			for ; i.Valid(); i.Next() {
				m := k.ExtractMembership(i.Key(), i.Value())
				require.Contains(t, test.storedMemberships, m)
			}
		})
	}
}

func TestKeeper_ExtractMembership(t *testing.T) {
	tests := []struct {
		name       string
		key        []byte
		value      []byte
		membership types.Membership
		mustFail   bool
	}{
		{
			"a good membership",
			[]byte{97, 99, 99, 114, 101, 100, 105, 116, 97, 116, 105, 111, 110, 115, 58, 115, 116, 111, 114, 97, 103, 101, 58, 20, 153, 39, 56, 31, 38, 42, 65, 168, 74, 73, 145, 237, 226, 147, 118, 104, 171, 0, 46, 239},
			[]byte{6, 98, 114, 111, 110, 122, 101},
			types.NewMembership(types.MembershipTypeBronze, testUser),
			false,
		},
		{
			"a badly serialized membership",
			[]byte{99, 99, 114, 101, 100, 105, 116, 97, 116, 105, 111, 110, 115, 58, 115, 116, 111, 114, 97, 103, 101, 58, 20, 153, 39, 56, 31, 38, 42, 65, 168, 74, 73, 145, 237, 226, 147, 118, 104, 171, 0, 46, 239},
			[]byte{6, 114, 111, 110, 122, 101},
			types.Membership{},
			true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			_, _, _, k := SetupTestInput()
			if !test.mustFail {
				m := k.ExtractMembership(test.key, test.value)
				require.Equal(t, test.membership, m)
			} else {
				require.Panics(t, func() {
					_ = k.ExtractMembership(test.key, test.value)
				})
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
			require.Equal(t, test.denom, string(store.Get([]byte(types.StableCreditsStoreKey))))
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

			require.Equal(t, test.denom, k.GetStableCreditsDenom(ctx))
		})
	}
}
