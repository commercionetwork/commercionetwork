package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/crypto/secp256k1"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DidDocument is the concrete serialization of the data model, according to a particular syntax.
// The Did Document contains attributes or claims about the DID Subject, and the DID itself is contained in
// the id property.
type DidDocument struct {
	Context string         `json:"@context"`
	ID      sdk.AccAddress `json:"id"`
	PubKeys PubKeys        `json:"publicKey"`
	Proof   Proof          `json:"proof"`
	Service Services       `json:"service"`
}

type Services []Service

func (s Services) Equals(other Services) bool {
	if len(s) != len(other) {
		return false
	}

	for key, value := range other {
		if !s[key].Equals(value) {
			return false
		}
	}

	return true
}

// Service represents a service type needed for DidDocument.
type Service struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

// Validate returns error when Service is not valid.
func (s Service) Validate() error {
	if s.Id == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "service field \"id\" is required")
	}

	if s.Type == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "service field \"type\" is required")
	}

	if s.ServiceEndpoint == "" {
		return sdkErr.Wrap(sdkErr.ErrInvalidRequest, "service field \"serviceEndpoint\" is required")
	}

	return nil
}

func (s Service) Equals(otherService Service) bool {
	return s.ServiceEndpoint == otherService.ServiceEndpoint &&
		s.Type == otherService.Type &&
		s.Id == otherService.Id
}

// didDocumentUnsigned is an intermediate type used to check for proof correctness
type didDocumentUnsigned struct {
	Context string         `json:"@context"`
	ID      sdk.AccAddress `json:"id"`
	PubKeys PubKeys        `json:"publicKey"`
}

// Equals returns true iff didDocument and other contain the same data
func (didDocument DidDocument) Equals(other DidDocument) bool {
	return didDocument.Context == other.Context &&
		didDocument.ID.Equals(other.ID) &&
		didDocument.PubKeys.Equals(other.PubKeys) &&
		didDocument.Proof.Equals(other.Proof) &&
		didDocument.Service.Equals(other.Service)
}

// Validate checks the data present inside this Did Document and returns an
// error if something is wrong
func (didDocument DidDocument) Validate() error {

	if didDocument.ID.Empty() {
		return sdkErr.Wrap(sdkErr.ErrInvalidAddress, (didDocument.ID.String()))
	}

	if didDocument.Context != "https://www.w3.org/ns/did/v1" {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, ("Invalid context, must be https://www.w3.org/ns/did/v1"))
	}

	for _, key := range didDocument.PubKeys {
		if err := key.Validate(); err != nil {
			return err
		}

		if !didDocument.ID.Equals(key.Controller) {
			return sdkErr.Wrap(sdkErr.ErrUnknownRequest, ("Public key controller must match did document id"))
		}
	}

	if !didDocument.PubKeys.HasVerificationAndSignatureKey() {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "specified public keys are not in the correct format")
	}

	if err := didDocument.Proof.Validate(); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("proof validation error: %s", err.Error()))
	}

	if err := didDocument.VerifyProof(); err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized, err.Error())
	}

	return nil
}

// VerifyProof verifies d's Proof against its content.
// The Proof is constructed as follows:
//  - let K be the Bech32 Account public key, embedded in the Proof "Verification Method" field
//  - let S be K transformed in a raw Secp256k1 public key
//  - let B be the SHA-512 (as defined in the FIPS 180-4) of the JSON representation of d minus the Proof field
//  - let L be the Proof Signature Value, decoded from Base64 encoding
// The Proof is verified if K.Verify(B, L) is verified.
func (didDocument DidDocument) VerifyProof() error {
	u := didDocumentUnsigned{
		Context: didDocument.Context,
		ID:      didDocument.ID,
		PubKeys: didDocument.PubKeys,
	}

	oProof := didDocument.Proof

	// get a public key object
	pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, oProof.VerificationMethod)
	if err != nil {
		return err
	}

	// get a seck256k1 public key
	sk := pk.(secp256k1.PubKeySecp256k1)

	// marshal u in json
	data, err := json.Marshal(u)
	if err != nil {
		return fmt.Errorf("could not marshal unsigned did document during proof verification: %w", err)
	}

	// get signature bytes from base64
	sigBytes, err := base64.StdEncoding.DecodeString(oProof.SignatureValue)
	if err != nil {
		return fmt.Errorf("could not decode base64 signature value: %w", err)
	}

	verified := sk.VerifyBytes(data, sigBytes)

	if !verified {
		return fmt.Errorf("proof signature verification failed")
	}

	return nil
}
