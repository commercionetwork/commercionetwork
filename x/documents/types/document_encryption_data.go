package types

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Equals returns true iff this and other contain the same data
func (data DocumentEncryptionData) Equals(other DocumentEncryptionData) bool {
	if len(data.Keys) != len(other.Keys) {
		return false
	}

	for index := range data.Keys {
		if !data.Keys[index].Equals(*other.Keys[index]) {
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
	for _, eData := range data.EncryptedData {
		if eData != "content" && eData != "content_uri" && eData != "metadata.content_uri" && eData != "metadata.schema.uri" {
			return errors.New("encrypted data not supported")
		}
	}

	return nil
}

// ContainsRecipient returns true iff data contains a key with recipient inside.
func (data DocumentEncryptionData) ContainsRecipient(recipient sdk.AccAddress) bool {
	for _, r := range data.Keys {
		if r.Recipient == recipient.String() {
			return true
		}
	}

	return false
}
