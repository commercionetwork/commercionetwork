package types

import (
	commons "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// GetSigners implements "github.com/cosmos/cosmos-sdk/types".Msg
func (msg *MsgSetDidDocument) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.DidDocument.ID)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes implements "github.com/cosmos/cosmos-sdk/types".Msg
func (msg *MsgSetDidDocument) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Route implements "github.com/cosmos/cosmos-sdk/types".Msg
func (msg *MsgSetDidDocument) Route() string {
	return RouterKey
}

// Type implements "github.com/cosmos/cosmos-sdk/types".Msg
func (msg *MsgSetDidDocument) Type() string {
	return MsgTypeSetDid
}

func validateContext(context []string) error {
	if !commons.Strings(context).IsSet() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid @context %s found elements with the same ID", context)
	}
	if len(context) == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid @context, must include %s", ContextDidV1)
	}
	if context[0] != ContextDidV1 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid @context, %s must be the first element", ContextDidV1)
	}

	return nil
}

func validateVerificationMethod(verificationMethod []*VerificationMethod, subject string) error {
	isVerificationMethodSet := func() bool {
		keys := []string{}
		for _, s := range verificationMethod {
			keys = append(keys, s.ID)
		}

		return commons.Strings(keys).IsSet()
	}
	if !isVerificationMethodSet() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid verificationMethod %s found elements with the same ID", verificationMethod)
	}
	var containsRsaSignature2018, containsRsaVerificationKey2018 bool
	for _, vm := range verificationMethod {
		if err := vm.isValid(subject); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid verificationMethod %s %e", vm, err)
		}
		if vm.Type == RsaSignature2018 {
			containsRsaSignature2018 = true
		}
		if vm.Type == RsaVerificationKey2018 {
			containsRsaVerificationKey2018 = true
		}
	}
	if !containsRsaSignature2018 || !containsRsaVerificationKey2018 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid verificationMethod %s types %s and %s are required", verificationMethod, RsaSignature2018, RsaVerificationKey2018)
	}

	return nil
}

func validateService(service []*Service) error {
	isServiceSet := func() bool {
		keys := []string{}
		for _, s := range service {
			keys = append(keys, s.ID)
		}

		return commons.Strings(keys).IsSet()
	}
	if !isServiceSet() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service %s found elements with the same ID", service)
	}
	for _, s := range service {
		if err := s.isValid(); err != nil {
			return err
		}
	}

	return nil
}

// ValidateBasic implements "github.com/cosmos/cosmos-sdk/types".Msg
func (msg *MsgSetDidDocument) ValidateBasic() error {

	if msg == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "msgSetDidDocument not defined")
	}

	// validate @context
	// @context The JSON-LD Context is either a string or a list containing any combination of strings and/or ordered maps.
	// The serialized value of @context MUST be the JSON String https://www.w3.org/ns/did/v1, or a JSON Array where the first item is the JSON String https://www.w3.org/ns/did/v1 and the subsequent items are serialized according to the JSON representation production rules.
	if err := validateContext(msg.DidDocument.Context); err != nil {
		return err
	}

	// validate ID
	if IsEmpty(msg.DidDocument.ID) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"id\" is required")
	}
	if err := isValidDidCom(msg.DidDocument.ID); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ID %s: %e", msg.DidDocument.ID, err)
	}

	// validate VerificationMethod
	if err := validateVerificationMethod(msg.DidDocument.VerificationMethod, msg.DidDocument.ID); err != nil {
		return err
	}

	// validate service
	// OPTIONAL
	// If present, the associated value MUST be a set of services, where each service is described by a map.
	// A conforming producer MUST NOT produce multiple service entries with the same id.
	if err := validateService(msg.DidDocument.Service); err != nil {
		return err
	}

	// validate authentication
	if !commons.Strings(msg.DidDocument.Authentication).IsSet() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid authentication %s found elements with the same ID", msg.DidDocument.Authentication)
	}
	for _, a := range msg.DidDocument.Authentication {
		if !msg.hasVerificationMethod(a) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid authentication: %s is not among the verification methods", a)
		}
	}

	// validate assertionMethod
	if !commons.Strings(msg.DidDocument.AssertionMethod).IsSet() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid assertionMethod %s found elements with the same ID", msg.DidDocument.AssertionMethod)
	}
	for _, am := range msg.DidDocument.AssertionMethod {
		if !msg.hasVerificationMethod(am) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid assertionMethod: %s is not among the verification methods", am)
		}
	}

	// validate capabilityDelegation
	if !commons.Strings(msg.DidDocument.CapabilityDelegation).IsSet() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid capabilityDelegation %s found elements with the same ID", msg.DidDocument.CapabilityDelegation)
	}
	for _, cd := range msg.DidDocument.CapabilityDelegation {
		if !msg.hasVerificationMethod(cd) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid capabilityDelegation: %s is not among the verification methods", cd)
		}
	}

	// validate capabilityInvocation
	if !commons.Strings(msg.DidDocument.CapabilityInvocation).IsSet() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid capabilityInvocation %s found elements with the same ID", msg.DidDocument.CapabilityInvocation)
	}
	for _, ci := range msg.DidDocument.CapabilityInvocation {
		if !msg.hasVerificationMethod(ci) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid capabilityInvocation: %s is not among the verification methods", ci)
		}
	}

	// validate keyAgreement
	if !commons.Strings(msg.DidDocument.KeyAgreement).IsSet() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid keyAgreement %s found elements with the same ID", msg.DidDocument.KeyAgreement)
	}
	for _, ka := range msg.DidDocument.KeyAgreement {
		if !msg.hasVerificationMethod(ka) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid keyAgreement: %s is not among the verification methods", ka)
		}
	}

	return nil
}

func (msg MsgSetDidDocument) hasVerificationMethod(id string) bool {
	for _, vm := range msg.DidDocument.VerificationMethod {
		// DID url
		if id == vm.ID {
			return true
		}
		// relative DID url
		if msg.DidDocument.ID+id == vm.ID {
			return true
		}
	}
	return false
}
