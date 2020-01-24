package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery
var documents = []types.Document{TestingDocument}

// ----------------------------------
// --- Documents
// ----------------------------------

func Test_queryGetReceivedDocuments_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	path := []string{types.QueryReceivedDocuments, TestingRecipient.String()}

	var actual []types.Document
	actualBz, err := querier(ctx, path, request)
	if err != nil {
		require.NoError(t, err)
	}
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, "[]", string(actualBz))
	require.Empty(t, actual)
}

func Test_queryGetReceivedDocuments_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	// Setup the store
	err := k.SaveDocument(ctx, TestingDocument)
	require.NoError(t, err)

	// Compose the path
	path := []string{types.QueryReceivedDocuments, TestingRecipient.String()}

	// Get the returned documents
	var actual []types.Document

	actualBz, err := querier(ctx, path, request)
	require.NoError(t, err)

	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, documents, actual)
}

func Test_queryGetSentDocuments_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	path := []string{types.QuerySentDocuments, TestingSender.String()}

	var actual []types.Document
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, "[]", string(actualBz))
	require.Empty(t, actual)
}

func Test_queryGetSentDocuments_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)
	//Setup the store
	err := k.SaveDocument(ctx, TestingDocument)
	require.NoError(t, err)

	// Compose the path
	path := []string{types.QuerySentDocuments, TestingSender.String()}

	// Get the returned documents
	var actual []types.Document
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, documents, actual)
}

// ---------------------------------
// --- Document receipts
// ---------------------------------

func Test_queryGetReceivedDocsReceipts_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	path := []string{types.QueryReceivedReceipts, TestingDocumentReceipt.Recipient.String(), ""}

	// Get the returned receipts
	var actual []types.DocumentReceipt
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, "[]", string(actualBz))
	require.Empty(t, actual)
}

func Test_queryGetReceivedDocsReceipts_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	require.NoError(t, k.SaveDocument(ctx, TestingDocument))

	tdr := TestingDocumentReceipt
	tdr.DocumentUUID = TestingDocument.UUID

	err := k.SaveReceipt(ctx, tdr)
	require.NoError(t, err)

	// Compose the path
	path := []string{types.QueryReceivedReceipts, TestingDocumentReceipt.Recipient.String(), ""}

	// Get the returned receipts
	querier := NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual []types.DocumentReceipt
	cdc.MustUnmarshalJSON(actualBz, &actual)

	expected := []types.DocumentReceipt{tdr}
	require.Equal(t, expected, actual)
}

func Test_queryGetReceivedDocsReceipts_WithDocUuid(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	require.NoError(t, k.SaveDocument(ctx, TestingDocument))

	tdr := TestingDocumentReceipt
	tdr.DocumentUUID = TestingDocument.UUID

	err := k.SaveReceipt(ctx, tdr)
	require.NoError(t, err)

	// Compose the path
	path := []string{types.QueryReceivedReceipts, TestingDocumentReceipt.Recipient.String(), tdr.DocumentUUID}

	// Get the returned receipts
	querier := NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual []types.DocumentReceipt
	cdc.MustUnmarshalJSON(actualBz, &actual)

	var expected = []types.DocumentReceipt{tdr}
	require.Equal(t, expected, actual)
}

func Test_queryGetSentDocsReceipts_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	path := []string{types.QuerySentReceipts, TestingDocumentReceipt.Sender.String()}

	var actual []types.DocumentReceipt
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, "[]", string(actualBz))
	require.Empty(t, actual)
}

func Test_queryGetSentDocsReceipts_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	require.NoError(t, k.SaveDocument(ctx, TestingDocument))

	tdr := TestingDocumentReceipt
	tdr.DocumentUUID = TestingDocument.UUID

	err := k.SaveReceipt(ctx, tdr)
	require.NoError(t, err)

	path := []string{types.QuerySentReceipts, TestingDocumentReceipt.Sender.String()}

	querier := NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual []types.DocumentReceipt
	cdc.MustUnmarshalJSON(actualBz, &actual)

	expected := []types.DocumentReceipt{tdr}
	require.Equal(t, expected, actual)
}

// ----------------------------------
// --- Document metadata schemes
// ----------------------------------

func Test_querySupportedMetadataSchemes_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	path := []string{types.QuerySupportedMetadataSchemes}

	var actual []types.MetadataSchema
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, "[]", string(actualBz))
	require.Empty(t, actual)
}

func Test_querySupportedMetadataSchemes_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	schemes := []types.MetadataSchema{
		{Type: "schema", SchemaURI: "https://example.com/schema", Version: "1.0.0"},
		{Type: "other-schema", SchemaURI: "https://example.com/other-schema", Version: "1.0.0"},
	}

	k.AddSupportedMetadataScheme(ctx, schemes[0])

	k.AddSupportedMetadataScheme(ctx, schemes[1])

	path := []string{types.QuerySupportedMetadataSchemes}

	var actual []types.MetadataSchema
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Contains(t, schemes, actual[0])
	require.Contains(t, schemes, actual[1])

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

	require.Equal(t, "[]", string(actualBz))
	require.Empty(t, actual)
}

func Test_queryTrustedMetadataProposers_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var querier = NewQuerier(k)

	k.AddTrustedSchemaProposer(ctx, TestingSender)
	k.AddTrustedSchemaProposer(ctx, TestingSender2)

	path := []string{types.QueryTrustedMetadataProposers}

	var actual []sdk.AccAddress
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Contains(t, actual, TestingSender)
	require.Contains(t, actual, TestingSender2)
}
