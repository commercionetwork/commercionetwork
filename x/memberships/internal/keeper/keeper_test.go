package keeper_test

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/nft"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_AssignMembership_InvalidType(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	invalidTypes := []string{"", "grn", "slver", "   ", "blck"}
	for _, test := range invalidTypes {
		_, err := k.AssignMembership(ctx, testUser, test)
		assert.NotNil(t, err)
	}
}

func TestKeeper_AssignMembership_NotExisting(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	tokenURI, err := k.AssignMembership(ctx, testUser, "black")
	assert.Nil(t, err)

	expectedID := "membership-" + testUser.String()
	assert.Equal(t, fmt.Sprintf("membership:black:%s", expectedID), tokenURI)
}

func TestKeeper_AssignMembership_Existing(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	memberships := []string{"black", "bronze", "silver", "gold", "black"}

	for _, membership := range memberships {
		tokenURI, err := k.AssignMembership(ctx, testUser, membership)
		assert.Nil(t, err)

		expectedID := "membership-" + testUser.String()
		expectedURI := fmt.Sprintf("membership:%s:%s", membership, expectedID)
		assert.Equal(t, expectedURI, tokenURI)
	}
}

func TestKeeper_RemoveMembership_NotExisting(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	deleted, err := k.RemoveMembership(ctx, testUser)
	assert.Nil(t, err)
	assert.True(t, deleted)
}

func TestKeeper_RemoveMembership_Existing(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	_, err := k.AssignMembership(ctx, testUser, "black")
	assert.Nil(t, err)

	deleted, err := k.RemoveMembership(ctx, testUser)
	assert.Nil(t, err)
	assert.True(t, deleted)

	_, found := k.GetMembership(ctx, testUser)
	assert.False(t, found)
}

func TestKeeper_GetMembership_NotExisting(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	_, _ = k.RemoveMembership(ctx, testUser)
	foundMembership, found := k.GetMembership(ctx, testUser)
	assert.Nil(t, foundMembership)
	assert.False(t, found)
}

func TestKeeper_GetMembership_Existing(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	membershipType := "black"
	_, err := k.AssignMembership(ctx, testUser, membershipType)
	assert.Nil(t, err)

	foundMembership, found := k.GetMembership(ctx, testUser)
	assert.True(t, found)
	assert.Equal(t, testUser, foundMembership.GetOwner())

	expectedID := "membership-" + testUser.String()
	assert.Equal(t, expectedID, foundMembership.GetID())
	assert.Equal(t, fmt.Sprintf("membership:%s:%s", membershipType, expectedID), foundMembership.GetTokenURI())
}

func TestKeeper_GetMembershipType(t *testing.T) {
	_, _, _, k := SetupTestInput()

	id := "123"
	membershipType := "black"
	membership := nft.NewBaseNFT(id, testUser, fmt.Sprintf("membership:%s:%s", membershipType, id))

	assert.Equal(t, membershipType, k.GetMembershipType(&membership))
}

func TestKeeper_GetMembershipsSet_EmptySet(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	set := k.GetMembershipsSet(ctx)
	assert.Empty(t, set)
}

func TestKeeper_GetMembershipsSet_FilledSet(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	first, _ := sdk.AccAddressFromBech32("cosmos18v6sv92yrxdxck4hvw0gyfccu7dggweyrsqcx9")
	_, _ = k.AssignMembership(ctx, first, "black")

	second, _ := sdk.AccAddressFromBech32("cosmos1e2zdv8l45mstf8uexgcyk54g87jmrx8sjn0g24")
	_, _ = k.AssignMembership(ctx, second, "silver")

	third, _ := sdk.AccAddressFromBech32("cosmos1e2cc3x2z7ku282kmh2x7jczylkajge0s2num6q")
	_, _ = k.AssignMembership(ctx, third, "bronze")

	set := k.GetMembershipsSet(ctx)

	firstMembership := types.Membership{Owner: first, MembershipType: "black"}
	secondMembership := types.Membership{Owner: second, MembershipType: "silver"}
	thirdMembership := types.Membership{Owner: third, MembershipType: "bronze"}

	assert.Len(t, set, 3)
	assert.Contains(t, set, firstMembership)
	assert.Contains(t, set, secondMembership)
	assert.Contains(t, set, thirdMembership)
}
