package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

var algorithms = map[string]int{
	"md5":     32,
	"sha-1":   40,
	"sha-224": 56,
	"sha-256": 64,
	"sha-384": 96,
	"sha-512": 128,
}

// DocumentChecksum represents the information related to the checksum of a document, if any
type DocumentChecksum struct {
	Value     string `json:"value"`
	Algorithm string `json:"algorithm"`
}

// Equals returns true iff this DocumentChecksum and other have the same contents
func (checksum DocumentChecksum) Equals(other DocumentChecksum) bool {
	return checksum.Value == other.Value &&
		checksum.Algorithm == other.Algorithm
}

// Validate returns an error if there is something wrong inside this DocumentChecksum
func (checksum DocumentChecksum) Validate() error {
	if len(strings.TrimSpace(checksum.Value)) == 0 {
		return errors.New("checksum value can't be empty")
	}
	if len(strings.TrimSpace(checksum.Algorithm)) == 0 {
		return errors.New("checksum algorithm can't be empty")
	}

	_, err := hex.DecodeString(checksum.Value)
	if err != nil {
		return errors.New("invalid checksum value (must be hex)")
	}

	algorithm := strings.ToLower(checksum.Algorithm)

	// Check that the algorithm is valid
	length, ok := algorithms[algorithm]
	if !ok {
		return fmt.Errorf("invalid algorithm type %s", algorithm)
	}

	// Check the validity of the checksum value
	if len(strings.TrimSpace(checksum.Value)) != length {
		return fmt.Errorf("invalid checksum length for algorithm %s", algorithm)
	}

	return nil
}
