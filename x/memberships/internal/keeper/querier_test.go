package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery

func TestQuerier_resolveIdentity_Existent(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	user := TestSignerAddress
	membershipType := "green"
	_, _ = k.AssignMembership(ctx, user, membershipType)

	path := []string{types.QueryGetMembership, user.String()}
	actual, _ := querier(ctx, path, request)

	expected, _ := codec.MarshalJSONIndent(cdc, MembershipResult{
		User:           user,
		MembershipType: membershipType,
	})
	assert.Equal(t, expected, actual)
}

func TestQuerier_ResolveIdentity_NonExistent(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	path := []string{types.QueryGetMembership, "nunu"}
	_, err := querier(ctx, path, request)
	assert.Error(t, err)
}
