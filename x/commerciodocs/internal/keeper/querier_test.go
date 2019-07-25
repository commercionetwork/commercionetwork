package keeper

import (
	"commercio-network/types"
	"commercio-network/x/commerciodocs"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

var querier = NewQuerier(commerciodocs.input.docsKeeper)
var request abci.RequestQuery

func Test_queryGetMetadata(t *testing.T) {
	path := []string{"metadata", "reference"}

	expected := MetadataResult{
		Document: commerciodocs.reference,
		Metadata: commerciodocs.metadata,
	}

	metadataStore := commerciodocs.input.ctx.KVStore(commerciodocs.input.docsKeeper.metadataStoreKey)
	metadataStore.Set([]byte(commerciodocs.reference), []byte(commerciodocs.metadata))

	res, _ := querier(commerciodocs.input.ctx, path, request)

	bz, _ := codec.MarshalJSONIndent(commerciodocs.input.docsKeeper.cdc, expected)

	assert.Equal(t, bz, res)
}

func Test_queryGetAuthorized(t *testing.T) {
	path := []string{"readers", "reference"}

	var readers = []types.Did{"reader1", "reader2"}

	expected := AuthorizedResult{
		Document: commerciodocs.reference,
		Readers:  readers,
	}

	readerStore := commerciodocs.input.ctx.KVStore(commerciodocs.input.docsKeeper.readersStoreKey)
	readerStore.Set([]byte(commerciodocs.reference), commerciodocs.input.cdc.MustMarshalBinaryBare(&readers))

	res, _ := querier(commerciodocs.input.ctx, path, request)

	bz, _ := codec.MarshalJSONIndent(commerciodocs.input.docsKeeper.cdc, expected)

	assert.Equal(t, bz, res)

}
