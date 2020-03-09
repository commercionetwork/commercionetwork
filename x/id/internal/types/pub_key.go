package types

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

// -------------
// --- PubKey
// -------------

// PubKey contains the information of a public key contained inside a Did Document
type PubKey struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Controller sdk.AccAddress `json:"controller"`
	PublicKey  string         `json:"publicKey"`
}

func NewPubKey(pubKeyID string, pubKeyType string, controller sdk.AccAddress, hexValue string) PubKey {
	return PubKey{
		ID:         pubKeyID,
		Type:       pubKeyType,
		Controller: controller,
		PublicKey:  hexValue,
	}
}

// Equals returns true iff pubKey and other contain the same data
func (pubKey PubKey) Equals(other PubKey) bool {
	return pubKey.ID == other.ID &&
		pubKey.Type == other.Type &&
		pubKey.Controller.Equals(other.Controller) &&
		pubKey.PublicKey == other.PublicKey
}

// Validate checks the data contained inside pubKey and returns an error if something is wrong
func (pubKey PubKey) Validate() error {

	regex, _ := regexp.Compile(fmt.Sprintf("^%s#keys-[0-9]+$", pubKey.Controller))
	if !regex.MatchString(pubKey.ID) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, (fmt.Sprintf("Invalid key id, must satisfy %s", regex)))
	}

	if err := keyTypeApproved(pubKey.Type); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
	}

	if pubKey.Type == KeyTypeRsaSignature || pubKey.Type == KeyTypeRsaVerification {
		err := validateRSAPubkey([]byte(pubKey.PublicKey))
		if err != nil {
			return err
		}
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

func (pubKeys PubKeys) HasVerificationAndSignatureKey() bool {
	hasVerification := false
	hasSignature := false

	for _, pk := range pubKeys {
		switch {
		case strings.HasSuffix(pk.ID, "#keys-1") && pk.Type == KeyTypeRsaVerification:
			hasVerification = true
		case strings.HasSuffix(pk.ID, "#keys-2") && pk.Type == KeyTypeRsaSignature:
			hasSignature = true
		}
	}

	return hasSignature && hasVerification
}

func validateRSAPubkey(key []byte) error {
	block, _ := pem.Decode(key)
	_, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}
	return nil
}

func keyTypeApproved(t string) error {
	switch t {
	case KeyTypeRsaVerification, KeyTypeRsaSignature, KeyTypeSecp256k1, KeyTypeEd25519:
		return nil
	default:
		return fmt.Errorf("key type %s not supported", t)
	}
}
