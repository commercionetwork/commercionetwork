package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// Equals returns true iff this dat and other contain the same data
func (key DocumentEncryptionKey) Equals(other DocumentEncryptionKey) bool {
	return key.Recipient == other.Recipient &&
		key.Value == other.Value
}

// Validate tries to validate all the data contained inside the given
// DocumentEncryptionKey and returns an error if something bad occurs
func (key DocumentEncryptionKey) Validate() error {
	if key.Recipient == "" {
		return fmt.Errorf("invalid address %s", key.Recipient)
	}

	if len(strings.TrimSpace(key.Value)) == 0 {
		return errors.New("encryption key value cannot be empty")
	}

	if _, err := hex.DecodeString(key.Value); err != nil {
		return errors.New("invalid encryption key value (must be hex)")
	}

	return nil
}
