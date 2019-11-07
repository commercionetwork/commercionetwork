package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// A proof on a DID Document is cryptographic proof of the integrity of the DID Document according to either:
// 1. The subject, or:
// 1. The controller, if present.
type Proof struct {
	Type           string    `json:"type"`
	Created        time.Time `json:"created"`
	Creator        string    `json:"creator"`
	SignatureValue string    `json:"signatureValue"`
}

func NewProof(proofType string, created time.Time, creator, signatureValue string) Proof {
	return Proof{
		Type:           proofType,
		Created:        created,
		Creator:        creator,
		SignatureValue: signatureValue,
	}
}

// Equals returns true iff proof and other contain the same data.
func (proof Proof) Equals(other Proof) bool {
	return proof.Type == other.Type &&
		proof.Created.Equal(other.Created) &&
		proof.Creator == other.Creator &&
		proof.SignatureValue == other.SignatureValue
}

// Validate checks for the content contained inside the proof and
// returns an error if something is invalid
func (proof Proof) Validate() sdk.Error {

	if proof.Type != "LinkedDataSignature2015" {
		return sdk.ErrUnknownRequest("Invalid proof type, must be LinkedDataSignature2015")
	}

	if proof.Created.IsZero() {
		return sdk.ErrUnknownRequest("Invalid proof creation type")
	}

	return nil
}
