package types_test

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestInvite_Empty(t *testing.T) {
	address, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	assert.True(t, types.Invite{}.Empty())
	assert.False(t, types.Invite{Sender: address}.Empty())
	assert.False(t, types.Invite{User: address}.Empty())
	assert.False(t, types.Invite{Rewarded: true}.Empty())
}

func TestInvite_Equals(t *testing.T) {
	user, _ := sdk.AccAddressFromBech32("cosmos1nm9lkhu4dufva9n8zt8q30yd5kuucp54kymqcn")
	sender, _ := sdk.AccAddressFromBech32("cosmos1007jzaanx5kmqnn3akgype2jseawfj80dne9t6")
	invite := types.NewInvite(sender, user)

	assert.False(t, invite.Equals(types.NewInvite(user, sender)))
	assert.False(t, invite.Equals(types.Invite{User: user, Sender: sender, Rewarded: true}))
	assert.True(t, invite.Equals(invite))
}
