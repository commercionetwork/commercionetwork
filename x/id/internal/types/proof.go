package types

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Proof struct {
	Type           string    `json:"type"`
	Created        time.Time `json:"created"`
	Creator        string    `json:"creator"`
	SignatureValue string    `json:"signatureValue"`
}

func (proof Proof) Equals(other Proof) bool {
	return proof.Type == other.Type &&
		proof.Created.Equal(other.Created) &&
		proof.Creator == other.Creator &&
		proof.SignatureValue == other.SignatureValue
}

func (proof Proof) Validate() sdk.Error {

	if proof.Type != "LinkedDataSignature2015" {
		return sdk.ErrUnknownRequest("Invalid proof type, must be LinkedDataSignature2015")
	}

	if proof.Created.IsZero() {
		return sdk.ErrUnknownRequest("Invalid proof creation type")
	}

	return nil
}
