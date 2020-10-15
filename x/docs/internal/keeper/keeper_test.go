package keeper

import (
	"testing"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"

	"github.com/commercionetwork/commercionetwork/x/docs/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// ----------------------------------
// --- Metadata schemes
// ----------------------------------

func TestKeeper_AddSupportedMetadataScheme(t *testing.T) {
	tests := []struct {
		name           string
		existingSchema []types.MetadataSchema
		newSchemas     []types.MetadataSchema
		correctType    bool
	}{
		{
			"no new schemas",
			[]types.MetadataSchema{
				{Type: "schema", SchemaURI: "https://example.com/schema", Version: "1.0.0"},
			},
			nil,
			true,
		},
		{
			"1 new schema",
			[]types.MetadataSchema{
				{Type: "schema", SchemaURI: "https://example.com/schema", Version: "1.0.0"},
			},
			[]types.MetadataSchema{
				{Type: "schema2", SchemaURI: "https://example.com/schema2", Version: "2.0.0"},
			},
			true,
		},
		{
			"2 new schemas",
			[]types.MetadataSchema{
				{Type: "schema", SchemaURI: "https://example.com/schema", Version: "1.0.0"},
			},
			[]types.MetadataSchema{
				{Type: "schema2", SchemaURI: "https://example.com/schema2", Version: "2.0.0"},
				{Type: "schema3", SchemaURI: "https://example.com/schema3", Version: "3.0.0"},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput()

			for _, pms := range tt.existingSchema {
				store := ctx.KVStore(k.StoreKey)
				msk := metadataSchemaKey(pms)
				store.Set(msk, k.cdc.MustMarshalBinaryBare(pms))
			}

			if tt.newSchemas == nil {
				for _, pms := range tt.existingSchema {
					supported := k.IsMetadataSchemeTypeSupported(ctx, pms.Type)
					require.Equal(t, tt.correctType, supported)
				}
				return
			}

			for _, nms := range tt.newSchemas {
				k.AddSupportedMetadataScheme(ctx, nms)
				supported := k.IsMetadataSchemeTypeSupported(ctx, nms.Type)
				require.Equal(t, true, supported)
			}

			stored := []types.MetadataSchema{}
			msi := k.SupportedMetadataSchemesIterator(ctx)
			defer msi.Close()

			for ; msi.Valid(); msi.Next() {
				m := types.MetadataSchema{}
				k.cdc.MustUnmarshalBinaryBare(msi.Value(), &m)

				stored = append(stored, m)
			}

			require.Equal(t, len(tt.newSchemas)+len(tt.existingSchema), len(stored))

			for _, nms := range tt.newSchemas {
				require.Contains(t, stored, nms)
			}
		})
	}
}

func TestKeeper_IsMetadataSchemeTypeSupported(t *testing.T) {
	tests := []struct {
		name                       string
		preexistantMetadataSchemes []types.MetadataSchema
		metadataSchemaPresent      bool
		metadataSchema             string
	}{
		{
			"schema not supported, no preexistant schemas",
			nil,
			false,
			"aSchema",
		},
		{
			"schema not supported, preexistant schemas",
			[]types.MetadataSchema{
				{Type: "schema", SchemaURI: "https://example.com/newSchema", Version: "1.0.0"},
			},
			false,
			"aSchema",
		},
		{
			"schema supported, preexistant schemas",
			[]types.MetadataSchema{
				{Type: "aSchema", SchemaURI: "https://example.com/newSchema", Version: "1.0.0"},
			},
			true,
			"aSchema",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput()
			for _, pms := range tt.preexistantMetadataSchemes {
				k.AddSupportedMetadataScheme(ctx, pms)
			}
			supported := k.IsMetadataSchemeTypeSupported(ctx, tt.metadataSchema)
			require.Equal(t, tt.metadataSchemaPresent, supported)
		})
	}
}

func TestKeeper_SupportedMetadataSchemesIterator(t *testing.T) {
	tests := []struct {
		name   string
		schema []types.MetadataSchema
	}{
		{
			"Empty list",
			[]types.MetadataSchema{},
		},
		{
			"1 element in list",
			[]types.MetadataSchema{
				{
					Type:      "schema",
					SchemaURI: "https://example.com/newSchema",
					Version:   "1.0.0",
				},
			},
		},
		{
			"2 elements in list",
			[]types.MetadataSchema{
				{
					Type:      "schema",
					SchemaURI: "https://example.com/newSchema",
					Version:   "1.0.0",
				},
				{
					Type:      "schema2",
					SchemaURI: "https://example.com/schema2",
					Version:   "2.0.0",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc, ctx, k := SetupTestInput()

			for _, ms := range tt.schema {
				store := ctx.KVStore(k.StoreKey)

				existingBz := cdc.MustMarshalBinaryBare(ms)
				store.Set(metadataSchemaKey(ms), existingBz)
			}

			result := []types.MetadataSchema{}
			smi := k.SupportedMetadataSchemesIterator(ctx)
			defer smi.Close()

			for ; smi.Valid(); smi.Next() {
				ms := types.MetadataSchema{}
				k.cdc.MustUnmarshalBinaryBare(smi.Value(), &ms)
				result = append(result, ms)
			}

			if tt.schema == nil {
				require.Empty(t, result)
			}

			require.Equal(t, len(tt.schema), len(result))

			for _, ms := range tt.schema {
				require.Contains(t, result, ms)
			}
		})
	}
}

// ----------------------------------
// --- Metadata schema proposers
// ----------------------------------

func TestKeeper_AddTrustedSchemaProposer(t *testing.T) {
	tests := []struct {
		name          string
		storedAddress sdk.AccAddress
		senderAddress sdk.AccAddress
	}{
		{
			"No stored address",
			nil,
			TestingSender,
		},
		{
			"1 element in stored address",
			TestingSender,
			TestingSender2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc, ctx, k := SetupTestInput()

			if tt.storedAddress != nil {
				store := ctx.KVStore(k.StoreKey)

				proposersBz := cdc.MustMarshalBinaryBare(&tt.storedAddress)
				store.Set(metadataSchemaProposerKey(tt.storedAddress), proposersBz)
			}

			k.AddTrustedSchemaProposer(ctx, tt.senderAddress)

			if tt.storedAddress == nil {
				ret := k.IsTrustedSchemaProposer(ctx, tt.senderAddress)
				require.True(t, ret)
				return
			}
			var stored []sdk.AccAddress

			tspi := k.TrustedSchemaProposersIterator(ctx)
			defer tspi.Close()

			for ; tspi.Valid(); tspi.Next() {
				p := sdk.AccAddress{}
				cdc.MustUnmarshalBinaryBare(tspi.Value(), &p)

				stored = append(stored, p)
			}

			require.Equal(t, 2, len(stored))
			require.Contains(t, stored, tt.storedAddress)
			require.Contains(t, stored, tt.senderAddress)

		})
	}
}

func TestKeeper_IsTrustedSchemaProposer(t *testing.T) {
	tests := []struct {
		name           string
		isEmpty        bool
		senderAddress  sdk.AccAddress
		senderAddress2 sdk.AccAddress
	}{
		{
			"Empty list",
			true,
			TestingSender,
			TestingSender2,
		},
		{
			"Existing list",
			false,
			TestingSender,
			TestingSender2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput()

			if tt.isEmpty {
				require.False(t, k.IsTrustedSchemaProposer(ctx, tt.senderAddress))
			} else {
				k.AddTrustedSchemaProposer(ctx, tt.senderAddress)
				require.True(t, k.IsTrustedSchemaProposer(ctx, tt.senderAddress))
			}

			require.False(t, k.IsTrustedSchemaProposer(ctx, tt.senderAddress2))
		})
	}
}

func TestKeeper_TrustedSchemaProposersIterator(t *testing.T) {
	tests := []struct {
		name            string
		senderAddresses []sdk.AccAddress
	}{
		{
			"Empty list",
			[]sdk.AccAddress{},
		},
		{
			"1 element in list",
			[]sdk.AccAddress{
				TestingSender,
			},
		},
		{
			"2 element in list",
			[]sdk.AccAddress{
				TestingSender,
				TestingSender2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc, ctx, k := SetupTestInput()

			for _, sa := range tt.senderAddresses {
				store := ctx.KVStore(k.StoreKey)

				proposersBz := cdc.MustMarshalBinaryBare(sa)
				store.Set(metadataSchemaProposerKey(sa), proposersBz)
			}

			result := []sdk.AccAddress{}
			tspi := k.TrustedSchemaProposersIterator(ctx)
			defer tspi.Close()
			for ; tspi.Valid(); tspi.Next() {
				ms := sdk.AccAddress{}
				k.cdc.MustUnmarshalBinaryBare(tspi.Value(), &ms)
				result = append(result, ms)
			}

			if tt.senderAddresses == nil {
				require.Empty(t, result)
			}

			require.Equal(t, len(tt.senderAddresses), len(result))

			for _, sa := range tt.senderAddresses {
				require.Contains(t, result, sa)
			}
		})
	}
}

// ----------------------------------
// --- Documents
// ----------------------------------

func TestKeeper_ShareDocument(t *testing.T) {
	var recipient sdk.AccAddress
	recipient, _ = sdk.AccAddressFromBech32("cosmos1h2z8u9294gtqmxlrnlyfueqysng3krh009fum7")

	tests := []struct {
		name           string
		storedDocument types.Document
		document       types.Document
		newRecipient   sdk.AccAddress
	}{
		{
			"No document in store",
			types.Document{},
			TestingDocument,
			nil,
		},
		{
			"One document in store, different recipient",
			TestingDocument,
			TestingDocument,
			recipient,
		},
		{
			"One document in store, different uuid",
			TestingDocument,
			types.Document{
				UUID:       TestingDocument.UUID + "new",
				ContentURI: TestingDocument.ContentURI,
				Metadata:   TestingDocument.Metadata,
				Checksum:   TestingDocument.Checksum,
				Sender:     TestingDocument.Sender,
				Recipients: TestingDocument.Recipients,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc, ctx, k := SetupTestInput()

			store := ctx.KVStore(k.StoreKey)

			if !tt.storedDocument.Equals(types.Document{}) {
				store.Set(getSentDocumentsIdsUUIDStoreKey(TestingSender, tt.storedDocument.UUID), cdc.MustMarshalBinaryBare(tt.storedDocument.UUID))
				store.Set(getReceivedDocumentsIdsUUIDStoreKey(TestingRecipient, tt.storedDocument.UUID), cdc.MustMarshalBinaryBare(tt.storedDocument.UUID))
			}

			if tt.newRecipient != nil {
				tt.document.Recipients = ctypes.Addresses{tt.newRecipient}
			}

			err := k.SaveDocument(ctx, tt.document)
			require.NoError(t, err)

			docsBz := store.Get(getDocumentStoreKey(tt.document.UUID))
			sentDocsBz := store.Get(getSentDocumentsIdsUUIDStoreKey(TestingSender, tt.document.UUID))
			receivedDocsBz := store.Get(getReceivedDocumentsIdsUUIDStoreKey(TestingRecipient, tt.document.UUID))

			if tt.newRecipient != nil {
				newReceivedDocsBz := store.Get(getReceivedDocumentsIdsUUIDStoreKey(tt.newRecipient, tt.document.UUID))

				var newReceivedDocs string
				cdc.MustUnmarshalBinaryBare(newReceivedDocsBz, &newReceivedDocs)
				require.Equal(t, tt.document.UUID, newReceivedDocs)
			}

			var stored types.Document
			cdc.MustUnmarshalBinaryBare(docsBz, &stored)
			require.Equal(t, stored, tt.document)

			var sentDocs, receivedDocs string
			cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
			cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)
			require.Equal(t, tt.document.UUID, sentDocs)
			require.Equal(t, tt.document.UUID, receivedDocs)

		})
	}
}

func TestKeeper_GetDocumentById(t *testing.T) {
	tests := []struct {
		name           string
		storedDocument types.Document
		wantedDoc      string
		wantErr        bool
	}{
		{
			"lookup on non existing document, empty store",
			types.Document{},
			TestingDocument.UUID,
			true,
		},
		{
			"lookup on non existing document, not empty store",
			TestingDocument,
			"",
			true,
		},
		{
			"lookup on existing document",
			TestingDocument,
			TestingDocument.UUID,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc, ctx, k := SetupTestInput()

			if tt.storedDocument.UUID != "" {
				store := ctx.KVStore(k.StoreKey)
				store.Set(getDocumentStoreKey(tt.storedDocument.UUID), cdc.MustMarshalBinaryBare(&tt.storedDocument))
			}

			doc, err := k.GetDocumentByID(ctx, tt.wantedDoc)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.storedDocument, doc)
			}
		})
	}
}

func TestKeeper_UserReceivedDocumentsIterator(t *testing.T) {
	tests := []struct {
		name      string
		recipient []byte
		docs      []types.Document
	}{
		{
			"no document in store",
			nil,
			[]types.Document{},
		},
		{
			"one document in store",
			TestingRecipient,
			[]types.Document{
				TestingDocument,
			},
		},
		{
			"multiple documents in store",
			TestingRecipient,
			[]types.Document{
				TestingDocument,
				{ // TestingDocument with different uuid
					UUID:           "uuid-2",
					Sender:         TestingDocument.Sender,
					Recipients:     TestingDocument.Recipients,
					Metadata:       TestingDocument.Metadata,
					ContentURI:     TestingDocument.ContentURI,
					Checksum:       TestingDocument.Checksum,
					EncryptionData: TestingDocument.EncryptionData,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc, ctx, k := SetupTestInput()

			store := ctx.KVStore(k.StoreKey)
			for _, document := range tt.docs {
				store.Set(getDocumentStoreKey(document.UUID), cdc.MustMarshalBinaryBare(document))
				store.Set(getReceivedDocumentsIdsUUIDStoreKey(tt.recipient, document.UUID), cdc.MustMarshalBinaryBare(document.UUID))
			}

			rdi := k.UserReceivedDocumentsIterator(ctx, tt.recipient)
			defer rdi.Close()

			documents := []types.Document{}
			for ; rdi.Valid(); rdi.Next() {
				id := ""
				k.cdc.MustUnmarshalBinaryBare(rdi.Value(), &id)
				doc, err := k.GetDocumentByID(ctx, id)
				require.NoError(t, err)

				documents = append(documents, doc)
			}

			require.Len(t, documents, len(tt.docs))
			for _, document := range tt.docs {
				require.Contains(t, documents, document)
			}
		})
	}
}

func TestKeeper_UserSentDocumentsIterator(t *testing.T) {
	tests := []struct {
		name   string
		sender []byte
		docs   []types.Document
	}{
		{
			"no document in store",
			TestingSender,
			[]types.Document{},
		},
		{
			"one document in store",
			TestingSender,
			[]types.Document{
				TestingDocument,
			},
		},
		{
			"multiple documents in store",
			TestingSender,
			[]types.Document{
				TestingDocument,
				{ // TestingDocument with different uuid
					UUID:           "uuid-2",
					Sender:         TestingDocument.Sender,
					Recipients:     TestingDocument.Recipients,
					Metadata:       TestingDocument.Metadata,
					ContentURI:     TestingDocument.ContentURI,
					Checksum:       TestingDocument.Checksum,
					EncryptionData: TestingDocument.EncryptionData,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc, ctx, k := SetupTestInput()

			store := ctx.KVStore(k.StoreKey)
			for _, document := range tt.docs {
				store.Set(getDocumentStoreKey(document.UUID), cdc.MustMarshalBinaryBare(document))
				store.Set(getSentDocumentsIdsUUIDStoreKey(tt.sender, document.UUID), cdc.MustMarshalBinaryBare(document.UUID))
			}

			documents := []types.Document{}
			di := k.UserSentDocumentsIterator(ctx, tt.sender)
			defer di.Close()

			for ; di.Valid(); di.Next() {
				id := ""
				k.cdc.MustUnmarshalBinaryBare(di.Value(), &id)
				doc, err := k.GetDocumentByID(ctx, id)
				require.NoError(t, err)

				documents = append(documents, doc)
			}

			require.Len(t, documents, len(tt.docs))
			for _, document := range tt.docs {
				require.Contains(t, documents, document)
			}
		})
	}
}

func TestKeeper_DocumentsIterator(t *testing.T) {
	tests := []struct {
		name string
		docs []types.Document
	}{
		{
			"no document in store",
			[]types.Document{},
		},
		{
			"one document in store",
			[]types.Document{
				TestingDocument,
			},
		},
		{
			"multiple documents in store",
			[]types.Document{
				TestingDocument,
				{ // TestingDocument with different uuid
					UUID:           "uuid-2",
					Sender:         TestingDocument.Sender,
					Recipients:     TestingDocument.Recipients,
					Metadata:       TestingDocument.Metadata,
					ContentURI:     TestingDocument.ContentURI,
					Checksum:       TestingDocument.Checksum,
					EncryptionData: TestingDocument.EncryptionData,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput()

			for _, document := range tt.docs {
				require.NoError(t, k.SaveDocument(ctx, document))
			}

			di := k.DocumentsIterator(ctx)
			defer di.Close()

			documents := []types.Document{}
			for ; di.Valid(); di.Next() {
				d := types.Document{}
				k.cdc.MustUnmarshalBinaryBare(di.Value(), &d)

				documents = append(documents, d)
			}

			require.Len(t, documents, len(tt.docs))
			for _, document := range tt.docs {
				require.Contains(t, documents, document)
			}

		})
	}
}

// ----------------------------------
// --- Document receipts
// ----------------------------------

func TestKeeper_SaveDocumentReceipt(t *testing.T) {
	tests := []struct {
		name       string
		document   types.Document
		receipt    types.DocumentReceipt
		newReceipt types.DocumentReceipt
	}{
		{
			"empty list",
			TestingDocument,
			TestingDocumentReceipt,
			types.DocumentReceipt{},
		},
		{
			"sent receipt already present",
			TestingDocument,
			TestingDocumentReceipt,
			types.DocumentReceipt{
				UUID:         TestingDocumentReceipt.UUID + "-new",
				Sender:       TestingSender,
				Recipient:    TestingDocumentReceipt.Recipient,
				TxHash:       TestingDocumentReceipt.TxHash,
				DocumentUUID: TestingDocument.UUID,
				Proof:        TestingDocumentReceipt.Proof,
			},
		},
		{
			"received receipt already present",
			TestingDocument,
			TestingDocumentReceipt,
			types.DocumentReceipt{
				UUID:         TestingDocumentReceipt.UUID + "-new",
				Sender:       TestingSender2,
				Recipient:    TestingDocumentReceipt.Recipient,
				TxHash:       TestingDocumentReceipt.TxHash,
				DocumentUUID: TestingDocument.UUID,
				Proof:        TestingDocumentReceipt.Proof,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc, ctx, k := SetupTestInput()

			require.NoError(t, k.SaveDocument(ctx, tt.document))

			tdr := tt.receipt
			tdr.DocumentUUID = tt.document.UUID
			require.NoError(t, k.SaveReceipt(ctx, tdr))

			store := ctx.KVStore(k.StoreKey)

			storedID := ""
			docReceiptBz := store.Get(getSentReceiptsIdsUUIDStoreKey(tt.receipt.Sender, tdr.DocumentUUID))
			cdc.MustUnmarshalBinaryBare(docReceiptBz, &storedID)

			stored, err := k.GetReceiptByID(ctx, storedID)
			require.NoError(t, err)

			require.Equal(t, stored, tdr)

			require.Error(t, k.SaveReceipt(ctx, tt.newReceipt))

			var storedSlice []types.DocumentReceipt
			si := k.UserSentReceiptsIterator(ctx, tt.receipt.Sender)

			defer si.Close()
			for ; si.Valid(); si.Next() {
				rid := ""
				k.cdc.MustUnmarshalBinaryBare(si.Value(), &rid)

				newReceipt, err := k.GetReceiptByID(ctx, rid)
				require.NoError(t, err)
				storedSlice = append(storedSlice, newReceipt)
			}

			require.Equal(t, 1, len(storedSlice))
			require.Contains(t, storedSlice, tdr)
			require.NotContains(t, storedSlice, tt.newReceipt)
		})
	}
}

func TestKeeper_SaveDocument(t *testing.T) {

	tests := []struct {
		name     string
		document types.Document
		wantErr  bool
	}{
		{
			"document UUID not specified",
			types.Document{},
			true,
		},
		{
			"document UUID empty",
			types.Document{
				UUID: "",
			},
			true,
		},
		{
			"duplicated document",
			TestingDocument,
			true,
		},
		{
			"UUID already in store",
			types.Document{
				UUID: "test-document-uuid",
			},
			true,
		},
		{
			"document's UUID not in store",
			types.Document{
				UUID:       "test-document-uuid_2",
				ContentURI: "https://example.com/document",
				Metadata: types.DocumentMetadata{
					ContentURI: "https://example.com/document/metadata",
					Schema: &types.DocumentMetadataSchema{
						URI:     "https://example.com/document/metadata/schema",
						Version: "1.0.0",
					},
				},
				Checksum: &types.DocumentChecksum{
					Value:     "93dfcaf3d923ec47edb8580667473987",
					Algorithm: "md5",
				},
				Sender:     TestingSender,
				Recipients: ctypes.Addresses{TestingRecipient},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k := SetupTestInput()

			store := ctx.KVStore(k.StoreKey)
			store.Set(getDocumentStoreKey(TestingDocument.UUID), k.cdc.MustMarshalBinaryBare(&TestingDocument))

			if tt.wantErr {
				require.Error(t, k.SaveDocument(ctx, tt.document))
			} else {
				require.NoError(t, k.SaveDocument(ctx, tt.document))
			}
		})
	}
}

func TestKeeper_SaveReceipt(t *testing.T) {
	tests := []struct {
		name            string
		documentReceipt types.DocumentReceipt
		wantErr         bool
	}{
		{
			"receipt UUID not specified",
			types.DocumentReceipt{
				DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
			},
			true,
		},
		{
			"receipt UUID empty",
			types.DocumentReceipt{
				UUID:         "",
				DocumentUUID: "6a2f41a3-c54c-fce8-32d2-0324e1c32e22",
			},
			true,
		},
		{
			"duplicated receipt",
			TestingDocumentReceipt,
			true,
		},
		{
			"UUID already in store",
			types.DocumentReceipt{
				UUID: "testing-document-receipt-uuid",
			},
			true,
		},
		{
			"receipt UUID not in store",
			types.DocumentReceipt{
				UUID:         TestingDocumentReceipt.UUID + "-new",
				Sender:       TestingDocumentReceipt.Sender,
				Recipient:    TestingDocumentReceipt.Recipient,
				TxHash:       TestingDocumentReceipt.TxHash,
				DocumentUUID: TestingDocument.UUID,
				Proof:        TestingDocumentReceipt.Proof,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc, ctx, k := SetupTestInput()

			store := ctx.KVStore(k.StoreKey)
			store.Set(getDocumentStoreKey(TestingDocument.UUID), k.cdc.MustMarshalBinaryBare(&TestingDocument))
			store.Set(getSentReceiptsIdsUUIDStoreKey(TestingDocumentReceipt.Sender, TestingDocumentReceipt.UUID), cdc.MustMarshalBinaryBare(TestingDocumentReceipt))

			if tt.wantErr {
				require.Error(t, k.SaveReceipt(ctx, tt.documentReceipt))
			} else {
				require.NoError(t, k.SaveReceipt(ctx, tt.documentReceipt))
			}
		})
	}
}

func TestKeeper_UserReceivedReceiptsIterator(t *testing.T) {
	tests := []struct {
		name            string
		documentReceipt []types.DocumentReceipt
	}{
		{
			"Empty list",
			[]types.DocumentReceipt{},
		},
		{
			"Filled list",
			[]types.DocumentReceipt{TestingDocumentReceipt},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cdc, ctx, k := SetupTestInput()

			for _, tdr := range tt.documentReceipt {
				store := ctx.KVStore(k.StoreKey)
				store.Set(getReceivedReceiptsIdsUUIDStoreKey(tdr.Recipient, tdr.UUID),
					cdc.MustMarshalBinaryBare(tdr.UUID))

				store.Set(getReceiptStoreKey(tdr.UUID), cdc.MustMarshalBinaryBare(tdr))
			}

			urri := k.UserReceivedReceiptsIterator(ctx, TestingDocumentReceipt.Recipient)
			defer urri.Close()

			receipts := []types.DocumentReceipt{}
			for ; urri.Valid(); urri.Next() {
				rid := ""
				k.cdc.MustUnmarshalBinaryBare(urri.Value(), &rid)

				r, err := k.GetReceiptByID(ctx, rid)
				require.NoError(t, err)

				receipts = append(receipts, r)
			}

			require.Equal(t, tt.documentReceipt, receipts)

		})
	}
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

			require.NoError(t, k.SaveDocument(ctx, tt.want))

			docKey := []byte{}

			di := k.DocumentsIterator(ctx)
			defer di.Close()
			for ; di.Valid(); di.Next() {
				docKey = di.Key()
			}

			extDoc, extUUID, extErr := k.ExtractDocument(ctx, docKey)

			if !tt.wantErr {
				require.NoError(t, extErr)
				require.Equal(t, tt.want, extDoc)
				require.Equal(t, tt.wantUUID, extUUID)
			} else {
				require.Error(t, extErr)
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

			require.Equal(t, tt.want, m)
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

			require.NoError(t, k.SaveDocument(ctx, tt.savedDocument))
			require.NoError(t, k.SaveReceipt(ctx, tt.want))

			recVal := []byte{}

			di, _ := k.ReceiptsIterators(ctx)
			defer di.Close()
			for ; di.Valid(); di.Next() {
				recVal = di.Value()
			}

			extDoc, extUUID, extErr := k.ExtractReceipt(ctx, recVal)

			if !tt.wantErr {
				require.NoError(t, extErr)
				require.Equal(t, tt.want, extDoc)
				require.Equal(t, tt.wantUUID, extUUID)
			} else {
				require.Error(t, extErr)
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

			require.Equal(t, tt.want, m)
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
				require.NoError(t, k.SaveDocument(ctx, tt.storedDocument))
			}

			if !tt.want.Equals(types.DocumentReceipt{}) {
				require.NoError(t, k.SaveReceipt(ctx, tt.want))
			}

			rr, err := k.GetReceiptByID(ctx, tt.want.UUID)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, rr)
			}
		})
	}
}
