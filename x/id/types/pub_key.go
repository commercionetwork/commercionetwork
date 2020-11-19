package types

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
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
	ID           string         `json:"id" example:"did:com:1tkgm3rra9cs3sfugjqdps30ujggf5klm425zvx#keys-1"`
	Type         string         `json:"type" example:"RsaVerificationKey2018"`
	Controller   sdk.AccAddress `json:"controller"`
	PublicKeyPem string         `json:"publicKeyPem" example:"-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvaM5rNKqd5sl1flSqRHg\nkKdGJzVcktZs0O1IO5A7TauzAtn0vRMr4moWYTn5nUCCiDFbTPoMyPp6tsaZScAD\nG9I7g4vK+/FcImcrdDdv9rjh1aGwkGK3AXUNEG+hkP+QsIBl5ORNSKn+EcdFmnUc\nzhNulA74zQ3xnz9cUtsPC464AWW0Yrlw40rJ/NmDYfepjYjikMVvJbKGzbN3Xwv7\nZzF4bPTi7giZlJuKbNUNTccPY/nPr5EkwZ5/cOZnAJGtmTtj0e0mrFTX8sMPyQx0\nO2uYM97z0SRkf8oeNQm+tyYbwGWY2TlCEXbvhP34xMaBTzWNF5+Z+FZi+UfPfVfK\nHQIDAQAB\n-----END PUBLIC KEY-----\n"`
}

func NewPubKey(pubKeyID string, pubKeyType string, controller sdk.AccAddress, hexValue string) PubKey {
	return PubKey{
		ID:           pubKeyID,
		Type:         pubKeyType,
		Controller:   controller,
		PublicKeyPem: hexValue,
	}
}

// Equals returns true iff pubKey and other contain the same data
func (pubKey PubKey) Equals(other PubKey) bool {
	return pubKey.ID == other.ID &&
		pubKey.Type == other.Type &&
		pubKey.Controller.Equals(other.Controller) &&
		pubKey.PublicKeyPem == other.PublicKeyPem
}

// Validate checks the data contained inside pubKey and returns an error if something is wrong
func (pubKey PubKey) Validate() error {

	if pubKey.Controller == nil {
		return errors.New("controller must be non-null")
	}

	regex, _ := regexp.Compile(fmt.Sprintf("^%s#keys-[0-9]+$", pubKey.Controller))
	if !regex.MatchString(pubKey.ID) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid key id, must satisfy %s", regex))
	}

	if err := keyTypeApproved(pubKey.Type); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, err.Error())
	}

	if pubKey.Type == KeyTypeRsaSignature || pubKey.Type == KeyTypeRsaVerification {
		err := validateRSAPubkey([]byte(pubKey.PublicKeyPem))
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

func (pubKeys PubKeys) noDuplicates() error {
	switch len(pubKeys) {
	case 0, 1:
		return nil
	case 2:
		if pubKeys[0].PublicKeyPem == pubKeys[1].PublicKeyPem {
			return errors.New("keys 0 and 1 are equal")
		}
	default:
		for i, k := range pubKeys {
			for l, m := range pubKeys {
				if i == l {
					continue
				}

				if k.PublicKeyPem == m.PublicKeyPem {
					return fmt.Errorf("keys %d and %d are equal", i, l)
				}

			}
		}
	}

	return nil
}

func validateRSAPubkey(key []byte) error {
	block, _ := pem.Decode(key)
	if block == nil {
		return errors.New("no valid PEM data found")
	}
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
