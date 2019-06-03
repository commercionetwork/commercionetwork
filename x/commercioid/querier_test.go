package commercioid

import (
	"commercio-network/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

var querier = NewQuerier(input.idKeeper)
var request abci.RequestQuery

func Test_queryResolveIdentity(t *testing.T) {
	path := []string{"identities", "newReader"}

	store := input.ctx.KVStore(input.idKeeper.identitiesStoreKey)
	store.Set([]byte(ownerIdentity), []byte(identityRef))

	actual, _ := querier(input.ctx, path, request)

	expected := IdentityResult{Did: ownerIdentity, DdoReference: identityRef}

	bz, _ := codec.MarshalJSONIndent(input.cdc, expected)

	assert.Equal(t, bz, actual)

}

func Test_queryResolveIdentity_unmarshalError(t *testing.T) {
	path := []string{"identities", "newReader"}

	_, err := querier(input.ctx, path, request)

	assert.Error(t, err)
}

func Test_queryGetConnections(t *testing.T) {
	path := []string{"connections", "newReader"}

	var userConnections = []types.Did{ownerIdentity}

	store := input.ctx.KVStore(input.idKeeper.connectionsStoreKey)
	store.Set([]byte(ownerIdentity), input.cdc.MustMarshalBinaryBare(&userConnections))

	actual, _ := querier(input.ctx, path, request)

	expected := ConnectionsResult{Did: ownerIdentity, Connections: userConnections}

	bz, _ := codec.MarshalJSONIndent(input.cdc, expected)

	assert.Equal(t, bz, actual)
}
