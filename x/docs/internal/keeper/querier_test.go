package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery
var documents = types.Documents{TestingDocument}

// ----------------------------------
// --- Documents
// ----------------------------------

func Test_queryGetReceivedDocuments_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	path := []string{types.QueryReceivedDocuments, TestingRecipient.String()}

	var actual types.Documents
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, "[]", string(actualBz))
	assert.Empty(t, actual)
}

func Test_queryGetReceivedDocuments_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	// Setup the store
	metadataStore := ctx.KVStore(k.StoreKey)
	documentIds := types.DocumentIDs{TestingDocument.UUID}
	metadataStore.Set(k.getReceivedDocumentsIdsStoreKey(TestingRecipient), cdc.MustMarshalBinaryBare(&documentIds))
	metadataStore.Set(k.getDocumentStoreKey(TestingDocument.UUID), cdc.MustMarshalBinaryBare(TestingDocument))

	// Compose the path
	path := []string{types.QueryReceivedDocuments, TestingRecipient.String()}

	// Get the returned documents
	var actual types.Documents
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, documents, actual)
}

func Test_queryGetSentDocuments_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	path := []string{types.QuerySentDocuments, TestingSender.String()}

	var actual types.Documents
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, "[]", string(actualBz))
	assert.Empty(t, actual)
}

func Test_queryGetSentDocuments_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)
	//Setup the store
	metadataStore := ctx.KVStore(k.StoreKey)
	documentIds := types.DocumentIDs{TestingDocument.UUID}
	metadataStore.Set(k.getSentDocumentsIdsStoreKey(TestingSender), cdc.MustMarshalBinaryBare(&documentIds))
	metadataStore.Set(k.getDocumentStoreKey(TestingDocument.UUID), cdc.MustMarshalBinaryBare(TestingDocument))

	// Compose the path
	path := []string{types.QuerySentDocuments, TestingSender.String()}

	// Get the returned documents
	var actual types.Documents
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, documents, actual)
}

// ---------------------------------
// --- Document receipts
// ---------------------------------

func Test_queryGetReceivedDocsReceipts_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	path := []string{types.QueryReceivedReceipts, TestingDocumentReceipt.Recipient.String(), ""}

	// Get the returned receipts
	var actual types.DocumentReceipts
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, "[]", string(actualBz))
	assert.Empty(t, actual)
}

func Test_queryGetReceivedDocsReceipts_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	//Setup the store
	store := ctx.KVStore(k.StoreKey)

	ids := types.DocumentReceiptsIDs{TestingDocumentReceipt.UUID}
	store.Set(k.getReceivedReceiptsIdsStoreKey(TestingDocumentReceipt.Recipient), cdc.MustMarshalBinaryBare(&ids))
	store.Set(k.getReceiptStoreKey(TestingDocumentReceipt.UUID), cdc.MustMarshalBinaryBare(&TestingDocumentReceipt))

	// Compose the path
	path := []string{types.QueryReceivedReceipts, TestingDocumentReceipt.Recipient.String(), ""}

	// Get the returned receipts
	querier := NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual types.DocumentReceipts
	cdc.MustUnmarshalJSON(actualBz, &actual)

	expected := types.DocumentReceipts{TestingDocumentReceipt}
	assert.Equal(t, expected, actual)
}

func Test_queryGetReceivedDocsReceipts_WithDocUuid(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	//Setup the store
	store := ctx.KVStore(k.StoreKey)

	var ids = types.DocumentReceiptsIDs{TestingDocumentReceipt.UUID}
	store.Set(k.getReceivedReceiptsIdsStoreKey(TestingDocumentReceipt.Recipient), cdc.MustMarshalBinaryBare(&ids))
	store.Set(k.getReceiptStoreKey(TestingDocumentReceipt.UUID), cdc.MustMarshalBinaryBare(&TestingDocumentReceipt))

	// Compose the path
	path := []string{types.QueryReceivedReceipts, TestingDocumentReceipt.Recipient.String(), TestingDocumentReceipt.DocumentUUID}

	// Get the returned receipts
	querier := NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual types.DocumentReceipts
	cdc.MustUnmarshalJSON(actualBz, &actual)

	var expected = types.DocumentReceipts{TestingDocumentReceipt}
	assert.Equal(t, expected, actual)
}

func Test_queryGetSentDocsReceipts_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	path := []string{types.QuerySentReceipts, TestingDocumentReceipt.Sender.String()}

	var actual types.DocumentReceipts
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, "[]", string(actualBz))
	assert.Empty(t, actual)
}

func Test_queryGetSentDocsReceipts_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	//Setup the store
	store := ctx.KVStore(k.StoreKey)

	var ids = types.DocumentReceiptsIDs{TestingDocumentReceipt.UUID}
	store.Set(k.getSentReceiptsIdsStoreKey(TestingDocumentReceipt.Sender), cdc.MustMarshalBinaryBare(&ids))
	store.Set(k.getReceiptStoreKey(TestingDocumentReceipt.UUID), cdc.MustMarshalBinaryBare(&TestingDocumentReceipt))

	path := []string{types.QuerySentReceipts, TestingDocumentReceipt.Sender.String()}

	querier := NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual types.DocumentReceipts
	cdc.MustUnmarshalJSON(actualBz, &actual)

	expected := types.DocumentReceipts{TestingDocumentReceipt}
	assert.Equal(t, expected, actual)
}

// ----------------------------------
// --- Document metadata schemes
// ----------------------------------

func Test_querySupportedMetadataSchemes_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	path := []string{types.QuerySupportedMetadataSchemes}

	var actual types.MetadataSchemes
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, "[]", string(actualBz))
	assert.Empty(t, actual)
}

func Test_querySupportedMetadataSchemes_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)
	//Setup the store
	store := ctx.KVStore(k.StoreKey)

	schemes := []types.MetadataSchema{
		{Type: "schema", SchemaURI: "https://example.com/schema", Version: "1.0.0"},
		{Type: "other-schema", SchemaURI: "https://example.com/other-schema", Version: "1.0.0"},
	}
	store.Set([]byte(types.SupportedMetadataSchemesStoreKey), cdc.MustMarshalBinaryBare(&schemes))

	path := []string{types.QuerySupportedMetadataSchemes}

	var actual []types.MetadataSchema
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, schemes, actual)
}

// -----------------------------------------
// --- Document metadata schemes proposers
// -----------------------------------------

func Test_queryTrustedMetadataProposers_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	path := []string{types.QueryTrustedMetadataProposers}

	var actual []sdk.AccAddress
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, "[]", string(actualBz))
	assert.Empty(t, actual)
}

func Test_queryTrustedMetadataProposers_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)
	//Setup the store
	proposers := []sdk.AccAddress{TestingSender, TestingSender2}

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.MetadataSchemaProposersStoreKey), cdc.MustMarshalBinaryBare(&proposers))

	path := []string{types.QueryTrustedMetadataProposers}

	var actual []sdk.AccAddress
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, proposers, actual)
}
