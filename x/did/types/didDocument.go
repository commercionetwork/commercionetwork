package types

import (
	commons "github.com/commercionetwork/commercionetwork/x/common/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	errorsmod "cosmossdk.io/errors"
)

func (ddo *DidDocument) Validate() error {

	if ddo == nil {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "didDocument not defined")
	}

	// validate @context
	// @context The JSON-LD Context is either a string or a list containing any combination of strings and/or ordered maps.
	// The serialized value of @context MUST be the JSON String https://www.w3.org/ns/did/v1, or a JSON Array where the first item is the JSON String https://www.w3.org/ns/did/v1 and the subsequent items are serialized according to the JSON representation production rules.
	if err := validateContextSlice(ddo.Context); err != nil {
		return err
	}

	// validate ID
	if IsEmpty(ddo.ID) {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"id\" is required")
	}
	if err := Validate(ddo.ID); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid ID %s: %e", ddo.ID, err)
	}

	// validate VerificationMethod
	if err := validateVerificationMethodSlice(ddo.VerificationMethod, ddo.ID); err != nil {
		return err
	}

	// validate service
	// OPTIONAL
	// If present, the associated value MUST be a set of services, where each service is described by a map.
	// A conforming producer MUST NOT produce multiple service entries with the same id.
	if err := validateServiceSlice(ddo.Service); err != nil {
		return err
	}

	// validate authentication
	if !commons.Strings(ddo.Authentication).IsSet() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid authentication %s found elements with the same ID", ddo.Authentication)
	}
	for _, a := range ddo.Authentication {
		if !ddo.hasVerificationMethod(a) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid authentication: %s is not among the verification methods", a)
		}
	}

	// validate assertionMethod
	if !commons.Strings(ddo.AssertionMethod).IsSet() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid assertionMethod %s found elements with the same ID", ddo.AssertionMethod)
	}
	for _, am := range ddo.AssertionMethod {
		if !ddo.hasVerificationMethod(am) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid assertionMethod: %s is not among the verification methods", am)
		}
	}

	// validate capabilityDelegation
	if !commons.Strings(ddo.CapabilityDelegation).IsSet() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid capabilityDelegation %s found elements with the same ID", ddo.CapabilityDelegation)
	}
	for _, cd := range ddo.CapabilityDelegation {
		if !ddo.hasVerificationMethod(cd) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid capabilityDelegation: %s is not among the verification methods", cd)
		}
	}

	// validate capabilityInvocation
	if !commons.Strings(ddo.CapabilityInvocation).IsSet() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid capabilityInvocation %s found elements with the same ID", ddo.CapabilityInvocation)
	}
	for _, ci := range ddo.CapabilityInvocation {
		if !ddo.hasVerificationMethod(ci) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid capabilityInvocation: %s is not among the verification methods", ci)
		}
	}

	// validate keyAgreement
	if !commons.Strings(ddo.KeyAgreement).IsSet() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid keyAgreement %s found elements with the same ID", ddo.KeyAgreement)
	}
	for _, ka := range ddo.KeyAgreement {
		if !ddo.hasVerificationMethod(ka) {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid keyAgreement: %s is not among the verification methods", ka)
		}
	}

	return nil
}

func validateContextSlice(context []string) error {
	if !commons.Strings(context).IsSet() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid @context %s found elements with the same ID", context)
	}
	if len(context) == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid @context, must include %s", ContextDidV1)
	}
	if context[0] != ContextDidV1 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid @context, %s must be the first element", ContextDidV1)
	}

	return nil
}

func validateVerificationMethodSlice(verificationMethod []*VerificationMethod, subject string) error {
	isVerificationMethodSet := func() bool {
		keys := []string{}
		for _, s := range verificationMethod {
			keys = append(keys, s.ID)
		}

		return commons.Strings(keys).IsSet()
	}
	if !isVerificationMethodSet() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid verificationMethod %s found elements with the same ID", verificationMethod)
	}
	var containsRsaSignature2018, containsRsaVerificationKey2018 bool
	for _, vm := range verificationMethod {
		if err := vm.Validate(subject); err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid verificationMethod %s %e", vm, err)
		}
		if vm.Type == RsaSignature2018 {
			containsRsaSignature2018 = true
		}
		if vm.Type == RsaVerificationKey2018 {
			containsRsaVerificationKey2018 = true
		}
	}
	if !containsRsaSignature2018 || !containsRsaVerificationKey2018 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid verificationMethod %s types %s and %s are required", verificationMethod, RsaSignature2018, RsaVerificationKey2018)
	}

	return nil
}

func validateServiceSlice(service []*Service) error {
	isServiceSet := func() bool {
		keys := []string{}
		for _, s := range service {
			keys = append(keys, s.ID)
		}

		return commons.Strings(keys).IsSet()
	}
	if !isServiceSet() {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service %s found elements with the same ID", service)
	}
	for _, s := range service {
		if err := s.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (ddo DidDocument) hasVerificationMethod(id string) bool {
	for _, vm := range ddo.VerificationMethod {
		// DID url
		if id == vm.ID {
			return true
		}
		// relative DID url
		if ddo.ID+id == vm.ID {
			return true
		}
	}
	return false
}
