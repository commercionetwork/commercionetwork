package keeper

/*
import (
	"commercio-network/types"
	"commercio-network/x/commercioid"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

var querier = types2.NewQuerier(commercioid.input.idKeeper)
var request abci.RequestQuery

func Test_queryResolveIdentity(t *testing.T) {
	path := []string{"identities", "newReader"}

	store := commercioid.input.ctx.KVStore(commercioid.input.idKeeper.identitiesStoreKey)
	store.Set([]byte(commercioid.ownerIdentity), []byte(commercioid.identityRef))

	actual, _ := querier(commercioid.input.ctx, path, request)

	expected := types2.IdentityResult{Did: commercioid.ownerIdentity, DdoReference: commercioid.identityRef}

	bz, _ := codec.MarshalJSONIndent(commercioid.input.cdc, expected)

	assert.Equal(t, bz, actual)

}

func Test_queryResolveIdentity_unmarshalError(t *testing.T) {
	path := []string{"identities", "nunu"}

	_, err := querier(commercioid.input.ctx, path, request)

	assert.Error(t, err)
}

func Test_queryGetConnections(t *testing.T) {
	path := []string{"connections", "newReader"}

	var userConnections = []types.Did{commercioid.ownerIdentity}

	store := commercioid.input.ctx.KVStore(commercioid.input.idKeeper.connectionsStoreKey)
	store.Set([]byte(commercioid.ownerIdentity), commercioid.input.cdc.MustMarshalBinaryBare(&userConnections))

	actual, _ := querier(commercioid.input.ctx, path, request)

	expected := types2.ConnectionsResult{Did: commercioid.ownerIdentity, Connections: userConnections}

	bz, _ := codec.MarshalJSONIndent(commercioid.input.cdc, expected)

	assert.Equal(t, bz, actual)
}

*/
