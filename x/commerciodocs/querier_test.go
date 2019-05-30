package commerciodocs

import (
	"commercio-network/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

var querier = NewQuerier(input.docsKeeper)
var request abci.RequestQuery

func Test_queryGetMetadata(t *testing.T) {
	path := []string{"metadata", "reference"}

	expected := MetadataResult{
		Document: reference,
		Metadata: metadata,
	}

	metadataStore := input.ctx.KVStore(input.docsKeeper.metadataStoreKey)
	metadataStore.Set([]byte(reference), []byte(metadata))

	res, _ := querier(input.ctx, path, request)

	bz, _ := codec.MarshalJSONIndent(input.docsKeeper.cdc, expected)

	assert.Equal(t, bz, res)
}

func Test_queryGetAuthorized(t *testing.T) {
	path := []string{"readers", "reference"}

	var readers = []types.Did{"reader1", "reader2"}

	expected := AuthorizedResult{
		Document: reference,
		Readers:  readers,
	}

	readerStore := input.ctx.KVStore(input.docsKeeper.readersStoreKey)
	readerStore.Set([]byte(reference), input.cdc.MustMarshalBinaryBare(&readers))

	res, _ := querier(input.ctx, path, request)

	bz, _ := codec.MarshalJSONIndent(input.docsKeeper.cdc, expected)

	assert.Equal(t, bz, res)

}
