package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery

func Test_queryResolveIdentity(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.IdentitiesStorePrefix+TestOwnerAddress.String()), []byte(TestDidDocumentUri))

	path := []string{"identities", TestOwnerAddress.String()}
	actual, _ := querier(ctx, path, request)

	expected, _ := codec.MarshalJSONIndent(cdc, IdentityResult{
		Did:          TestOwnerAddress,
		DdoReference: TestDidDocumentUri,
	})
	assert.Equal(t, expected, actual)
}

func Test_queryResolveIdentity_nonExistentIdentity(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)
	path := []string{"identities", "nunu"}
	_, err := querier(ctx, path, request)
	assert.Error(t, err)
}
