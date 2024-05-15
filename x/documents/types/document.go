package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	errors "cosmossdk.io/errors"
	uuid "github.com/satori/go.uuid"
)

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

	if doc.Metadata != nil {
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
	if _, err := sdk.AccAddressFromBech32(doc.Sender); err != nil {
		return errors.Wrap(sdkErr.ErrInvalidAddress, doc.Sender)
	}

	if len(doc.Recipients) == 0 {
		return errors.Wrap(sdkErr.ErrInvalidAddress, "Recipients cannot be empty")
	}

	for _, recipient := range doc.Recipients {
		if _, err := sdk.AccAddressFromBech32(recipient); err != nil {
			return errors.Wrap(sdkErr.ErrInvalidAddress, recipient)
		}
	}

	if !validateUUID(doc.UUID) {
		return errors.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("Invalid document UUID: %s", doc.UUID))
	}

	if err := doc.Metadata.Validate(); err != nil {
		return errors.Wrap(sdkErr.ErrInvalidRequest, err.Error())
	}

	if doc.Checksum != nil {
		if err := doc.Checksum.Validate(); err != nil {
			return errors.Wrap(sdkErr.ErrInvalidRequest, err.Error())
		}
	}

	if doc.EncryptionData != nil {
		if err := doc.EncryptionData.Validate(); err != nil {
			return errors.Wrap(sdkErr.ErrInvalidRequest, err.Error())
		}

		// check that each document recipient have some encrypted data
		for _, recipient := range doc.Recipients {
			recipientAccAddr, _ := sdk.AccAddressFromBech32(recipient)
			if !doc.EncryptionData.ContainsRecipient(recipientAccAddr) {
				errMsg := fmt.Sprintf(
					"%s is a recipient inside the document but not in the encryption data",
					recipient,
				)
				return errors.Wrap(sdkErr.ErrInvalidAddress, errMsg)
			}
		}

		// check that there are no spurious encryption data recipients not present
		// in the document recipient list
		recipients := make(map[string]struct{})
		for _, recipient := range doc.Recipients {
			recipients[recipient] = struct{}{}
		}
		for _, encAdd := range doc.EncryptionData.Keys {
			if _, found := recipients[encAdd.Recipient]; !found {
				return errors.Wrap(sdkErr.ErrInvalidRequest, fmt.Sprintf("the recipient %s is inside encryption data but not along the recipients", encAdd.Recipient))
			}
		}

		for _, fieldName := range doc.EncryptionData.EncryptedData {
			switch fieldName {
			case "content_uri":
				if doc.ContentURI == "" {
					return errors.Wrap(sdkErr.ErrInvalidRequest, "field ContentUri marked as encrypted but not present in document")
				}
			// case "metadata.content_uri":
			// 	if doc.Metadata == nil || doc.Metadata.ContentURI == "" {
			// 		return errors.Wrap(sdkErr.ErrInvalidRequest, "field Metadata.ContentURI marked as encrypted but not present in document")
			// 	}
			case "metadata.schema.uri":
				if doc.Metadata == nil || doc.Metadata.Schema == nil || doc.Metadata.Schema.URI == "" {
					return errors.Wrap(sdkErr.ErrInvalidRequest, "field Metadata.Schema.URI marked as encrypted but not present in document")
				}
			}
		}

	}

	if doc.DoSign != nil {
		if doc.Checksum == nil {
			return errors.Wrap(
				sdkErr.ErrInvalidRequest,
				"field \"checksum\" not present in document, but required when using do_sign",
			)
		}

		if doc.ContentURI == "" {
			return errors.Wrap(
				sdkErr.ErrInvalidRequest,
				"field \"content_uri\" not present in document, but required when using do_sign",
			)
		}

		if err := SdnData(doc.DoSign.SdnData).Validate(); err != nil {
			return errors.Wrap(sdkErr.ErrInvalidRequest,
				err.Error(),
			)
		}
	}

	if err := doc.lengthLimits(); err != nil {
		return errors.Wrap(sdkErr.ErrInvalidRequest,
			err.Error(),
		)
	}

	return nil
}
