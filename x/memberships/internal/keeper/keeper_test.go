package keeper

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_AddTrustedMinter(t *testing.T) {
	_, ctx, k := SetupTestInput()
	membershipsStore := ctx.KVStore(k.StoreKey)
	storeData := membershipsStore.Get([]byte(types.TrustedMinterPrefix + TestSignerAddress.String()))
	assert.Nil(t, storeData)

	k.AddTrustedMinter(ctx, TestSignerAddress)

	afterOpLen := membershipsStore.Get([]byte(types.TrustedMinterPrefix + TestSignerAddress.String()))
	assert.Equal(t, TestSignerAddress.Bytes(), afterOpLen)
}

func TestKeeper_GetTrustedMinters(t *testing.T) {
	_, ctx, k := SetupTestInput()
	membershipsStore := ctx.KVStore(k.StoreKey)
	membershipsStore.Set([]byte(types.TrustedMinterPrefix+TestSignerAddress.String()), TestSignerAddress.Bytes())
	membershipsStore.Set([]byte(types.TrustedMinterPrefix+TestUserAddress.String()), TestUserAddress.Bytes())

	minters := k.GetTrustedMinters(ctx)

	assert.True(t, minters.Contains(TestSignerAddress))
	assert.True(t, minters.Contains(TestUserAddress))
}

func TestKeeper_getMembershipTokenId(t *testing.T) {
	_, _, k := SetupTestInput()
	actual := k.getMembershipTokenId(TestSignerAddress)
	assert.Equal(t, "membership-cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0", actual)
}

func TestKeeper_getMembershipUri(t *testing.T) {
	_, _, k := SetupTestInput()
	id := k.getMembershipTokenId(TestSignerAddress)
	actual := k.getMembershipUri("black", id)
	assert.Equal(t, fmt.Sprintf("membership:black:%s", id), actual)
}

func TestKeeper_AssignMembership_InvalidType(t *testing.T) {
	_, ctx, k := SetupTestInput()
	invalidTypes := []string{"", "grn", "slver", "   ", "blck"}
	for _, test := range invalidTypes {
		_, err := k.AssignMembership(ctx, TestSignerAddress, test)
		assert.NotNil(t, err)
	}
}

func TestKeeper_AssignMembership_NotExisting(t *testing.T) {
	_, ctx, k := SetupTestInput()
	tokenURI, err := k.AssignMembership(ctx, TestSignerAddress, "green")
	assert.Nil(t, err)

	expectedId := k.getMembershipTokenId(TestSignerAddress)
	assert.Equal(t, fmt.Sprintf("membership:green:%s", expectedId), tokenURI)
}

func TestKeeper_AssignMembership_Existing(t *testing.T) {
	_, ctx, k := SetupTestInput()
	memberships := []string{"green", "bronze", "silver", "gold", "black"}

	for _, membership := range memberships {
		tokenURI, err := k.AssignMembership(ctx, TestSignerAddress, membership)
		assert.Nil(t, err)

		expectedId := k.getMembershipTokenId(TestSignerAddress)
		expectedURI := k.getMembershipUri(membership, expectedId)
		assert.Equal(t, expectedURI, tokenURI)
	}
}

func TestKeeper_RemoveMembership_NotExisting(t *testing.T) {
	_, ctx, k := SetupTestInput()

	deleted, err := k.RemoveMembership(ctx, TestSignerAddress)
	assert.Nil(t, err)
	assert.True(t, deleted)
}

func TestKeeper_RemoveMembership_Existing(t *testing.T) {
	_, ctx, k := SetupTestInput()
	_, err := k.AssignMembership(ctx, TestSignerAddress, "green")
	assert.Nil(t, err)

	deleted, err := k.RemoveMembership(ctx, TestSignerAddress)
	assert.Nil(t, err)
	assert.True(t, deleted)

	_, found := k.GetMembership(ctx, TestSignerAddress)
	assert.False(t, found)
}

func TestKeeper_GetMembership_NotExisting(t *testing.T) {
	_, ctx, k := SetupTestInput()
	_, _ = k.RemoveMembership(ctx, TestSignerAddress)
	foundMembership, found := k.GetMembership(ctx, TestSignerAddress)
	assert.Nil(t, foundMembership)
	assert.False(t, found)
}

func TestKeeper_GetMembership_Existing(t *testing.T) {
	_, ctx, k := SetupTestInput()
	membershipType := "green"
	_, err := k.AssignMembership(ctx, TestSignerAddress, membershipType)
	assert.Nil(t, err)

	foundMembership, found := k.GetMembership(ctx, TestSignerAddress)
	assert.True(t, found)
	assert.Equal(t, TestSignerAddress, foundMembership.GetOwner())

	expectedId := k.getMembershipTokenId(TestSignerAddress)
	assert.Equal(t, expectedId, foundMembership.GetID())
	assert.Equal(t, k.getMembershipUri(membershipType, expectedId), foundMembership.GetTokenURI())
}

func TestKeeper_GetMembershipType(t *testing.T) {
	_, _, k := SetupTestInput()

	id := "123"
	membershipType := "green"
	membership := nft.NewBaseNFT(id, TestSignerAddress, k.getMembershipUri(membershipType, id))

	assert.Equal(t, membershipType, k.GetMembershipType(&membership))
}

// ----------------------
// --- Genesis utils
// ----------------------

func clearNFTs(t *testing.T) {
	_, ctx, k := SetupTestInput()
	if collection, found := k.NftKeeper.GetCollection(ctx, types.NftDenom); found {
		for _, membershipNft := range collection.NFTs {
			err := k.NftKeeper.DeleteNFT(ctx, types.NftDenom, membershipNft.GetID())
			assert.Nil(t, err)
		}
	}
}

func TestKeeper_GetMembershipsSet_EmptySet(t *testing.T) {
	//clearNFTs(t)
	_, ctx, k := SetupTestInput()

	set := k.GetMembershipsSet(ctx)
	assert.Empty(t, set)
}

func TestKeeper_GetMembershipsSet_FilledSet(t *testing.T) {
	//clearNFTs(t)
	_, ctx, k := SetupTestInput()

	first, err := sdk.AccAddressFromBech32("cosmos18v6sv92yrxdxck4hvw0gyfccu7dggweyrsqcx9")
	second, err := sdk.AccAddressFromBech32("cosmos1e2zdv8l45mstf8uexgcyk54g87jmrx8sjn0g24")
	third, err := sdk.AccAddressFromBech32("cosmos1e2cc3x2z7ku282kmh2x7jczylkajge0s2num6q")
	assert.Nil(t, err)

	_, err = k.AssignMembership(ctx, first, "green")
	_, err = k.AssignMembership(ctx, second, "silver")
	_, err = k.AssignMembership(ctx, third, "bronze")

	set := k.GetMembershipsSet(ctx)

	firstMembership := types.Membership{Owner: first, MembershipType: "green"}
	secondMembership := types.Membership{Owner: second, MembershipType: "silver"}
	thirdMembership := types.Membership{Owner: third, MembershipType: "bronze"}

	assert.Equal(t, 3, len(set))
	assert.Contains(t, set, firstMembership)
	assert.Contains(t, set, secondMembership)
	assert.Contains(t, set, thirdMembership)
}
