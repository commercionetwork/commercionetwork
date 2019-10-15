package keeper

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/stretchr/testify/assert"
)

// ---------------------
// --- Invites
// ---------------------

func TestKeeper_InviteUser_NoInvite(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)

	err := k.InviteUser(ctx, TestUser, TestUser2)
	assert.Nil(t, err)

	var invite types.Invite
	accreditationBz := store.Get(k.getInviteStoreKey(TestUser))
	cdc.MustUnmarshalBinaryBare(accreditationBz, &invite)

	assert.Equal(t, TestUser, invite.User)
	assert.Equal(t, TestUser2, invite.Sender)
	assert.False(t, invite.Rewarded)
}

func TestKeeper_InviteUser_ExistentInvite(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	existingAccredit := types.Invite{User: TestUser, Sender: TestUser2, Rewarded: false}

	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getInviteStoreKey(TestUser), cdc.MustMarshalBinaryBare(existingAccredit))

	err := k.InviteUser(ctx, TestUser, TestUser2)
	assert.NotNil(t, err)
}

func TestKeeper_GetInvite_NoInvite(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()

	_, found := k.GetInvite(ctx, TestUser)
	assert.False(t, found)
}

func TestKeeper_GetInvite_ExistingInvite(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	expected := types.Invite{User: TestUser, Sender: TestUser2, Rewarded: false}

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
	_, ctx, _, _, k := GetTestInput()

	invites := k.GetInvites(ctx)
	assert.Empty(t, invites)
}

func TestKeeper_GetInvites_NonEmptyList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	inv1 := types.Invite{Sender: TestUser2, User: TestUser, Rewarded: false}
	inv2 := types.Invite{Sender: TestUser2, User: TestUser2, Rewarded: false}

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

// ---------------------
// --- Signers
// ---------------------

func TestKeeper_AddTrustedServiceProvider_EmptyList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	k.AddTrustedServiceProvider(ctx, TestTsp)

	var signers ctypes.Addresses
	signersBz := store.Get([]byte(types.TrustedSignersStoreKey))
	cdc.MustUnmarshalBinaryBare(signersBz, &signers)

	assert.Equal(t, 1, len(signers))
	assert.Contains(t, signers, TestTsp)
}

func TestKeeper_AddTrustedServiceProvider_ExistingList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	signers := ctypes.Addresses{TestTsp}
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&signers))

	k.AddTrustedServiceProvider(ctx, TestUser)

	var actual ctypes.Addresses
	actualBz := store.Get([]byte(types.TrustedSignersStoreKey))
	cdc.MustUnmarshalBinaryBare(actualBz, &actual)

	assert.Equal(t, 2, len(actual))
	assert.Contains(t, actual, TestTsp)
	assert.Contains(t, actual, TestUser)
}

func TestKeeper_GetTrustedServiceProviders_EmptyList(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)

	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	signers := k.GetTrustedServiceProviders(ctx)

	assert.Empty(t, signers)
}

func TestKeeper_GetTrustedServiceProviders_ExistingList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	signers := ctypes.Addresses{TestTsp, TestUser, TestUser2}
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&signers))

	actual := k.GetTrustedServiceProviders(ctx)
	assert.Equal(t, 3, len(actual))
	assert.Contains(t, actual, TestTsp)
	assert.Contains(t, actual, TestUser)
	assert.Contains(t, actual, TestUser2)
}

func TestKeeper_IsTrustedServiceProvider_EmptyList(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	assert.False(t, k.IsTrustedServiceProvider(ctx, TestTsp))
	assert.False(t, k.IsTrustedServiceProvider(ctx, TestUser))
	assert.False(t, k.IsTrustedServiceProvider(ctx, TestUser2))
}

func TestKeeper_IsTrustedServiceProvider_ExistingList(t *testing.T) {
	cdc, ctx, _, _, k := GetTestInput()

	store := ctx.KVStore(k.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {

	}

	signers := ctypes.Addresses{TestUser, TestTsp}
	store.Set([]byte(types.TrustedSignersStoreKey), cdc.MustMarshalBinaryBare(&signers))

	assert.True(t, k.IsTrustedServiceProvider(ctx, TestUser))
	assert.True(t, k.IsTrustedServiceProvider(ctx, TestTsp))
	assert.False(t, k.IsTrustedServiceProvider(ctx, TestUser2))
}

// -------------------
// --- Memberships
// -------------------

func TestKeeper_getMembershipTokenId(t *testing.T) {
	_, _, _, _, k := GetTestInput()
	actual := k.getMembershipTokenId(TestUser)
	assert.Equal(t, fmt.Sprintf("membership-%s", TestUser.String()), actual)
}

func TestKeeper_getMembershipUri(t *testing.T) {
	_, _, _, _, k := GetTestInput()
	id := k.getMembershipTokenId(TestUser)
	actual := k.getMembershipUri("black", id)
	assert.Equal(t, fmt.Sprintf("membership:black:%s", id), actual)
}

func TestKeeper_AssignMembership_InvalidType(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()
	invalidTypes := []string{"", "grn", "slver", "   ", "blck"}
	for _, test := range invalidTypes {
		_, err := k.AssignMembership(ctx, TestUser, test)
		assert.NotNil(t, err)
	}
}

func TestKeeper_AssignMembership_NotExisting(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()
	tokenURI, err := k.AssignMembership(ctx, TestUser, "black")
	assert.Nil(t, err)

	expectedId := k.getMembershipTokenId(TestUser)
	assert.Equal(t, fmt.Sprintf("membership:black:%s", expectedId), tokenURI)
}

func TestKeeper_AssignMembership_Existing(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()
	memberships := []string{"black", "bronze", "silver", "gold", "black"}

	for _, membership := range memberships {
		tokenURI, err := k.AssignMembership(ctx, TestUser, membership)
		assert.Nil(t, err)

		expectedId := k.getMembershipTokenId(TestUser)
		expectedURI := k.getMembershipUri(membership, expectedId)
		assert.Equal(t, expectedURI, tokenURI)
	}
}

func TestKeeper_RemoveMembership_NotExisting(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()

	deleted, err := k.RemoveMembership(ctx, TestUser)
	assert.Nil(t, err)
	assert.True(t, deleted)
}

func TestKeeper_RemoveMembership_Existing(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()
	_, err := k.AssignMembership(ctx, TestUser, "black")
	assert.Nil(t, err)

	deleted, err := k.RemoveMembership(ctx, TestUser)
	assert.Nil(t, err)
	assert.True(t, deleted)

	_, found := k.GetMembership(ctx, TestUser)
	assert.False(t, found)
}

func TestKeeper_GetMembership_NotExisting(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()
	_, _ = k.RemoveMembership(ctx, TestUser)
	foundMembership, found := k.GetMembership(ctx, TestUser)
	assert.Nil(t, foundMembership)
	assert.False(t, found)
}

func TestKeeper_GetMembership_Existing(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()
	membershipType := "black"
	_, err := k.AssignMembership(ctx, TestUser, membershipType)
	assert.Nil(t, err)

	foundMembership, found := k.GetMembership(ctx, TestUser)
	assert.True(t, found)
	assert.Equal(t, TestUser, foundMembership.GetOwner())

	expectedId := k.getMembershipTokenId(TestUser)
	assert.Equal(t, expectedId, foundMembership.GetID())
	assert.Equal(t, k.getMembershipUri(membershipType, expectedId), foundMembership.GetTokenURI())
}

func TestKeeper_GetMembershipType(t *testing.T) {
	_, _, _, _, k := GetTestInput()

	id := "123"
	membershipType := "black"
	membership := nft.NewBaseNFT(id, TestUser, k.getMembershipUri(membershipType, id))

	assert.Equal(t, membershipType, k.GetMembershipType(&membership))
}

func TestKeeper_GetMembershipsSet_EmptySet(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()

	set := k.GetMembershipsSet(ctx)
	assert.Empty(t, set)
}

func TestKeeper_GetMembershipsSet_FilledSet(t *testing.T) {
	_, ctx, _, _, k := GetTestInput()

	first, err := sdk.AccAddressFromBech32("cosmos18v6sv92yrxdxck4hvw0gyfccu7dggweyrsqcx9")
	second, err := sdk.AccAddressFromBech32("cosmos1e2zdv8l45mstf8uexgcyk54g87jmrx8sjn0g24")
	third, err := sdk.AccAddressFromBech32("cosmos1e2cc3x2z7ku282kmh2x7jczylkajge0s2num6q")
	assert.Nil(t, err)

	_, err = k.AssignMembership(ctx, first, "black")
	_, err = k.AssignMembership(ctx, second, "silver")
	_, err = k.AssignMembership(ctx, third, "bronze")

	set := k.GetMembershipsSet(ctx)

	firstMembership := types.Membership{Owner: first, MembershipType: "black"}
	secondMembership := types.Membership{Owner: second, MembershipType: "silver"}
	thirdMembership := types.Membership{Owner: third, MembershipType: "bronze"}

	assert.Equal(t, 3, len(set))
	assert.Contains(t, set, firstMembership)
	assert.Contains(t, set, secondMembership)
	assert.Contains(t, set, thirdMembership)
}
