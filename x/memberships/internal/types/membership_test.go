package types_test

import (
	"strings"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestMembership_Equals(t *testing.T) {
	owner, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	membership := types.NewMembership(types.MembershipTypeBronze, owner)

	assert.False(t, membership.Equals(types.NewMembership(types.MembershipTypeSilver, membership.Owner)))
	assert.False(t, membership.Equals(types.NewMembership(types.MembershipTypeGold, membership.Owner)))
	assert.False(t, membership.Equals(types.NewMembership(types.MembershipTypeBlack, membership.Owner)))
	assert.False(t, membership.Equals(types.NewMembership(membership.MembershipType, sdk.AccAddress{})))
	assert.True(t, membership.Equals(membership))
}

func TestIsMembershipTypeValid(t *testing.T) {
	assert.True(t, types.IsMembershipTypeValid(types.MembershipTypeBronze))
	assert.True(t, types.IsMembershipTypeValid(types.MembershipTypeSilver))
	assert.True(t, types.IsMembershipTypeValid(types.MembershipTypeGold))
	assert.True(t, types.IsMembershipTypeValid(types.MembershipTypeBlack))
	assert.False(t, types.IsMembershipTypeValid(strings.ToUpper(types.MembershipTypeBronze)))
}

func TestCanUpgrade(t *testing.T) {
	assert.False(t, types.CanUpgrade(types.MembershipTypeBronze, types.MembershipTypeBronze))
	assert.True(t, types.CanUpgrade(types.MembershipTypeBronze, types.MembershipTypeSilver))
	assert.True(t, types.CanUpgrade(types.MembershipTypeBronze, types.MembershipTypeGold))
	assert.True(t, types.CanUpgrade(types.MembershipTypeBronze, types.MembershipTypeBlack))

	assert.False(t, types.CanUpgrade(types.MembershipTypeSilver, types.MembershipTypeBronze))
	assert.False(t, types.CanUpgrade(types.MembershipTypeSilver, types.MembershipTypeSilver))
	assert.True(t, types.CanUpgrade(types.MembershipTypeSilver, types.MembershipTypeGold))
	assert.True(t, types.CanUpgrade(types.MembershipTypeSilver, types.MembershipTypeBlack))

	assert.False(t, types.CanUpgrade(types.MembershipTypeGold, types.MembershipTypeBronze))
	assert.False(t, types.CanUpgrade(types.MembershipTypeGold, types.MembershipTypeSilver))
	assert.False(t, types.CanUpgrade(types.MembershipTypeGold, types.MembershipTypeGold))
	assert.True(t, types.CanUpgrade(types.MembershipTypeGold, types.MembershipTypeBlack))

	assert.False(t, types.CanUpgrade(types.MembershipTypeBlack, types.MembershipTypeBronze))
	assert.False(t, types.CanUpgrade(types.MembershipTypeBlack, types.MembershipTypeSilver))
	assert.False(t, types.CanUpgrade(types.MembershipTypeBlack, types.MembershipTypeGold))
	assert.False(t, types.CanUpgrade(types.MembershipTypeBlack, types.MembershipTypeBlack))
}

func TestMemberships_AppendIfMissing(t *testing.T) {
	owner, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	membership1 := types.NewMembership(types.MembershipTypeBronze, owner)
	membership2 := types.NewMembership(types.MembershipTypeSilver, owner)

	tests := []struct {
		name             string
		memberships      types.Memberships
		membership       types.Membership
		shouldBeAppended bool
	}{
		{
			name:             "Membership is appended to empty slice",
			memberships:      types.Memberships{},
			membership:       membership1,
			shouldBeAppended: true,
		},
		{
			name:             "Membership is appended to existing list",
			memberships:      types.Memberships{membership1},
			membership:       membership2,
			shouldBeAppended: true,
		},
		{
			name:             "Membership is not appended when existing",
			memberships:      types.Memberships{membership1, membership2},
			membership:       membership1,
			shouldBeAppended: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			result, appended := test.memberships.AppendIfMissing(test.membership)
			assert.Equal(t, test.shouldBeAppended, appended)
			assert.Contains(t, result, test.membership)
		})
	}
}
