package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

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
	if len(checksum.Value) == 0 {
		return errors.New("checksum value can't be empty")
	}
	if len(checksum.Algorithm) == 0 {
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
		return errors.New(fmt.Sprintf("invalid algorithm type %s", algorithm))
	}

	// Check the validity of the checksum value
	if len(checksum.Value) != length {
		return errors.New(fmt.Sprintf("invalid checksum length for algorithm %s", algorithm))
	}

	return nil
}
