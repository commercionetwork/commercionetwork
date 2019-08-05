package keeper

import (
	"commercio-network/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

var querier = NewQuerier(testUtils.idKeeper)
var request abci.RequestQuery

func Test_queryResolveIdentity(t *testing.T) {
	path := []string{"identities", "newReader"}

	store := testUtils.ctx.KVStore(testUtils.idKeeper.identitiesStoreKey)
	store.Set([]byte(testOwnerIdentity), []byte(testIdentityRef))

	actual, _ := querier(testUtils.ctx, path, request)

	expected := IdentityResult{Did: testOwnerIdentity, DdoReference: testIdentityRef}

	bz, _ := codec.MarshalJSONIndent(testUtils.cdc, expected)

	assert.Equal(t, bz, actual)

}

func Test_queryResolveIdentity_unmarshalError(t *testing.T) {
	path := []string{"identities", "nunu"}

	_, err := querier(testUtils.ctx, path, request)

	assert.Error(t, err)
}

func Test_queryGetConnections(t *testing.T) {
	path := []string{"connections", "newReader"}

	var userConnections = []types.Did{testOwnerIdentity}

	store := testUtils.ctx.KVStore(testUtils.idKeeper.connectionsStoreKey)
	store.Set([]byte(testOwnerIdentity), testUtils.cdc.MustMarshalBinaryBare(&userConnections))

	actual, _ := querier(testUtils.ctx, path, request)

	expected := ConnectionsResult{Did: testOwnerIdentity, Connections: userConnections}

	bz, _ := codec.MarshalJSONIndent(testUtils.cdc, expected)

	assert.Equal(t, bz, actual)
}
