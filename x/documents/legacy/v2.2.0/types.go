// DONTCOVER
// nolint
package v2_2_0

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	uuid "github.com/satori/go.uuid"

	"encoding/hex"
	"errors"

	"strings"

	"github.com/commercionetwork/commercionetwork/x/common/types"
)

const (
	ModuleName = "docs"

	SdnDataCommonName   = "common_name"
	SdnDataSurname      = "surname"
	SdnDataSerialNumber = "serial_number"
	SdnDataGivenName    = "given_name"
	SdnDataOrganization = "organization"
	SdnDataCountry      = "country"

	InputStringSep = ","
)

var algorithms = map[string]int{
	"md5":     32,
	"sha-1":   40,
	"sha-224": 56,
	"sha-256": 64,
	"sha-384": 96,
	"sha-512": 128,
}

var validSdnData = map[string]struct{}{
	SdnDataCommonName:   {},
	SdnDataSurname:      {},
	SdnDataSerialNumber: {},
	SdnDataGivenName:    {},
	SdnDataOrganization: {},
	SdnDataCountry:      {},
}

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

// Validate checks that the SdnData is valid, only accepts value included in
// validSdnData.
func (s SdnData) Validate() error {
	for _, val := range s {
		if _, ok := validSdnData[val]; !ok {
			return fmt.Errorf("sdn_data value \"%s\" is not supported", val)
		}
	}

	return nil
}

// NewSdnDataFromString generates a SdnData struct based on the input string.
// The input string expects a comma-separated value as:
// "common_name,surname,serial_number"
// If empty string is provided, a SdnData with default value will be provided. Default : "serial_number".
func NewSdnDataFromString(input string) (SdnData, error) {
	if input == "" {
		return SdnData{SdnDataSerialNumber}, nil
	}

	var split SdnData = strings.Split(input, InputStringSep)
	err := split.Validate()
	if err != nil {
		return SdnData{}, err
	}

	return split, nil
}

// Equals returns true iff this dat and other contain the same data
func (key DocumentEncryptionKey) Equals(other DocumentEncryptionKey) bool {
	return key.Recipient.Equals(other.Recipient) &&
		key.Value == other.Value
}

// Validate tries to validate all the data contained inside the given
// DocumentEncryptionKey and returns an error if something bad occurs
func (key DocumentEncryptionKey) Validate() error {
	if key.Recipient.Empty() {
		return fmt.Errorf("invalid address %s", key.Recipient.String())
	}

	if len(strings.TrimSpace(key.Value)) == 0 {
		return errors.New("encryption key value cannot be empty")
	}

	if _, err := hex.DecodeString(key.Value); err != nil {
		return errors.New("invalid encryption key value (must be hex)")
	}

	return nil
}

// Equals returns true iff this dat and other contain the same data
func (data DocumentEncryptionData) Equals(other DocumentEncryptionData) bool {
	if len(data.Keys) != len(other.Keys) {
		return false
	}

	for index := range data.Keys {
		if !data.Keys[index].Equals(other.Keys[index]) {
			return false
		}
	}

	if len(data.EncryptedData) != len(other.EncryptedData) {
		return false
	}

	for index := range data.EncryptedData {
		if data.EncryptedData[index] != other.EncryptedData[index] {
			return false
		}
	}

	return true
}

// Validate tries to validate all the data contained inside the given
// DocumentEncryptionData and returns an error if something is wrong
func (data DocumentEncryptionData) Validate() error {

	if len(data.Keys) == 0 {
		return errors.New("encryption data keys cannot be empty")
	}

	// Validate the keys
	for _, key := range data.Keys {
		if err := key.Validate(); err != nil {
			return err
		}
	}

	// Validate the encrypted data
	for _, data := range data.EncryptedData {
		if data != "content" && data != "content_uri" && data != "metadata.content_uri" && data != "metadata.schema.uri" {
			return errors.New("encrypted data not supported")
		}
	}

	return nil
}

// ContainsRecipient returns true iff data contains a key with recipient inside.
func (data DocumentEncryptionData) ContainsRecipient(recipient sdk.AccAddress) bool {
	for _, r := range data.Keys {
		if r.Recipient.Equals(recipient) {
			return true
		}
	}

	return false
}

func (metaSchema DocumentMetadataSchema) Equals(metSchema2 DocumentMetadataSchema) bool {
	return metaSchema.URI == metSchema2.URI &&
		metaSchema.Version == metSchema2.Version
}

// Equals returns true iff this metadata and other contain the same data
func (metadata DocumentMetadata) Equals(other DocumentMetadata) bool {
	if metadata.ContentURI != other.ContentURI {
		return false
	}

	if metadata.Schema != nil && other.Schema != nil {
		return metadata.Schema.Equals(*other.Schema)
	}

	return metadata.Schema == other.Schema
}

// Validate tries to validate all the data contained inside the given
// DocumentMetadata and returns an error if something is wrong
func (metadata DocumentMetadata) Validate() error {
	if len(strings.TrimSpace(metadata.ContentURI)) == 0 {
		return errors.New("metadata.content_uri can't be empty")
	}

	if (metadata.Schema == nil) && len(strings.TrimSpace(metadata.SchemaType)) == 0 {
		return errors.New("either metadata.schema or metadata.schema_type must be defined")
	}

	if metadata.Schema != nil {
		if len(strings.TrimSpace(metadata.Schema.URI)) == 0 {
			return errors.New("metadata.schema.uri can't be empty")
		}
		if len(strings.TrimSpace(metadata.Schema.Version)) == 0 {
			return errors.New("metadata.schema.version can't be empty")
		}
	}
	return nil
}

// Equals returns true iff this DocumentChecksum and other have the same contents
func (checksum DocumentChecksum) Equals(other DocumentChecksum) bool {
	return checksum.Value == other.Value &&
		checksum.Algorithm == other.Algorithm
}

// Validate returns an error if there is something wrong inside this DocumentChecksum
func (checksum DocumentChecksum) Validate() error {
	if len(strings.TrimSpace(checksum.Value)) == 0 {
		return errors.New("checksum value can't be empty")
	}
	if len(strings.TrimSpace(checksum.Algorithm)) == 0 {
		return errors.New("checksum algorithm can't be empty")
	}

	_, err := hex.DecodeString(checksum.Value)
	if err != nil {
		return errors.New("invalid checksum value (must be hex)")
	}

	algorithm := strings.ToLower(checksum.Algorithm)

	// Check that the algorithm is valid
	length, ok := algorithms[algorithm]
	if !ok {
		return fmt.Errorf("invalid checksum algorithm type %s", algorithm)
	}

	// Check the validity of the checksum value
	if len(strings.TrimSpace(checksum.Value)) != length {
		return fmt.Errorf("invalid checksum length for algorithm %s", algorithm)
	}

	return nil
}

// validateUUID returns true when uuidStr is a valid UUID, false otherwise.
func validateUUID(uuidStr string) bool {
	_, err := uuid.FromString(uuidStr)

	// when err is nil, uuidStr is a valid UUID
	return err == nil
}

func (doc Document) lengthLimits() error {
	e := func(fieldName string, maxLen int) error {
		return fmt.Errorf("%s content can't be longer than %d bytes", fieldName, maxLen)
	}

	if len(doc.ContentURI) > 512 {
		return e("content_uri", 512)
	}

	if len(doc.Metadata.ContentURI) > 512 {
		return e("metadata.content_uri", 512)
	}

	if s := doc.Metadata.Schema; s != nil {
		if len(s.URI) > 512 {
			return e("metadata.schema.uri", 512)
		}
		if len(s.Version) > 32 {
			return e("metadata.schema.version", 32)
		}
	}

	if len(doc.Metadata.SchemaType) > 512 {
		return e("metadata.schema_type", 512)
	}

	if doc.EncryptionData != nil {
		for i, key := range doc.EncryptionData.Keys {
			if len(key.Value) > 512 {
				return e(fmt.Sprintf("encryption key #%d", i), 512)
			}
		}
	}

	if ds := doc.DoSign; ds != nil {
		if len(ds.VcrID) > 64 {
			return e("do_sign.vcr_id", 64)
		}

		if len(ds.CertificateProfile) > 32 {
			return e("do_sign.certificate_profile", 32)
		}
	}

	return nil
}

func (doc Document) Validate() error {
	if doc.Sender.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, doc.Sender.String())
	}

	if doc.Recipients.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Recipients cannot be empty")
	}

	for _, recipient := range doc.Recipients {
		if recipient.Empty() {
			return sdkErr.Wrap(sdkErr.ErrInvalidAddress, recipient.String())
		}
	}

	if !validateUUID(doc.UUID) {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("Invalid document UUID: %s", doc.UUID))
	}

	err := doc.Metadata.Validate()
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}

	if doc.Checksum != nil {
		err = doc.Checksum.Validate()
		if err != nil {
			return sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
		}
	} else {
		errMsg := fmt.Sprintf(
			"checksum can't be empty",
		)
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, errMsg)
	}

	if doc.EncryptionData != nil {
		err = doc.EncryptionData.Validate()
		if err != nil {
			return sdkErr.Wrap(sdkErr.ErrInvalidRequest, err.Error())
		}
	}

	if doc.EncryptionData != nil {

		// check that each document recipient have some encrypted data
		for _, recipient := range doc.Recipients {
			if !doc.EncryptionData.ContainsRecipient(recipient) {
				errMsg := fmt.Sprintf(
					"%s is a recipient inside the document but not in the encryption data",
					recipient.String(),
				)
				return sdkErr.Wrap(sdkErr.ErrInvalidAddress, errMsg)
			}
		}

		// check that there are no spurious encryption data recipients not present
		// in the document recipient list
		for _, encAdd := range doc.EncryptionData.Keys {
			if !doc.Recipients.Contains(encAdd.Recipient) {
				errMsg := fmt.Sprintf(
					"%s is a recipient inside encryption data but not inside the message",
					encAdd.Recipient.String(),
				)
				return sdkErr.Wrap(sdkErr.ErrInvalidAddress, errMsg)
			}
		}

		// Check that the `encrypted_data' field name is actually present in doc
		fNotPresent := func(s string) error {
			return sdkErr.Wrap(sdkErr.ErrInvalidRequest,
				fmt.Sprintf("field \"%s\" not present in document, but marked as encrypted", s),
			)
		}

		for _, fieldName := range doc.EncryptionData.EncryptedData {
			switch fieldName {
			case "content_uri":
				if doc.ContentURI == "" {
					return fNotPresent("content_uri")
				}
			case "metadata.schema.uri":
				if doc.Metadata.Schema == nil || doc.Metadata.Schema.URI == "" {
					return fNotPresent("metadata.schema.uri")
				}
			}
		}

	}

	if doc.DoSign != nil {
		if doc.Checksum == nil {
			return sdkErr.Wrap(
				sdkErr.ErrInvalidRequest,
				"field \"checksum\" not present in document, but required when using do_sign",
			)
		}

		if doc.ContentURI == "" {
			return sdkErr.Wrap(
				sdkErr.ErrInvalidRequest,
				"field \"content_uri\" not present in document, but required when using do_sign",
			)
		}

		err := doc.DoSign.SdnData.Validate()
		if err != nil {
			return sdkErr.Wrap(
				sdkErr.ErrInvalidRequest,
				err.Error(),
			)
		}
	}

	if err := doc.lengthLimits(); err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest,
			err.Error(),
		)
	}

	return nil
}
