// DONTCOVER
// nolint
package v1_3_0

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	uuid "github.com/satori/go.uuid"

	"github.com/commercionetwork/commercionetwork/x/common/types"
	v120docs "github.com/commercionetwork/commercionetwork/x/documents/legacy/v1.2.0"
)

// Migrate accepts exported genesis state from v1.1.0 and migrates it to v1.2.0
// genesis state. This migration changes the data that are saved for each document
// removing the metadata proof
func Migrate(oldGenState v120docs.GenesisState) GenesisState {

	documents := Documents{}
	var receipts []DocumentReceipt
	for _, userData := range oldGenState.UsersData {

		// Migrate just the sent documents, as the received ones are the same just in a different spot
		for _, sentDoc := range userData.SentDocuments {
			document := migrateDocument(sentDoc)
			document.Sender = userData.User
			document.Recipients = findDocumentRecipients(document, oldGenState.UsersData)
			documents = documents.AppendIfMissingID(document)
		}

		// Migrate only the sent receipts, as the received ones are the same just in a different spot
		for index, sentReceipt := range userData.SentReceipts {
			receipt := migrateReceipt(index, sentReceipt)
			receipts = append(receipts, receipt)
		}
	}

	supportedMetadataSchemes := make([]MetadataSchema, len(oldGenState.SupportedMetadataSchemes))
	for i, schema := range oldGenState.SupportedMetadataSchemes {
		supportedMetadataSchemes[i] = MetadataSchema{
			Type:      schema.Type,
			SchemaURI: schema.SchemaURI,
			Version:   schema.Version,
		}
	}

	return GenesisState{
		Documents:                      documents,
		Receipts:                       receipts,
		SupportedMetadataSchemes:       supportedMetadataSchemes,
		TrustedMetadataSchemaProposers: oldGenState.TrustedMetadataSchemaProposers,
	}
}

// migrateDocuments migrates a single v1.2.x document into a 1.3.0 document
func migrateDocument(doc v120docs.Document) Document {

	// Convert the metadata schemes
	var metadataSchema *DocumentMetadataSchema
	if doc.Metadata.Schema != nil {
		metadataSchema = &DocumentMetadataSchema{
			URI:     doc.Metadata.Schema.URI,
			Version: doc.Metadata.Schema.Version,
		}
	}

	// Convert the encryption data
	var encryptionData *DocumentEncryptionData
	if doc.EncryptionData != nil {

		// Convert encryption keys
		keys := make([]DocumentEncryptionKey, len(doc.EncryptionData.Keys))
		for i, key := range doc.EncryptionData.Keys {
			keys[i] = DocumentEncryptionKey{
				Recipient: key.Recipient,
				Value:     key.Value,
			}
		}

		encryptionData = &DocumentEncryptionData{
			Keys:          nil,
			EncryptedData: doc.EncryptionData.EncryptedData,
		}
	}

	// Return a new document
	return Document{
		UUID: doc.UUID,
		Metadata: DocumentMetadata{
			ContentURI: doc.Metadata.ContentURI,
			SchemaType: doc.Metadata.SchemaType,
			Schema:     metadataSchema,
		},
		ContentURI: doc.ContentURI,
		Checksum: &DocumentChecksum{
			Value:     doc.Checksum.Value,
			Algorithm: doc.Checksum.Algorithm,
		},
		EncryptionData: encryptionData,
	}
}

// findDocumentRecipients returns the list of all the sdk.AccAddress that are the
// recipients of the given document
func findDocumentRecipients(document Document, userData []v120docs.UserDocumentsData) types.Addresses {
	var recipients types.Addresses

	// Iterate over all the users' received documents searching for one with the same uuid
	for _, data := range userData {
		for _, receivedDoc := range data.ReceivedDocuments {
			if receivedDoc.UUID == document.UUID {
				recipients = appendIfMissing(recipients, data.User)
			}
		}
	}

	return recipients
}

// appendIfMissing returns a new sdk.AccAddress list that is made of the addresses list and the given address
// if such address does not exist inside the list. Otherwise the original list is returned
func appendIfMissing(addresses types.Addresses, address sdk.AccAddress) types.Addresses {
	for _, a := range addresses {
		if a.Equals(address) {
			return addresses
		}
	}

	return append(addresses, address)
}

// migrateReceipts migrates a v1.2.0 document receipt into a v1.3.0 document receipt
func migrateReceipt(index int, receipt v120docs.DocumentReceipt) DocumentReceipt {
	ns, _ := uuid.FromString("cfbb5b51-6ac0-43b0-8e09-022236285e31")
	return DocumentReceipt{
		UUID:         uuid.NewV3(ns, strconv.Itoa(index)).String(),
		Sender:       receipt.Sender,
		Recipient:    receipt.Recipient,
		TxHash:       receipt.TxHash,
		DocumentUUID: receipt.DocumentUUID,
		Proof:        receipt.Proof,
	}
}
