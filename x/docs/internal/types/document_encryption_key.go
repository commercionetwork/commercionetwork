package types

import (
	"encoding/base64"
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
	Encoding  string         `json:"encoding"`  // Encoding used to write the above value
}

// Equals returns true iff this dat and other contain the same data
func (key DocumentEncryptionKey) Equals(other DocumentEncryptionKey) bool {
	return key.Recipient.Equals(other.Recipient) &&
		key.Value == other.Value &&
		key.Encoding == other.Encoding
}

// Validate tries to validate all the data contained inside the given
// DocumentEncryptionKey and returns an error if something bad occurs
// TODO: Test this
func (key DocumentEncryptionKey) Validate() error {
	if key.Recipient.Empty() {
		return errors.New(fmt.Sprintf("invalid address %s", key.Recipient.String()))
	}

	if len(strings.TrimSpace(key.Value)) == 0 {
		return errors.New("encryption key value cannot be empty")
	}

	if len(strings.TrimSpace(key.Encoding)) == 0 {
		return errors.New("encryption key encoding cannot be empty")
	}

	encoding := strings.ToLower(key.Encoding)
	if encoding != "base64" && encoding != "hex" {
		return errors.New("encryption key encoding method unknown")
	}

	if encoding == "hex" {
		if _, err := hex.DecodeString(key.Value); err != nil {
			return err
		}
	}

	if encoding == "base64" {
		if _, err := base64.StdEncoding.DecodeString(key.Value); err != nil {
			return err
		}
	}

	return nil
}
