package types

import (
	"encoding/hex"
	"strings"
)

// ValidateHex returns true if the given value is a valid hexadecimal-encoded string, false otherwise
func ValidateHex(value string) bool {

	depositProof := strings.TrimSpace(value)
	if len(depositProof) == 0 {
		return false
	}

	if _, err := hex.DecodeString(depositProof); err != nil {
		return false
	}

	return true
}
