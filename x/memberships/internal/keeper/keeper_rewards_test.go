package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_DepositIntoPool_EmptyPool(t *testing.T) {
	cdc, ctx, bankK, _, k := GetTestInput()

	coins := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000)))
	_ = bankK.SetCoins(ctx, TestUser, coins)

	store := ctx.KVStore(k.StoreKey)

	deposit := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
	err := k.DepositIntoPool(ctx, TestUser, deposit)
	assert.Nil(t, err)

	var pool sdk.Coins
	poolBz := store.Get([]byte(types.LiquidityPoolStoreKey))
	cdc.MustUnmarshalBinaryBare(poolBz, &pool)
	assert.Equal(t, deposit, pool)
}

func TestKeeper_DepositIntoPool_ExistingPool(t *testing.T) {
	cdc, ctx, bankK, _, k := GetTestInput()

	coins := sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000)))
	_ = bankK.SetCoins(ctx, TestUser, coins)

	pool := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.LiquidityPoolStoreKey), cdc.MustMarshalBinaryBare(&pool))

	addition := sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000)))
	err := k.DepositIntoPool(ctx, TestUser, addition)
	assert.Nil(t, err)

	var actual sdk.Coins
	actualBz := store.Get([]byte(types.LiquidityPoolStoreKey))
	cdc.MustUnmarshalBinaryBare(actualBz, &actual)

	expected := sdk.NewCoins(
		sdk.NewCoin("uatom", sdk.NewInt(100)),
		sdk.NewCoin("ucommercio", sdk.NewInt(1000)),
	)
	assert.Equal(t, expected, actual)
}

func TestKeeper_SetPoolFunds_EmptyPool(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)

	deposit := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
	k.SetPoolFunds(ctx, deposit)

	var pool sdk.Coins
	poolBz := store.Get([]byte(types.LiquidityPoolStoreKey))
	cdc.MustUnmarshalBinaryBare(poolBz, &pool)
	assert.Equal(t, deposit, pool)
}

func TestKeeper_SetPoolFunds_ExistingPool(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	pool := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.LiquidityPoolStoreKey), cdc.MustMarshalBinaryBare(&pool))

	addition := sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000)))
	k.SetPoolFunds(ctx, addition)

	var actual sdk.Coins
	actualBz := store.Get([]byte(types.LiquidityPoolStoreKey))
	cdc.MustUnmarshalBinaryBare(actualBz, &actual)

	expected := sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000)))
	assert.Equal(t, expected, actual)
}

func TestKeeper_GetPoolFunds_EmptyPool(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()

	pool := k.GetPoolFunds(ctx)

	assert.Empty(t, pool)
}

func TestKeeper_GetPoolFunds_ExistingPool(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	expected := sdk.NewCoins(
		sdk.NewCoin("uatom", sdk.NewInt(100)),
		sdk.NewCoin("ucommercio", sdk.NewInt(1000)),
	)

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.LiquidityPoolStoreKey), cdc.MustMarshalBinaryBare(&expected))

	pool := k.GetPoolFunds(ctx)

	assert.Equal(t, expected, pool)
}

func TestKeeper_DistributeReward_InviteSenderWithoutMembership(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()

	// Setup
	invite := types.Invite{User: TestUser, Sender: testInviteSender, Rewarded: false}
	k.SaveInvite(ctx, invite)

	// Test
	err := k.DistributeReward(ctx, invite, types.MembershipTypeBronze)
	assert.Error(t, err)
	assert.Contains(t, err.Result().Log, "Invite sender does not have a membership")
}

func TestKeeper_DistributeReward_InviteRecipientWithoutMembership(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()

	// Setup
	invite := types.Invite{User: TestUser, Sender: testInviteSender, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, invite.Sender, types.MembershipTypeBlack)

	// Test
	err := k.DistributeReward(ctx, invite, types.MembershipTypeBronze)
	assert.Error(t, err)
	assert.Contains(t, err.Result().Log, "Invite recipient does not have a membership")
}

func TestKeeper_DistributeReward_InviteRecipientWrongMembershipType(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()

	// Setup
	invite := types.Invite{User: TestUser, Sender: testInviteSender, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, invite.Sender, types.MembershipTypeBlack)
	_, _ = k.AssignMembership(ctx, invite.User, types.MembershipTypeGold)

	// Test
	err := k.DistributeReward(ctx, invite, types.MembershipTypeBronze)
	assert.Error(t, err)
	assert.Contains(t, err.Result().Log, "Invite recipient membership is not the same as the bought membership")
}

func TestKeeper_DistributeReward_EnoughPoolAmount(t *testing.T) {
	_, ctx, bankK, _, k := GetTestInput()

	// Setup
	invite := types.Invite{User: TestUser, Sender: testInviteSender, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, invite.Sender, types.MembershipTypeBlack)
	_, _ = k.AssignMembership(ctx, invite.User, types.MembershipTypeGold)

	poolFunds := sdk.NewCoins(sdk.NewInt64Coin(TestStableCreditsDenom, 1000000000000))
	k.SetPoolFunds(ctx, poolFunds)

	err := k.DistributeReward(ctx, invite, types.MembershipTypeGold)
	assert.NoError(t, err)

	expectedRewards := sdk.NewCoins(sdk.NewInt64Coin(TestStableCreditsDenom, 2250000000))
	expectedRemainingPool := poolFunds.Sub(expectedRewards)

	assert.Equal(t, expectedRewards, bankK.GetCoins(ctx, invite.Sender))
	assert.Equal(t, expectedRemainingPool, k.GetPoolFunds(ctx))

	storedInvite, _ := k.GetInvite(ctx, invite.User)
	assert.True(t, storedInvite.Rewarded)
}

func TestKeeper_DistributeReward_InsufficientPoolFundsGreaterThanZero(t *testing.T) {
	_, ctx, bankK, _, k := GetTestInput()

	// Setup
	invite := types.Invite{User: TestUser, Sender: testInviteSender, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, invite.Sender, types.MembershipTypeBlack)
	_, _ = k.AssignMembership(ctx, invite.User, types.MembershipTypeGold)

	poolFunds := sdk.NewCoins(sdk.NewInt64Coin(TestStableCreditsDenom, 1000000))
	k.SetPoolFunds(ctx, poolFunds)

	err := k.DistributeReward(ctx, invite, types.MembershipTypeGold)
	assert.NoError(t, err)

	assert.Equal(t, poolFunds, bankK.GetCoins(ctx, invite.Sender))
	assert.Empty(t, k.GetPoolFunds(ctx))

	storedInvite, _ := k.GetInvite(ctx, invite.User)
	assert.True(t, storedInvite.Rewarded)
}

func TestKeeper_DistributeReward_InsufficientPoolFundsZero(t *testing.T) {
	_, ctx, bankK, _, k := GetTestInput()

	// Setup
	invite := types.Invite{User: TestUser, Sender: testInviteSender, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, invite.Sender, types.MembershipTypeBlack)
	_, _ = k.AssignMembership(ctx, invite.User, types.MembershipTypeGold)
	k.SetPoolFunds(ctx, nil)

	err := k.DistributeReward(ctx, invite, types.MembershipTypeGold)
	assert.NoError(t, err)

	assert.Empty(t, bankK.GetCoins(ctx, invite.Sender))
	assert.Empty(t, k.GetPoolFunds(ctx))

	storedInvite, _ := k.GetInvite(ctx, invite.User)
	assert.True(t, storedInvite.Rewarded)
}
