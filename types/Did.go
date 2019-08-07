package types

import (
	"errors"
	"strings"
)

type Did string

const DidPrefix = "did:commercio:"

func (did Did) Validate() error {
	if !strings.HasPrefix(string(did), DidPrefix) {
		return errors.New("invalid did prefix")
	}

	return nil
}
