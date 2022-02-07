package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	uuid "github.com/satori/go.uuid"
)

// Equals returns true when doc equals other, false otherwise.
func (doc Document) Equals(other Document) bool {
	metadata := false
	if doc.Metadata == nil && other.Metadata == nil {
		metadata = true
	} else if (doc.Metadata == nil && other.Metadata != nil) || (doc.Metadata != nil && other.Metadata == nil) {
		metadata = false
	} else {
		metadata = doc.Metadata.Equals(*other.Metadata)
	}

	validContent := doc.UUID == other.UUID &&
		doc.ContentURI == other.ContentURI && metadata

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

// Validate certify whether doc is a valid Document instance or not.
// It returns an error with the validation failure motivation when the validation process
// fails.
func (doc Document) Validate() error {
	if doc.Sender == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, doc.Sender)
	}

	if len(doc.Recipients) == 0 {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Recipients cannot be empty")
	}

	for _, recipient := range doc.Recipients {
		_, err := sdk.AccAddressFromBech32(recipient)
		if err != nil {
			return sdkErr.Wrap(sdkErr.ErrInvalidAddress, recipient)
		}
	}

	if !validateUUID(doc.UUID) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid document UUID: %s", doc.UUID))
	}

	err := doc.Metadata.Validate()
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
	}

	if doc.Checksum != nil {
		err = doc.Checksum.Validate()
		if err != nil {
			return sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
		}
	}

	if doc.EncryptionData != nil {
		err = doc.EncryptionData.Validate()
		if err != nil {
			return sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
		}

		// check that each document recipient have some encrypted data
		for _, recipient := range doc.Recipients {
			recipientAccAddr, _ := sdk.AccAddressFromBech32(recipient)
			if !doc.EncryptionData.ContainsRecipient(recipientAccAddr) {
				errMsg := fmt.Sprintf(
					"%s is a recipient inside the document but not in the encryption data",
					recipient,
				)
				return sdkErr.Wrap(sdkErr.ErrInvalidAddress, errMsg)
			}
		}
		/*
			// check that there are no spurious encryption data recipients not present
			// in the document recipient list
			for _, encAdd := range doc.EncryptionData.Keys {
				if !doc.Recipients.Contains(encAdd.Recipient) {
					errMsg := fmt.Sprintf(
						"%s is a recipient inside encryption data but not inside the message",
						encAdd.Recipient,
					)
					return sdkErr.Wrap(sdkErr.ErrInvalidAddress, errMsg)
				}
			}*/

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
				// case "metadata.schema.uri":
				// 	// already checked in validate of metadata
				// 	if doc.Metadata.Schema == nil || doc.Metadata.Schema.URI == "" {
				// 		return fNotPresent("metadata.schema.uri")
				// 	}
			}
		}

	}

	if doc.DoSign != nil {
		if doc.Checksum == nil {
			return sdkErr.Wrap(
				sdkErr.ErrUnknownRequest,
				"field \"checksum\" not present in document, but required when using do_sign",
			)
		}

		if doc.ContentURI == "" {
			return sdkErr.Wrap(
				sdkErr.ErrUnknownRequest,
				"field \"content_uri\" not present in document, but required when using do_sign",
			)
		}
		/*
			err := *(doc.DoSign).SdnData.Validate()
			if err != nil {
				return sdkErr.Wrap(
					sdkErr.ErrUnknownRequest,
					err.Error(),
				)
			}*/
	}

	if err := doc.lengthLimits(); err != nil {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest,
			err.Error(),
		)
	}

	return nil
}
