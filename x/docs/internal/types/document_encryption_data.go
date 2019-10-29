package types

import (
	"errors"
)

// DocumentEncryptionData contains the data that are related to the way
// that a document's contents or other data have been encrypted
type DocumentEncryptionData struct {
	Keys          []DocumentEncryptionKey `json:"keys"`           // contains the keys used to encrypt the data
	EncryptedData []string                `json:"encrypted_data"` // contains the list of data that have been encrypted
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
