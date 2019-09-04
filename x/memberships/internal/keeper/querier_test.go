package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var querier = NewQuerier(TestUtils.MembershipKeeper)
var request abci.RequestQuery

func TestQuerier_resolveIdentity_Existent(t *testing.T) {
	user := TestSignerAddress
	membershipType := "green"
	_, _ = TestUtils.MembershipKeeper.AssignMembership(TestUtils.Ctx, user, membershipType)

	path := []string{types.QueryGetMembership, user.String()}
	actual, _ := querier(TestUtils.Ctx, path, request)

	expected, _ := codec.MarshalJSONIndent(TestUtils.Cdc, MembershipResult{
		User:           user,
		MembershipType: membershipType,
	})
	assert.Equal(t, expected, actual)
}

func TestQuerier_ResolveIdentity_NonExistent(t *testing.T) {
	path := []string{types.QueryGetMembership, "nunu"}
	_, err := querier(TestUtils.Ctx, path, request)
	assert.Error(t, err)
}
