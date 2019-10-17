package types

import (
	"encoding/hex"
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PubKeys []PubKey

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

func (pubKeys PubKeys) FindById(id string) (PubKey, bool) {
	for _, key := range pubKeys {
		if key.Id == id {
			return key, true
		}
	}

	return PubKey{}, false
}

type PubKey struct {
	Id           string         `json:"id"`
	Type         string         `json:"type"`
	Controller   sdk.AccAddress `json:"controller"`
	PublicKeyHex string         `json:"publicKeyHex"`
}

func (pubKey PubKey) Equals(other PubKey) bool {
	return pubKey.Id == other.Id &&
		pubKey.Type == other.Type &&
		pubKey.Controller.Equals(other.Controller) &&
		pubKey.PublicKeyHex == other.PublicKeyHex
}

func (pubKey PubKey) Validate() sdk.Error {

	regex, _ := regexp.Compile(fmt.Sprintf("^%s#keys-^[0-9]+$", pubKey.Controller.String()))
	if !regex.MatchString(pubKey.Id) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("Invalid key id, must satisfy %s", regex.String()))
	}

	if pubKey.Type != KeyTypeRsa && pubKey.Type != KeyTypeSecp256k1 && pubKey.Type != KeyTypeEd25519 {
		msg := fmt.Sprintf("Invalid key type, must be either %s, %s or %s",
			KeyTypeRsa, KeyTypeSecp256k1, KeyTypeEd25519)
		return sdk.ErrUnknownRequest(msg)
	}

	if _, err := hex.DecodeString(pubKey.PublicKeyHex); err != nil {
		return sdk.ErrUnknownRequest("Invalid publicKeyHex value")
	}

	return nil
}
