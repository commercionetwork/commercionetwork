package keeper

import (
	"commercio-network/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

var querier = NewQuerier(testUtils.docsKeeper)
var request abci.RequestQuery

func Test_queryGetMetadata(t *testing.T) {
	path := []string{"testMetadata", "testReference"}

	expected := MetadataResult{
		Document: testReference,
		Metadata: testMetadata,
	}

	metadataStore := testUtils.ctx.KVStore(testUtils.docsKeeper.metadataStoreKey)
	metadataStore.Set([]byte(testReference), []byte(testMetadata))

	res, _ := querier(testUtils.ctx, path, request)

	bz, _ := codec.MarshalJSONIndent(testUtils.docsKeeper.cdc, expected)

	assert.Equal(t, bz, res)
}

func Test_queryGetAuthorized(t *testing.T) {
	path := []string{"readers", "testReference"}

	var readers = []types.Did{"reader1", "reader2"}

	expected := AuthorizedResult{
		Document: testReference,
		Readers:  readers,
	}

	readerStore := testUtils.ctx.KVStore(testUtils.docsKeeper.readersStoreKey)
	readerStore.Set([]byte(testReference), testUtils.cdc.MustMarshalBinaryBare(&readers))

	res, _ := querier(testUtils.ctx, path, request)

	bz, _ := codec.MarshalJSONIndent(testUtils.docsKeeper.cdc, expected)

	assert.Equal(t, bz, res)

}
