package keeper

import (
	"github.com/commercionetwork/commercionetwork/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

var querier = NewQuerier(TestUtils.DocsKeeper)
var request abci.RequestQuery

func Test_queryGetMetadata(t *testing.T) {
	path := []string{"TestMetadata", "TestReference"}

	expected := MetadataResult{
		Document: TestReference,
		Metadata: TestMetadata,
	}

	metadataStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.metadataStoreKey)
	metadataStore.Set([]byte(TestReference), []byte(TestMetadata))

	res, _ := querier(TestUtils.Ctx, path, request)

	bz, _ := codec.MarshalJSONIndent(TestUtils.DocsKeeper.cdc, expected)

	assert.Equal(t, bz, res)
}

func Test_queryGetAuthorized(t *testing.T) {
	path := []string{"readers", "TestReference"}

	var readers = []types.Did{"reader1", "reader2"}

	expected := AuthorizedResult{
		Document: TestReference,
		Readers:  readers,
	}

	readerStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.readersStoreKey)
	readerStore.Set([]byte(TestReference), TestUtils.Cdc.MustMarshalBinaryBare(&readers))

	res, _ := querier(TestUtils.Ctx, path, request)

	bz, _ := codec.MarshalJSONIndent(TestUtils.DocsKeeper.cdc, expected)

	assert.Equal(t, bz, res)

}
