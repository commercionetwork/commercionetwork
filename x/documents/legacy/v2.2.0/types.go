// DONTCOVER
// nolint
package v2_2_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/common/types"
)

const (
	ModuleName = "docs"
)

// ---------------
// --- Genesis
// ---------------

// v2.2.0 documents genesis state
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
	Sender         sdk.AccAddress          `json:"sender" swaggertype:"string" example:"did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf"`
	Recipients     types.Addresses         `json:"recipients" swaggertype:"array,string" example:"did:com:12p24st9asf394jv04e8sxrl9c384jjqwejv0gf"`
	UUID           string                  `json:"uuid" swaggertype:"string" example:"d0f6c692-506f-4bd7-bdf4-f6693633d1da"`
	Metadata       DocumentMetadata        `json:"metadata"`
	ContentURI     string                  `json:"content_uri,omitempty"`     // Optional
	Checksum       *DocumentChecksum       `json:"checksum,omitempty"`        // Optional
	EncryptionData *DocumentEncryptionData `json:"encryption_data,omitempty"` // Optional
	DoSign         *DocumentDoSign         `json:"do_sign,omitempty"`         // Optional
}

// DocumentMetadata represents the information about the metadata associated to a document
type DocumentMetadata struct {
	ContentURI string                  `json:"content_uri"`
	SchemaType string                  `json:"schema_type,omitempty"` // Optional - Either this or schema must be defined
	Schema     *DocumentMetadataSchema `json:"schema,omitempty"`      // Optional - Either this or schema_type must be defined
}

type DocumentMetadataSchema struct {
	URI     string `json:"uri"`
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

// DocumentDoSign represents the optional DoSign value inside a Document.
type DocumentDoSign struct {
	StorageURI         string  `json:"storage_uri"`
	SignerInstance     string  `json:"signer_instance"`
	SdnData            SdnData `json:"sdn_data"`
	VcrID              string  `json:"vcr_id"`
	CertificateProfile string  `json:"certificate_profile"`
}

// SdnData represents the SdnData value inside a DocumentDoSign struct.
type SdnData []string

// ----------------------
// --- Document receipt
// ---------------------

type DocumentReceipt struct {
	UUID         string         `json:"uuid"`
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
