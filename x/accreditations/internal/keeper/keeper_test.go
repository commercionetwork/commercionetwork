package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// ---------------------
// --- Invites
// ---------------------

func TestKeeper_InviteUser_NoInvite(t *testing.T) {
	ctx, cdc, _, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)

	err := k.InviteUser(ctx, TestUser, TestInviteSender)
	assert.Nil(t, err)

	var invite types.Invite
	accreditationBz := store.Get(k.getInviteStoreKey(TestUser))
	cdc.MustUnmarshalBinaryBare(accreditationBz, &invite)

	assert.Equal(t, TestUser, invite.User)
	assert.Equal(t, TestInviteSender, invite.Sender)
	assert.False(t, invite.Rewarded)
}

func TestKeeper_InviteUser_ExistentInvite(t *testing.T) {
	ctx, cdc, _, _, _, k := GetTestInput()

	existingAccredit := types.Invite{User: TestUser, Sender: TestInviteSender, Rewarded: false}

	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getInviteStoreKey(TestUser), cdc.MustMarshalBinaryBare(existingAccredit))

	err := k.InviteUser(ctx, TestUser, TestInviteSender)
	assert.NotNil(t, err)
}

func TestKeeper_GetInvite_NoInvite(t *testing.T) {
	ctx, _, _, _, _, k := GetTestInput()

	_, found := k.GetInvite(ctx, TestUser)
	assert.False(t, found)
}

func TestKeeper_GetInvite_ExistingInvite(t *testing.T) {
	ctx, cdc, _, _, _, k := GetTestInput()

	expected := types.Invite{User: TestUser, Sender: TestInviteSender, Rewarded: false}

	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getInviteStoreKey(TestUser), cdc.MustMarshalBinaryBare(expected))

	stored, found := k.GetInvite(ctx, TestUser)
	assert.True(t, found)
	assert.Equal(t, expected, stored)
}

// ---------------------
// --- Accreditations
// ---------------------

func TestKeeper_GetInvites_EmptyList(t *testing.T) {
	ctx, _, _, _, _, k := GetTestInput()

	invites := k.GetInvites(ctx)
	assert.Empty(t, invites)
}

func TestKeeper_GetInvites_NonEmptyList(t *testing.T) {
	ctx, cdc, _, _, _, k := GetTestInput()

	inv1 := types.Invite{Sender: TestInviteSender, User: TestUser, Rewarded: false}
	inv2 := types.Invite{Sender: TestInviteSender, User: TestUser2, Rewarded: false}

	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getInviteStoreKey(TestUser), cdc.MustMarshalBinaryBare(inv1))
	store.Set(k.getInviteStoreKey(TestUser2), cdc.MustMarshalBinaryBare(inv2))

	invites := k.GetInvites(ctx)
	assert.Equal(t, 2, len(invites))
	assert.Contains(t, invites, inv1)
	assert.Contains(t, invites, inv2)
}

// ---------------------
// --- Pool
// ---------------------

func TestKeeper_DepositIntoPool_EmptyPool(t *testing.T) {
	ctx, cdc, _, bankKeeper, _, accreditationKeeper := GetTestInput()

	coins := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000)))
	_ = bankKeeper.SetCoins(ctx, TestUser, coins)

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	deposit := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
	err := accreditationKeeper.DepositIntoPool(ctx, TestUser, deposit)
	assert.Nil(t, err)

	var pool sdk.Coins
	poolBz := store.Get([]byte(types.LiquidityPoolStoreKey))
	cdc.MustUnmarshalBinaryBare(poolBz, &pool)
	assert.Equal(t, deposit, pool)
}

func TestKeeper_DepositIntoPool_ExistingPool(t *testing.T) {
	ctx, cdc, _, bankKeeper, _, accreditationKeeper := GetTestInput()

	coins := sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000)))
	_ = bankKeeper.SetCoins(ctx, TestUser, coins)

	pool := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Set([]byte(types.LiquidityPoolStoreKey), cdc.MustMarshalBinaryBare(&pool))

	addition := sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000)))
	err := accreditationKeeper.DepositIntoPool(ctx, TestUser, addition)
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
	ctx, cdc, _, _, _, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	deposit := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
	accreditationKeeper.SetPoolFunds(ctx, deposit)

	var pool sdk.Coins
	poolBz := store.Get([]byte(types.LiquidityPoolStoreKey))
	cdc.MustUnmarshalBinaryBare(poolBz, &pool)
	assert.Equal(t, deposit, pool)
}

func TestKeeper_SetPoolFunds_ExistingPool(t *testing.T) {
	ctx, cdc, _, _, _, accreditationKeeper := GetTestInput()

	pool := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Set([]byte(types.LiquidityPoolStoreKey), cdc.MustMarshalBinaryBare(&pool))

	addition := sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000)))
	accreditationKeeper.SetPoolFunds(ctx, addition)

	var actual sdk.Coins
	actualBz := store.Get([]byte(types.LiquidityPoolStoreKey))
	cdc.MustUnmarshalBinaryBare(actualBz, &actual)

	expected := sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000)))
	assert.Equal(t, expected, actual)
}

func TestKeeper_GetPoolFunds_EmptyPool(t *testing.T) {
	ctx, _, _, _, _, accreditationKeeper := GetTestInput()

	pool := accreditationKeeper.GetPoolFunds(ctx)

	assert.Empty(t, pool)
}

func TestKeeper_GetPoolFunds_ExistingPool(t *testing.T) {
	ctx, cdc, _, _, _, accreditationKeeper := GetTestInput()

	expected := sdk.NewCoins(
		sdk.NewCoin("uatom", sdk.NewInt(100)),
		sdk.NewCoin("ucommercio", sdk.NewInt(1000)),
	)

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Set([]byte(types.LiquidityPoolStoreKey), cdc.MustMarshalBinaryBare(&expected))

	pool := accreditationKeeper.GetPoolFunds(ctx)

	assert.Equal(t, expected, pool)
}

// ---------------------
// --- Signers
// ---------------------

func TestKeeper_AddTrustedServiceProvider_EmptyList(t *testing.T) {
	ctx, cdc, _, _, _, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	accreditationKeeper.AddTrustedServiceProvider(ctx, TestTsp)

	var signers ctypes.Addresses
	signersBz := store.Get([]byte(types.TrustedSignersStoreKey))
	cdc.MustUnmarshalBinaryBare(signersBz, &signers)

	assert.Equal(t, 1, len(signers))
	assert.Contains(t, signers, TestTsp)
}

func TestKeeper_AddTrustedServiceProvider_ExistingList(t *testing.T) {
	ctx, cdc, _, _, _, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	signers := ctypes.Addresses{TestTsp}
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&signers))

	accreditationKeeper.AddTrustedServiceProvider(ctx, TestUser)

	var actual ctypes.Addresses
	actualBz := store.Get([]byte(types.TrustedSignersStoreKey))
	cdc.MustUnmarshalBinaryBare(actualBz, &actual)

	assert.Equal(t, 2, len(actual))
	assert.Contains(t, actual, TestTsp)
	assert.Contains(t, actual, TestUser)
}

func TestKeeper_GetTrustedServiceProviders_EmptyList(t *testing.T) {
	ctx, _, _, _, _, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	signers := accreditationKeeper.GetTrustedServiceProviders(ctx)

	assert.Empty(t, signers)
}

func TestKeeper_GetTrustedServiceProviders_ExistingList(t *testing.T) {
	ctx, cdc, _, _, _, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	signers := ctypes.Addresses{TestTsp, TestUser, TestInviteSender}
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&signers))

	actual := accreditationKeeper.GetTrustedServiceProviders(ctx)
	assert.Equal(t, 3, len(actual))
	assert.Contains(t, actual, TestTsp)
	assert.Contains(t, actual, TestUser)
	assert.Contains(t, actual, TestInviteSender)
}

func TestKeeper_IsTrustedServiceProvider_EmptyList(t *testing.T) {
	ctx, _, _, _, _, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	assert.False(t, accreditationKeeper.IsTrustedServiceProvider(ctx, TestTsp))
	assert.False(t, accreditationKeeper.IsTrustedServiceProvider(ctx, TestUser))
	assert.False(t, accreditationKeeper.IsTrustedServiceProvider(ctx, TestInviteSender))
}

func TestKeeper_IsTrustedServiceProvider_ExistingList(t *testing.T) {
	ctx, cdc, _, _, _, accreditationKeeper := GetTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	signers := ctypes.Addresses{TestUser, TestTsp}
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&signers))

	assert.True(t, accreditationKeeper.IsTrustedServiceProvider(ctx, TestUser))
	assert.True(t, accreditationKeeper.IsTrustedServiceProvider(ctx, TestTsp))
	assert.False(t, accreditationKeeper.IsTrustedServiceProvider(ctx, TestInviteSender))
}
