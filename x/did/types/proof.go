package types

import (
	"encoding/base64"
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const ProofPurposeAuthentication = "authentication"

// Validate checks for the content contained inside the proof and
// returns an error if something is invalid
func (proof Proof) Validate() error {
	// proof is empty
	if proof == (Proof{}) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "empty proof")
	}

	if proof.Type != KeyTypeSecp256k12019 {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("Invalid proof type, must be %s", KeyTypeSecp256k12019))
	}

	created, _ := time.Parse(time.RFC3339, proof.Created)
	if created.IsZero() {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Invalid proof creation time")
	}

	if proof.ProofPurpose != ProofPurposeAuthentication {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "proof purpose must be \"authentication\"")
	}

	controller, err := sdk.AccAddressFromBech32(proof.Controller)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "invalid controller, must be a valid bech32-encoded address")
	}

	// decode the bech32 public key
	ppk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, proof.VerificationMethod)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "invalid verification method, must be a bech32-encoded public key")
	}

	ppkAddress, err := sdk.AccAddressFromHex(ppk.Address().String())
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("could not derive AccAddress from verification method: %s", err))
	}

	if !controller.Equals(ppkAddress) {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "verification method-derived AccAddress differs from controller")
	}

	_, err = base64.StdEncoding.DecodeString(proof.SignatureValue)
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "signature value must be base64 encoded")
	}
	return nil
}
