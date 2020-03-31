package types

import (
	"fmt"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DidDocument is the concrete serialization of the data model, according to a particular syntax.
// The Did Document contains attributes or claims about the DID Subject, and the DID itself is contained in
// the id property.
type DidDocument struct {
	Context        string         `json:"@context"`
	ID             sdk.AccAddress `json:"id"`
	PubKeys        PubKeys        `json:"publicKey"`
	Authentication ctypes.Strings `json:"authentication"`
	Proof          Proof          `json:"proof"`
	Services       Services       `json:"service"`
}

// Equals returns true iff didDocument and other contain the same data
func (didDocument DidDocument) Equals(other DidDocument) bool {
	return didDocument.Context == other.Context &&
		didDocument.ID.Equals(other.ID) &&
		didDocument.PubKeys.Equals(other.PubKeys) &&
		didDocument.Authentication.Equals(other.Authentication) &&
		didDocument.Proof.Equals(other.Proof) &&
		didDocument.Services.Equals(other.Services)
}

// Validate checks the data present inside this Did Document and returns an
// error if something is wrong
func (didDocument DidDocument) Validate() error {

	if didDocument.ID.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, didDocument.ID.String())
	}

	if didDocument.Context != "https://www.w3.org/ns/did/v1" {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid context, must be https://www.w3.org/ns/did/v1")
	}

	if len(didDocument.PubKeys) != 3 {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Field publicKey must have length of 3")
	}

	// -----------------------------------
	// --- Validate the authentication
	// -----------------------------------

	if len(didDocument.Authentication) != 1 {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Array authentication cannot have more than one item")
	}

	authKey, found := didDocument.PubKeys.FindByID(didDocument.Authentication[0])
	if !found {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Authentication key not found inside publicKey array")
	}

	if authKey.Type != KeyTypeSecp256k1 && authKey.Type != KeyTypeEd25519 {
		msg := fmt.Sprintf("Authentication key type must be either %s or %s", KeyTypeSecp256k1, KeyTypeEd25519)
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, msg)
	}

	// --------------------------------
	// --- Validate public keys
	// --------------------------------

	for _, key := range didDocument.PubKeys {
		if err := key.Validate(); err != nil {
			return err
		}

		if !didDocument.ID.Equals(key.Controller) {
			return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Public key controller must match did document id")
		}
	}

	// ------------------------------------------
	// --- Validate the proof creator value
	// ------------------------------------------

	if didDocument.Proof.Creator != didDocument.Authentication[0] {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid proof key, must be the authentication one")
	}

	// -----------------------------
	// --- Validate the services
	// -----------------------------

	for _, service := range didDocument.Services {
		if err := service.Validate(); err != nil {
			return err
		}
	}

	return nil
}
