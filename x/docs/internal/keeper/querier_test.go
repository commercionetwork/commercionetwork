package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var querier = NewQuerier(TestUtils.DocsKeeper)
var request abci.RequestQuery

func Test_queryGetReceivedDocuments(t *testing.T) {
	documents := []types.Document{TestingDocument}

	// Setup the store
	metadataStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	metadataStore.Set(
		[]byte(ReceivedDocumentsPrefix+TestingDocument.Recipient.String()),
		TestUtils.Cdc.MustMarshalBinaryBare(&documents),
	)

	// Compose the path
	path := []string{"received", TestingDocument.Recipient.String()}

	// Get the returned documents
	var actual []types.Document
	actualBz, _ := querier(TestUtils.Ctx, path, request)
	TestUtils.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, documents, actual)
}
