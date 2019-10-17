package keeper

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/stretchr/testify/assert"
)

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
