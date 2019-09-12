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
	membershipsStore := TestUtils.Ctx.KVStore(TestUtils.MembershipKeeper.StoreKey)
	storeData := membershipsStore.Get([]byte(types.TrustedMinterPrefix + TestSignerAddress.String()))
	assert.Nil(t, storeData)

	TestUtils.MembershipKeeper.AddTrustedMinter(TestUtils.Ctx, TestSignerAddress)

	afterOpLen := membershipsStore.Get([]byte(types.TrustedMinterPrefix + TestSignerAddress.String()))
	assert.Equal(t, TestSignerAddress.Bytes(), afterOpLen)
}

func TestKeeper_GetTrustedMinters(t *testing.T) {
	membershipsStore := TestUtils.Ctx.KVStore(TestUtils.MembershipKeeper.StoreKey)
	membershipsStore.Set([]byte(types.TrustedMinterPrefix+TestSignerAddress.String()), TestSignerAddress.Bytes())
	membershipsStore.Set([]byte(types.TrustedMinterPrefix+TestUserAddress.String()), TestUserAddress.Bytes())

	minters := TestUtils.MembershipKeeper.GetTrustedMinters(TestUtils.Ctx)

	assert.True(t, minters.Contains(TestSignerAddress))
	assert.True(t, minters.Contains(TestUserAddress))
}

func TestKeeper_getMembershipTokenId(t *testing.T) {
	actual := TestUtils.MembershipKeeper.getMembershipTokenId(TestSignerAddress)
	assert.Equal(t, "membership-cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0", actual)
}

func TestKeeper_getMembershipUri(t *testing.T) {
	id := TestUtils.MembershipKeeper.getMembershipTokenId(TestSignerAddress)
	actual := TestUtils.MembershipKeeper.getMembershipUri("black", id)
	assert.Equal(t, fmt.Sprintf("membership:black:%s", id), actual)
}

func TestKeeper_AssignMembership_InvalidType(t *testing.T) {
	invalidTypes := []string{"", "grn", "slver", "   ", "blck"}
	for _, test := range invalidTypes {
		_, err := TestUtils.MembershipKeeper.AssignMembership(TestUtils.Ctx, TestSignerAddress, test)
		assert.NotNil(t, err)
	}
}

func TestKeeper_AssignMembership_NotExisting(t *testing.T) {
	tokenURI, err := TestUtils.MembershipKeeper.AssignMembership(TestUtils.Ctx, TestSignerAddress, "green")
	assert.Nil(t, err)

	expectedId := TestUtils.MembershipKeeper.getMembershipTokenId(TestSignerAddress)
	assert.Equal(t, fmt.Sprintf("membership:green:%s", expectedId), tokenURI)
}

func TestKeeper_AssignMembership_Existing(t *testing.T) {
	memberships := []string{"green", "bronze", "silver", "gold", "black"}

	for _, membership := range memberships {
		tokenURI, err := TestUtils.MembershipKeeper.AssignMembership(TestUtils.Ctx, TestSignerAddress, membership)
		assert.Nil(t, err)

		expectedId := TestUtils.MembershipKeeper.getMembershipTokenId(TestSignerAddress)
		expectedURI := TestUtils.MembershipKeeper.getMembershipUri(membership, expectedId)
		assert.Equal(t, expectedURI, tokenURI)
	}
}

func TestKeeper_RemoveMembership_NotExisting(t *testing.T) {
	// pre-delete to clean up the cache
	_, _ = TestUtils.MembershipKeeper.RemoveMembership(TestUtils.Ctx, TestSignerAddress)

	// real deletion
	deleted, err := TestUtils.MembershipKeeper.RemoveMembership(TestUtils.Ctx, TestSignerAddress)
	assert.Nil(t, err)
	assert.True(t, deleted)
}

func TestKeeper_RemoveMembership_Existing(t *testing.T) {
	_, err := TestUtils.MembershipKeeper.AssignMembership(TestUtils.Ctx, TestSignerAddress, "green")
	assert.Nil(t, err)

	deleted, err := TestUtils.MembershipKeeper.RemoveMembership(TestUtils.Ctx, TestSignerAddress)
	assert.Nil(t, err)
	assert.True(t, deleted)

	_, found := TestUtils.MembershipKeeper.GetMembership(TestUtils.Ctx, TestSignerAddress)
	assert.False(t, found)
}

func TestKeeper_GetMembership_NotExisting(t *testing.T) {
	_, _ = TestUtils.MembershipKeeper.RemoveMembership(TestUtils.Ctx, TestSignerAddress)
	foundMembership, found := TestUtils.MembershipKeeper.GetMembership(TestUtils.Ctx, TestSignerAddress)
	assert.Nil(t, foundMembership)
	assert.False(t, found)
}

func TestKeeper_GetMembership_Existing(t *testing.T) {
	membershipType := "green"
	_, err := TestUtils.MembershipKeeper.AssignMembership(TestUtils.Ctx, TestSignerAddress, membershipType)
	assert.Nil(t, err)

	foundMembership, found := TestUtils.MembershipKeeper.GetMembership(TestUtils.Ctx, TestSignerAddress)
	assert.True(t, found)
	assert.Equal(t, TestSignerAddress, foundMembership.GetOwner())

	expectedId := TestUtils.MembershipKeeper.getMembershipTokenId(TestSignerAddress)
	assert.Equal(t, expectedId, foundMembership.GetID())
	assert.Equal(t, TestUtils.MembershipKeeper.getMembershipUri(membershipType, expectedId), foundMembership.GetTokenURI())
}

func TestKeeper_GetMembershipType(t *testing.T) {
	id := "123"
	membershipType := "green"
	membership := nft.NewBaseNFT(id, TestSignerAddress, TestUtils.MembershipKeeper.getMembershipUri(membershipType, id))

	assert.Equal(t, membershipType, TestUtils.MembershipKeeper.GetMembershipType(&membership))
}

// ----------------------
// --- Genesis utils
// ----------------------

func clearNFTs(t *testing.T) {
	if collection, found := TestUtils.MembershipKeeper.NftKeeper.GetCollection(TestUtils.Ctx, types.NftDenom); found {
		for _, membershipNft := range collection.NFTs {
			err := TestUtils.MembershipKeeper.NftKeeper.DeleteNFT(TestUtils.Ctx, types.NftDenom, membershipNft.GetID())
			assert.Nil(t, err)
		}
	}
}

func TestKeeper_GetMembershipsSet_EmptySet(t *testing.T) {
	clearNFTs(t)

	set := TestUtils.MembershipKeeper.GetMembershipsSet(TestUtils.Ctx)
	assert.Empty(t, set)
}

func TestKeeper_GetMembershipsSet_FilledSet(t *testing.T) {
	clearNFTs(t)

	first, err := sdk.AccAddressFromBech32("cosmos18v6sv92yrxdxck4hvw0gyfccu7dggweyrsqcx9")
	second, err := sdk.AccAddressFromBech32("cosmos1e2zdv8l45mstf8uexgcyk54g87jmrx8sjn0g24")
	third, err := sdk.AccAddressFromBech32("cosmos1e2cc3x2z7ku282kmh2x7jczylkajge0s2num6q")
	assert.Nil(t, err)

	_, err = TestUtils.MembershipKeeper.AssignMembership(TestUtils.Ctx, first, "green")
	_, err = TestUtils.MembershipKeeper.AssignMembership(TestUtils.Ctx, second, "silver")
	_, err = TestUtils.MembershipKeeper.AssignMembership(TestUtils.Ctx, third, "bronze")

	set := TestUtils.MembershipKeeper.GetMembershipsSet(TestUtils.Ctx)

	firstMembership := types.Membership{Owner: first, MembershipType: "green"}
	secondMembership := types.Membership{Owner: second, MembershipType: "silver"}
	thirdMembership := types.Membership{Owner: third, MembershipType: "bronze"}

	assert.Equal(t, 3, len(set))
	assert.Contains(t, set, firstMembership)
	assert.Contains(t, set, secondMembership)
	assert.Contains(t, set, thirdMembership)
}
