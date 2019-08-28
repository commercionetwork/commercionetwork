package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/types"
	keys "github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var querier = NewQuerier(TestUtils.DocsKeeper)
var request abci.RequestQuery
var documents = []types.Document{TestingDocument}

func Test_queryGetReceivedDocuments(t *testing.T) {
	// Setup the store
	metadataStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	metadataStore.Set(
		[]byte(keys.ReceivedDocumentsPrefix+TestingDocument.Recipient.String()),
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

func Test_queryGetSentDocuments(t *testing.T) {
	// Setup the store
	metadataStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	metadataStore.Set(
		[]byte(keys.SentDocumentsPrefix+TestingDocument.Sender.String()),
		TestUtils.Cdc.MustMarshalBinaryBare(&documents),
	)

	// Compose the path
	path := []string{"sent", TestingDocument.Sender.String()}

	// Get the returned documents
	var actual []types.Document
	actualBz, _ := querier(TestUtils.Ctx, path, request)
	TestUtils.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, documents, actual)
}

// ----------------------------------
// --- DocumentReceipt
// ----------------------------------

func TestKeeper_GetUserReceivedReceipts(t *testing.T) {
	receipts := []types.DocumentReceipt{TestingDocumentReceipt}

	//Setup the store
	receiptStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	receiptStore.Set([]byte(keys.DocumentReceiptPrefix+TestingDocumentReceipt.Recipient.String()),
		TestUtils.Cdc.MustMarshalBinaryBare(&receipts))

	// Compose the path
	path := []string{"receipts", TestingDocumentReceipt.Recipient.String()}

	//Get the returned receipts
	var actual []types.DocumentReceipt
	actualBz, _ := querier(TestUtils.Ctx, path, request)
	TestUtils.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, receipts, actual)
}

func TestKeeper_GetReceiptByDocumentUuid(t *testing.T) {
	receipts := []types.DocumentReceipt{TestingDocumentReceipt}

	//Setup the store
	receiptStore := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	receiptStore.Set([]byte(keys.DocumentReceiptPrefix+TestingDocumentReceipt.Recipient.String()),
		TestUtils.Cdc.MustMarshalBinaryBare(&receipts))

	// Compose the path
	path := []string{"receipt", TestingDocumentReceipt.Recipient.String(), TestingDocumentReceipt.Uuid}

	//Get the returned receipts
	var actual types.DocumentReceipt
	actualBz, _ := querier(TestUtils.Ctx, path, request)
	TestUtils.Cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, TestingDocumentReceipt, actual)
}
