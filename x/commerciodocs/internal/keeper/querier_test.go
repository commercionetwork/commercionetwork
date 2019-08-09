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

func Test_queryGetReceivedDocuments(t *testing.T) {
	documents := []types.Document{TestingDocument}

	metadataStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	metadataStore.Set(
		[]byte(ReceivedDocumentsPrefix+TestingDocument.Recipient.String()),
		codec.MustMarshalJSONIndent(TestUtils.Cdc, documents),
	)

	path := []string{"received", TestingDocument.Recipient.String()}
	actual, _ := querier(TestUtils.Ctx, path, request)

	expected, _ := codec.MarshalJSONIndent(TestUtils.Cdc, documents)
	assert.Equal(t, expected, actual)
}

//func Test_queryGetAuthorized(t *testing.T) {
//	path := []string{"readers", "TestReference"}
//
//	var readers = []types.Did{"reader1", "reader2"}
//
//	expected := AuthorizedResult{
//		Document: TestReference,
//		Readers:  readers,
//	}
//
//	readerStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.readersStoreKey)
//	readerStore.Set([]byte(TestReference), TestUtils.Cdc.MustMarshalBinaryBare(&readers))
//
//	res, _ := querier(TestUtils.Ctx, path, request)
//
//	bz, _ := codec.MarshalJSONIndent(TestUtils.DocsKeeper.cdc, expected)
//
//	assert.Equal(t, bz, res)
//
//}
