package types

import (
	commons "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (msg *MsgSetDidDocument) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.ID)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetDidDocument) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetDidDocument) Route() string {
	return RouterKey
}

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

func validateVerificationMethod(verificationMethod []*VerificationMethod) error {
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
		if err := vm.isValid(); err != nil {
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

func (msg *MsgSetDidDocument) ValidateBasic() error {

	if msg == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "msgSetDidDocument not defined")
	}

	// validate @context
	// @context The JSON-LD Context is either a string or a list containing any combination of strings and/or ordered maps.
	// The serialized value of @context MUST be the JSON String https://www.w3.org/ns/did/v1, or a JSON Array where the first item is the JSON String https://www.w3.org/ns/did/v1 and the subsequent items are serialized according to the JSON representation production rules.
	if err := validateContext(msg.Context); err != nil {
		return err
	}

	// validate ID
	if IsEmpty(msg.ID) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"id\" is required")
	}
	if err := isValidDidCom(msg.ID); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ID %s: %e", msg.ID, err)
	}

	// validate VerificationMethod
	if err := validateVerificationMethod(msg.VerificationMethod); err != nil {
		return err
	}

	// validate service
	// OPTIONAL
	// If present, the associated value MUST be a set of services, where each service is described by a map.
	// A conforming producer MUST NOT produce multiple service entries with the same id.
	if err := validateService(msg.Service); err != nil {
		return err
	}

	// validate authentication
	if !commons.Strings(msg.Authentication).IsSet() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid authentication %s found elements with the same ID", msg.Authentication)
	}
	for _, a := range msg.Authentication {
		if !msg.HasVerificationMethod(a) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid authentication: %s is not among the verification methods", a)
		}
	}

	// validate assertionMethod
	if !commons.Strings(msg.AssertionMethod).IsSet() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid assertionMethod %s found elements with the same ID", msg.AssertionMethod)
	}
	for _, am := range msg.AssertionMethod {
		if !msg.HasVerificationMethod(am) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid assertionMethod: %s is not among the verification methods", am)
		}
	}

	// validate capabilityDelegation
	if !commons.Strings(msg.CapabilityDelegation).IsSet() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid capabilityDelegation %s found elements with the same ID", msg.CapabilityDelegation)
	}
	for _, cd := range msg.CapabilityDelegation {
		if !msg.HasVerificationMethod(cd) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid capabilityDelegation: %s is not among the verification methods", cd)
		}
	}

	// validate capabilityInvocation
	if !commons.Strings(msg.CapabilityInvocation).IsSet() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid capabilityInvocation %s found elements with the same ID", msg.CapabilityInvocation)
	}
	for _, ci := range msg.CapabilityInvocation {
		if !msg.HasVerificationMethod(ci) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid capabilityInvocation: %s is not among the verification methods", ci)
		}
	}

	// validate keyAgreement
	if !commons.Strings(msg.KeyAgreement).IsSet() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid keyAgreement %s found elements with the same ID", msg.KeyAgreement)
	}
	for _, ka := range msg.KeyAgreement {
		if !msg.HasVerificationMethod(ka) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid keyAgreement: %s is not among the verification methods", ka)
		}
	}

	return nil
}

func (msg MsgSetDidDocument) HasVerificationMethod(id string) bool {
	for _, vm := range msg.VerificationMethod {
		// DID url
		if id == vm.ID {
			return true
		}
		// relative DID url
		if msg.ID+id == vm.ID {
			return true
		}
	}
	return false
}
