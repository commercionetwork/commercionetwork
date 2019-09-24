package types

import (
	"encoding/hex"
	"errors"
	"strings"
)

// DidDocument represents the data related to a single Did Document
type DidDocument struct {
	Uri         string `json:"uri"`
	ContentHash string `json:"content_hash"`
}

// Equals returns true iff this didDocument and other contain the same data
func (didDocument DidDocument) Equals(other DidDocument) bool {
	return didDocument.Uri == other.Uri &&
		didDocument.ContentHash == other.Uri
}

// Validate checks the data present inside this Did Document and returns an
// error if something is wrong
func (didDocument DidDocument) Validate() error {
	if len(strings.TrimSpace(didDocument.Uri)) == 0 {
		return errors.New("did document uri cannot be empty")
	}

	if len(strings.TrimSpace(didDocument.ContentHash)) == 0 {
		return errors.New("did document content hash cannot be empty")
	}

	if _, err := hex.DecodeString(didDocument.ContentHash); err != nil {
		return errors.New("did document content hash must be a valid hex string")
	}

	if len(strings.TrimSpace(didDocument.ContentHash)) != 64 {
		return errors.New("did document content hash is not a valid SHA256 hash")
	}

	return nil
}
