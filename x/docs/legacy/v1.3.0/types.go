// DONTCOVER
// nolint
package v1_3_0

import (
	"github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "docs"
)

// ---------------
// --- Genesis
// ---------------

// v1.3.0 docs genesis state
type GenesisState struct {
	Documents                      []Document        `json:"documents"`
	Receipts                       []DocumentReceipt `json:"receipts"`
	SupportedMetadataSchemes       []MetadataSchema  `json:"supported_metadata_schemes"`
	TrustedMetadataSchemaProposers []sdk.AccAddress  `json:"trusted_metadata_schema_proposers"`
}

// -----------------
// --- Document
// -----------------

type Document struct {
	Uuid           string                  `json:"uuid"`
	Metadata       DocumentMetadata        `json:"metadata"`
	ContentUri     string                  `json:"content_uri"`     // Optional
	Checksum       *DocumentChecksum       `json:"checksum"`        // Optional
	EncryptionData *DocumentEncryptionData `json:"encryption_data"` // Optional
	Sender         sdk.AccAddress          `json:"sender"`
	Recipients     types.Addresses         `json:"recipients"`
}

type DocumentMetadata struct {
	ContentUri string                  `json:"content_uri"`
	SchemaType string                  `json:"schema_type"` // Optional - Either this or schema must be defined
	Schema     *DocumentMetadataSchema `json:"schema"`      // Optional - Either this or schema_type must be defined
}

type DocumentMetadataSchema struct {
	Uri     string `json:"uri"`
	Version string `json:"version"`
}

type DocumentChecksum struct {
	Value     string `json:"value"`
	Algorithm string `json:"algorithm"`
}

type DocumentEncryptionData struct {
	Keys          []DocumentEncryptionKey `json:"keys"`           // contains the keys used to encrypt the data
	EncryptedData []string                `json:"encrypted_data"` // contains the list of data that have been encrypted
}

type DocumentEncryptionKey struct {
	Recipient sdk.AccAddress `json:"recipient"` // Recipient that should use this data
	Value     string         `json:"value"`     // Value of the key that should be used. This is encrypted with the recipient's public key
}

// ----------------------
// --- Document receipt
// ---------------------

type DocumentReceipt struct {
	Uuid         string         `json:"uuid"`
	Sender       sdk.AccAddress `json:"sender"`
	Recipient    sdk.AccAddress `json:"recipient"`
	TxHash       string         `json:"tx_hash"`
	DocumentUuid string         `json:"document_uuid"`
	Proof        string         `json:"proof"`
}

// ------------------------
// --- Metadata schemes
// -------------------------

type MetadataSchema struct {
	Type      string `json:"type"`
	SchemaUri string `json:"schema_uri"`
	Version   string `json:"version"`
}
