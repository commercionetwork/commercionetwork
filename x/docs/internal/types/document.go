package types

import (
	"github.com/commercionetwork/commercionetwork/x/docs"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DocumentMetadataSchema represents the information about the schema that should be used in order to
// validate the metadata associated with a document.
type DocumentMetadataSchema struct {
	Uri     string `json:"uri"`
	Version string `json:"version"`
}

func (metaSchema DocumentMetadataSchema) Equals(metSchema2 DocumentMetadataSchema) bool {
	return metaSchema.Uri == metSchema2.Uri &&
		metaSchema.Version == metSchema2.Version
}

// DocumentMetadata represents the information about the metadata associated to a document
type DocumentMetadata struct {
	ContentUri string                 `json:"content_uri"`
	Schema     DocumentMetadataSchema `json:"schema"`
	Proof      string                 `json:"proof"`
}

func (docMeta DocumentMetadata) Equals(docMeta2 DocumentMetadata) bool {
	return docMeta.ContentUri == docMeta2.ContentUri &&
		docMeta.Proof == docMeta2.Proof &&
		docMeta.Schema.Equals(docMeta2.Schema)
}

// DocumentChecksum represents the information related to the checksum of a document, if any
type DocumentChecksum struct {
	Value     string `json:"value"`
	Algorithm string `json:"algorithm"`
}

func (checksum DocumentChecksum) Equals(checksum2 DocumentChecksum) bool {
	return checksum.Value == checksum2.Value &&
		checksum.Algorithm == checksum2.Algorithm
}

// Document contains the generic information about a single document which has been sent from a user to another user.
// It contains the information about its content, its associated metadata and the related checksum.
// In order to be valid, a document must have a non-empty and unique UUID and a valid metadata information.
// Both the content and the checksum information are optional.
type Document struct {
	Sender     sdk.AccAddress   `json:"sender"`
	Recipient  sdk.AccAddress   `json:"recipient"`
	Uuid       string           `json:"uuid"`
	Metadata   DocumentMetadata `json:"metadata"`
	ContentUri string           `json:"content_uri"` // Optional
	Checksum   DocumentChecksum `json:"checksum"`    // Optional
}

func (doc Document) Equals(doc2 Document) bool {
	return doc.Sender.Equals(doc2.Sender) &&
		doc.Recipient.Equals(doc2.Recipient) &&
		doc.Uuid == doc2.Uuid &&
		doc.ContentUri == doc2.ContentUri &&
		doc.Metadata.Equals(doc2.Metadata) &&
		doc.Checksum.Equals(doc2.Checksum)
}

type Documents []Document

func (documents Documents) AppendIfMissing(i docs.Document) []docs.Document {
	for _, ele := range documents {
		if ele.Equals(i) {
			return documents
		}
	}
	return append(documents, i)
}
