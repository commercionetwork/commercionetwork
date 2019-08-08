package keeper

import (
	"github.com/commercionetwork/commercionetwork/x/commercioid/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

var querier = NewQuerier(TestUtils.IdKeeper)
var request abci.RequestQuery

func Test_queryResolveIdentity(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.IdKeeper.StoreKey)
	store.Set([]byte(types.IdentitiesStorePrefix+TestOwnerAddress.String()), []byte(TestDidDocumentReference))

	path := []string{"identities", TestOwnerAddress.String()}
	actual, _ := querier(TestUtils.Ctx, path, request)

	expected, _ := codec.MarshalJSONIndent(TestUtils.Cdc, IdentityResult{
		Did:          TestOwnerAddress,
		DdoReference: TestDidDocumentReference,
	})
	assert.Equal(t, expected, actual)
}

func Test_queryResolveIdentity_nonExistentIdentity(t *testing.T) {
	path := []string{"identities", "nunu"}
	_, err := querier(TestUtils.Ctx, path, request)
	assert.Error(t, err)
}
