package keeper

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/stretchr/testify/assert"
)

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
