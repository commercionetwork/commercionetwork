package keeper_test

import (
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/memberships/types"
)

func TestKeeper_DepositIntoPool(t *testing.T) {
	tests := []struct {
		name         string
		existingPool sdk.Coins
		user         sdk.AccAddress
		userAmt      sdk.Coins
		deposit      sdk.Coins
		expectedPool sdk.Coins
		expectedUser sdk.Coins
		error        error
	}{
		{
			name:         "Empty deposit pool is incremented properly",
			existingPool: sdk.NewCoins(),
			user:         testUser,
			userAmt:      sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(100))),
			deposit:      sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(100))),
			expectedPool: sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(100))),
			expectedUser: sdk.Coins{},
		},
		{
			name:         "Existing deposit pool in incremented properly",
			user:         testUser,
			userAmt:      sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000))),
			existingPool: sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(100))),
			deposit:      sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000))),
			expectedPool: sdk.NewCoins(
				sdk.NewCoin("ucommercio", sdk.NewInt(1100)),
			),
			expectedUser: sdk.Coins{},
		},
		{
			name:         "Credits fails if user has not enough money",
			user:         testUser,
			userAmt:      sdk.NewCoins(),
			existingPool: sdk.NewCoins(),
			deposit:      sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000))),
			expectedPool: sdk.NewCoins(),
			expectedUser: sdk.NewCoins(),
			error:        sdkErr.Wrap(sdkErr.ErrInsufficientFunds, "insufficient account funds;  < 1000ucommercio"),
		},
		{
			name:         "deposit fails because not expressed in ucommercio",
			user:         testUser,
			userAmt:      sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100))),
			existingPool: sdk.NewCoins(),
			deposit:      sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(99))),
			expectedPool: sdk.NewCoins(),
			expectedUser: sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100))),
			error:        sdkErr.Wrap(sdkErr.ErrInsufficientFunds, "deposit into membership pool can only be expressed in ucommercio"),
		},
	}

	for _, test := range tests {
		ctx, bk, _, k := SetupTestInput()

		k.SupplyKeeper.SetSupply(ctx, supply.NewSupply(test.userAmt))
		_ = bk.SetCoins(ctx, test.user, test.userAmt)

		_ = k.SupplyKeeper.MintCoins(ctx, types.ModuleName, test.existingPool)

		err := k.DepositIntoPool(ctx, test.user, test.deposit)
		if test.error != nil {
			require.Equal(t, test.error.Error(), err.Error())
		} else {
			require.NoError(t, err)
		}

		require.True(t, test.expectedPool.IsEqual(k.GetPoolFunds(ctx)))
		require.True(t, test.expectedUser.IsEqual(bk.GetCoins(ctx, test.user)))
	}
}

func TestKeeper_GetPoolFunds(t *testing.T) {
	tests := []struct {
		name         string
		existingPool sdk.Coins
	}{
		{
			name:         "Empty pool is returned properly",
			existingPool: sdk.Coins{},
		},
		{
			name: "Non empty pool is returned properly",
			existingPool: sdk.NewCoins(
				sdk.NewCoin("uatom", sdk.NewInt(100)),
				sdk.NewCoin("ucommercio", sdk.NewInt(1000)),
			),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			_ = k.SupplyKeeper.MintCoins(ctx, types.ModuleName, test.existingPool)

			actual := k.GetPoolFunds(ctx)
			require.True(t, test.existingPool.IsEqual(actual))
		})
	}
}

func TestKeeper_DistributeReward(t *testing.T) {
	tests := []struct {
		name                      string
		inviteSenderMembership    string
		inviteRecipientMembership string
		invite                    types.Invite
		user                      sdk.AccAddress
		poolFunds                 sdk.Coins
		expectedInviteSenderAmt   sdk.Coins
		expectedPoolAmt           sdk.Coins
		error                     error
	}{
		{
			name:                   "Invite recipient without membership returns error",
			inviteSenderMembership: types.MembershipTypeBlack,
			invite:                 types.NewInvite(testInviteSender, testUser, "black"),
			user:                   testUser,
			error:                  sdkErr.Wrap(sdkErr.ErrUnauthorized, "Invite recipient does not have a membership"),
		},
		{
			name:                      "Insufficient pool funds greater than zero gives all reward available",
			inviteSenderMembership:    types.MembershipTypeBlack,
			inviteRecipientMembership: types.MembershipTypeGold,
			invite:                    types.NewInvite(testInviteSender, testUser, "black"),
			user:                      testUser,
			poolFunds:                 sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1000000)),
			expectedInviteSenderAmt:   sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1000000)),
			expectedPoolAmt:           sdk.Coins{},
		},
		{
			name:                      "Empty pool funds does not distribute anything",
			inviteSenderMembership:    types.MembershipTypeBlack,
			inviteRecipientMembership: types.MembershipTypeGold,
			invite:                    types.NewInvite(testInviteSender, testUser, "black"),
			user:                      testUser,
			poolFunds:                 sdk.Coins{},
			expectedInviteSenderAmt:   sdk.Coins{},
			expectedPoolAmt:           sdk.Coins{},
		},
		{
			name:                      "Enough pool funds",
			inviteSenderMembership:    types.MembershipTypeBlack,
			inviteRecipientMembership: types.MembershipTypeGold,
			invite:                    types.NewInvite(testInviteSender, testUser, "black"),
			user:                      testUser,
			poolFunds:                 sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1000000000000)),
			expectedInviteSenderAmt:   sdk.NewCoins(sdk.NewInt64Coin(testDenom, 2250000000)),
			expectedPoolAmt:           sdk.NewCoins(sdk.NewInt64Coin(testDenom, 997750000000)),
		},
		{
			name: "Invite is already rewarded",
			invite: types.Invite{
				Sender:           testInviteSender,
				SenderMembership: "black",
				User:             testUser,
				Status:           types.InviteStatusRewarded,
			},
			poolFunds:       sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1000000000000)),
			expectedPoolAmt: sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1000000000000)),
		},
		{
			name: "Invite is invalid",
			invite: types.Invite{
				Sender:           testInviteSender,
				SenderMembership: "black",
				User:             testUser,
				Status:           types.InviteStatusInvalid,
			},
			poolFunds:       sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1000000000000)),
			expectedPoolAmt: sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1000000000000)),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, bk, _, k := SetupTestInput()

			k.SaveInvite(ctx, test.invite)
			_ = k.AssignMembership(ctx, test.invite.Sender, test.inviteSenderMembership)
			_ = k.AssignMembership(ctx, test.invite.User, test.inviteRecipientMembership)

			_ = k.SupplyKeeper.MintCoins(ctx, types.ModuleName, test.poolFunds)

			err := k.DistributeReward(ctx, test.invite)
			if test.error != nil {
				require.Equal(t, test.error.Error(), err.Error())
			} else {
				require.NoError(t, err)
			}

			if test.error == nil && test.invite.Status != types.InviteStatusInvalid {
				storedInvite, _ := k.GetInvite(ctx, test.invite.User)
				require.Equal(t, storedInvite.Status, types.InviteStatusRewarded)
			}

			asd := test.expectedPoolAmt.IsEqual(k.GetPoolFunds(ctx))
			require.True(t, asd)
			require.True(t, test.expectedInviteSenderAmt.IsEqual(bk.GetCoins(ctx, test.invite.Sender)))
		})
	}
}
