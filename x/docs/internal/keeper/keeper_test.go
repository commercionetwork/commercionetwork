package keeper

import (
	"testing"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// ----------------------------------
// --- Metadata schemes
// ----------------------------------

func TestKeeper_AddSupportedMetadataScheme_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	//Setup the store
	store := ctx.KVStore(k.StoreKey)

	schema := types.MetadataSchema{Type: "schema", SchemaUri: "https://example.com/schema", Version: "1.0.0"}
	k.AddSupportedMetadataScheme(ctx, schema)

	var stored types.MetadataSchemes
	storedBz := store.Get([]byte(types.SupportedMetadataSchemesStoreKey))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, 1, len(stored))
	assert.Contains(t, stored, schema)
}

func TestKeeper_AddSupportedMetadataScheme_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	//Setup the store
	store := ctx.KVStore(k.StoreKey)

	existingSchema := types.MetadataSchema{Type: "schema", SchemaUri: "https://example.com/newSchema", Version: "1.0.0"}
	existing := []types.MetadataSchema{existingSchema}
	existingBz := cdc.MustMarshalBinaryBare(&existing)
	store.Set([]byte(types.SupportedMetadataSchemesStoreKey), existingBz)

	newSchema := types.MetadataSchema{Type: "schema2", SchemaUri: "https://example.com/schema2", Version: "2.0.0"}
	k.AddSupportedMetadataScheme(ctx, newSchema)

	var stored types.MetadataSchemes
	storedBz := store.Get([]byte(types.SupportedMetadataSchemesStoreKey))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, 2, len(stored))
	assert.Contains(t, stored, existingSchema)
	assert.Contains(t, stored, newSchema)
}

func TestKeeper_IsMetadataSchemeTypeSupported_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()

	assert.False(t, k.IsMetadataSchemeTypeSupported(ctx, "schema"))
	assert.False(t, k.IsMetadataSchemeTypeSupported(ctx, "schema2"))
	assert.False(t, k.IsMetadataSchemeTypeSupported(ctx, "non-existent"))
}

func TestKeeper_IsMetadataSchemeTypeSupported_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	existingSchema := types.MetadataSchema{Type: "schema", SchemaUri: "https://example.com/newSchema", Version: "1.0.0"}
	existing := []types.MetadataSchema{existingSchema}
	existingBz := cdc.MustMarshalBinaryBare(&existing)
	store.Set([]byte(types.SupportedMetadataSchemesStoreKey), existingBz)

	assert.True(t, k.IsMetadataSchemeTypeSupported(ctx, "schema"))
	assert.False(t, k.IsMetadataSchemeTypeSupported(ctx, "schema2"))
	assert.False(t, k.IsMetadataSchemeTypeSupported(ctx, "any-schema"))
}

func TestKeeper_GetSupportedMetadataSchemes_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()

	result := k.GetSupportedMetadataSchemes(ctx)

	assert.Empty(t, result)
}

func TestKeeper_GetSupportedMetadataSchemes_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	existingSchema := types.MetadataSchema{Type: "schema", SchemaUri: "https://example.com/newSchema", Version: "1.0.0"}
	existing := types.MetadataSchemes{existingSchema}
	existingBz := cdc.MustMarshalBinaryBare(&existing)
	store.Set([]byte(types.SupportedMetadataSchemesStoreKey), existingBz)

	actual := k.GetSupportedMetadataSchemes(ctx)

	assert.Equal(t, existing, actual)
}

// ----------------------------------
// --- Metadata schema proposers
// ----------------------------------

func TestKeeper_AddTrustedSchemaProposer_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	k.AddTrustedSchemaProposer(ctx, TestingSender)

	var proposers []sdk.AccAddress
	proposersBz := store.Get([]byte(types.MetadataSchemaProposersStoreKey))
	cdc.MustUnmarshalBinaryBare(proposersBz, &proposers)

	assert.Equal(t, 1, len(proposers))
	assert.Contains(t, proposers, TestingSender)
}

func TestKeeper_AddTrustedSchemaProposer_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	existing := []sdk.AccAddress{TestingSender}
	proposersBz := cdc.MustMarshalBinaryBare(&existing)
	store.Set([]byte(types.MetadataSchemaProposersStoreKey), proposersBz)

	k.AddTrustedSchemaProposer(ctx, TestingSender2)

	var stored []sdk.AccAddress
	storedBz := store.Get([]byte(types.MetadataSchemaProposersStoreKey))
	cdc.MustUnmarshalBinaryBare(storedBz, &stored)

	assert.Equal(t, 2, len(stored))
	assert.Contains(t, stored, TestingSender)
	assert.Contains(t, stored, TestingSender2)
}

func TestKeeper_IsTrustedSchemaProposer_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()

	assert.False(t, k.IsTrustedSchemaProposer(ctx, TestingSender))
	assert.False(t, k.IsTrustedSchemaProposer(ctx, TestingSender2))
}

func TestKeeper_IsTrustedSchemaProposerExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	existing := []sdk.AccAddress{TestingSender}
	proposersBz := cdc.MustMarshalBinaryBare(&existing)
	store.Set([]byte(types.MetadataSchemaProposersStoreKey), proposersBz)

	assert.True(t, k.IsTrustedSchemaProposer(ctx, TestingSender))
	assert.False(t, k.IsTrustedSchemaProposer(ctx, TestingSender2))
}

func TestKeeper_GetTrustedSchemaProposers_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()

	stored := k.GetTrustedSchemaProposers(ctx)

	assert.Empty(t, stored)
}

func TestKeeper_GetTrustedSchemaProposers_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	existing := ctypes.Addresses{TestingSender}
	proposersBz := cdc.MustMarshalBinaryBare(&existing)
	store.Set([]byte(types.MetadataSchemaProposersStoreKey), proposersBz)

	stored := k.GetTrustedSchemaProposers(ctx)

	assert.Equal(t, existing, stored)
}

// ----------------------------------
// --- Documents
// ----------------------------------

func TestKeeper_ShareDocument_EmptyLists(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	err := k.ShareDocument(ctx, TestingSender, ctypes.Addresses{TestingRecipient}, TestingDocument)
	assert.Nil(t, err)

	docsBz := store.Get(k.getDocumentStoreKey(TestingDocument.Uuid))
	sentDocsBz := store.Get(k.getSentDocumentsStoreKey(TestingSender))
	receivedDocsBz := store.Get(k.getReceivedDocumentsStoreKey(TestingRecipient))

	var stored types.Document
	cdc.MustUnmarshalBinaryBare(docsBz, &stored)
	assert.Equal(t, stored, TestingDocument)

	var sentDocs, receivedDocs types.DocumentIds
	cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
	cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)

	assert.Equal(t, 1, len(sentDocs))
	assert.Contains(t, sentDocs, TestingDocument.Uuid)

	assert.Equal(t, 1, len(receivedDocs))
	assert.Contains(t, receivedDocs, TestingDocument.Uuid)
}

func TestKeeper_ShareDocument_ExistingDocument(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	store.Set(k.getDocumentStoreKey(TestingDocument.Uuid), cdc.MustMarshalBinaryBare(TestingDocument))

	documentsIds := types.DocumentIds{TestingDocument.Uuid}
	store.Set(k.getSentDocumentsStoreKey(TestingSender), cdc.MustMarshalBinaryBare(&documentsIds))
	store.Set(k.getReceivedDocumentsStoreKey(TestingRecipient), cdc.MustMarshalBinaryBare(&documentsIds))

	err := k.ShareDocument(ctx, TestingSender, ctypes.Addresses{TestingRecipient}, TestingDocument)
	assert.NotNil(t, err)

	sentDocsBz := store.Get(k.getSentDocumentsStoreKey(TestingSender))
	receivedDocsBz := store.Get(k.getReceivedDocumentsStoreKey(TestingRecipient))

	var sentDocs, receivedDocs types.DocumentIds
	cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
	cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)

	assert.Equal(t, 1, len(sentDocs))
	assert.Contains(t, sentDocs, TestingDocument.Uuid)

	assert.Equal(t, 1, len(receivedDocs))
	assert.Contains(t, receivedDocs, TestingDocument.Uuid)
}

func TestKeeper_ShareDocument_SameDocumentDifferentRecipient(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	documentsIds := types.DocumentIds{TestingDocument.Uuid}

	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getSentDocumentsStoreKey(TestingSender), cdc.MustMarshalBinaryBare(&documentsIds))
	store.Set(k.getReceivedDocumentsStoreKey(TestingRecipient), cdc.MustMarshalBinaryBare(&documentsIds))

	newRecipient, _ := sdk.AccAddressFromBech32("cosmos1h2z8u9294gtqmxlrnlyfueqysng3krh009fum7")
	newDocument := types.Document{
		ContentUri: TestingDocument.ContentUri,
		Metadata:   TestingDocument.Metadata,
		Checksum:   TestingDocument.Checksum,
	}
	err := k.ShareDocument(ctx, TestingSender, ctypes.Addresses{newRecipient}, newDocument)
	assert.Nil(t, err)

	sentDocsBz := store.Get(k.getSentDocumentsStoreKey(TestingSender))
	receivedDocsBz := store.Get(k.getReceivedDocumentsStoreKey(TestingRecipient))
	newReceivedDocsBz := store.Get(k.getReceivedDocumentsStoreKey(newRecipient))

	var sentDocs, receivedDocs, newReceivedDocs types.DocumentIds
	cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
	cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)
	cdc.MustUnmarshalBinaryBare(newReceivedDocsBz, &newReceivedDocs)

	assert.Equal(t, 1, len(sentDocs))
	assert.Contains(t, sentDocs, TestingDocument.Uuid)

	assert.Equal(t, 1, len(receivedDocs))
	assert.Contains(t, receivedDocs, TestingDocument.Uuid)

	assert.Equal(t, 1, len(newReceivedDocs))
	assert.Contains(t, newReceivedDocs, newDocument.Uuid)
}

func TestKeeper_ShareDocument_SameDocumentDifferentUuid(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	documentsIds := types.DocumentIds{TestingDocument.Uuid}

	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getSentDocumentsStoreKey(TestingSender), cdc.MustMarshalBinaryBare(&documentsIds))
	store.Set(k.getReceivedDocumentsStoreKey(TestingRecipient), cdc.MustMarshalBinaryBare(&documentsIds))

	newDocument := types.Document{
		Uuid:       TestingDocument.Uuid + "new",
		ContentUri: TestingDocument.ContentUri,
		Metadata:   TestingDocument.Metadata,
		Checksum:   TestingDocument.Checksum,
	}
	err := k.ShareDocument(ctx, TestingSender, ctypes.Addresses{TestingRecipient}, newDocument)
	assert.Nil(t, err)

	sentDocsBz := store.Get(k.getSentDocumentsStoreKey(TestingSender))
	receivedDocsBz := store.Get(k.getReceivedDocumentsStoreKey(TestingRecipient))

	var sentDocs, receivedDocs types.DocumentIds
	cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
	cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)

	assert.Equal(t, 2, len(sentDocs))
	assert.Contains(t, sentDocs, TestingDocument.Uuid)

	assert.Equal(t, 2, len(receivedDocs))
	assert.Contains(t, receivedDocs, TestingDocument.Uuid)
}

func TestKeeper_GetUserReceivedDocuments_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()

	receivedDocs, err := k.GetUserReceivedDocuments(ctx, TestingRecipient)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(receivedDocs))
}

func TestKeeper_GetUserReceivedDocuments_NonEmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	documents := types.Documents{TestingDocument}
	documentIds := types.DocumentIds{TestingDocument.Uuid}

	store.Set(k.getDocumentStoreKey(TestingDocument.Uuid), cdc.MustMarshalBinaryBare(TestingDocument))
	store.Set(k.getReceivedDocumentsStoreKey(TestingRecipient), cdc.MustMarshalBinaryBare(&documentIds))

	receivedDocs, err := k.GetUserReceivedDocuments(ctx, TestingRecipient)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(receivedDocs))
	assert.Equal(t, documents, receivedDocs)
}

func TestKeeper_GetUserSentDocuments_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()

	sentDocuments, err := k.GetUserSentDocuments(ctx, TestingSender)
	assert.Nil(t, err)

	assert.Equal(t, 0, len(sentDocuments))
}

func TestKeeper_GetUserSentDocuments_NonEmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	documents := types.Documents{TestingDocument}
	documentIds := types.DocumentIds{TestingDocument.Uuid}
	store.Set(k.getDocumentStoreKey(TestingDocument.Uuid), cdc.MustMarshalBinaryBare(TestingDocument))
	store.Set(k.getSentDocumentsStoreKey(TestingSender), cdc.MustMarshalBinaryBare(&documentIds))

	sentDocuments, err := k.GetUserSentDocuments(ctx, TestingSender)
	assert.Nil(t, err)

	assert.Equal(t, 1, len(sentDocuments))
	assert.Equal(t, documents, sentDocuments)
}

// ----------------------------------
// --- Document receipts
// ----------------------------------

func TestKeeper_SendDocumentReceipt_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	k.SendDocumentReceipt(ctx, TestingDocumentReceipt)

	var stored types.DocumentReceipts
	docReceiptBz := store.Get(k.getSentReceiptsStoreKey(TestingDocumentReceipt.Sender))
	cdc.MustUnmarshalBinaryBare(docReceiptBz, &stored)

	assert.Equal(t, 1, len(stored))
	assert.Equal(t, types.DocumentReceipts{TestingDocumentReceipt}, stored)
}

func TestKeeper_SendDocumentReceipt_ExistingReceipt(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var existing = types.DocumentReceipts{TestingDocumentReceipt}

	store := ctx.KVStore(k.StoreKey)
	store.Set(
		k.getSentReceiptsStoreKey(TestingDocumentReceipt.Sender),
		cdc.MustMarshalBinaryBare(&existing),
	)

	k.SendDocumentReceipt(ctx, TestingDocumentReceipt)

	var stored types.DocumentReceipts
	docReceiptBz := store.Get(k.getSentReceiptsStoreKey(TestingDocumentReceipt.Sender))
	cdc.MustUnmarshalBinaryBare(docReceiptBz, &stored)

	assert.Equal(t, 1, len(stored))
	assert.Equal(t, existing, stored)
}

func TestKeeper_GetUserReceivedReceipts_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()

	receipts := k.GetUserReceivedReceipts(ctx, TestingDocumentReceipt.Recipient)

	assert.Empty(t, receipts)
}

func TestKeeper_GetUserReceivedReceipts_FilledList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var existing = types.DocumentReceipts{TestingDocumentReceipt}

	store := ctx.KVStore(k.StoreKey)
	store.Set(
		k.getReceivedReceiptsStoreKey(TestingDocumentReceipt.Recipient),
		cdc.MustMarshalBinaryBare(&existing),
	)

	actualReceipts := k.GetUserReceivedReceipts(ctx, TestingDocumentReceipt.Recipient)

	assert.Equal(t, existing, actualReceipts)
}

func TestKeeper_GetUserReceivedReceiptsForDocument_UuidNotFound(t *testing.T) {
	_, ctx, k := SetupTestInput()
	receipt := k.GetUserReceivedReceiptsForDocument(ctx, TestingDocumentReceipt.Recipient, "111")
	assert.Empty(t, receipt)
}

func TestKeeper_GetUserReceivedReceiptsForDocument_UuidFound(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var TestingDocumentReceipt2 = types.DocumentReceipt{
		Sender:       TestingSender2,
		Recipient:    TestingDocumentReceipt.Recipient,
		TxHash:       TestingDocumentReceipt.TxHash,
		DocumentUuid: TestingDocumentReceipt.DocumentUuid,
		Proof:        TestingDocumentReceipt.Proof,
	}

	stored := types.DocumentReceipts{TestingDocumentReceipt, TestingDocumentReceipt2}

	store := ctx.KVStore(k.StoreKey)
	store.Set(
		k.getReceivedReceiptsStoreKey(TestingDocumentReceipt.Recipient),
		cdc.MustMarshalBinaryBare(&stored),
	)

	actual := k.GetUserReceivedReceiptsForDocument(
		ctx,
		TestingDocumentReceipt.Recipient,
		TestingDocumentReceipt.DocumentUuid,
	)

	assert.Equal(t, stored, actual)
}

// ----------------------------------
// --- Genesis utils
// ----------------------------------

func TestKeeper_GetUsersSet_FilledSet(t *testing.T) {
	_, ctx, k := SetupTestInput()
	_ = k.ShareDocument(ctx, TestingSender, ctypes.Addresses{TestingRecipient}, TestingDocument)
	k.SendDocumentReceipt(ctx, TestingDocumentReceipt)

	users, err := k.GetUsersSet(ctx)

	assert.Nil(t, err)
	assert.Contains(t, users, TestingSender)
	assert.Contains(t, users, TestingRecipient)
	assert.Contains(t, users, TestingDocumentReceipt.Sender)
	assert.Contains(t, users, TestingDocumentReceipt.Recipient)
}

func TestKeeper_GetUsersSet_EmptySet(t *testing.T) {
	_, ctx, k := SetupTestInput()
	users, err := k.GetUsersSet(ctx)

	assert.Nil(t, err)
	assert.Empty(t, users)
}

func TestKeeper_SetUserDocuments(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	documents := types.Documents{TestingDocument}

	k.SetUserDocuments(ctx, TestingSender, documents, types.Documents{})
	k.SetUserDocuments(ctx, TestingRecipient, types.Documents{}, documents)

	sentBz := store.Get(k.getSentDocumentsStoreKey(TestingSender))
	receivedBz := store.Get(k.getReceivedDocumentsStoreKey(TestingRecipient))

	var sentDocuments, receivedDocuments types.Documents
	cdc.MustUnmarshalBinaryBare(sentBz, &sentDocuments)
	cdc.MustUnmarshalBinaryBare(receivedBz, &receivedDocuments)

	assert.Equal(t, documents, sentDocuments)
	assert.Equal(t, documents, receivedDocuments)
}

func TestKeeper_SetUserReceipts(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	receipts := types.DocumentReceipts{TestingDocumentReceipt}

	k.SetUserReceipts(ctx, TestingDocumentReceipt.Sender, receipts, types.DocumentReceipts{})
	k.SetUserReceipts(ctx, TestingDocumentReceipt.Recipient, types.DocumentReceipts{}, receipts)

	sentBz := store.Get(k.getSentReceiptsStoreKey(TestingDocumentReceipt.Sender))
	receivedBz := store.Get(k.getReceivedReceiptsStoreKey(TestingDocumentReceipt.Recipient))

	var sentReceipts, receivedReceipts types.DocumentReceipts
	cdc.MustUnmarshalBinaryBare(sentBz, &sentReceipts)
	cdc.MustUnmarshalBinaryBare(receivedBz, &receivedReceipts)

	assert.Equal(t, receipts, sentReceipts)
	assert.Equal(t, receipts, receivedReceipts)
}
