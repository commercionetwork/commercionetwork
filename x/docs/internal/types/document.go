package types

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	uuid "github.com/satori/go.uuid"
)

// Document contains the generic information about a single document which has been sent from a user to another user.
// It contains the information about its content, its associated metadata and the related checksum.
// In order to be valid, a document must have a non-empty and unique UUID and a valid metadata information.
// Both the content and the checksum information are optional.
type Document struct {
	Sender         sdk.AccAddress          `json:"sender"`
	Recipients     types.Addresses         `json:"recipients"`
	UUID           string                  `json:"uuid"`
	Metadata       DocumentMetadata        `json:"metadata"`
	ContentURI     string                  `json:"content_uri"`     // Optional
	Checksum       *DocumentChecksum       `json:"checksum"`        // Optional
	EncryptionData *DocumentEncryptionData `json:"encryption_data"` // Optional
	DoSign         *DocumentDoSign         `json:"do_sign"`         // Optional
}

// Equals returns true when doc equals other, false otherwise.
func (doc Document) Equals(other Document) bool {
	validContent := doc.UUID == other.UUID &&
		doc.ContentURI == other.ContentURI &&
		doc.Metadata.Equals(other.Metadata)

	var validChecksum bool
	if doc.Checksum != nil && other.Checksum != nil {
		validChecksum = doc.Checksum.Equals(*other.Checksum)
	} else {
		validChecksum = doc.Checksum == other.Checksum
	}

	var validEncryptionData bool
	if doc.EncryptionData != nil && other.EncryptionData != nil {
		validEncryptionData = doc.EncryptionData.Equals(*other.EncryptionData)
	} else {
		validEncryptionData = doc.EncryptionData == other.EncryptionData
	}

	return validContent && validChecksum && validEncryptionData
}

// validateUUID returns true when uuidStr is a valid UUID, false otherwise.
func validateUUID(uuidStr string) bool {
	_, err := uuid.FromString(uuidStr)

	// when err is nil, uuidStr is a valid UUID
	return err == nil
}

// Validate certify whether doc is a valid Document instance or not.
// It returns an error with the validation failure motivation when the validation process
// fails.
func (doc Document) Validate() error {
	if doc.Sender.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, (doc.Sender.String()))
	}

	if doc.Recipients.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, ("Recipients cannot be empty"))
	}

	for _, recipient := range doc.Recipients {
		if recipient.Empty() {
			return sdkErr.Wrap(sdkErr.ErrInvalidAddress, (recipient.String()))
		}
	}

	if !validateUUID(doc.UUID) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, (fmt.Sprintf("Invalid document UUID: %s", doc.UUID)))
	}

	err := doc.Metadata.Validate()
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, (err.Error()))
	}

	if doc.Checksum != nil {
		err = doc.Checksum.Validate()
		if err != nil {
			return sdkErr.Wrap(sdkErr.ErrUnknownRequest, (err.Error()))
		}
	}

	if doc.EncryptionData != nil {
		err = doc.EncryptionData.Validate()
		if err != nil {
			return sdkErr.Wrap(sdkErr.ErrUnknownRequest, (err.Error()))
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
				return sdkErr.Wrap(sdkErr.ErrInvalidAddress, (errMsg))
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
				return sdkErr.Wrap(sdkErr.ErrInvalidAddress, (errMsg))
			}
		}

		// Check that the `encrypted_data' field name is actually present in doc
		fNotPresent := func(s string) error {
			return sdkErr.Wrap(sdkErr.ErrUnknownRequest,
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
				sdkErr.ErrUnknownRequest,
				"field \"checksum\" not present in document, but marked do_sign",
			)
		}

		if doc.ContentURI == "" {
			return sdkErr.Wrap(
				sdkErr.ErrUnknownRequest,
				"field \"content_uri\" not present in document, but marked do_sign",
			)
		}
	}

	return nil
}
