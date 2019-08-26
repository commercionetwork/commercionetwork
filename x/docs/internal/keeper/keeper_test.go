package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_ShareDocument_EmptyLists(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)

	TestUtils.DocsKeeper.ShareDocument(TestUtils.Ctx, TestingDocument)

	sentDocsBz := store.Get([]byte(SentDocumentsPrefix + TestingSender.String()))
	receivedDocsBz := store.Get([]byte(ReceivedDocumentsPrefix + TestingRecipient.String()))

	var sentDocs, receivedDocs []types.Document
	TestUtils.Cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
	TestUtils.Cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)

	assert.Equal(t, 1, len(sentDocs))
	assert.Contains(t, sentDocs, TestingDocument)

	assert.Equal(t, 1, len(receivedDocs))
	assert.Equal(t, sentDocs, receivedDocs)
}

func TestKeeper_ShareDocument_ExistingDocument(t *testing.T) {
	documents := []types.Document{TestingDocument}
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set([]byte(SentDocumentsPrefix+TestingSender.String()), TestUtils.Cdc.MustMarshalBinaryBare(documents))
	store.Set([]byte(ReceivedDocumentsPrefix+TestingRecipient.String()), TestUtils.Cdc.MustMarshalBinaryBare(documents))

	TestUtils.DocsKeeper.ShareDocument(TestUtils.Ctx, TestingDocument)

	sentDocsBz := store.Get([]byte(SentDocumentsPrefix + TestingSender.String()))
	receivedDocsBz := store.Get([]byte(ReceivedDocumentsPrefix + TestingRecipient.String()))

	var sentDocs, receivedDocs []types.Document
	TestUtils.Cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
	TestUtils.Cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)

	assert.Equal(t, 1, len(sentDocs))
	assert.Contains(t, sentDocs, TestingDocument)

	assert.Equal(t, 1, len(receivedDocs))
	assert.Equal(t, sentDocs, receivedDocs)
}

func TestKeeper_ShareDocument_SameInfoDifferentRecipient(t *testing.T) {
	documents := []types.Document{TestingDocument}
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set([]byte(SentDocumentsPrefix+TestingSender.String()), TestUtils.Cdc.MustMarshalBinaryBare(documents))
	store.Set([]byte(ReceivedDocumentsPrefix+TestingRecipient.String()), TestUtils.Cdc.MustMarshalBinaryBare(documents))

	newRecipient, _ := sdk.AccAddressFromBech32("cosmos1h2z8u9294gtqmxlrnlyfueqysng3krh009fum7")
	newDocument := types.Document{
		Sender:     TestingDocument.Sender,
		Recipient:  newRecipient,
		ContentUri: TestingDocument.ContentUri,
		Metadata:   TestingDocument.Metadata,
		Checksum:   TestingDocument.Checksum,
	}
	TestUtils.DocsKeeper.ShareDocument(TestUtils.Ctx, newDocument)

	sentDocsBz := store.Get([]byte(SentDocumentsPrefix + TestingSender.String()))
	receivedDocsBz := store.Get([]byte(ReceivedDocumentsPrefix + TestingRecipient.String()))
	newReceivedDocsBz := store.Get([]byte(ReceivedDocumentsPrefix + newRecipient.String()))

	var sentDocs, receivedDocs, newReceivedDocs []types.Document
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
	store.Delete([]byte(ReceivedDocumentsPrefix + TestingRecipient.String()))

	receivedDocs := TestUtils.DocsKeeper.GetUserReceivedDocuments(TestUtils.Ctx, TestingDocument.Sender)
	assert.Equal(t, 0, len(receivedDocs))
}

func TestKeeper_GetUserReceivedDocuments_NonEmptyList(t *testing.T) {
	documents := []types.Document{TestingDocument}
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set([]byte(ReceivedDocumentsPrefix+TestingRecipient.String()), TestUtils.Cdc.MustMarshalBinaryBare(documents))

	receivedDocs := TestUtils.DocsKeeper.GetUserReceivedDocuments(TestUtils.Ctx, TestingRecipient)

	assert.Equal(t, 1, len(receivedDocs))
	assert.Equal(t, documents, receivedDocs)
}

func TestKeeper_GetUserSentDocuments_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Delete([]byte(SentDocumentsPrefix + TestingSender.String()))

	sentDocuments := TestUtils.DocsKeeper.GetUserSentDocuments(TestUtils.Ctx, TestingDocument.Sender)
	assert.Equal(t, 0, len(sentDocuments))
}

func TestKeeper_GetUserSentDocuments_NonEmptyList(t *testing.T) {
	documents := []types.Document{TestingDocument}
	store := TestUtils.Ctx.KVStore(TestUtils.DocsKeeper.StoreKey)
	store.Set([]byte(SentDocumentsPrefix+TestingSender.String()), TestUtils.Cdc.MustMarshalBinaryBare(documents))

	sentDocuments := TestUtils.DocsKeeper.GetUserSentDocuments(TestUtils.Ctx, TestingSender)

	assert.Equal(t, 1, len(sentDocuments))
	assert.Equal(t, documents, sentDocuments)
}
