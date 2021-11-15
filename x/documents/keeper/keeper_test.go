package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func setupKeeper(t testing.TB) (*Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, sdk.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, sdk.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	keeper := NewKeeper(
		codec.NewProtoCodec(registry), storeKey, memStoreKey,
		//nil, nil, nil,
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return keeper, ctx
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
		newRecipient   string
	}{
		{
			"No document in store",
			types.Document{},
			TestingDocument,
			"",
		},
		{
			"One document in store, different recipient",
			TestingDocument,
			TestingDocument,
			recipient.String(),
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
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			store := ctx.KVStore(k.storeKey)

			if !tt.storedDocument.Equals(types.Document{}) {
				store.Set(getSentDocumentsIdsUUIDStoreKey(TestingSender, tt.storedDocument.UUID), []byte(tt.storedDocument.UUID))        //k.cdc.MustMarshalBinaryBare(tt.storedDocument.UUID))
				store.Set(getReceivedDocumentsIdsUUIDStoreKey(TestingRecipient, tt.storedDocument.UUID), []byte(tt.storedDocument.UUID)) // k.cdc.MustMarshalBinaryBare(tt.storedDocument.UUID))
			}

			if tt.newRecipient != "" {
				tt.document.Recipients = append([]string{}, tt.newRecipient)
			}

			err := k.SaveDocument(ctx, tt.document)
			require.NoError(t, err)

			docsBz := store.Get(getDocumentStoreKey(tt.document.UUID))
			sentDocsBz := store.Get(getSentDocumentsIdsUUIDStoreKey(TestingSender, tt.document.UUID))
			receivedDocsBz := store.Get(getReceivedDocumentsIdsUUIDStoreKey(TestingRecipient, tt.document.UUID))

			if tt.newRecipient != "" {
				recipientAddr, _ := sdk.AccAddressFromBech32(tt.newRecipient)
				newReceivedDocsBz := store.Get(getReceivedDocumentsIdsUUIDStoreKey(recipientAddr, tt.document.UUID))

				newReceivedDocs := string(newReceivedDocsBz)
				//k.cdc.MustUnmarshalBinaryBare(newReceivedDocsBz, &newReceivedDocs)
				require.Equal(t, tt.document.UUID, newReceivedDocs)
			}

			var stored types.Document
			k.cdc.MustUnmarshalBinaryBare(docsBz, &stored)
			require.Equal(t, stored, tt.document)

			sentDocs := string(sentDocsBz)
			receivedDocs := string(receivedDocsBz)
			//k.cdc.MustUnmarshalBinaryBare(sentDocsBz, &sentDocs)
			//k.cdc.MustUnmarshalBinaryBare(receivedDocsBz, &receivedDocs)
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
			k, ctx := setupKeeper(t)

			if tt.storedDocument.UUID != "" {
				store := ctx.KVStore(k.storeKey)
				//store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.DocumentKey))
				store.Set(getDocumentStoreKey(tt.storedDocument.UUID), k.cdc.MustMarshalBinaryBare(&tt.storedDocument))
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
			k, ctx := setupKeeper(t)

			store := ctx.KVStore(k.storeKey)
			for _, document := range tt.docs {
				store.Set(getDocumentStoreKey(document.UUID), k.cdc.MustMarshalBinaryBare(&document))
				//store.Set(getReceivedDocumentsIdsUUIDStoreKey(tt.recipient, document.UUID), k.cdc.MustMarshalBinaryBare(document.UUID))
				store.Set(getReceivedDocumentsIdsUUIDStoreKey(tt.recipient, document.UUID), []byte(document.UUID))
			}

			rdi := k.UserReceivedDocumentsIterator(ctx, tt.recipient)
			defer rdi.Close()

			documents := []types.Document{}
			for ; rdi.Valid(); rdi.Next() {
				id := string(rdi.Value())
				//k.cdc.MustUnmarshalBinaryBare(rdi.Value(), &id)
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
			k, ctx := setupKeeper(t)

			store := ctx.KVStore(k.storeKey)
			for _, document := range tt.docs {
				store.Set(getDocumentStoreKey(document.UUID), k.cdc.MustMarshalBinaryBare(&document))
				//store.Set(getSentDocumentsIdsUUIDStoreKey(tt.sender, document.UUID), k.cdc.MustMarshalBinaryBare(document.UUID))
				store.Set(getSentDocumentsIdsUUIDStoreKey(tt.sender, document.UUID), []byte(document.UUID))
			}

			documents := []types.Document{}
			di := k.UserSentDocumentsIterator(ctx, tt.sender)
			defer di.Close()

			for ; di.Valid(); di.Next() {
				id := string(di.Value())
				//k.cdc.MustUnmarshalBinaryBare(di.Value(), &id)
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
			k, ctx := setupKeeper(t)

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
				Sender:       TestingSender.String(),
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
				Sender:       TestingSender2.String(),
				Recipient:    TestingDocumentReceipt.Recipient,
				TxHash:       TestingDocumentReceipt.TxHash,
				DocumentUUID: TestingDocument.UUID,
				Proof:        TestingDocumentReceipt.Proof,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			require.NoError(t, k.SaveDocument(ctx, tt.document))

			tdr := tt.receipt
			tdr.DocumentUUID = tt.document.UUID
			require.NoError(t, k.SaveReceipt(ctx, tdr))

			store := ctx.KVStore(k.storeKey)

			senderAccadrr, _ := sdk.AccAddressFromBech32(tdr.Sender)
			docReceiptBz := store.Get(getSentReceiptsIdsUUIDStoreKey(senderAccadrr, tdr.DocumentUUID))
			storedID := string(docReceiptBz)
			//k.cdc.MustUnmarshalBinaryBare(docReceiptBz, &storedID)

			stored, err := k.GetReceiptByID(ctx, storedID)
			require.NoError(t, err)

			require.Equal(t, stored, tdr)

			require.Error(t, k.SaveReceipt(ctx, tt.newReceipt))

			var storedSlice []types.DocumentReceipt
			senderAccadrr, _ = sdk.AccAddressFromBech32(tt.receipt.Sender)
			si := k.UserSentReceiptsIterator(ctx, senderAccadrr)

			defer si.Close()
			for ; si.Valid(); si.Next() {
				rid := string(si.Value())
				//k.cdc.MustUnmarshalBinaryBare(si.Value(), &rid)

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
				Metadata: &types.DocumentMetadata{
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
				Sender:     TestingSender.String(),
				Recipients: append([]string{}, TestingRecipient.String()),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			store := ctx.KVStore(k.storeKey)
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
			k, ctx := setupKeeper(t)

			store := ctx.KVStore(k.storeKey)
			store.Set(getDocumentStoreKey(TestingDocument.UUID), k.cdc.MustMarshalBinaryBare(&TestingDocument))
			senderAccadrr, _ := sdk.AccAddressFromBech32(TestingDocumentReceipt.Sender)
			store.Set(getSentReceiptsIdsUUIDStoreKey(senderAccadrr, TestingDocumentReceipt.UUID), k.cdc.MustMarshalBinaryBare(&TestingDocumentReceipt))

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
			k, ctx := setupKeeper(t)

			for _, tdr := range tt.documentReceipt {
				store := ctx.KVStore(k.storeKey)
				recipientAccAdrr, _ := sdk.AccAddressFromBech32(tdr.Recipient)
				store.Set(getReceivedReceiptsIdsUUIDStoreKey(recipientAccAdrr, tdr.UUID), []byte(tdr.UUID))

				store.Set(getReceiptStoreKey(tdr.UUID), k.cdc.MustMarshalBinaryBare(&tdr))
			}
			recipientAccAdrr, _ := sdk.AccAddressFromBech32(TestingDocumentReceipt.Recipient)
			urri := k.UserReceivedReceiptsIterator(ctx, recipientAccAdrr)
			defer urri.Close()

			receipts := []types.DocumentReceipt{}
			for ; urri.Valid(); urri.Next() {
				rid := string(urri.Value())
				//k.cdc.MustUnmarshalBinaryBare(urri.Value(), &rid)

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
			k, ctx := setupKeeper(t)

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
			k, ctx := setupKeeper(t)

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
			k, ctx := setupKeeper(t)

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
