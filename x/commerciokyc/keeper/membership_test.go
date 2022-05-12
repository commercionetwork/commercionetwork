package keeper

import (
	"fmt"
	"testing"
	"time"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
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
		expectedUserBalance  sdk.Coins
		expectedPoolBalance  sdk.Coins
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
		{
			name:      "Memberships matrix option reward not exists",
			invite:    types.Invite{Sender: testTsp.String(), SenderMembership: "bold", User: testUser2.String(), Status: uint64(types.InviteStatusPending)},
			mustError: true,
			//expectedError: sdkErr.Wrap(sdkErr.ErrInvalidRequest, "Invalid reward options"),
		},
		{
			name:      "Pool has zero tokens",
			invite:    types.Invite{Sender: testTsp.String(), SenderMembership: "gold", User: testUser2.String(), Status: uint64(types.InviteStatusPending)},
			pool:      coins,
			mustError: true,
			//expectedError: sdkErr.Wrap(sdkErr.ErrUnauthorized, "ABR pool has zero tokens"),
		},
		// TODO reward return different amount then expected. Amount of current pool
		/*{
			name:      "Account has not sufficient funds (pool is small then expected reward)",
			invite:    types.Invite{Sender: testTsp.String(), SenderMembership: "gold", User: testUser2.String(), Status: uint64(types.InviteStatusPending)},
			pool:      sdk.NewCoins(sdk.NewCoin(stakeDenom, sdk.NewInt(10000000))),
			mustError: true,
			//expectedError: sdkErr.Wrap(sdkErr.ErrUnauthorized, "ABR pool has zero tokens"),
			// could not move collateral amount to module account, 43478261ucommercio is smaller than 100000001ucommercio: insufficient funds
		},*/
		{
			name:                 "Account correctly rewarded",
			invite:               types.Invite{Sender: testTsp.String(), SenderMembership: "gold", User: testUser2.String(), Status: uint64(types.InviteStatusPending)},
			pool:                 sdk.NewCoins(sdk.NewCoin(stakeDenom, sdk.NewInt(1000000000000))),
			expectedInviteStatus: int64(types.InviteStatusRewarded),
			expectedUserBalance:  sdk.NewCoins(sdk.NewCoin(stableCreditDenom, sdk.NewInt(1750000000))),
			expectedPoolBalance:  sdk.NewCoins(sdk.NewCoin(stakeDenom, sdk.NewInt(998775000000))),
			mustError:            false,
		},
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
		userAccAddr, _ := sdk.AccAddressFromBech32(test.invite.User)
		//senderMembershipType := test.invite.SenderMembership
		//senderBalance := bk.GetAllBalances(ctx, senderAccAddr)
		//userBalance := bk.GetAllBalances(ctx, userAccAddr)
		k.SaveInvite(ctx, test.invite)

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
			inviteRet, _ := k.GetInvite(ctx, userAccAddr)
			require.Equal(t, test.expectedInviteStatus, int64(inviteRet.Status))
		}

		if test.expectedUserBalance != nil {
			senderBalance := bk.GetAllBalances(ctx, senderAccAddr)
			require.Equal(t, test.expectedUserBalance, senderBalance)
		}

		if test.expectedPoolBalance != nil {
			currentPoolBalance := k.GetLiquidityPoolAmount(ctx)
			require.Equal(t, test.expectedPoolBalance, currentPoolBalance)
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
				k.cdc.MustUnmarshal(i.Value(), &m)
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
			expectedExpiration: currentTime.Add(secondsPerYear),
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

func TestIsValidMembership(t *testing.T) {
	type args struct {
		expiredAt time.Time
		mt        string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Expired membership return error",
			args: args{
				expiredAt: time.Now(),
				mt:        types.MembershipTypeBronze,
			},
			want: false,
		},
		{
			name: "Black membership is always valid",
			args: args{
				expiredAt: time.Now(),
				mt:        types.MembershipTypeBlack,
			},
			want: true,
		},
		{
			name: "Valid membership",
			args: args{
				expiredAt: time.Now().Add(time.Hour * 24),
				mt:        types.MembershipTypeBlack,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, _ := SetupTestInput()
			if got := IsValidMembership(ctx, tt.args.expiredAt, tt.args.mt); got != tt.want {
				t.Errorf("IsValidMembership() = %v, want %v", got, tt.want)
			}
		})
	}
}
