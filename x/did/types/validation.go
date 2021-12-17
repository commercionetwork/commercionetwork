package types

import (
	commons "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func isValidDidCom(did string) error {
	_, err := sdk.AccAddressFromBech32(did)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid ID address (%s), %e", did, err)
	}

	return nil
}

// func (ddo DidDocument) Validate() error {

// 	// validate ID
// 	if err := isValidDidCom(ddo.ID); err != nil {
// 		return err
// 	}

// 	// validate @context
// 	// @context The JSON-LD Context is either a string or a list containing any combination of strings and/or ordered maps.
// 	// The serialized value of @context MUST be the JSON String https://www.w3.org/ns/did/v1, or a JSON Array where the first item is the JSON String https://www.w3.org/ns/did/v1 and the subsequent items are serialized according to the JSON representation production rules.
// 	if commons.Strings(ddo.Context).Contains(ContextDidV1) {
// 		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid @context, must include %s", ContextDidV1)
// 	}

// 	// validate VerificationMethod
// 	for _, vm := range ddo.VerificationMethod {
// 		if err := vm.Validate(); err != nil {
// 			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid VerificationMethod %s %e", vm, err)
// 		}
// 	}

// 	// validate service
// 	// OPTIONAL
// 	// TODO
// 	// If present, the associated value MUST be a set of services, where each service is described by a map.
// 	// A conforming producer MUST NOT produce multiple service entries with the same id.
// 	for _, s := range ddo.Service {
// 		err := s.Validate()
// 		if err != nil {
// 			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Service %s %e", s, err)
// 		}
// 	}

// 	// validate Authentication
// 	for _, a := range ddo.Authentication {
// 		if err := a.Validate(); err != nil {
// 			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Autentication %s %e", a, err)
// 		}
// 	}

// 	// validate AssertionMethod
// 	for _, am := range ddo.AssertionMethod {
// 		err := am.Validate()
// 		if err != nil {
// 			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Autentication %s %e", am, err)
// 		}
// 	}

// 	// validate CapabilityDelegation
// 	for _, cd := range ddo.CapabilityDelegation {
// 		if err := cd.Validate(); err != nil {
// 			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid CapabilityDelegation %s %e", cd, err)
// 		}
// 	}

// 	// validate CapabilityInvocation
// 	for _, ci := range ddo.CapabilityInvocation {
// 		if err := ci.Validate(); err != nil {
// 			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid CapabilityInvocation %s %e", ci, err)
// 		}
// 	}

// 	// validate KeyAgreement
// 	for _, ka := range ddo.KeyAgreement {
// 		if err := ka.Validate(); err != nil {
// 			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid KeyAgreement %s %e", ka, err)
// 		}
// 	}

// 	return nil
// }

func (s *Service) Validate() error {
	if s == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service is not defined")
	}

	// validate id
	// Required
	// A string that conforms to the rules of RFC3986 for URIs.
	if IsEmpty(s.ID) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"id\" is required")
	}
	if !IsValidRFC3986Uri(s.ID) {
		return sdkerrors.Wrap(ErrInvalidRFC3986UriFormat, "service field \"id\" must conform to the rules of RFC3986 for URIs")
	}

	// validate type
	// Required
	// A string.
	// W3C recommendation: In order to maximize interoperability, the service type and its associated properties SHOULD be registered in the DID Specification Registries
	if IsEmpty(s.Type) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"type\" is required")
	}

	// validate serviceEndpoint
	// Required
	// A string that conforms to the rules of RFC3986 for URIs.
	if IsEmpty(s.ServiceEndpoint) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"serviceEndpoint\" is required")
	}
	if !IsValidRFC3986Uri(s.ServiceEndpoint) {
		return sdkerrors.Wrap(ErrInvalidRFC3986UriFormat, "service field \"serviceEndpoint\" must conform to the rules of RFC3986 for URIs")
	}

	return nil
}

func (v *VerificationMethod) Validate() error {
	if v == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod is not defined")
	}

	// validate ID
	// Required
	// A string that conforms to the rules for DID URL.
	if IsEmpty(v.ID) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"id\" is required")
	}
	if !IsValidDIDURL(v.ID) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"id\" must conform to the rules for DID URL")
	}

	// validate type
	// Required
	// A string that references exactly one verification method type.
	// W3C recommendation: In order to maximize global interoperability, the verification method type SHOULD be registered in the DID Specification Registries -> https://www.w3.org/TR/did-spec-registries/#verification-method-types
	if IsEmpty(v.Type) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"type\" is required")
	}
	if commons.Strings(verificationMethodTypes).Contains(v.Controller) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"type\" is not supported")
	}

	// validate controller
	// Required
	// A string that conforms to the rules of DID Syntax.
	if IsEmpty(v.Controller) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"controller\" is required")
	}
	if !IsValidDID(v.Controller) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"controller\" must conform to the rules of DID Syntax")
	}

	// validate publicKeyMultibase
	// A verification method MUST NOT contain multiple verification material properties for the same material. For example, expressing key material in a verification method using both publicKeyJwk and publicKeyMultibase at the same time is prohibited.
	// -> using only publicKeyMultibase
	if IsEmpty(v.PublicKeyMultibase) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "verificationMethod field \"publicKeyMultibase\" is required")
	}

	return nil
}
