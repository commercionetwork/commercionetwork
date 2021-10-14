package types

import (
	//"encoding/base64"
	//"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSetIdentity{}

func NewMsgSetIdentity(context string, ID string, pubkeys []*PubKey, service []*Service) *MsgSetIdentity {
	return &MsgSetIdentity{
		Context: context,
		ID:      ID,
		PubKeys: pubkeys,
		Service: service,
	}
}

func (msg *MsgSetIdentity) Route() string {
	return RouterKey
}

func (msg *MsgSetIdentity) Type() string {
	return "SetIdentity"
}

func (msg *MsgSetIdentity) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.ID)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetIdentity) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetIdentity) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.ID)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Context != ContextDidV1 {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "Invalid context, must be https://www.w3.org/ns/did/v1")
	}

	controller, _ := sdk.AccAddressFromBech32(msg.ID)

	for _, key := range msg.PubKeys {
		if err := key.Validate(); err != nil {
			return err
		}
		keycontroller, _ := sdk.AccAddressFromBech32(key.Controller)
		if !controller.Equals(keycontroller) {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "Public key controller must match did document id")
		}
	}

	var pubKeys PubKeys
	pubKeys = msg.PubKeys
	if err := pubKeys.noDuplicates(); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if !pubKeys.HasVerificationAndSignatureKey() {
		return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "specified public keys are not in the correct format")
	}
	/*
		if msg.Proof == nil {
			return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "proof not provided")
		}

		if err := msg.Proof.Validate(); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, fmt.Sprintf("proof validation error: %s", err.Error()))
		}
	*/
	// we have some service, we should validate 'em
	if msg.Service != nil {
		for i, service := range msg.Service {
			err := service.Validate()
			if err != nil {
				return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, fmt.Sprintf("service %d validation failed: %w", i, err))
			}
		}
	}
	/*
		if err := msg.VerifyProof(); err != nil {
			return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, err.Error())
		}*/

	if err := msg.lengthLimits(); err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	return nil
}

// VerifyProof verifies d's Proof against its content.
// The Proof is constructed as follows:
//  - let K be the Bech32 Account public key, embedded in the Proof "Verification Method" field
//  - let S be K transformed in a raw Secp256k1 public key
//  - let B be the SHA-256 (as defined in the FIPS 180-4) of the JSON representation of d minus the Proof field
//  - let L be the Proof Signature Value, decoded from Base64 encoding
// The Proof is verified if K.Verify(B, L) is verified.
/*func (msg *MsgSetIdentity) VerifyProof() error {
	//u := DidDocumentUnsigned(msg)
	u := msg
	// Explicitly zero out the Proof field.
	//
	// Here we leverage the fact that encoding/json do not encode nil pointers,
	// effectively giving us DidDocument-(Proof field).
	u.Proof = nil

	oProof := msg.Proof

	// get a public key object
	pk, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, oProof.VerificationMethod)
	if err != nil {
		return err
	}
	// get a seck256k1 public key
	// ********** TO CONTROL *********************
	// TYPE is not supported
	//sk := pk.(secp256k1.PubKeySecp256k1)
	// ********** TO CONTROL *********************

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

	//verified := sk.VerifyBytes(data, sigBytes)
	verified := pk.VerifySignature(data, sigBytes)

	if !verified {
		return fmt.Errorf("proof signature verification failed")
	}

	return nil
}
*/
func (msg *MsgSetIdentity) lengthLimits() error {
	e := func(fieldName string, maxLen int) error {
		return fmt.Errorf("%s content can't be longer than %d bytes", fieldName, maxLen)
	}

	for i, s := range msg.Service {
		if len(s.ID) > 64 {
			return e(fmt.Sprintf("service.%d.id", i), 64)
		}

		if len(s.Type) > 64 {
			return e(fmt.Sprintf("service.%d.type", i), 64)
		}

		if len(s.ServiceEndpoint) > 512 {
			return e(fmt.Sprintf("service.%d.serviceEndpoint", i), 512)
		}
	}

	return nil
}
