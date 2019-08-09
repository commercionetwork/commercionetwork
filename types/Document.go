package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type DocumentMetadataSchema struct {
	Uri     string `json:"uri"`
	Version string `json:"version"`
}

type DocumentMetadata struct {
	ContentUri string                 `json:"content_uri"`
	Schema     DocumentMetadataSchema `json:"schema"`
	Proof      string                 `json:"proof"`
}

type DocumentChecksum struct {
	Value     string `json:"value"`
	Algorithm string `json:"algorithm"`
}

type Document struct {
	Sender     sdk.AccAddress   `json:"sender"`
	Recipient  sdk.AccAddress   `json:"recipient"`
	Uuid       string           `json:"uuid"`
	ContentUri string           `json:"content_uri"`
	Metadata   DocumentMetadata `json:"metadata"`
	Checksum   DocumentChecksum `json:"checksum"`
}

//TODO don't know if its correct to let these functions here, i think that they are like class method so it seems
// correct to me
func (checksum DocumentChecksum) Equals(checksum2 DocumentChecksum) bool {
	if checksum.Value != checksum2.Value {
		return false
	}
	if checksum.Algorithm != checksum2.Algorithm {
		return false
	}
	return true
}

func (metaSchema DocumentMetadataSchema) Equals(metSchema2 DocumentMetadataSchema) bool {
	if metaSchema.Uri != metSchema2.Uri {
		return false
	}
	if metaSchema.Version != metSchema2.Version {
		return false
	}
	return true
}

func (docMeta DocumentMetadata) Equals(docMeta2 DocumentMetadata) bool {
	if docMeta.ContentUri != docMeta2.ContentUri {
		return false
	}
	if docMeta.Proof != docMeta2.Proof {
		return false
	}
	if !docMeta.Schema.Equals(docMeta2.Schema) {
		return false
	}
	return true
}

func (doc Document) Equals(doc2 Document) bool {

	if !doc.Sender.Equals(doc2.Sender) {
		return false
	}
	if !doc.Recipient.Equals(doc2.Recipient) {
		return false
	}
	if doc.Uuid != doc2.Uuid {
		return false
	}
	if doc.ContentUri != doc2.ContentUri {
		return false
	}
	if !doc.Metadata.Equals(doc2.Metadata) {
		return false
	}
	if !doc.Checksum.Equals(doc2.Checksum) {
		return false
	}

	return true
}
