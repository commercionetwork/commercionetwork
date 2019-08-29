package keeper

import (
	"testing"

	keys "github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_ShareDocument_EmptyLists(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)

	TestUtils.DocsKeeper.ShareDocument(TestUtils.Ctx, TestingDocument)

	sentDocsBz := store.Get([]byte(keys.SentDocumentsPrefix + TestingSender.String()))
	receivedDocsBz := store.Get([]byte(keys.ReceivedDocumentsPrefix + TestingRecipient.String()))

	var sentDocs, receivedDocs []keys.Document
	TestUtils.Cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
	TestUtils.Cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)

	assert.Equal(t, 1, len(sentDocs))
	assert.Contains(t, sentDocs, TestingDocument)

	assert.Equal(t, 1, len(receivedDocs))
	assert.Equal(t, sentDocs, receivedDocs)
}

func TestKeeper_ShareDocument_ExistingDocument(t *testing.T) {
	documents := []keys.Document{TestingDocument}
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set([]byte(keys.SentDocumentsPrefix+TestingSender.String()), TestUtils.Cdc.MustMarshalBinaryBare(&documents))
	store.Set([]byte(keys.ReceivedDocumentsPrefix+TestingRecipient.String()), TestUtils.Cdc.MustMarshalBinaryBare(&documents))

	TestUtils.DocsKeeper.ShareDocument(TestUtils.Ctx, TestingDocument)

	sentDocsBz := store.Get([]byte(keys.SentDocumentsPrefix + TestingSender.String()))
	receivedDocsBz := store.Get([]byte(keys.ReceivedDocumentsPrefix + TestingRecipient.String()))

	var sentDocs, receivedDocs []keys.Document
	TestUtils.Cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
	TestUtils.Cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)

	assert.Equal(t, 1, len(sentDocs))
	assert.Contains(t, sentDocs, TestingDocument)

	assert.Equal(t, 1, len(receivedDocs))
	assert.Equal(t, sentDocs, receivedDocs)
}

func TestKeeper_ShareDocument_SameInfoDifferentRecipient(t *testing.T) {
	documents := []keys.Document{TestingDocument}
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set([]byte(keys.SentDocumentsPrefix+TestingSender.String()), TestUtils.Cdc.MustMarshalBinaryBare(&documents))
	store.Set([]byte(keys.ReceivedDocumentsPrefix+TestingRecipient.String()), TestUtils.Cdc.MustMarshalBinaryBare(&documents))

	newRecipient, _ := sdk.AccAddressFromBech32("cosmos1h2z8u9294gtqmxlrnlyfueqysng3krh009fum7")
	newDocument := keys.Document{
		Sender:     TestingDocument.Sender,
		Recipient:  newRecipient,
		ContentUri: TestingDocument.ContentUri,
		Metadata:   TestingDocument.Metadata,
		Checksum:   TestingDocument.Checksum,
	}
	TestUtils.DocsKeeper.ShareDocument(TestUtils.Ctx, newDocument)

	sentDocsBz := store.Get([]byte(keys.SentDocumentsPrefix + TestingSender.String()))
	receivedDocsBz := store.Get([]byte(keys.ReceivedDocumentsPrefix + TestingRecipient.String()))
	newReceivedDocsBz := store.Get([]byte(keys.ReceivedDocumentsPrefix + newRecipient.String()))

	var sentDocs, receivedDocs, newReceivedDocs []keys.Document
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
	store.Delete([]byte(keys.ReceivedDocumentsPrefix + TestingRecipient.String()))

	receivedDocs := TestUtils.DocsKeeper.GetUserReceivedDocuments(TestUtils.Ctx, TestingDocument.Sender)
	assert.Equal(t, 0, len(receivedDocs))
}

func TestKeeper_GetUserReceivedDocuments_NonEmptyList(t *testing.T) {
	documents := []keys.Document{TestingDocument}
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set([]byte(keys.ReceivedDocumentsPrefix+TestingRecipient.String()), TestUtils.Cdc.MustMarshalBinaryBare(&documents))
	receivedDocs := TestUtils.DocsKeeper.GetUserReceivedDocuments(TestUtils.Ctx, TestingRecipient)

	assert.Equal(t, 1, len(receivedDocs))
	assert.Equal(t, documents, receivedDocs)
}

func TestKeeper_GetUserSentDocuments_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete([]byte(keys.SentDocumentsPrefix + TestingSender.String()))

	sentDocuments := TestUtils.DocsKeeper.GetUserSentDocuments(TestUtils.Ctx, TestingDocument.Sender)
	assert.Equal(t, 0, len(sentDocuments))
}

func TestKeeper_GetUserSentDocuments_NonEmptyList(t *testing.T) {
	documents := []keys.Document{TestingDocument}
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set([]byte(keys.SentDocumentsPrefix+TestingSender.String()), TestUtils.Cdc.MustMarshalBinaryBare(&documents))

	sentDocuments := TestUtils.DocsKeeper.GetUserSentDocuments(TestUtils.Ctx, TestingSender)

	assert.Equal(t, 1, len(sentDocuments))
	assert.Equal(t, documents, sentDocuments)
}

// ----------------------------------
// --- DocumentReceipt
// ----------------------------------

func TestKeeper_ShareDocumentReceipt_EmptyList(t *testing.T) {
	TestUtils.DocsKeeper.ShareDocumentReceipt(TestUtils.Ctx, TestingDocumentReceipt)

	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	docReceiptBz := store.Get([]byte(keys.DocumentReceiptPrefix + TestingDocumentReceipt.Uuid +
		TestingDocumentReceipt.Recipient.String()))

	var actual keys.DocumentReceipt

	TestUtils.Cdc.MustUnmarshalBinaryBare(docReceiptBz, &actual)

	assert.Equal(t, TestingDocumentReceipt, actual)
}

func TestKeeper_ShareDocumentReceipt_ExistingReceipt(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)

	store.Set([]byte(keys.DocumentReceiptPrefix+TestingDocumentReceipt.Uuid+TestingDocumentReceipt.Recipient.String()),
		TestUtils.Cdc.MustMarshalBinaryBare(TestingDocumentReceipt))

	TestUtils.DocsKeeper.ShareDocumentReceipt(TestUtils.Ctx, TestingDocumentReceipt)

	var counter = 0
	iterator := sdk.KVStorePrefixIterator(store, []byte(keys.DocumentReceiptPrefix))
	for ; iterator.Valid(); iterator.Next() {
		counter++
	}

	assert.Equal(t, 1, counter)
}

func TestKeeper_GetUserReceivedReceipts_EmptyList(t *testing.T) {

	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete([]byte(keys.DocumentReceiptPrefix + TestingDocumentReceipt.Uuid + TestingDocumentReceipt.Recipient.String()))

	receipts := TestUtils.DocsKeeper.GetUserReceivedReceipts(TestUtils.Ctx, TestingDocumentReceipt.Recipient)

	assert.Empty(t, receipts)
}

func TestKeeper_GetUserReceivedReceipts_FilledList(t *testing.T) {
	var TestingDocumentReceipt2 = keys.DocumentReceipt{
		Sender:    TestingRecipient,
		Recipient: TestingSender,
		TxHash:    "txHash",
		Uuid:      "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
		Proof:     "proof",
	}

	var expectedReceipts = []keys.DocumentReceipt{TestingDocumentReceipt}

	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set([]byte(keys.DocumentReceiptPrefix+TestingDocumentReceipt.Uuid+TestingDocumentReceipt.Recipient.String()),
		TestUtils.Cdc.MustMarshalBinaryBare(&TestingDocumentReceipt))
	store.Set([]byte(keys.DocumentReceiptPrefix+TestingDocumentReceipt2.Uuid+TestingDocumentReceipt2.Recipient.String()),
		TestUtils.Cdc.MustMarshalBinaryBare(&TestingDocumentReceipt2))

	actualReceipts := TestUtils.DocsKeeper.GetUserReceivedReceipts(TestUtils.Ctx, TestingDocumentReceipt.Recipient)

	assert.Equal(t, expectedReceipts, actualReceipts)
}

func TestKeeper_GetReceiptByDocumentUuid_UuidNotFound(t *testing.T) {
	receipt := TestUtils.DocsKeeper.GetReceiptByDocumentUuid(TestUtils.Ctx, TestingDocumentReceipt.Recipient, "111")
	assert.Empty(t, receipt)
}

func TestKeeper_GetReceiptByDocumentUuid_UuidFound(t *testing.T) {
	var TestingDocumentReceipt2 = keys.DocumentReceipt{
		Sender:    TestingRecipient,
		Recipient: TestingSender,
		TxHash:    "txHash",
		Uuid:      "6a2f41a3-c54c-fce8-32d2-0324e1c32e25",
		Proof:     "proof",
	}

	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set([]byte(keys.DocumentReceiptPrefix+TestingDocumentReceipt.Uuid+TestingDocumentReceipt.Recipient.String()),
		TestUtils.Cdc.MustMarshalBinaryBare(TestingDocumentReceipt))
	store.Set([]byte(keys.DocumentReceiptPrefix+TestingDocumentReceipt2.Uuid+TestingDocumentReceipt2.Recipient.String()),
		TestUtils.Cdc.MustMarshalBinaryBare(TestingDocumentReceipt2))

	actual := TestUtils.DocsKeeper.GetReceiptByDocumentUuid(TestUtils.Ctx, TestingDocumentReceipt.Recipient,
		TestingDocumentReceipt.Uuid)
	actual2 := TestUtils.DocsKeeper.GetReceiptByDocumentUuid(TestUtils.Ctx, TestingDocumentReceipt2.Recipient,
		TestingDocumentReceipt2.Uuid)

	assert.Equal(t, TestingDocumentReceipt, actual)
	assert.Equal(t, TestingDocumentReceipt2, actual2)
}
