package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// -------------------
// --- Credential
// -------------------

type Credential struct {
	Timestamp int64          `json:"timestamp"` // Block height at which the credential has been inserted
	User      sdk.AccAddress `json:"user"`
	Verifier  sdk.AccAddress `json:"verifier"`
}

func NewCredential(user, verifier sdk.AccAddress, blockHeight int64) Credential {
	return Credential{
		Timestamp: blockHeight,
		User:      user,
		Verifier:  verifier,
	}
}

func (c Credential) Equals(other Credential) bool {
	return c.Timestamp == other.Timestamp &&
		c.User.Equals(other.User) &&
		c.Verifier.Equals(other.Verifier)
}

// -------------------
// --- Credentials
// -------------------

// Credentials represent a slice of Credential objects
type Credentials []Credential

// Contains returns true of the given credentials is contained inside the credentials slice
func (credentials Credentials) Contains(credential Credential) bool {
	for _, c := range credentials {
		if c.Equals(credential) {
			return true
		}
	}
	return false
}

// AppendIfMissing returns a new Credentials object containing the given credential.
// If the credential has been appended because previously missing, returns true. Otherwise returns false.
func (credentials Credentials) AppendIfMissing(credential Credential) (Credentials, bool) {
	if credentials.Contains(credential) {
		return credentials, false
	}
	return append(credentials, credential), true
}
