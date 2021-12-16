package keeper_test

import (
	"fmt"
	"testing"
	"time"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func TestKeeper_AssignMembership(t *testing.T) {
	tests := []struct {
		name               string
		existingMembership string
		membershipType     string
		expectedNotExists  bool
		userIsTsp          bool
		user               sdk.AccAddress
		tsp                sdk.AccAddress
		expiredAt          time.Time
		error              error
	}{
		{
			name:           "Invalid membership type returns error",
			membershipType: "grn",
			user:           testUser,
			tsp:            testTsp,
			expiredAt:      testExpiration,
			error:          sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid membership type: grn"),
		},
		{
			name:           "Membership with invalid expired date",
			user:           testUser,
			tsp:            testTsp,
			expiredAt:      testExpirationNegative,
			membershipType: types.MembershipTypeBronze,
			error:          sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid expiry date: %s is before current block time", testExpirationNegative)),
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
			expiredAt:      testExpiration,
			membershipType: types.MembershipTypeBronze,
		},
		{
			name:               "Existing membership is replaced",
			user:               testUser,
			tsp:                testTsp,
			expiredAt:          testExpiration,
			existingMembership: types.MembershipTypeBronze,
			membershipType:     types.MembershipTypeGold,
		},
		{
			name:               "Cannot assign tsp membership",
			userIsTsp:          true,
			user:               testUser,
			tsp:                testTsp,
			expiredAt:          testExpiration,
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
				err := k.AssignMembership(ctx, test.user, test.existingMembership, test.tsp, test.expiredAt)
				require.NoError(t, err)
			}

			if test.userIsTsp {
				k.AddTrustedServiceProvider(ctx, test.user)
			}

			err := k.AssignMembership(ctx, test.user, test.membershipType, test.tsp, test.expiredAt)

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

func TestKeeper_DeleteMembership(t *testing.T) {
	tests := []struct {
		name       string
		membership types.Membership
		tsp        sdk.AccAddress
		mustError  bool
	}{
		{
			name:       "Non existing membership throws an error",
			membership: types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testExpiration),
			mustError:  true,
		},
		{
			name:       "Existing membership is removed properly",
			membership: types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testExpiration),
			mustError:  false,
		},
		{
			name:       "Tsp membership cannot be removed",
			membership: types.NewMembership(types.MembershipTypeBlack, testTsp, testUser, testExpiration),
			tsp:        testUser,
			mustError:  true,
		},
	}

	for _, test := range tests {
		ctx, _, _, k := SetupTestInput()

		// if the test should not throw an error when removing, we must add
		// a membership first
		membershipOwner, _ := sdk.AccAddressFromBech32(test.membership.Owner)
		membershipTspAddress, _ := sdk.AccAddressFromBech32(test.membership.TspAddress)

		if !test.mustError {
			_ = k.AssignMembership(ctx, membershipOwner, test.membership.MembershipType, membershipTspAddress, *test.membership.ExpiryAt)
		}

		if test.tsp != nil {
			k.AddTrustedServiceProvider(ctx, test.tsp)
		}

		err := k.DeleteMembership(ctx, membershipOwner)
		if !test.mustError {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}

		_, err = k.GetMembership(ctx, membershipOwner)
		require.Error(t, err)
	}
}
func TestKeeper_DistributeReward(t *testing.T) {
	coins := sdk.Coins{}
	tests := []struct {
		name                 string
		invite               types.Invite
		pool                 sdk.Coins
		expectedGovBalance   sdk.Coins
		expectedUserBalance  sdk.Coins
		expectedInviteStatus int64
		mustError            bool
		expectedError        error
	}{
		{
			name:                 "Invite status is invalid",
			invite:               types.Invite{Sender: testInviteSender.String(), SenderMembership: types.MembershipTypeGold, User: testUser.String(), Status: uint64(types.InviteStatusRewarded)},
			expectedInviteStatus: int64(types.InviteStatusRewarded),
			mustError:            false,
		},
		{
			name:      "Invite sender does no have a membership",
			invite:    types.Invite{Sender: testInviteSender.String(), SenderMembership: types.MembershipTypeGold, User: testUser2.String(), Status: uint64(types.InviteStatusPending)},
			mustError: true,
		},
		{
			name:      "Invite recipient does not have a membership",
			invite:    types.Invite{Sender: testUser2.String(), SenderMembership: types.MembershipTypeGold, User: testUser.String(), Status: uint64(types.InviteStatusPending)},
			mustError: true,
		},
		/*{
			name:          "Memberships matrix option reward not exists",
			invite:        types.Invite{Sender: testTsp.String(), SenderMembership: "bold", User: testUser2.String(), Status: uint64(types.InviteStatusPending)},
			mustError:     true,
			expectedError: sdkErr.Wrap(sdkErr.ErrInvalidRequest, "Invalid reward options"),
		},*/
		{
			name:      "Pool has zero tokens",
			invite:    types.Invite{Sender: testTsp.String(), SenderMembership: "gold", User: testUser2.String(), Status: uint64(types.InviteStatusPending)},
			pool:      coins,
			mustError: true,
			//expectedError: sdkErr.Wrap(sdkErr.ErrUnauthorized, "ABR pool has zero tokens"),
		},
		{
			name:      "Account has not sufficient funds (pool is small then expected reward)",
			invite:    types.Invite{Sender: testTsp.String(), SenderMembership: "gold", User: testUser2.String(), Status: uint64(types.InviteStatusPending)},
			pool:      sdk.NewCoins(sdk.NewCoin(testDenom, sdk.NewInt(100000000))),
			mustError: true,
			//expectedError: sdkErr.Wrap(sdkErr.ErrUnauthorized, "ABR pool has zero tokens"),
			// could not move collateral amount to module account, 43478261ucommercio is smaller than 100000001ucommercio: insufficient funds
		},
		// TODO: correct reward fail on open position
		/*{
			name:                 "Account correctly rewarded",
			invite:               types.Invite{Sender: testTsp.String(), SenderMembership: "gold", User: testUser2.String(), Status: uint64(types.InviteStatusPending)},
			pool:                 sdk.NewCoins(sdk.NewCoin(testDenom, sdk.NewInt(1000000000000))),
			expectedInviteStatus: int64(types.InviteStatusRewarded),
			mustError:            false,
		},*/
	}
	for _, test := range tests {
		ctx, bk, _, k := SetupTestInput()

		err := k.AssignMembership(ctx, testUser2, types.MembershipTypeGold, testTsp, testExpiration)
		require.NoError(t, err)
		err = k.AssignMembership(ctx, testTsp, types.MembershipTypeBlack, testTsp, testExpiration)
		require.NoError(t, err)
		err = k.SetLiquidityPoolToAccount(ctx, test.pool)
		require.NoError(t, err)
		senderAccAddr, _ := sdk.AccAddressFromBech32(test.invite.Sender)
		//userAccAddr, _ := sdk.AccAddressFromBech32(test.invite.User)
		//senderMembershipType := test.invite.SenderMembership
		senderBalance := bk.GetAllBalances(ctx, senderAccAddr)
		//userBalance := bk.GetAllBalances(ctx, userAccAddr)

		err = k.DistributeReward(ctx, test.invite)
		if !test.mustError {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
			if test.expectedError != nil {
				require.Equal(t, test.expectedError, err)
			}
		}

		if test.expectedInviteStatus > 0 {
			require.Equal(t, test.expectedInviteStatus, int64(test.invite.Status))
		}

		if test.expectedUserBalance != nil {
			require.Equal(t, test.expectedUserBalance, senderBalance)
		}

	}

}

func TestKeeper_GetMembership(t *testing.T) {
	tests := []struct {
		name                   string
		existingMembershipType string
		user                   sdk.AccAddress
		tsp                    sdk.AccAddress
		expiration             time.Time
		expectedError          error
		expectedMembership     types.Membership
	}{
		{
			name:       "Non existing membership is returned properly",
			user:       testUser,
			tsp:        testTsp,
			expiration: testExpiration,
			expectedError: sdkErr.Wrap(sdkErr.ErrUnknownRequest,
				fmt.Sprintf("membership not found for user \"%s\"", testUser.String()),
			),
		},
		{
			name:                   "Existing membership is returned properly",
			existingMembershipType: types.MembershipTypeBronze,
			user:                   testUser,
			tsp:                    testTsp,
			expiration:             testExpiration,
			expectedError:          nil,
			expectedMembership: types.Membership{
				Owner:          testUser.String(),
				TspAddress:     testTsp.String(),
				MembershipType: types.MembershipTypeBronze,
				ExpiryAt:       &testExpiration,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			_ = k.AssignMembership(ctx, test.user, test.existingMembershipType, test.tsp, test.expiration)

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
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testExpiration),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testExpiration),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, m := range test.storedMemberships {
				mOwner, _ := sdk.AccAddressFromBech32(m.Owner)
				mTspAddress, _ := sdk.AccAddressFromBech32(m.TspAddress)
				err := k.AssignMembership(ctx, mOwner, m.MembershipType, mTspAddress, *m.ExpiryAt)
				require.NoError(t, err)
			}
			ms := k.GetMemberships(ctx)
			for _, mg := range ms {
				require.Contains(t, test.storedMemberships, *mg)
			}
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
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testExpiration),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testExpiration),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, m := range test.storedMemberships {
				mOwner, _ := sdk.AccAddressFromBech32(m.Owner)
				mTspAddress, _ := sdk.AccAddressFromBech32(m.TspAddress)
				err := k.AssignMembership(ctx, mOwner, m.MembershipType, mTspAddress, *m.ExpiryAt)
				require.NoError(t, err)
			}
			i := k.MembershipIterator(ctx)
			for ; i.Valid(); i.Next() {
				//m := k.ExtractMembership(i.Value())
				var m types.Membership
				k.Cdc.MustUnmarshalBinaryBare(i.Value(), &m)
				require.Contains(t, test.storedMemberships, m)
			}
		})
	}
}

func TestKeeper_ComputeExpiryHeight(t *testing.T) {
	currentTime := time.Now()
	tests := []struct {
		name               string
		expectedExpiration time.Time
		curTime            time.Time
	}{
		{
			name:               "Compute expiry",
			expectedExpiration: currentTime.Add(SecondsPerYear),
			curTime:            currentTime,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			_, _, _, k := SetupTestInput()

			computedHeight := k.ComputeExpiryHeight(test.curTime)
			require.Equal(t, test.expectedExpiration, computedHeight)

		})
	}
}

/*
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
			types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testExpiration),
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
*/

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
				types.NewMembership(types.MembershipTypeBronze, testUser, testUser2, testExpiration),
				types.NewMembership(types.MembershipTypeGold, testUser2, testUser, testExpiration),
			},
			expetedMemberships: types.Memberships{},
		},
		{
			name: "Tsp with some memberships return set with correct items",
			tsp:  testTsp,
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testExpiration),
				types.NewMembership(types.MembershipTypeGold, testUser2, testUser, testExpiration),
			},
			expetedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testExpiration),
			},
		},
		{
			name: "Tsp with all memberships return all set",
			tsp:  testTsp,
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testExpiration),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testExpiration),
			},
			expetedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testExpiration),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testExpiration),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, m := range test.storedMemberships {
				mOwner, _ := sdk.AccAddressFromBech32(m.Owner)
				mTspAddress, _ := sdk.AccAddressFromBech32(m.TspAddress)
				err := k.AssignMembership(ctx, mOwner, m.MembershipType, mTspAddress, *m.ExpiryAt)
				require.NoError(t, err)
			}
			ms := k.GetTspMemberships(ctx, test.tsp)
			require.Equal(t, test.expetedMemberships, ms)
		})
	}
}

func TestKeeper_ExportMemberships(t *testing.T) {
	tests := []struct {
		name                string
		storedMemberships   types.Memberships
		expectedMemberships types.Memberships
	}{
		{
			name:                "Empty set is returned properly",
			storedMemberships:   types.Memberships{},
			expectedMemberships: types.Memberships{},
		},

		{
			name: "All memberships return all set",
			storedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testExpiration),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testExpiration),
			},
			expectedMemberships: types.Memberships{
				types.NewMembership(types.MembershipTypeBronze, testUser, testTsp, testExpiration),
				types.NewMembership(types.MembershipTypeGold, testUser2, testTsp, testExpiration),
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			for _, m := range test.storedMemberships {
				mOwner, _ := sdk.AccAddressFromBech32(m.Owner)
				mTspAddress, _ := sdk.AccAddressFromBech32(m.TspAddress)
				err := k.AssignMembership(ctx, mOwner, m.MembershipType, mTspAddress, *m.ExpiryAt)
				require.NoError(t, err)
			}
			ms := k.ExportMemberships(ctx)
			require.Equal(t, test.expectedMemberships, ms)
		})
	}
}

func TestKeeper_RemoveExpiredMemberships(t *testing.T) {
	curRemoveTime := time.Now().AddDate(1, 0, 1)
	curTime := testExpiration.AddDate(1, 0, 0)
	curTimePlusYear := curRemoveTime.UTC().Add(SecondsPerYear)
	tests := []struct {
		name                string
		storedMemberships   []*types.Membership
		expectedMemberships []*types.Membership
	}{
		{
			name:                "Empty set is properly stored",
			storedMemberships:   []*types.Membership{},
			expectedMemberships: []*types.Membership{},
		},
		{
			name: "All expired memberships are properly removed",
			storedMemberships: []*types.Membership{
				&types.Membership{testUser.String(), testUser2.String(), types.MembershipTypeBronze, &testExpiration},
				&types.Membership{testUser2.String(), testUser.String(), types.MembershipTypeGold, &testExpiration},
			},
			expectedMemberships: []*types.Membership{},
		},
		{
			name: "Some memberships expired are properly removed",
			storedMemberships: []*types.Membership{
				&types.Membership{testUser.String(), testUser2.String(), types.MembershipTypeBronze, &testExpiration},
				&types.Membership{testUser2.String(), testTsp.String(), types.MembershipTypeBronze, &testExpiration},
				&types.Membership{testUser3.String(), testUser.String(), types.MembershipTypeGold, &curTime},
			},
			expectedMemberships: []*types.Membership{
				&types.Membership{testUser3.String(), testUser.String(), types.MembershipTypeGold, &curTime},
			},
		},
		{
			name: "All memberships not expired are properly left on store",
			storedMemberships: []*types.Membership{
				&types.Membership{testUser.String(), testTsp.String(), types.MembershipTypeBronze, &curTime},
				&types.Membership{testUser2.String(), testTsp.String(), types.MembershipTypeGold, &curTime},
			},
			expectedMemberships: []*types.Membership{
				&types.Membership{testUser.String(), testTsp.String(), types.MembershipTypeBronze, &curTime},
				&types.Membership{testUser2.String(), testTsp.String(), types.MembershipTypeGold, &curTime},
			},
		},
		{
			name: "Black membership left in store and renewed",
			storedMemberships: []*types.Membership{
				&types.Membership{testUser.String(), testTsp.String(), types.MembershipTypeBlack, &testExpiration},
				&types.Membership{testUser2.String(), testTsp.String(), types.MembershipTypeGold, &testExpiration},
			},
			expectedMemberships: []*types.Membership{
				&types.Membership{testUser.String(), testTsp.String(), types.MembershipTypeBlack, &curTimePlusYear},
			},
		},
	}

	for _, test := range tests {

		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			ctx = ctx.WithBlockTime(curRemoveTime)
			for _, m := range test.storedMemberships {
				mOwner, _ := sdk.AccAddressFromBech32(m.Owner)
				mTspAddress, _ := sdk.AccAddressFromBech32(m.TspAddress)
				err := k.AssignMembership(ctx, mOwner, m.MembershipType, mTspAddress, *m.ExpiryAt)
				require.NoError(t, err)
			}
			err2 := k.RemoveExpiredMemberships(ctx)
			require.NoError(t, err2)
			ms := k.GetMemberships(ctx)
			require.Equal(t, test.expectedMemberships, ms)
		})
	}
}
