package types

import (
	"errors"
	"regexp"
	"strings"
)

// Document contains the generic information about a single document which has been sent from a user to another user.
// It contains the information about its content, its associated metadata and the related checksum.
// In order to be valid, a document must have a non-empty and unique UUID and a valid metadata information.
// Both the content and the checksum information are optional.
type Document struct {
	Uuid           string                  `json:"uuid"`
	Metadata       DocumentMetadata        `json:"metadata"`
	ContentUri     string                  `json:"content_uri"`     // Optional
	Checksum       *DocumentChecksum       `json:"checksum"`        // Optional
	EncryptionData *DocumentEncryptionData `json:"encryption_data"` // Optional
}

// TODO: Test
func (doc Document) Equals(other Document) bool {
	validContent := doc.Uuid == other.Uuid &&
		doc.ContentUri == other.ContentUri &&
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

func validateUuid(uuid string) bool {
	regex := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	return regex.MatchString(uuid)
}

func (doc Document) Validate() error {
	if !validateUuid(doc.Uuid) {
		return errors.New("invalid document UUID")
	}
	if len(strings.TrimSpace(doc.ContentUri)) == 0 {
		return errors.New("document content Uri can't be empty")
	}

	err := doc.Metadata.Validate()
	if err != nil {
		return err
	}

	if doc.Checksum != nil {
		err = (*doc.Checksum).Validate()
		if err != nil {
			return err
		}
	}

	if doc.EncryptionData != nil {
		err = (*doc.EncryptionData).Validate()
		if err != nil {
			return err
		}
	}

	return nil
}

type Documents []Document

func (documents Documents) AppendIfMissing(i Document) []Document {
	for _, ele := range documents {
		if ele.Equals(i) {
			return documents
		}
	}
	return append(documents, i)
}
