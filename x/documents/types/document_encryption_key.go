package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DocumentEncryptionKey contains the data related to a specific encryption key
// used to encrypt some document's data when sending it to a specific recipient
type DocumentEncryptionKey struct {
	Recipient sdk.AccAddress `json:"recipient"` // Recipient that should use this data
	Value     string         `json:"value"`     // Value of the key that should be used. This is encrypted with the recipient's public key
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
