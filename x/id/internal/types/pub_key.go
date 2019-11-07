package types

import (
	"encoding/hex"
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// -------------
// --- PubKey
// -------------

// PubKey contains the information of a public key contained inside a Did Document
type PubKey struct {
	ID           string         `json:"id"`
	Type         string         `json:"type"`
	Controller   sdk.AccAddress `json:"controller"`
	PublicKeyHex string         `json:"publicKeyHex"`
}

func NewPubKey(pubKeyID string, pubKeyType string, controller sdk.AccAddress, hexValue string) PubKey {
	return PubKey{
		ID:           pubKeyID,
		Type:         pubKeyType,
		Controller:   controller,
		PublicKeyHex: hexValue,
	}
}

// Equals returns true iff pubKey and other contain the same data
func (pubKey PubKey) Equals(other PubKey) bool {
	return pubKey.ID == other.ID &&
		pubKey.Type == other.Type &&
		pubKey.Controller.Equals(other.Controller) &&
		pubKey.PublicKeyHex == other.PublicKeyHex
}

// Validate checks the data contained inside pubKey and returns an error if something is wrong
func (pubKey PubKey) Validate() sdk.Error {

	regex, _ := regexp.Compile(fmt.Sprintf("^%s#keys-[0-9]+$", pubKey.Controller))
	if !regex.MatchString(pubKey.ID) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid key id, must satisfy %s", regex))
	}

	if pubKey.Type != KeyTypeRsa && pubKey.Type != KeyTypeSecp256k1 && pubKey.Type != KeyTypeEd25519 {
		msg := fmt.Sprintf("Invalid key type, must be either %s, %s or %s", KeyTypeRsa, KeyTypeSecp256k1, KeyTypeEd25519)
		return sdk.ErrUnknownRequest(msg)
	}

	if _, err := hex.DecodeString(pubKey.PublicKeyHex); err != nil {
		return sdk.ErrUnknownRequest("Invalid publicKeyHex value")
	}

	return nil
}

// --------------
// --- PubKeys
// --------------

// PubKeys represents a slice of PubKey objects
type PubKeys []PubKey

// Equals returns true iff pubKeys contains the same data as other in the same order
func (pubKeys PubKeys) Equals(other PubKeys) bool {
	if len(pubKeys) != len(other) {
		return false
	}

	for index, key := range pubKeys {
		if !key.Equals(other[index]) {
			return false
		}
	}

	return true
}

// FindByID returns the key having the given id present inside the pubKeys object
// If no key has been found, returns false
func (pubKeys PubKeys) FindByID(id string) (PubKey, bool) {
	for _, key := range pubKeys {
		if key.ID == id {
			return key, true
		}
	}

	return PubKey{}, false
}
