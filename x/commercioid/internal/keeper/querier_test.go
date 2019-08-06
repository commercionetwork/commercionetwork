package keeper

import (
	"github.com/commercionetwork/commercionetwork/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

var querier = NewQuerier(TestUtils.IdKeeper)
var request abci.RequestQuery

func Test_queryResolveIdentity(t *testing.T) {
	path := []string{"identities", "newReader"}

	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.identitiesStoreKey)
	store.Set([]byte(TestOwnerIdentity), []byte(TestIdentityRef))

	actual, _ := querier(TestUtils.Ctx, path, request)

	expected := IdentityResult{Did: TestOwnerIdentity, DdoReference: TestIdentityRef}

	bz, _ := codec.MarshalJSONIndent(TestUtils.Cdc, expected)

	assert.Equal(t, bz, actual)

}

func Test_queryResolveIdentity_unmarshalError(t *testing.T) {
	path := []string{"identities", "nunu"}

	_, err := querier(TestUtils.Ctx, path, request)

	assert.Error(t, err)
}

func Test_queryGetConnections(t *testing.T) {
	path := []string{"connections", "newReader"}

	var userConnections = []types.Did{TestOwnerIdentity}

	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.connectionsStoreKey)
	store.Set([]byte(TestOwnerIdentity), TestUtils.Cdc.MustMarshalBinaryBare(&userConnections))

	actual, _ := querier(TestUtils.Ctx, path, request)

	expected := ConnectionsResult{Did: TestOwnerIdentity, Connections: userConnections}

	bz, _ := codec.MarshalJSONIndent(TestUtils.Cdc, expected)

	assert.Equal(t, bz, actual)
}
