package types

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (pubKey PubKey) Validate() error {

	if pubKey.Controller == "" {
		return errors.New("controller must be non-null")
	}

	regex, _ := regexp.Compile(fmt.Sprintf("^%s#keys-[0-9]+$", pubKey.Controller))
	if !regex.MatchString(pubKey.ID) {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Invalid key id, must satisfy %s", regex))
	}

	if err := keyTypeApproved(pubKey.Type); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, err.Error())
	}

	if pubKey.Type == KeyTypeRsaSignature || pubKey.Type == KeyTypeRsaVerification {
		err := validateRSAPubkey([]byte(pubKey.PublicKeyPem))
		if err != nil {
			return err
		}
	}

	return nil
}

func keyTypeApproved(t string) error {
	switch t {
	case KeyTypeRsaVerification, KeyTypeRsaSignature, KeyTypeSecp256k1, KeyTypeEd25519, KeyTypeBls12381G1Key2020:
		return nil
	default:
		return fmt.Errorf("key type %s not supported", t)
	}
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

// Equals returns true iff pubKey and other contain the same data
func (pubKey PubKey) Equals(other PubKey) bool {
	controller, _ := sdk.AccAddressFromBech32(pubKey.Controller)
	otherController, _ := sdk.AccAddressFromBech32(other.Controller)
	return pubKey.ID == other.ID &&
		pubKey.Type == other.Type &&
		controller.Equals(otherController) &&
		pubKey.PublicKeyPem == other.PublicKeyPem
}

type PubKeys []*PubKey

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
