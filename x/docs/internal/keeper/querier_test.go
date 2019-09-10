package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var querier = NewQuerier(TestUtils.DocsKeeper)
var request abci.RequestQuery
var documents = types.Documents{TestingDocument}

// ----------------------------------
// --- Documents
// ----------------------------------

func Test_queryGetReceivedDocuments(t *testing.T) {
	// Setup the store
	metadataStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	metadataStore.Set(
		[]byte(types.ReceivedDocumentsPrefix+TestingDocument.Recipient.String()),
		TestUtils.Cdc.MustMarshalBinaryBare(&documents),
	)

	// Compose the path
	path := []string{types.QueryReceivedDocuments, TestingDocument.Recipient.String()}

	// Get the returned documents
	var actual types.Documents
	actualBz, _ := querier(TestUtils.Ctx, path, request)
	TestUtils.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, documents, actual)
}

func Test_queryGetSentDocuments(t *testing.T) {
	// Setup the store
	metadataStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	metadataStore.Set(
		[]byte(types.SentDocumentsPrefix+TestingDocument.Sender.String()),
		TestUtils.Cdc.MustMarshalBinaryBare(&documents),
	)

	// Compose the path
	path := []string{types.QuerySentDocuments, TestingDocument.Sender.String()}

	// Get the returned documents
	var actual types.Documents
	actualBz, _ := querier(TestUtils.Ctx, path, request)
	TestUtils.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, documents, actual)
}

// ---------------------------------
// --- Document receipts
// ---------------------------------

func Test_queryGetReceivedDocsReceipts(t *testing.T) {
	// Setup the store
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete(TestUtils.DocsKeeper.getReceivedReceiptsStoreKey(TestingDocumentReceipt.Recipient))

	var stored = types.DocumentReceipts{TestingDocumentReceipt}
	store.Set(
		TestUtils.DocsKeeper.getReceivedReceiptsStoreKey(TestingDocumentReceipt.Recipient),
		TestUtils.Cdc.MustMarshalBinaryBare(&stored),
	)

	// Compose the path
	path := []string{types.QueryReceivedReceipts, TestingDocumentReceipt.Recipient.String(), ""}

	// Get the returned receipts
	var actual types.DocumentReceipts
	actualBz, _ := querier(TestUtils.Ctx, path, request)
	TestUtils.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, stored, actual)
}

func Test_queryGetReceivedDocsReceipts_WithDocUuid(t *testing.T) {
	// Setup the store
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete(TestUtils.DocsKeeper.getReceivedReceiptsStoreKey(TestingDocumentReceipt.Recipient))

	var stored = types.DocumentReceipts{TestingDocumentReceipt}
	store.Set(
		TestUtils.DocsKeeper.getReceivedReceiptsStoreKey(TestingDocumentReceipt.Recipient),
		TestUtils.Cdc.MustMarshalBinaryBare(&stored),
	)

	// Compose the path
	path := []string{types.QueryReceivedReceipts, TestingDocumentReceipt.Recipient.String(), TestingDocumentReceipt.DocumentUuid}

	// Get the returned receipts
	var actual types.DocumentReceipts
	actualBz, _ := querier(TestUtils.Ctx, path, request)
	TestUtils.Cdc.MustUnmarshalJSON(actualBz, &actual)

	var expected = types.DocumentReceipts{TestingDocumentReceipt}
	assert.Equal(t, expected, actual)
}
