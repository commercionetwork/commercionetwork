package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery

func TestQuerier_resolveIdentity_Existent(t *testing.T) {
	cdc, ctx, _, _, k := SetupTestInput()
	var querier = NewQuerier(k)

	user := TestUserAddress
	membershipType := "bronze"
	_, _ = k.AssignMembership(ctx, user, membershipType)

	path := []string{types.QueryGetMembership, user.String()}
	bz, _ := querier(ctx, path, request)

	var actual MembershipResult
	cdc.MustUnmarshalJSON(bz, &actual)

	expected := MembershipResult{User: user, MembershipType: membershipType}
	assert.Equal(t, expected, actual)
}

func TestQuerier_ResolveIdentity_NonExistent(t *testing.T) {
	_, ctx, _, _, k := SetupTestInput()
	var querier = NewQuerier(k)

	path := []string{types.QueryGetMembership, "nunu"}
	_, err := querier(ctx, path, request)
	assert.Error(t, err)
}
