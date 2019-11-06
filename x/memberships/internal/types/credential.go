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

type Credentials []Credential

func (credentials Credentials) Contains(credential Credential) bool {
	for _, c := range credentials {
		if c.Equals(credential) {
			return true
		}
	}
	return false
}

func (credentials Credentials) AppendIfMissing(credential Credential) (newList Credentials, edited bool) {
	if credentials.Contains(credential) {
		return credentials, false
	}
	return append(credentials, credential), true
}
