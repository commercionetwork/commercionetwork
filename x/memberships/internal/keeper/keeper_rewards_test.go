package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_DepositIntoPool_EmptyPool(t *testing.T) {
	ctx, bankK, _, k := SetupTestInput()

	coins := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000)))
	_ = bankK.SetCoins(ctx, testUser, coins)

	deposit := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
	err := k.DepositIntoPool(ctx, testUser, deposit)
	assert.Nil(t, err)
	assert.Equal(t, deposit, k.GetPoolFunds(ctx))
}

func TestKeeper_DepositIntoPool_ExistingPool(t *testing.T) {
	ctx, bankK, _, k := SetupTestInput()

	coins := sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000)))
	_ = bankK.SetCoins(ctx, testUser, coins)

	pool := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
	_ = k.supplyKeeper.MintCoins(ctx, types.ModuleName, pool)

	addition := sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000)))
	err := k.DepositIntoPool(ctx, testUser, addition)
	assert.Nil(t, err)

	expected := sdk.NewCoins(
		sdk.NewCoin("uatom", sdk.NewInt(100)),
		sdk.NewCoin("ucommercio", sdk.NewInt(1000)),
	)
	assert.Equal(t, expected, k.GetPoolFunds(ctx))
}

func TestKeeper_GetPoolFunds_EmptyPool(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	pool := k.GetPoolFunds(ctx)

	assert.Empty(t, pool)
}

func TestKeeper_GetPoolFunds_ExistingPool(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	expected := sdk.NewCoins(
		sdk.NewCoin("uatom", sdk.NewInt(100)),
		sdk.NewCoin("ucommercio", sdk.NewInt(1000)),
	)
	_ = k.supplyKeeper.MintCoins(ctx, types.ModuleName, expected)

	pool := k.GetPoolFunds(ctx)
	assert.Equal(t, expected, pool)
}

func TestKeeper_DistributeReward_InviteSenderWithoutMembership(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	// Setup
	invite := types.Invite{User: testUser, Sender: testInviteSender, Rewarded: false}
	k.SaveInvite(ctx, invite)

	// Test
	err := k.DistributeReward(ctx, invite, types.MembershipTypeBronze)
	assert.Error(t, err)
	assert.Contains(t, err.Result().Log, "Invite sender does not have a membership")
}

func TestKeeper_DistributeReward_InviteRecipientWithoutMembership(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	// Setup
	invite := types.Invite{User: testUser, Sender: testInviteSender, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, invite.Sender, types.MembershipTypeBlack)

	// Test
	err := k.DistributeReward(ctx, invite, types.MembershipTypeBronze)
	assert.Error(t, err)
	assert.Contains(t, err.Result().Log, "Invite recipient does not have a membership")
}

func TestKeeper_DistributeReward_InviteRecipientWrongMembershipType(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	// Setup
	invite := types.Invite{User: testUser, Sender: testInviteSender, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, invite.Sender, types.MembershipTypeBlack)
	_, _ = k.AssignMembership(ctx, invite.User, types.MembershipTypeGold)

	// Test
	err := k.DistributeReward(ctx, invite, types.MembershipTypeBronze)
	assert.Error(t, err)
	assert.Contains(t, err.Result().Log, "Invite recipient membership is not the same as the bought membership")
}

func TestKeeper_DistributeReward_EnoughPoolAmount(t *testing.T) {
	ctx, bankK, _, k := SetupTestInput()

	// Setup
	invite := types.Invite{User: testUser, Sender: testInviteSender, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, invite.Sender, types.MembershipTypeBlack)
	_, _ = k.AssignMembership(ctx, invite.User, types.MembershipTypeGold)

	poolFunds := sdk.NewCoins(sdk.NewInt64Coin(testStableCreditsDenom, 1000000000000))
	_ = k.supplyKeeper.MintCoins(ctx, types.ModuleName, poolFunds)

	err := k.DistributeReward(ctx, invite, types.MembershipTypeGold)
	assert.NoError(t, err)

	expectedRewards := sdk.NewCoins(sdk.NewInt64Coin(testStableCreditsDenom, 2250000000))
	expectedRemainingPool := poolFunds.Sub(expectedRewards)

	assert.Equal(t, expectedRewards, bankK.GetCoins(ctx, invite.Sender))
	assert.Equal(t, expectedRemainingPool, k.GetPoolFunds(ctx))

	storedInvite, _ := k.GetInvite(ctx, invite.User)
	assert.True(t, storedInvite.Rewarded)
}

func TestKeeper_DistributeReward_InsufficientPoolFundsGreaterThanZero(t *testing.T) {
	ctx, bankK, _, k := SetupTestInput()

	// Setup
	invite := types.Invite{User: testUser, Sender: testInviteSender, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, invite.Sender, types.MembershipTypeBlack)
	_, _ = k.AssignMembership(ctx, invite.User, types.MembershipTypeGold)

	poolFunds := sdk.NewCoins(sdk.NewInt64Coin(testStableCreditsDenom, 1000000))
	_ = k.supplyKeeper.MintCoins(ctx, types.ModuleName, poolFunds)

	err := k.DistributeReward(ctx, invite, types.MembershipTypeGold)
	assert.NoError(t, err)

	assert.Equal(t, poolFunds, bankK.GetCoins(ctx, invite.Sender))
	assert.Empty(t, k.GetPoolFunds(ctx))

	storedInvite, _ := k.GetInvite(ctx, invite.User)
	assert.True(t, storedInvite.Rewarded)
}

func TestKeeper_DistributeReward_InsufficientPoolFundsZero(t *testing.T) {
	ctx, bankK, _, k := SetupTestInput()

	// Setup
	invite := types.Invite{User: testUser, Sender: testInviteSender, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, invite.Sender, types.MembershipTypeBlack)
	_, _ = k.AssignMembership(ctx, invite.User, types.MembershipTypeGold)

	err := k.DistributeReward(ctx, invite, types.MembershipTypeGold)
	assert.NoError(t, err)

	assert.Empty(t, bankK.GetCoins(ctx, invite.Sender))
	assert.Empty(t, k.GetPoolFunds(ctx))

	storedInvite, _ := k.GetInvite(ctx, invite.User)
	assert.True(t, storedInvite.Rewarded)
}
