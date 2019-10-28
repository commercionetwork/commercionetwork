// DONTCOVER
// nolint
package v1_1_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "docs"
)

// v1.1.0 docs genesis state
type GenesisState struct {
	UsersData                      []UserData       `json:"users_data"`
	SupportedMetadataSchemes       []MetadataSchema `json:"supported_metadata_schemes"`
	TrustedMetadataSchemaProposers []sdk.AccAddress `json:"trusted_metadata_schema_proposers"`
}

type UserData struct {
	User              sdk.AccAddress    `json:"user"`
	SentDocuments     []Document        `json:"sent_documents"`
	ReceivedDocuments []Document        `json:"received_documents"`
	SentReceipts      []DocumentReceipt `json:"sent_receipts"`
	ReceivedReceipts  []DocumentReceipt `json:"received_receipts"`
}

// -----------------
// --- Document
// -----------------

type Document struct {
	UUID       string           `json:"uuid"`
	Metadata   DocumentMetadata `json:"metadata"`
	ContentURI string           `json:"content_uri"`
	Checksum   DocumentChecksum `json:"checksum"`
	Sender     sdk.AccAddress   `json:"sender"`
	Recipient  sdk.AccAddress   `json:"recipient"`
}

type DocumentMetadata struct {
	ContentURI string                  `json:"content_uri"`
	SchemaType string                  `json:"schema_type"` // Optional - Either this or schema must be defined
	Schema     *DocumentMetadataSchema `json:"schema"`      // Optional - Either this or schema_type must be defined
}

type DocumentMetadataSchema struct {
	URI     string `json:"uri"`
	Version string `json:"version"`
}

type DocumentChecksum struct {
	Value     string `json:"value"`
	Algorithm string `json:"algorithm"`
}

// -----------------------
// --- Document receipt
// -----------------------

type DocumentReceipt struct {
	Sender       sdk.AccAddress `json:"sender"`
	Recipient    sdk.AccAddress `json:"recipient"`
	TxHash       string         `json:"tx_hash"`
	DocumentUUID string         `json:"document_uuid"`
	Proof        string         `json:"proof"`
}

// ------------------------
// --- Metadata schemes
// -------------------------

type MetadataSchema struct {
	Type      string `json:"type"`
	SchemaURI string `json:"schema_uri"`
	Version   string `json:"version"`
}
