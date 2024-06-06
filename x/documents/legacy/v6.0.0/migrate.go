package v600

import (
	"time"

	v300 "github.com/commercionetwork/commercionetwork/x/documents/legacy/v3.0.0"
	"github.com/commercionetwork/commercionetwork/x/documents/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateStore(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	oldDocsIterator := sdk.KVStorePrefixIterator(store, []byte(v300.DocumentStorePrefix))
	oldDocReceiptsIterator := sdk.KVStorePrefixIterator(store, []byte(types.ReceiptsStorePrefix))

	defer oldDocsIterator.Close()
	for ; oldDocsIterator.Valid(); oldDocsIterator.Next() {
		var oldDocument v300.Document
		cdc.MustUnmarshal(oldDocsIterator.Value(), &oldDocument)
		
		document := migrateDocument(oldDocument)
		store.Set(oldDocsIterator.Key(), cdc.MustMarshal(document))
	}

	defer oldDocReceiptsIterator.Close()
	for ; oldDocReceiptsIterator.Valid(); oldDocReceiptsIterator.Next() {
		var oldDocumentReceipt v300.DocumentReceipt
		cdc.MustUnmarshal(oldDocReceiptsIterator.Value(), &oldDocumentReceipt)
		
		documentReceipt := migrateDocumentReceipt(oldDocumentReceipt)
		store.Set(oldDocReceiptsIterator.Key(), cdc.MustMarshal(documentReceipt))
	}

	return nil
}

// migrateDocuments migrates a single v3.0.0 document into a 6.0.0 document
func migrateDocument(doc v300.Document) *types.Document {
	var encryptionData *types.DocumentEncryptionData
	if doc.EncryptionData != nil {

		// Convert encryption keys
		keys := make([]*types.DocumentEncryptionKey, len(doc.EncryptionData.Keys))
		for i, key := range doc.EncryptionData.Keys {
			keys[i] = &types.DocumentEncryptionKey{
				Recipient: key.Recipient,
				Value:     key.Value,
			}
		}

		encryptionData = &types.DocumentEncryptionData{
			Keys:          keys,
			EncryptedData: doc.EncryptionData.EncryptedData,
		}
	}

	var metadata *types.DocumentMetadata
	if doc.Metadata != nil {
		metadata = &types.DocumentMetadata{
			ContentURI: doc.Metadata.ContentURI,
		}
		if doc.Metadata.Schema != nil {
			metadata.Schema = &types.DocumentMetadataSchema{
				URI: doc.Metadata.Schema.URI,
				Version: doc.Metadata.Schema.Version,
			}
		}
	}

	var checksum *types.DocumentChecksum
	if doc.Checksum != nil {
		checksum = &types.DocumentChecksum{
			Value: doc.Checksum.Value,
			Algorithm: doc.Checksum.Algorithm,
		}
	}
	
	var doSign *types.DocumentDoSign
	if doc.DoSign != nil {
		doSign = &types.DocumentDoSign{
			StorageURI: doc.DoSign.StorageURI,
			SignerInstance: doc.DoSign.SignerInstance,
			SdnData: doc.DoSign.SdnData,
			VcrID: doc.DoSign.VcrID,
			CertificateProfile: doc.DoSign.CertificateProfile,
		}
	}
	//Set the time to the symbolic unix epoch 01/01/1970
	var timestamp time.Time = time.Unix(0, 0).UTC()

	// Return a new document
	return &types.Document{
		Sender:     doc.Sender,
		Recipients: doc.Recipients,
		UUID:       doc.UUID,
		Metadata: metadata,
		ContentURI:     doc.ContentURI,
		Checksum:      checksum,
		EncryptionData: encryptionData,
		DoSign: doSign,
		Timestamp: &timestamp,
	}
}


func migrateDocumentReceipt(receipt v300.DocumentReceipt) *types.DocumentReceipt {
	//Set the time to the symbolic unix epoch 01/01/1970
	var timestamp time.Time = time.Unix(0, 0).UTC()

	return &types.DocumentReceipt{
		UUID:         receipt.UUID,
		Sender:       receipt.Sender,
		Recipient:    receipt.Recipient,
		TxHash:       receipt.TxHash,
		DocumentUUID: receipt.DocumentUUID,
		Proof:        receipt.Proof,
		Timestamp:	  &timestamp,
	}
}