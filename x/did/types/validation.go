package types

import (
	commons "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// const didKeyPrefix = "did:com:"

func isValidDid(did string) error {
	_, err := sdk.AccAddressFromBech32(did)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid ID address (%s)", did, err)
	}

	return nil
}

func (ddo DidDocument) Validate() error {

	// validate ID
	if err := isValidDid(ddo.ID); err != nil {
		return err
	}

	// validate Context
	// @context The JSON-LD Context is either a string or a list containing any combination of strings and/or ordered maps.
	if commons.Strings(ddo.Context).Contains(ContextDidV1) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Context, must include %s", ContextDidV1)
	}

	// validate VerificationMethod
	for _, vm := range ddo.VerificationMethod {
		if err := vm.Validate(); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid VerificationMethod %s %e", vm, err)
		}
	}

	// validate Service
	for _, s := range ddo.Service {
		err := s.Validate()
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Service %s %e", s, err)
		}
	}

	// validate Authentication
	for _, a := range ddo.Authentication {
		if err := a.Validate(); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Autentication %s %e", a, err)
		}
	}

	// validate AssertionMethod
	for _, am := range ddo.AssertionMethod {
		err := am.Validate()
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Autentication %s %e", am, err)
		}
	}

	// validate CapabilityDelegation
	for _, cd := range ddo.CapabilityDelegation {
		if err := cd.Validate(); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid CapabilityDelegation %s %e", cd, err)
		}
	}

	// validate CapabilityInvocation
	for _, ci := range ddo.CapabilityInvocation {
		if err := ci.Validate(); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid CapabilityInvocation %s %e", ci, err)
		}
	}

	// validate KeyAgreement
	for _, ka := range ddo.KeyAgreement {
		if err := ka.Validate(); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid KeyAgreement %s %e", ka, err)
		}
	}

	return nil
}

func (s Service) Validate() error {
	if s.ID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"id\" is required")
	}

	if s.Type == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"type\" is required")
	}

	if s.ServiceEndpoint == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"serviceEndpoint\" is required")
	}

	return nil
}

func (v VerificationMethod) Validate() error {
	if v.ID == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"id\" is required")
	}

	if v.Type == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"type\" is required")
	}

	if v.Controller == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "service field \"serviceEndpoint\" is required")
	}

	// TODO

	return nil
}
