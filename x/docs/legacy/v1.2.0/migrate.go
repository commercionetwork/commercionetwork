package v1_2_0

import (
	v110docs "github.com/commercionetwork/commercionetwork/x/docs/legacy/v1.1.0"
)

// Migrate accepts exported genesis state from v1.1.0 and migrates it to v1.2.0
// genesis state. This migration changes the data that are saved for each document
// removing the metadata proof
func Migrate(oldGenState v110docs.GenesisState) GenesisState {

	usersData := make([]UserDocumentsData, len(oldGenState.UsersData))
	for i, userData := range oldGenState.UsersData {
		usersData[i] = UserDocumentsData{
			User:              userData.User,
			SentDocuments:     migrateDocuments(userData.SentDocuments),
			ReceivedDocuments: migrateDocuments(userData.ReceivedDocuments),
			SentReceipts:      migrateReceipts(userData.SentReceipts),
			ReceivedReceipts:  migrateReceipts(userData.ReceivedReceipts),
		}
	}

	supportedMetadataSchemes := make([]MetadataSchema, len(oldGenState.SupportedMetadataSchemes))
	for i, schema := range oldGenState.SupportedMetadataSchemes {
		supportedMetadataSchemes[i] = MetadataSchema{
			Type:      schema.Type,
			SchemaUri: schema.SchemaUri,
			Version:   schema.Version,
		}
	}

	return GenesisState{
		UsersData:                      usersData,
		SupportedMetadataSchemes:       supportedMetadataSchemes,
		TrustedMetadataSchemaProposers: oldGenState.TrustedMetadataSchemaProposers,
	}
}

// migrateDocuments migrates a list of v1.1.0 documents into a list of v1.2.0 documents
func migrateDocuments(oldDocs []v110docs.Document) []Document {
	documents := make([]Document, len(oldDocs))
	for i, doc := range oldDocs {

		var metadataSchema *DocumentMetadataSchema
		if doc.Metadata.Schema != nil {
			metadataSchema = &DocumentMetadataSchema{
				Uri:     doc.Metadata.Schema.Uri,
				Version: doc.Metadata.Schema.Version,
			}
		}

		documents[i] = Document{
			Uuid: doc.Uuid,
			Metadata: DocumentMetadata{
				ContentUri: doc.Metadata.ContentUri,
				SchemaType: doc.Metadata.SchemaType,
				Schema:     metadataSchema,
			},
			ContentUri: doc.ContentUri,
			Checksum: &DocumentChecksum{
				Value:     doc.Checksum.Value,
				Algorithm: doc.Checksum.Algorithm,
			},
			EncryptionData: nil,
		}
	}

	return documents
}

// migrateReceipts migrates a list of v1.1.0 document receipts into a list of v1.2.0 document receipts
func migrateReceipts(oldReceipts []v110docs.DocumentReceipt) []DocumentReceipt {
	documentReceipts := make([]DocumentReceipt, len(oldReceipts))
	for i, receipt := range oldReceipts {
		documentReceipts[i] = DocumentReceipt{
			Sender:       receipt.Sender,
			Recipient:    receipt.Recipient,
			TxHash:       receipt.TxHash,
			DocumentUuid: receipt.DocumentUuid,
			Proof:        receipt.Proof,
		}
	}

	return documentReceipts
}
