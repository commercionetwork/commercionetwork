package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

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
// --- DocumentReceipt
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
