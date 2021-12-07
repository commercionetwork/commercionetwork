package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	commons "github.com/commercionetwork/commercionetwork/x/common/types"
)

func NewMsgSetDidDocument(context string, ID string) *MsgSetDidDocument {
	return &MsgSetDidDocument{
		&DidDocument{Context: []string{context}, ID: ID},
	}

}

func (msg *MsgSetDidDocument) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.DidDocument.ID)
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

func (msg *MsgSetDidDocument) ValidateBasic() error {

	ddo := msg.DidDocument

	// validate ID
	_, err := sdk.AccAddressFromBech32(ddo.ID)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid ID address (%s)", err)
	}

	// TODO: check signer is the same as ID

	// validate Context
	if commons.Strings(ddo.Context).Contains(ContextDidV1) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Context, must include %s", ContextDidV1)
	}

	// validate VerificationMethod
	for _, vm := range ddo.VerificationMethod {
		if err = vm.Validate(); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid VerificationMethod %s %e", vm, err)
		}
	}

	// validate Service
	for _, s := range ddo.Service {
		err = s.Validate()
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Service %s %e", s, err)
		}
	}

	// validate Authentication
	for _, a := range ddo.Authentication {
		if err = a.Validate(); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Autentication %s %e", a, err)
		}
	}

	// validate AssertionMethod
	for _, am := range ddo.AssertionMethod {
		err = am.Validate()
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid Autentication %s %e", am, err)
		}
	}

	// validate CapabilityDelegation
	for _, cd := range ddo.CapabilityDelegation {
		if err = cd.Validate(); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid CapabilityDelegation %s %e", cd, err)
		}
	}

	// validate CapabilityInvocation
	for _, ci := range ddo.CapabilityInvocation {
		if err = ci.Validate(); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid CapabilityInvocation %s %e", ci, err)
		}
	}

	// validate KeyAgreement
	for _, ka := range ddo.KeyAgreement {
		if err = ka.Validate(); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "invalid KeyAgreement %s %e", ka, err)
		}
	}

	// var pubKeys PubKeys
	// pubKeys = msg.PubKeys
	// if err := pubKeys.noDuplicates(); err != nil {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	// }

	// if !pubKeys.HasVerificationAndSignatureKey() {
	// 	return sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "specified public keys are not in the correct format")
	// }

	// if err := msg.lengthLimits(); err != nil {
	// 	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	// }

	return nil
}
