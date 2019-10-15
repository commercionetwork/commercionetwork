package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// -------------------
// --- Credential
// -------------------

type Credential struct {
	Timestamp time.Time      `json:"timestamp"`
	User      sdk.AccAddress `json:"user"`
	Verifier  sdk.AccAddress `json:"verifier"`
}

func (c Credential) Equals(other Credential) bool {
	return c.Timestamp.Equal(other.Timestamp) &&
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
