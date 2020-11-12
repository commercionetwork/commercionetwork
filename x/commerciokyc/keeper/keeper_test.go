package keeper

import (
	"fmt"
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func TestKeeper_BuyMembership(t *testing.T) {
	tests := []struct {
		name           string
		membershipType string
		user           sdk.AccAddress
		tsp            sdk.AccAddress
		bankAmount     sdk.Coins
		height         int64
		error          error
	}{
		{
			name:           "Invalid membership type black returns error",
			membershipType: types.MembershipTypeBlack,
			user:           testUser,
			tsp:            testTsp,
			height:         testHeight,
			bankAmount:     sdk.NewCoins(sdk.NewInt64Coin(stableCreditDenom, 10000000000)),
			error:          sdkErr.Wrap(sdkErr.ErrInvalidAddress, "cannot buy black membership"),
		},
		{
			name:           "Valid membership type returns no errors",
			membershipType: types.MembershipTypeBronze,
			user:           testUser,
			tsp:            testTsp,
			bankAmount:     sdk.NewCoins(sdk.NewInt64Coin(stableCreditDenom, 10000000000)),
			height:         testHeight,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, bk, _, k := SetupTestInput()
			_ = bk.SetCoins(ctx, test.tsp, test.bankAmount)
			err := k.BuyMembership(ctx, test.user, test.membershipType, test.tsp, test.height)

			if err != nil && test.error != nil {
				require.Equal(t, test.error.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeper_AssignMembership(t *testing.T) {
	tests := []struct {
		name               string
		existingMembership string
		membershipType     string
		expectedNotExists  bool
		userIsTsp          bool
		user               sdk.AccAddress
		tsp                sdk.AccAddress
		height             int64
		error              error
	}{
		{
			name:           "Invalid membership type returns error",
			membershipType: "grn",
			user:           testUser,
			tsp:            testTsp,
			height:         testHeight,
			error:          sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid membership type: grn"),
		},
		{
			name:           "Membership with invalid height zero or negative",
			user:           testUser,
			tsp:            testTsp,
			height:         testHeightNegative,
			membershipType: types.MembershipTypeBronze,
			error:          sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid expiry height: -1"),
		},
		/*{
			name:               "Invalid tsp",
			user:               testUser,
			tsp:                testUser,
			height:             testHeight,
			existingMembership: types.MembershipTypeBronze,
			error:              sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid expiry height: -1"),
		},*/
		{
			name:           "Non existing membership is properly saved",
			user:           testUser,
			tsp:            testTsp,
			height:         testHeight,
			membershipType: types.MembershipTypeBronze,
		},
		{
			name:               "Existing membership is replaced",
			user:               testUser,
			tsp:                testTsp,
			height:             testHeight,
			existingMembership: types.MembershipTypeBronze,
			membershipType:     types.MembershipTypeGold,
		},
		{
			name:               "Cannot assign tsp membership",
			userIsTsp:          true,
			user:               testUser,
			tsp:                testTsp,
			height:             testHeight,
			existingMembership: types.MembershipTypeBlack,
			membershipType:     types.MembershipTypeGold,
			error:              sdkErr.Wrap(sdkErr.ErrUnauthorized, "account \""+testUser.String()+"\" is a Trust Service Provider: remove from tsps list before"),
		},
		/*{
			name:               "Assign \"none\" membership type to delete membership",
			user:               testUser,
			tsp:                testTsp,
			height:             testHeight,
			existingMembership: types.MembershipTypeBronze,
			membershipType:     types.MembershipTypeNone,
			expectedNotExists:  true,
			error:              sdkErr.Wrap(sdkErr.ErrUnknownRequest, "membership not found for user \"cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae\""),
		},*/
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			if len(test.existingMembership) != 0 {
				err := k.AssignMembership(ctx, test.user, test.existingMembership, test.tsp, test.height)
				require.NoError(t, err)
			}

			if test.userIsTsp {
				k.AddTrustedServiceProvider(ctx, test.user)
			}

			err := k.AssignMembership(ctx, test.user, test.membershipType, test.tsp, test.height)

			if err != nil {
				if test.expectedNotExists {
					_, err2 := k.GetMembership(ctx, test.user)
					require.Equal(t, test.error.Error(), err2.Error())
				}
				require.Equal(t, test.error.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestKeeper_ComputeExpiryHeight(t *testing.T) {
	tests := []struct {
		name              string
		expectedNumBlocks int64
		curHeight         int64
	}{
		{
			name:              "Compute expiry",
			expectedNumBlocks: yearBlocks + testHeight,
			curHeight:         testHeight,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			_, _, _, k := SetupTestInput()

			computedHeight := k.ComputeExpiryHeight(test.curHeight)
			require.Equal(t, test.expectedNumBlocks, computedHeight)

		})
	}
}

func TestKeeper_GetMembership(t *testing.T) {
	tests := []struct {
		name                   string
		existingMembershipType string
		user                   sdk.AccAddress
		tsp                    sdk.AccAddress
		height                 int64
		expectedError          error
		expectedMembership     types.Membership
	}{
		{
			name:   "Non existing membership is returned properly",
			user:   testUser,
			tsp:    testTsp,
			height: testHeight,
			expectedError: sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf("membership not found for user \"%s\"", testUser.String()),
			),
		},
		{
			name:                   "Existing membership is returned properly",
			existingMembershipType: types.MembershipTypeBronze,
			user:                   testUser,
			tsp:                    testTsp,
			height:                 testHeight,
			expectedError:          nil,
			expectedMembership: types.Membership{
				Owner:          testUser,
				TspAddress:     testTsp,
				MembershipType: types.MembershipTypeBronze,
				ExpiryAt:       testHeight,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			_ = k.AssignMembership(ctx, test.user, test.existingMembershipType, test.tsp, test.height)

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

func TestKeeper_RemoveMembership(t *testing.T) {
	tests := []struct {
		name       string
		membership types.Membership
		tsp        sdk.AccAddress
		mustError  bool
	}{
		{
			name:       "Non existing membership throws an error",
			membership: types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testHeight),
			mustError:  true,
		},
		{
			name:       "Existing membership is removed properly",
			membership: types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testHeight),
			mustError:  false,
		},
		{
			name:       "Tsp membership cannot be removed",
			membership: types.NewMembership(types.MembershipTypeBlack, testTsp, testUser, testHeight),
			tsp:        testUser,
			mustError:  true,
		},
	}

	for _, test := range tests {
		ctx, _, _, k := SetupTestInput()

		// if the test should not throw an error when removing, we must add
		// a membership first
		if !test.mustError {
			_ = k.AssignMembership(ctx, test.membership.Owner, test.membership.MembershipType, test.membership.TspAddress, test.membership.ExpiryAt)
		}

		if test.tsp != nil {
			k.AddTrustedServiceProvider(ctx, test.tsp)
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
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testHeight),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testHeight),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, m := range test.storedMemberships {
				err := k.AssignMembership(ctx, m.Owner, m.MembershipType, m.TspAddress, m.ExpiryAt)
				require.NoError(t, err)
			}
			i := k.MembershipIterator(ctx)
			for ; i.Valid(); i.Next() {
				m := k.ExtractMembership(i.Value())
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
			[]byte{10, 20, 153, 39, 56, 31, 38, 42, 65, 168, 74, 73, 145, 237, 226, 147, 118, 104, 171, 0, 46, 239, 18, 20, 251, 182, 16, 225, 99, 30, 161, 9, 143, 124, 32, 185, 74, 85, 162, 77, 31, 208, 217, 44, 26, 6, 98, 114, 111, 110, 122, 101, 32, 10},
			types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testHeight),
			false,
		},
		{
			// Da correggere
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
				m := k.ExtractMembership(test.value)
				require.Equal(t, test.membership, m)
			} else {
				require.Panics(t, func() {
					_ = k.ExtractMembership(test.value)
				})
			}
		})
	}
}

func TestKeeper_GetMemberships(t *testing.T) {
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
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testHeight),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testHeight),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, m := range test.storedMemberships {
				err := k.AssignMembership(ctx, m.Owner, m.MembershipType, m.TspAddress, m.ExpiryAt)
				require.NoError(t, err)
			}
			ms := k.GetMemberships(ctx)
			for _, mg := range ms {
				require.Contains(t, test.storedMemberships, mg)
			}
		})
	}
}

func TestKeeper_GetTspMemberships(t *testing.T) {
	tests := []struct {
		name               string
		tsp                sdk.AccAddress
		storedMemberships  types.Memberships
		expetedMemberships types.Memberships
	}{
		{
			name:               "Empty set is returned properly",
			tsp:                testTsp,
			storedMemberships:  types.Memberships{},
			expetedMemberships: types.Memberships{},
		},
		{
			name: "Tsp without membership return empty set properly ",
			tsp:  testTsp,
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testUser2, testHeight),
				types.NewMembership(types.MembershipTypeGold, testUser2, testUser, testHeight),
			},
			expetedMemberships: types.Memberships{},
		},
		{
			name: "Tsp with some memberships return set with correct items",
			tsp:  testTsp,
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testHeight),
				types.NewMembership(types.MembershipTypeGold, testUser2, testUser, testHeight),
			},
			expetedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testHeight),
			},
		},
		{
			name: "Tsp with all memberships return all set",
			tsp:  testTsp,
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testHeight),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testHeight),
			},
			expetedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testHeight),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testHeight),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, m := range test.storedMemberships {
				err := k.AssignMembership(ctx, m.Owner, m.MembershipType, m.TspAddress, m.ExpiryAt)
				require.NoError(t, err)
			}
			ms := k.GetTspMemberships(ctx, test.tsp)
			require.Equal(t, test.expetedMemberships, ms)
		})
	}
}

func TestKeeper_GetExportMemberships(t *testing.T) {
	tests := []struct {
		name                string
		blockHeight         int64
		storedMemberships   types.Memberships
		expectedMemberships types.Memberships
	}{
		{
			name:                "Empty set is returned properly",
			storedMemberships:   types.Memberships{},
			expectedMemberships: types.Memberships{},
		},
		{
			name:        "All expired membership return empty set",
			blockHeight: int64(11),
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testUser2, testHeight),
				types.NewMembership(types.MembershipTypeGold, testUser2, testUser, testHeight),
			},
			expectedMemberships: types.Memberships{},
		},
		{
			name:        "Some memberships expired return set with correct items",
			blockHeight: int64(11),
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, int64(12)),
				types.NewMembership(types.MembershipTypeBronze, testUser2, testTsp, int64(11)),
				types.NewMembership(types.MembershipTypeGold, testUser3, testUser, testHeight),
			},
			expectedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, int64(1)),
			},
		},
		{
			name:        "All memberships not expired return all set",
			blockHeight: int64(9),
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testHeight),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testHeight),
			},
			expectedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, int64(1)),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, int64(1)),
			},
		},
		{
			name:        "Black membership returned in set although it has expired",
			blockHeight: int64(11),
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBlack, testUser, testTsp, testHeight),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testHeight),
			},
			expectedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBlack, testUser, testTsp, int64(-1)),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, m := range test.storedMemberships {
				err := k.AssignMembership(ctx, m.Owner, m.MembershipType, m.TspAddress, m.ExpiryAt)
				require.NoError(t, err)
			}
			ms := k.GetExportMemberships(ctx, test.blockHeight)
			require.Equal(t, test.expectedMemberships, ms)
		})
	}
}

func TestKeeper_RemoveExpiredMemberships(t *testing.T) {
	tests := []struct {
		name                string
		storedMemberships   types.Memberships
		expectedMemberships types.Memberships
	}{
		{
			name:                "Empty set is properly stored",
			storedMemberships:   types.Memberships{},
			expectedMemberships: types.Memberships{},
		},
		{
			name: "All expired memberships are properly removed",
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testUser2, int64(1)),
				types.NewMembership(types.MembershipTypeGold, testUser2, testUser, int64(1)),
			},
			expectedMemberships: types.Memberships{},
		},
		{
			name: "Some memberships expired are properly removed",
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, int64(1)),
				types.NewMembership(types.MembershipTypeBronze, testUser2, testTsp, int64(1)),
				types.NewMembership(types.MembershipTypeGold, testUser3, testUser, testHeight),
			},
			expectedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeGold, testUser3, testUser, testHeight),
			},
		},
		{
			name: "All memberships not expired are properly left on store",
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testHeight),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testHeight),
			},
			expectedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testHeight),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testHeight),
			},
		},
		{
			name: "Black membership left in store and renewed",
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBlack, testUser, testTsp, int64(1)),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, int64(1)),
			},
			expectedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBlack, testUser, testTsp, yearBlocks),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			for _, m := range test.storedMemberships {
				err := k.AssignMembership(ctx, m.Owner, m.MembershipType, m.TspAddress, m.ExpiryAt)
				require.NoError(t, err)
			}
			err2 := k.RemoveExpiredMemberships(ctx)
			require.NoError(t, err2)
			ms := k.GetMemberships(ctx)
			require.Equal(t, test.expectedMemberships, ms)
		})
	}
}
