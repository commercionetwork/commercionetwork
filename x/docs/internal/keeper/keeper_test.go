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
	_, ctx, k := SetupTestInput()

	schema := types.MetadataSchema{Type: "schema", SchemaURI: "https://example.com/schema", Version: "1.0.0"}
	k.AddSupportedMetadataScheme(ctx, schema)

	ret := k.IsMetadataSchemeTypeSupported(ctx, schema.Type)
	assert.True(t, ret)
}

func TestKeeper_AddSupportedMetadataScheme_ExistingList(t *testing.T) {
	_, ctx, k := SetupTestInput()
	//Setup the store

	existingSchema := types.MetadataSchema{Type: "schema", SchemaURI: "https://example.com/newSchema", Version: "1.0.0"}
	k.AddSupportedMetadataScheme(ctx, existingSchema)

	newSchema := types.MetadataSchema{Type: "schema2", SchemaURI: "https://example.com/schema2", Version: "2.0.0"}
	k.AddSupportedMetadataScheme(ctx, newSchema)

	stored := []types.MetadataSchema{}
	msi := k.SupportedMetadataSchemesIterator(ctx)
	defer msi.Close()

	for ; msi.Valid(); msi.Next() {
		m := types.MetadataSchema{}
		k.cdc.MustUnmarshalBinaryBare(msi.Value(), &m)

		stored = append(stored, m)
	}

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
	_, ctx, k := SetupTestInput()

	existingSchema := types.MetadataSchema{Type: "schema", SchemaURI: "https://example.com/newSchema", Version: "1.0.0"}
	k.AddSupportedMetadataScheme(ctx, existingSchema)

	assert.True(t, k.IsMetadataSchemeTypeSupported(ctx, "schema"))
	assert.False(t, k.IsMetadataSchemeTypeSupported(ctx, "schema2"))
	assert.False(t, k.IsMetadataSchemeTypeSupported(ctx, "any-schema"))
}

func TestKeeper_SupportedMetadataSchemesIterator_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()

	result := []types.MetadataSchema{}
	smi := k.SupportedMetadataSchemesIterator(ctx)
	defer smi.Close()

	for ; smi.Valid(); smi.Next() {
		ms := types.MetadataSchema{}
		k.cdc.MustUnmarshalBinaryBare(smi.Value(), &ms)
		result = append(result, ms)
	}

	assert.Empty(t, result)
}

func TestKeeper_SupportedMetadataSchemesIterator_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	existingSchema := types.MetadataSchema{Type: "schema", SchemaURI: "https://example.com/newSchema", Version: "1.0.0"}
	existingBz := cdc.MustMarshalBinaryBare(existingSchema)
	store.Set(metadataSchemaKey(existingSchema), existingBz)

	result := []types.MetadataSchema{}
	smi := k.SupportedMetadataSchemesIterator(ctx)
	defer smi.Close()

	for ; smi.Valid(); smi.Next() {
		ms := types.MetadataSchema{}
		k.cdc.MustUnmarshalBinaryBare(smi.Value(), &ms)
		result = append(result, ms)
	}

	assert.Equal(t, []types.MetadataSchema{existingSchema}, result)
}

// ----------------------------------
// --- Metadata schema proposers
// ----------------------------------

func TestKeeper_AddTrustedSchemaProposer_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()

	k.AddTrustedSchemaProposer(ctx, TestingSender)
	ret := k.IsTrustedSchemaProposer(ctx, TestingSender)
	assert.True(t, ret)
}

func TestKeeper_AddTrustedSchemaProposer_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	proposersBz := cdc.MustMarshalBinaryBare(&TestingSender)
	store.Set(metadataSchemaProposerKey(TestingSender), proposersBz)

	k.AddTrustedSchemaProposer(ctx, TestingSender2)

	var stored []sdk.AccAddress

	tspi := k.TrustedSchemaProposersIterator(ctx)
	defer tspi.Close()

	for ; tspi.Valid(); tspi.Next() {
		p := sdk.AccAddress{}
		cdc.MustUnmarshalBinaryBare(tspi.Value(), &p)

		stored = append(stored, p)
	}

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
	_, ctx, k := SetupTestInput()

	k.AddTrustedSchemaProposer(ctx, TestingSender)

	assert.True(t, k.IsTrustedSchemaProposer(ctx, TestingSender))
	assert.False(t, k.IsTrustedSchemaProposer(ctx, TestingSender2))
}

func TestKeeper_TrustedSchemaProposersIterator_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()

	result := []sdk.AccAddress{}
	tspi := k.TrustedSchemaProposersIterator(ctx)
	defer tspi.Close()

	for ; tspi.Valid(); tspi.Next() {
		ms := sdk.AccAddress{}
		k.cdc.MustUnmarshalBinaryBare(tspi.Value(), &ms)
		result = append(result, ms)
	}

	assert.Empty(t, result)
}

func TestKeeper_TrustedSchemaProposersIterator_ExistingList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	proposersBz := cdc.MustMarshalBinaryBare(TestingSender)
	store.Set(metadataSchemaProposerKey(TestingSender), proposersBz)

	result := []sdk.AccAddress{}
	tspi := k.TrustedSchemaProposersIterator(ctx)
	defer tspi.Close()

	for ; tspi.Valid(); tspi.Next() {
		ms := sdk.AccAddress{}
		k.cdc.MustUnmarshalBinaryBare(tspi.Value(), &ms)
		result = append(result, ms)
	}

	assert.Equal(t, []sdk.AccAddress{TestingSender}, result)
}

// ----------------------------------
// --- Documents
// ----------------------------------

func TestKeeper_ShareDocument_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	err := k.SaveDocument(ctx, TestingDocument)
	assert.NoError(t, err)

	docsBz := store.Get(getDocumentStoreKey(TestingDocument.UUID))
	sentDocsBz := store.Get(getSentDocumentsIdsUUIDStoreKey(TestingSender, TestingDocument.UUID))
	receivedDocsBz := store.Get(getReceivedDocumentsIdsUUIDStoreKey(TestingRecipient, TestingDocument.UUID))

	var stored types.Document
	cdc.MustUnmarshalBinaryBare(docsBz, &stored)
	assert.Equal(t, stored, TestingDocument)

	var sentDocs, receivedDocs string
	cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
	cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)

	assert.Equal(t, TestingDocument.UUID, sentDocs)

	assert.Equal(t, TestingDocument.UUID, receivedDocs)
}

func TestKeeper_ShareDocument_ExistingDocument(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	store.Set(getDocumentStoreKey(TestingDocument.UUID), cdc.MustMarshalBinaryBare(TestingDocument))

	err := k.SaveDocument(ctx, TestingDocument)
	assert.NotNil(t, err)
}

func TestKeeper_ShareDocument_ExistingDocument_DifferentRecipient(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)
	store.Set(getSentDocumentsIdsUUIDStoreKey(TestingSender, TestingDocument.UUID), cdc.MustMarshalBinaryBare(TestingDocument.UUID))
	store.Set(getReceivedDocumentsIdsUUIDStoreKey(TestingRecipient, TestingDocument.UUID), cdc.MustMarshalBinaryBare(TestingDocument.UUID))

	newRecipient, _ := sdk.AccAddressFromBech32("cosmos1h2z8u9294gtqmxlrnlyfueqysng3krh009fum7")
	newDocument := types.Document{
		UUID:       TestingDocument.UUID,
		ContentURI: TestingDocument.ContentURI,
		Metadata:   TestingDocument.Metadata,
		Checksum:   TestingDocument.Checksum,
		Sender:     TestingDocument.Sender,
		Recipients: ctypes.Addresses{newRecipient},
	}
	err := k.SaveDocument(ctx, newDocument)
	assert.Nil(t, err)

	sentDocsBz := store.Get(getSentDocumentsIdsUUIDStoreKey(TestingSender, TestingDocument.UUID))
	receivedDocsBz := store.Get(getReceivedDocumentsIdsUUIDStoreKey(TestingRecipient, TestingDocument.UUID))
	newReceivedDocsBz := store.Get(getReceivedDocumentsIdsUUIDStoreKey(newRecipient, TestingDocument.UUID))

	var sentDocs, receivedDocs, newReceivedDocs string
	cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
	cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)
	cdc.MustUnmarshalBinaryBare(newReceivedDocsBz, &newReceivedDocs)

	assert.Equal(t, TestingDocument.UUID, sentDocs)

	assert.Equal(t, TestingDocument.UUID, receivedDocs)

	assert.Equal(t, newDocument.UUID, newReceivedDocs)
}

func TestKeeper_ShareDocument_ExistingDocument_DifferentUuid(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)
	store.Set(getSentDocumentsIdsUUIDStoreKey(TestingSender, TestingDocument.UUID), cdc.MustMarshalBinaryBare(TestingDocument.UUID))
	store.Set(getReceivedDocumentsIdsUUIDStoreKey(TestingRecipient, TestingDocument.UUID), cdc.MustMarshalBinaryBare(TestingDocument.UUID))

	newDocument := types.Document{
		UUID:       TestingDocument.UUID + "new",
		ContentURI: TestingDocument.ContentURI,
		Metadata:   TestingDocument.Metadata,
		Checksum:   TestingDocument.Checksum,
		Recipients: TestingDocument.Recipients,
		Sender:     TestingDocument.Sender,
	}
	err := k.SaveDocument(ctx, newDocument)
	assert.Nil(t, err)

	sentDocsBz := store.Get(getSentDocumentsIdsUUIDStoreKey(TestingSender, TestingDocument.UUID))
	receivedDocsBz := store.Get(getReceivedDocumentsIdsUUIDStoreKey(TestingRecipient, TestingDocument.UUID))

	var sentDocs, receivedDocs string
	cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
	cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)

	assert.Equal(t, TestingDocument.UUID, sentDocs)

	assert.Equal(t, TestingDocument.UUID, receivedDocs)
}

func TestKeeper_GetDocumentById_NonExisting(t *testing.T) {
	_, ctx, k := SetupTestInput()
	_, err := k.GetDocumentByID(ctx, "non-existing")
	assert.Error(t, err)
}

func TestKeeper_GetDocumentById_Existing(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)
	store.Set(getDocumentStoreKey(TestingDocument.UUID), cdc.MustMarshalBinaryBare(&TestingDocument))

	doc, err := k.GetDocumentByID(ctx, TestingDocument.UUID)
	assert.NoError(t, err)
	assert.Equal(t, TestingDocument, doc)
}

func TestKeeper_UserReceivedDocumentsIterator_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()

	rdi := k.UserReceivedDocumentsIterator(ctx, TestingRecipient)
	defer rdi.Close()

	docs := []types.Document{}
	for ; rdi.Valid(); rdi.Next() {
		doc, err := k.GetDocumentByID(ctx, string(rdi.Value()))
		assert.NoError(t, err)

		docs = append(docs, doc)
	}

	assert.Empty(t, docs)
}

func TestKeeper_UserReceivedDocumentsIterator_NonEmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	store.Set(getDocumentStoreKey(TestingDocument.UUID), cdc.MustMarshalBinaryBare(TestingDocument))
	store.Set(getReceivedDocumentsIdsUUIDStoreKey(TestingRecipient, TestingDocument.UUID), cdc.MustMarshalBinaryBare(TestingDocument.UUID))

	rdi := k.UserReceivedDocumentsIterator(ctx, TestingRecipient)
	defer rdi.Close()

	docs := []types.Document{}
	for ; rdi.Valid(); rdi.Next() {
		id := ""
		k.cdc.MustUnmarshalBinaryBare(rdi.Value(), &id)
		doc, err := k.GetDocumentByID(ctx, id)
		assert.NoError(t, err)

		docs = append(docs, doc)
	}

	assert.Equal(t, 1, len(docs))
	assert.Equal(t, []types.Document{TestingDocument}, docs)
}

func TestKeeper_UserSentDocumentsIterator_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()

	docs := []types.Document{}
	sdi := k.UserSentDocumentsIterator(ctx, TestingSender)
	defer sdi.Close()

	for ; sdi.Valid(); sdi.Next() {
		id := ""
		k.cdc.MustUnmarshalBinaryBare(sdi.Value(), &id)
		doc, err := k.GetDocumentByID(ctx, id)
		assert.NoError(t, err)

		docs = append(docs, doc)
	}

	assert.Empty(t, docs)
}

func TestKeeper_UserSentDocumentsIterator_NonEmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	store.Set(getDocumentStoreKey(TestingDocument.UUID), cdc.MustMarshalBinaryBare(TestingDocument))
	store.Set(getSentDocumentsIdsUUIDStoreKey(TestingRecipient, TestingDocument.UUID), cdc.MustMarshalBinaryBare(TestingDocument.UUID))

	rdi := k.UserSentDocumentsIterator(ctx, TestingRecipient)
	defer rdi.Close()

	docs := []types.Document{}
	for ; rdi.Valid(); rdi.Next() {
		id := ""
		k.cdc.MustUnmarshalBinaryBare(rdi.Value(), &id)
		doc, err := k.GetDocumentByID(ctx, id)
		assert.NoError(t, err)

		docs = append(docs, doc)
	}

	assert.Equal(t, 1, len(docs))
	assert.Equal(t, []types.Document{TestingDocument}, docs)
}

func TestKeeper_DocumentsIterator_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()
	di := k.DocumentsIterator(ctx)
	defer di.Close()

	documents := []types.Document{}
	for ; di.Valid(); di.Next() {
		d := types.Document{}
		k.cdc.MustUnmarshalBinaryBare(di.Value(), &d)

		documents = append(documents, d)
	}

	assert.Empty(t, documents)
}

func TestKeeper_DocumentsIterator_ExistingList(t *testing.T) {
	_, ctx, k := SetupTestInput()

	doc1 := TestingDocument
	doc2 := types.Document{
		UUID:           "uuid-2",
		Sender:         TestingDocument.Sender,
		Recipients:     TestingDocument.Recipients,
		Metadata:       TestingDocument.Metadata,
		ContentURI:     TestingDocument.ContentURI,
		Checksum:       TestingDocument.Checksum,
		EncryptionData: TestingDocument.EncryptionData,
	}
	assert.NoError(t, k.SaveDocument(ctx, doc1))
	assert.NoError(t, k.SaveDocument(ctx, doc2))

	di := k.DocumentsIterator(ctx)
	defer di.Close()

	docs := []types.Document{}
	for ; di.Valid(); di.Next() {
		d := types.Document{}
		k.cdc.MustUnmarshalBinaryBare(di.Value(), &d)

		docs = append(docs, d)
	}

	assert.Len(t, docs, 2)
	assert.Contains(t, docs, doc1)
	assert.Contains(t, docs, doc2)
}

// ----------------------------------
// --- Document receipts
// ----------------------------------

func TestKeeper_SaveDocumentReceipt_EmptyList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)

	assert.NoError(t, k.SaveDocument(ctx, TestingDocument))

	tdr := TestingDocumentReceipt
	tdr.DocumentUUID = TestingDocument.UUID
	assert.NoError(t, k.SaveReceipt(ctx, tdr))

	storedID := ""
	docReceiptBz := store.Get(getSentReceiptsIdsUUIDStoreKey(TestingDocumentReceipt.Sender, tdr.DocumentUUID))
	cdc.MustUnmarshalBinaryBare(docReceiptBz, &storedID)

	stored, err := k.GetReceiptByID(ctx, storedID)
	assert.NoError(t, err)

	assert.Equal(t, stored, tdr)
}

func TestKeeper_SaveDocumentReceipt_ExistingReceipt(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)
	store.Set(getSentReceiptsIdsUUIDStoreKey(TestingDocumentReceipt.Sender, TestingDocumentReceipt.UUID), cdc.MustMarshalBinaryBare(TestingDocumentReceipt))

	assert.Error(t, k.SaveReceipt(ctx, TestingDocumentReceipt))
}

func TestKeeper_SaveDocumentReceipt_ExistingReceipt_DifferentUuid(t *testing.T) {
	_, ctx, k := SetupTestInput()

	assert.NoError(t, k.SaveDocument(ctx, TestingDocument))

	oldReceipt := TestingDocumentReceipt
	oldReceipt.DocumentUUID = TestingDocument.UUID

	newReceipt := types.DocumentReceipt{
		UUID:         TestingDocumentReceipt.UUID + "-new",
		Sender:       TestingDocumentReceipt.Sender,
		Recipient:    TestingDocumentReceipt.Recipient,
		TxHash:       TestingDocumentReceipt.TxHash,
		DocumentUUID: TestingDocument.UUID,
		Proof:        TestingDocumentReceipt.Proof,
	}

	assert.NoError(t, k.SaveReceipt(ctx, oldReceipt))
	assert.Error(t, k.SaveReceipt(ctx, newReceipt))

	var stored []types.DocumentReceipt
	si := k.UserSentReceiptsIterator(ctx, TestingDocumentReceipt.Sender)
	defer si.Close()
	for ; si.Valid(); si.Next() {
		rid := ""
		k.cdc.MustUnmarshalBinaryBare(si.Value(), &rid)

		newReceipt, err := k.GetReceiptByID(ctx, rid)
		assert.NoError(t, err)
		stored = append(stored, newReceipt)
	}

	assert.Equal(t, 1, len(stored))
	assert.Contains(t, stored, oldReceipt)
	assert.NotContains(t, stored, newReceipt)
}

func TestKeeper_UserReceivedReceiptsIterator_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()

	urri := k.UserReceivedReceiptsIterator(ctx, TestingDocumentReceipt.Recipient)
	defer urri.Close()

	receipts := []types.DocumentReceipt{}
	for ; urri.Valid(); urri.Next() {
		rid := ""
		k.cdc.MustUnmarshalBinaryBare(urri.Value(), &rid)

		r, err := k.GetReceiptByID(ctx, rid)
		assert.NoError(t, err)

		receipts = append(receipts, r)
	}

	assert.Empty(t, receipts)
}

func TestKeeper_UserReceivedReceiptsIterator_FilledList(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)
	store.Set(getReceivedReceiptsIdsUUIDStoreKey(TestingDocumentReceipt.Recipient, TestingDocumentReceipt.UUID),
		cdc.MustMarshalBinaryBare(TestingDocumentReceipt.UUID))

	store.Set(getReceiptStoreKey(TestingDocumentReceipt.UUID), cdc.MustMarshalBinaryBare(TestingDocumentReceipt))

	urri := k.UserReceivedReceiptsIterator(ctx, TestingDocumentReceipt.Recipient)
	defer urri.Close()

	receipts := []types.DocumentReceipt{}
	for ; urri.Valid(); urri.Next() {
		rid := ""
		k.cdc.MustUnmarshalBinaryBare(urri.Value(), &rid)

		r, err := k.GetReceiptByID(ctx, rid)
		assert.NoError(t, err)

		receipts = append(receipts, r)
	}

	expected := []types.DocumentReceipt{TestingDocumentReceipt}

	assert.Equal(t, expected, receipts)
}

func TestKeeper_ExtractDocument(t *testing.T) {
	tests := []struct {
		name     string
		want     types.Document
		wantUUID string
		wantErr  bool
	}{
		{
			"stored document",
			TestingDocument,
			TestingDocument.UUID,
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput()

			assert.NoError(t, k.SaveDocument(ctx, tt.want))

			docKey := []byte{}

			di := k.DocumentsIterator(ctx)
			defer di.Close()
			for ; di.Valid(); di.Next() {
				docKey = di.Key()
			}

			extDoc, extUUID, extErr := k.ExtractDocument(ctx, docKey)

			if !tt.wantErr {
				assert.NoError(t, extErr)
				assert.Equal(t, tt.want, extDoc)
				assert.Equal(t, tt.wantUUID, extUUID)
			} else {
				assert.Error(t, extErr)
			}
		})
	}
}

func TestKeeper_ExtractMetadataSchema(t *testing.T) {
	tests := []struct {
		name string
		want types.MetadataSchema
	}{
		{
			"stored metadataSchema",
			types.MetadataSchema{Type: "ms"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput()
			k.AddSupportedMetadataScheme(ctx, tt.want)

			ki := k.SupportedMetadataSchemesIterator(ctx)
			defer ki.Close()

			mIterVal := []byte{}

			for ; ki.Valid(); ki.Next() {
				mIterVal = ki.Value()
			}

			m := k.ExtractMetadataSchema(mIterVal)

			assert.Equal(t, tt.want, m)
		})
	}
}

func TestKeeper_ExtractReceipt(t *testing.T) {
	r := TestingDocumentReceipt
	r.DocumentUUID = TestingDocument.UUID

	tests := []struct {
		name          string
		savedDocument types.Document
		want          types.DocumentReceipt
		wantUUID      string
		wantErr       bool
	}{
		{
			"stored receipt",
			TestingDocument,
			r,
			r.UUID,
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput()

			assert.NoError(t, k.SaveDocument(ctx, tt.savedDocument))
			assert.NoError(t, k.SaveReceipt(ctx, tt.want))

			recVal := []byte{}

			di, _ := k.ReceiptsIterators(ctx)
			defer di.Close()
			for ; di.Valid(); di.Next() {
				recVal = di.Value()
			}

			extDoc, extUUID, extErr := k.ExtractReceipt(ctx, recVal)

			if !tt.wantErr {
				assert.NoError(t, extErr)
				assert.Equal(t, tt.want, extDoc)
				assert.Equal(t, tt.wantUUID, extUUID)
			} else {
				assert.Error(t, extErr)
			}
		})
	}
}

func TestKeeper_ExtractTrustedSchemaProposer(t *testing.T) {
	tests := []struct {
		name string
		want sdk.AccAddress
	}{
		{
			"stored trusted schema proposer",
			TestingSender,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput()
			k.AddTrustedSchemaProposer(ctx, tt.want)

			ki := k.TrustedSchemaProposersIterator(ctx)
			defer ki.Close()

			mIterVal := []byte{}

			for ; ki.Valid(); ki.Next() {
				mIterVal = ki.Value()
			}

			m := k.ExtractTrustedSchemaProposer(mIterVal)

			assert.Equal(t, tt.want, m)
		})
	}
}

func TestKeeper_GetReceiptByID(t *testing.T) {
	r := TestingDocumentReceipt
	r.DocumentUUID = TestingDocument.UUID

	tests := []struct {
		name           string
		storedDocument types.Document
		want           types.DocumentReceipt
		wantErr        bool
	}{
		{
			"lookup on existing receipt",
			TestingDocument,
			r,
			false,
		},
		{
			"lookup on non existing receipt",
			types.Document{},
			types.DocumentReceipt{},
			true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput()

			if !tt.storedDocument.Equals(types.Document{}) {
				assert.NoError(t, k.SaveDocument(ctx, tt.storedDocument))
			}

			if !tt.want.Equals(types.DocumentReceipt{}) {
				assert.NoError(t, k.SaveReceipt(ctx, tt.want))
			}

			rr, err := k.GetReceiptByID(ctx, tt.want.UUID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, rr)
			}
		})
	}
}
