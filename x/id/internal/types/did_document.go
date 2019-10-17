package types

import (
	"fmt"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DidDocument represents the data related to a single Did Document
type DidDocument struct {
	Context        string         `json:"@context"`
	Id             sdk.AccAddress `json:"id"`
	PubKeys        PubKeys        `json:"publicKey"`
	Authentication ctypes.Strings `json:"authentication"`
	Proof          Proof          `json:"proof"`
}

// Equals returns true iff this didDocument and other contain the same data
func (didDocument DidDocument) Equals(other DidDocument) bool {
	return didDocument.Context == other.Context &&
		didDocument.Id.Equals(other.Id) &&
		didDocument.PubKeys.Equals(other.PubKeys) &&
		didDocument.Authentication.Equals(other.Authentication) &&
		didDocument.Proof.Equals(other.Proof)
}

// Validate checks the data present inside this Did Document and returns an
// error if something is wrong
func (didDocument DidDocument) Validate() sdk.Error {

	if didDocument.Context != "https://www.w3.org/2019/did/v1" {
		return sdk.ErrUnknownRequest("Invalid context, must be https://www.w3.org/2019/did/v1")
	}

	if len(didDocument.PubKeys) != 3 {
		return sdk.ErrUnknownRequest("Field publicKey must have length of 3")
	}

	// -----------------------------------
	// --- Validate the authentication
	// -----------------------------------

	if len(didDocument.Authentication) != 1 {
		return sdk.ErrUnknownRequest("Array authentication cannot have more than one item")
	}

	authKey, found := didDocument.PubKeys.FindById(didDocument.Authentication[0])
	if !found {
		return sdk.ErrUnknownRequest("Authentication key not found inside publicKey array")
	}

	if authKey.Type != KeyTypeSecp256k1 && authKey.Type != KeyTypeEd25519 {
		msg := fmt.Sprintf("Authentication key type must be either %s or %s", KeyTypeSecp256k1, KeyTypeEd25519)
		return sdk.ErrUnknownRequest(msg)
	}

	// --------------------------------
	// --- Validate public keys
	// --------------------------------

	for _, key := range didDocument.PubKeys {
		if err := key.Validate(); err != nil {
			return err
		}

		if !didDocument.Id.Equals(key.Controller) {
			return sdk.ErrUnknownRequest("Public key controller must match did document id")
		}
	}

	// ------------------------------------------
	// --- Validate the proof creator value
	// ------------------------------------------

	if didDocument.Proof.Creator != didDocument.Authentication[0] {
		return sdk.ErrUnknownRequest("Invalid proof key, must be the authentication one")
	}

	return nil
}
