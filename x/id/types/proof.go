package types

import (
	"encoding/base64"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
)

const ProofPurposeAuthentication = "authentication"

// A proof on a DID Document is cryptographic proof of the integrity of the DID Document according to either:
// 1. The subject, or:
// 1. The controller, if present.
type Proof struct {
	Type               string    `json:"type" example:"EcdsaSecp256k1VerificationKey2019"`
	Created            time.Time `json:"created" example:"2020-04-22T04:23:50.73112321Z"`
	ProofPurpose       string    `json:"proofPurpose" example:"authentication"`
	Controller         string    `json:"controller" example:"did:com:1tkgm3rra9cs3sfugjqdps30ujggf5klm425zvx"`
	VerificationMethod string    `json:"verificationMethod" example:"did:com:pub1addwnpepqt6lnn5v0c3rys49v5v9f4kvcchehnu7kyk8t8vce5lsxfy7e2pxwyvmf6t"`
	SignatureValue     string    `json:"signatureValue" example:"nIgRvObXlF2OIbktZcQJw0UU7zDEku8cEBq7194YOjhEvD5wBZ+TcNu9GNRZucC6OyuplHfK6uo57+3lVQbpgA=="`
}

func NewProof(proofType string, created time.Time, proofPurpose, controller, verificationMethod, signatureValue string) Proof {
	return Proof{
		Type:               proofType,
		Created:            created,
		ProofPurpose:       proofPurpose,
		Controller:         controller,
		VerificationMethod: verificationMethod,
		SignatureValue:     signatureValue,
	}
}

// Equals returns true iff proof and other contain the same data.
func (proof Proof) Equals(other Proof) bool {
	return proof.Type == other.Type &&
		proof.Created.Equal(other.Created) &&
		proof.ProofPurpose == other.ProofPurpose &&
		proof.Controller == other.Controller &&
		proof.VerificationMethod == other.VerificationMethod &&
		proof.SignatureValue == other.SignatureValue
}

// Validate checks for the content contained inside the proof and
// returns an error if something is invalid
func (proof Proof) Validate() error {
	// proof is empty
	if proof == (Proof{}) {
		return sdkErr.Wrap(sdkErr.ErrUnauthorized, "empty proof")
	}

	if proof.Type != KeyTypeSecp256k12019 {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid proof type, must be %s", KeyTypeSecp256k12019))
	}

	if proof.Created.IsZero() {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid proof creation time")
	}

	if proof.ProofPurpose != ProofPurposeAuthentication {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "proof purpose must be \"authentication\"")
	}

	controller, err := sdk.AccAddressFromBech32(proof.Controller)
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "invalid controller, must be a valid bech32-encoded address")
	}

	// decode the bech32 public key
	ppk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, proof.VerificationMethod)
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "invalid verification method, must be a bech32-encoded public key")
	}

	ppkAddress, err := sdk.AccAddressFromHex(ppk.Address().String())
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("could not derive AccAddress from verification method: %s", err))
	}

	if !controller.Equals(ppkAddress) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "verification method-derived AccAddress differs from controller")
	}

	_, err = base64.StdEncoding.DecodeString(proof.SignatureValue)
	if err != nil {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, "signature value must be base64 encoded")
	}
	return nil
}
