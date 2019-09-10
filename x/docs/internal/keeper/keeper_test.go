package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------
// --- Metadata schemes
// ----------------------------------

func TestKeeper_AddSupportedMetadataScheme_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete([]byte(types.SupportedMetadataSchemesStoreKey))

	schema := types.MetadataSchema{Type: "schema", SchemaUri: "https://example.com/schema", Version: "1.0.0"}
	TestUtils.DocsKeeper.AddSupportedMetadataScheme(TestUtils.Ctx, schema)

	var stored types.MetadataSchemes
	storedBz := store.Get([]byte(types.SupportedMetadataSchemesStoreKey))
	TestUtils.Cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, 1, len(stored))
	assert.Contains(t, stored, schema)
}

func TestKeeper_AddSupportedMetadataScheme_ExistingList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)

	existingSchema := types.MetadataSchema{Type: "schema", SchemaUri: "https://example.com/newSchema", Version: "1.0.0"}
	existing := []types.MetadataSchema{existingSchema}
	existingBz := TestUtils.Cdc.MustMarshalBinaryBare(&existing)
	store.Set([]byte(types.SupportedMetadataSchemesStoreKey), existingBz)

	newSchema := types.MetadataSchema{Type: "schema2", SchemaUri: "https://example.com/schema2", Version: "2.0.0"}
	TestUtils.DocsKeeper.AddSupportedMetadataScheme(TestUtils.Ctx, newSchema)

	var stored types.MetadataSchemes
	storedBz := store.Get([]byte(types.SupportedMetadataSchemesStoreKey))
	TestUtils.Cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, 2, len(stored))
	assert.Contains(t, stored, existingSchema)
	assert.Contains(t, stored, newSchema)
}

func TestKeeper_IsMetadataSchemeTypeSupported_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete([]byte(types.SupportedMetadataSchemesStoreKey))

	assert.False(t, TestUtils.DocsKeeper.IsMetadataSchemeTypeSupported(TestUtils.Ctx, "schema"))
	assert.False(t, TestUtils.DocsKeeper.IsMetadataSchemeTypeSupported(TestUtils.Ctx, "schema2"))
	assert.False(t, TestUtils.DocsKeeper.IsMetadataSchemeTypeSupported(TestUtils.Ctx, "non-existent"))
}

func TestKeeper_IsMetadataSchemeTypeSupported_ExistingList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)

	existingSchema := types.MetadataSchema{Type: "schema", SchemaUri: "https://example.com/newSchema", Version: "1.0.0"}
	existing := []types.MetadataSchema{existingSchema}
	existingBz := TestUtils.Cdc.MustMarshalBinaryBare(&existing)
	store.Set([]byte(types.SupportedMetadataSchemesStoreKey), existingBz)

	assert.True(t, TestUtils.DocsKeeper.IsMetadataSchemeTypeSupported(TestUtils.Ctx, "schema"))
	assert.False(t, TestUtils.DocsKeeper.IsMetadataSchemeTypeSupported(TestUtils.Ctx, "schema2"))
	assert.False(t, TestUtils.DocsKeeper.IsMetadataSchemeTypeSupported(TestUtils.Ctx, "any-schema"))
}

func TestKeeper_GetSupportedMetadataSchemes_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete([]byte(types.SupportedMetadataSchemesStoreKey))

	result := TestUtils.DocsKeeper.GetSupportedMetadataSchemes(TestUtils.Ctx)

	assert.Empty(t, result)
}

func TestKeeper_GetSupportedMetadataSchemes_ExistingList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)

	existingSchema := types.MetadataSchema{Type: "schema", SchemaUri: "https://example.com/newSchema", Version: "1.0.0"}
	existing := types.MetadataSchemes{existingSchema}
	existingBz := TestUtils.Cdc.MustMarshalBinaryBare(&existing)
	store.Set([]byte(types.SupportedMetadataSchemesStoreKey), existingBz)

	actual := TestUtils.DocsKeeper.GetSupportedMetadataSchemes(TestUtils.Ctx)

	assert.Equal(t, existing, actual)
}

// ----------------------------------
// --- Metadata schema proposers
// ----------------------------------

func TestKeeper_AddTrustedSchemaProposer_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete([]byte(types.MetadataSchemaProposersStoreKey))

	TestUtils.DocsKeeper.AddTrustedSchemaProposer(TestUtils.Ctx, TestingSender)

	var proposers []sdk.AccAddress
	proposersBz := store.Get([]byte(types.MetadataSchemaProposersStoreKey))
	TestUtils.Cdc.MustUnmarshalBinaryBare(proposersBz, &proposers)

	assert.Equal(t, 1, len(proposers))
	assert.Contains(t, proposers, TestingSender)
}

func TestKeeper_AddTrustedSchemaProposer_ExistingList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)

	existing := []sdk.AccAddress{TestingSender}
	proposersBz := TestUtils.Cdc.MustMarshalBinaryBare(&existing)
	store.Set([]byte(types.MetadataSchemaProposersStoreKey), proposersBz)

	TestUtils.DocsKeeper.AddTrustedSchemaProposer(TestUtils.Ctx, TestingSender2)

	var stored []sdk.AccAddress
	storedBz := store.Get([]byte(types.MetadataSchemaProposersStoreKey))
	TestUtils.Cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, 2, len(stored))
	assert.Contains(t, stored, TestingSender)
	assert.Contains(t, stored, TestingSender2)
}

func TestKeeper_IsTrustedSchemaProposer_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete([]byte(types.MetadataSchemaProposersStoreKey))

	assert.False(t, TestUtils.DocsKeeper.IsTrustedSchemaProposer(TestUtils.Ctx, TestingSender))
	assert.False(t, TestUtils.DocsKeeper.IsTrustedSchemaProposer(TestUtils.Ctx, TestingSender2))
}

func TestKeeper_IsTrustedSchemaProposerExistingList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)

	existing := []sdk.AccAddress{TestingSender}
	proposersBz := TestUtils.Cdc.MustMarshalBinaryBare(&existing)
	store.Set([]byte(types.MetadataSchemaProposersStoreKey), proposersBz)

	assert.True(t, TestUtils.DocsKeeper.IsTrustedSchemaProposer(TestUtils.Ctx, TestingSender))
	assert.False(t, TestUtils.DocsKeeper.IsTrustedSchemaProposer(TestUtils.Ctx, TestingSender2))
}

func TestKeeper_GetTrustedSchemaProposers_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete([]byte(types.MetadataSchemaProposersStoreKey))

	stored := TestUtils.DocsKeeper.GetTrustedSchemaProposers(TestUtils.Ctx)

	assert.Empty(t, stored)
}

func TestKeeper_GetTrustedSchemaProposers_ExistingList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)

	existing := types.Addresses{TestingSender}
	proposersBz := TestUtils.Cdc.MustMarshalBinaryBare(&existing)
	store.Set([]byte(types.MetadataSchemaProposersStoreKey), proposersBz)

	stored := TestUtils.DocsKeeper.GetTrustedSchemaProposers(TestUtils.Ctx)

	assert.Equal(t, existing, stored)
}

// ----------------------------------
// --- Documents
// ----------------------------------

func TestKeeper_ShareDocument_EmptyLists(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)

	TestUtils.DocsKeeper.ShareDocument(TestUtils.Ctx, TestingDocument)

	sentDocsBz := store.Get([]byte(types.SentDocumentsPrefix + TestingSender.String()))
	receivedDocsBz := store.Get([]byte(types.ReceivedDocumentsPrefix + TestingRecipient.String()))

	var sentDocs, receivedDocs types.Documents
	TestUtils.Cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
	TestUtils.Cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)

	assert.Equal(t, 1, len(sentDocs))
	assert.Contains(t, sentDocs, TestingDocument)

	assert.Equal(t, 1, len(receivedDocs))
	assert.Equal(t, sentDocs, receivedDocs)
}

func TestKeeper_ShareDocument_ExistingDocument(t *testing.T) {
	documents := types.Documents{TestingDocument}
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set([]byte(types.SentDocumentsPrefix+TestingSender.String()), TestUtils.Cdc.MustMarshalBinaryBare(&documents))
	store.Set([]byte(types.ReceivedDocumentsPrefix+TestingRecipient.String()), TestUtils.Cdc.MustMarshalBinaryBare(&documents))

	TestUtils.DocsKeeper.ShareDocument(TestUtils.Ctx, TestingDocument)

	sentDocsBz := store.Get([]byte(types.SentDocumentsPrefix + TestingSender.String()))
	receivedDocsBz := store.Get([]byte(types.ReceivedDocumentsPrefix + TestingRecipient.String()))

	var sentDocs, receivedDocs types.Documents
	TestUtils.Cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
	TestUtils.Cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)

	assert.Equal(t, 1, len(sentDocs))
	assert.Contains(t, sentDocs, TestingDocument)

	assert.Equal(t, 1, len(receivedDocs))
	assert.Equal(t, sentDocs, receivedDocs)
}

func TestKeeper_ShareDocument_SameInfoDifferentRecipient(t *testing.T) {
	documents := types.Documents{TestingDocument}

	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set(
		TestUtils.DocsKeeper.getSentDocumentsStoreKey(TestingDocument.Sender),
		TestUtils.Cdc.MustMarshalBinaryBare(&documents),
	)
	store.Set(
		TestUtils.DocsKeeper.getReceivedDocumentsStoreKey(TestingDocument.Recipient),
		TestUtils.Cdc.MustMarshalBinaryBare(&documents),
	)

	newRecipient, _ := sdk.AccAddressFromBech32("cosmos1h2z8u9294gtqmxlrnlyfueqysng3krh009fum7")
	newDocument := types.Document{
		Sender:     TestingDocument.Sender,
		Recipient:  newRecipient,
		ContentUri: TestingDocument.ContentUri,
		Metadata:   TestingDocument.Metadata,
		Checksum:   TestingDocument.Checksum,
	}
	TestUtils.DocsKeeper.ShareDocument(TestUtils.Ctx, newDocument)

	sentDocsBz := store.Get(TestUtils.DocsKeeper.getSentDocumentsStoreKey(TestingDocument.Sender))
	receivedDocsBz := store.Get(TestUtils.DocsKeeper.getReceivedDocumentsStoreKey(TestingDocument.Recipient))
	newReceivedDocsBz := store.Get(TestUtils.DocsKeeper.getReceivedDocumentsStoreKey(newRecipient))

	var sentDocs, receivedDocs, newReceivedDocs types.Documents
	TestUtils.Cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
	TestUtils.Cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)
	TestUtils.Cdc.MustUnmarshalBinaryBare(newReceivedDocsBz, &newReceivedDocs)

	assert.Equal(t, 2, len(sentDocs))
	assert.Contains(t, sentDocs, TestingDocument)
	assert.Contains(t, sentDocs, newDocument)

	assert.Equal(t, 1, len(receivedDocs))
	assert.Contains(t, receivedDocs, TestingDocument)

	assert.Equal(t, 1, len(newReceivedDocs))
	assert.Contains(t, newReceivedDocs, newDocument)
}

func TestKeeper_GetUserReceivedDocuments_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete([]byte(types.ReceivedDocumentsPrefix + TestingRecipient.String()))

	receivedDocs := TestUtils.DocsKeeper.GetUserReceivedDocuments(TestUtils.Ctx, TestingDocument.Sender)
	assert.Equal(t, 0, len(receivedDocs))
}

func TestKeeper_GetUserReceivedDocuments_NonEmptyList(t *testing.T) {
	documents := types.Documents{TestingDocument}
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set([]byte(types.ReceivedDocumentsPrefix+TestingRecipient.String()), TestUtils.Cdc.MustMarshalBinaryBare(&documents))
	receivedDocs := TestUtils.DocsKeeper.GetUserReceivedDocuments(TestUtils.Ctx, TestingRecipient)

	assert.Equal(t, 1, len(receivedDocs))
	assert.Equal(t, documents, receivedDocs)
}

func TestKeeper_GetUserSentDocuments_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete([]byte(types.SentDocumentsPrefix + TestingSender.String()))

	sentDocuments := TestUtils.DocsKeeper.GetUserSentDocuments(TestUtils.Ctx, TestingDocument.Sender)
	assert.Equal(t, 0, len(sentDocuments))
}

func TestKeeper_GetUserSentDocuments_NonEmptyList(t *testing.T) {
	documents := types.Documents{TestingDocument}
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set([]byte(types.SentDocumentsPrefix+TestingSender.String()), TestUtils.Cdc.MustMarshalBinaryBare(&documents))

	sentDocuments := TestUtils.DocsKeeper.GetUserSentDocuments(TestUtils.Ctx, TestingSender)

	assert.Equal(t, 1, len(sentDocuments))
	assert.Equal(t, documents, sentDocuments)
}

// ----------------------------------
// --- Document receipts
// ----------------------------------

func TestKeeper_SendDocumentReceipt_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete(TestUtils.DocsKeeper.getSentReceiptsStoreKey(TestingDocumentReceipt.Sender))

	TestUtils.DocsKeeper.SendDocumentReceipt(TestUtils.Ctx, TestingDocumentReceipt)

	var stored types.DocumentReceipts
	docReceiptBz := store.Get(TestUtils.DocsKeeper.getSentReceiptsStoreKey(TestingDocumentReceipt.Sender))
	TestUtils.Cdc.MustUnmarshalBinaryBare(docReceiptBz, &stored)

	assert.Equal(t, 1, len(stored))
	assert.Equal(t, types.DocumentReceipts{TestingDocumentReceipt}, stored)
}

func TestKeeper_SendDocumentReceipt_ExistingReceipt(t *testing.T) {
	var existing = types.DocumentReceipts{TestingDocumentReceipt}

	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set(
		TestUtils.DocsKeeper.getSentReceiptsStoreKey(TestingDocumentReceipt.Sender),
		TestUtils.Cdc.MustMarshalBinaryBare(&existing),
	)

	TestUtils.DocsKeeper.SendDocumentReceipt(TestUtils.Ctx, TestingDocumentReceipt)

	var stored types.DocumentReceipts
	docReceiptBz := store.Get(TestUtils.DocsKeeper.getSentReceiptsStoreKey(TestingDocumentReceipt.Sender))
	TestUtils.Cdc.MustUnmarshalBinaryBare(docReceiptBz, &stored)

	assert.Equal(t, 1, len(stored))
	assert.Equal(t, existing, stored)
}

func TestKeeper_GetUserReceivedReceipts_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete(TestUtils.DocsKeeper.getReceivedReceiptsStoreKey(TestingDocumentReceipt.Recipient))

	receipts := TestUtils.DocsKeeper.GetUserReceivedReceipts(TestUtils.Ctx, TestingDocumentReceipt.Recipient)

	assert.Empty(t, receipts)
}

func TestKeeper_GetUserReceivedReceipts_FilledList(t *testing.T) {
	var existing = types.DocumentReceipts{TestingDocumentReceipt}

	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set(
		TestUtils.DocsKeeper.getReceivedReceiptsStoreKey(TestingDocumentReceipt.Recipient),
		TestUtils.Cdc.MustMarshalBinaryBare(&existing),
	)

	actualReceipts := TestUtils.DocsKeeper.GetUserReceivedReceipts(TestUtils.Ctx, TestingDocumentReceipt.Recipient)

	assert.Equal(t, existing, actualReceipts)
}

func TestKeeper_GetUserReceivedReceiptsForDocument_UuidNotFound(t *testing.T) {
	receipt := TestUtils.DocsKeeper.GetUserReceivedReceiptsForDocument(TestUtils.Ctx, TestingDocumentReceipt.Recipient, "111")
	assert.Empty(t, receipt)
}

func TestKeeper_GetUserReceivedReceiptsForDocument_UuidFound(t *testing.T) {
	var TestingDocumentReceipt2 = types.DocumentReceipt{
		Sender:       TestingSender2,
		Recipient:    TestingDocumentReceipt.Recipient,
		TxHash:       TestingDocumentReceipt.TxHash,
		DocumentUuid: TestingDocumentReceipt.DocumentUuid,
		Proof:        TestingDocumentReceipt.Proof,
	}

	stored := types.DocumentReceipts{TestingDocumentReceipt, TestingDocumentReceipt2}

	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set(
		TestUtils.DocsKeeper.getReceivedReceiptsStoreKey(TestingDocumentReceipt.Recipient),
		TestUtils.Cdc.MustMarshalBinaryBare(&stored),
	)

	actual := TestUtils.DocsKeeper.GetUserReceivedReceiptsForDocument(
		TestUtils.Ctx,
		TestingDocumentReceipt.Recipient,
		TestingDocumentReceipt.DocumentUuid,
	)

	assert.Equal(t, stored, actual)
}

// ----------------------------------
// --- Genesis utils
// ----------------------------------

func TestKeeper_GetUsersSet_FilledSet(t *testing.T) {
	TestUtils.DocsKeeper.ShareDocument(TestUtils.Ctx, TestingDocument)
	TestUtils.DocsKeeper.SendDocumentReceipt(TestUtils.Ctx, TestingDocumentReceipt)

	users, err := TestUtils.DocsKeeper.GetUsersSet(TestUtils.Ctx)

	assert.Nil(t, err)
	assert.Contains(t, users, TestingDocument.Sender)
	assert.Contains(t, users, TestingDocument.Recipient)
	assert.Contains(t, users, TestingDocumentReceipt.Sender)
	assert.Contains(t, users, TestingDocumentReceipt.Recipient)
}

func TestKeeper_GetUsersSet_EmptySet(t *testing.T) {
	// Cleanup the store
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	users, err := TestUtils.DocsKeeper.GetUsersSet(TestUtils.Ctx)

	assert.Nil(t, err)
	assert.Empty(t, users)
}

func TestKeeper_SetUserDocuments(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete(TestUtils.DocsKeeper.getSentDocumentsStoreKey(TestingDocument.Sender))
	store.Delete(TestUtils.DocsKeeper.getReceivedDocumentsStoreKey(TestingDocument.Recipient))

	documents := types.Documents{TestingDocument}

	TestUtils.DocsKeeper.SetUserDocuments(TestUtils.Ctx, TestingDocument.Sender, documents, types.Documents{})
	TestUtils.DocsKeeper.SetUserDocuments(TestUtils.Ctx, TestingDocument.Recipient, types.Documents{}, documents)

	sentBz := store.Get(TestUtils.DocsKeeper.getSentDocumentsStoreKey(TestingDocument.Sender))
	receivedBz := store.Get(TestUtils.DocsKeeper.getReceivedDocumentsStoreKey(TestingDocument.Recipient))

	var sentDocuments, receivedDocuments types.Documents
	TestUtils.Cdc.MustUnmarshalBinaryBare(sentBz, &sentDocuments)
	TestUtils.Cdc.MustUnmarshalBinaryBare(receivedBz, &receivedDocuments)

	assert.Equal(t, documents, sentDocuments)
	assert.Equal(t, documents, receivedDocuments)
}

func TestKeeper_SetUserReceipts(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete(TestUtils.DocsKeeper.getSentReceiptsStoreKey(TestingDocumentReceipt.Sender))
	store.Delete(TestUtils.DocsKeeper.getReceivedReceiptsStoreKey(TestingDocumentReceipt.Recipient))

	receipts := types.DocumentReceipts{TestingDocumentReceipt}

	TestUtils.DocsKeeper.SetUserReceipts(TestUtils.Ctx, TestingDocumentReceipt.Sender, receipts, types.DocumentReceipts{})
	TestUtils.DocsKeeper.SetUserReceipts(TestUtils.Ctx, TestingDocumentReceipt.Recipient, types.DocumentReceipts{}, receipts)

	sentBz := store.Get(TestUtils.DocsKeeper.getSentReceiptsStoreKey(TestingDocumentReceipt.Sender))
	receivedBz := store.Get(TestUtils.DocsKeeper.getReceivedReceiptsStoreKey(TestingDocumentReceipt.Recipient))

	var sentReceipts, receivedReceipts types.DocumentReceipts
	TestUtils.Cdc.MustUnmarshalBinaryBare(sentBz, &sentReceipts)
	TestUtils.Cdc.MustUnmarshalBinaryBare(receivedBz, &receivedReceipts)

	assert.Equal(t, receipts, sentReceipts)
	assert.Equal(t, receipts, receivedReceipts)
}
