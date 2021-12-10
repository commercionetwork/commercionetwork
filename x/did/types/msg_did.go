package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewMsgSetDidDocument(context string, ID string) *MsgSetDidDocument {
	return &MsgSetDidDocument{
		Context:              []string{},
		ID:                   "",
		VerificationMethod:   []*VerificationMethod{},
		Service:              []*Service{},
		Authentication:       []*VerificationMethod{},
		AssertionMethod:      []*VerificationMethod{},
		CapabilityDelegation: []*VerificationMethod{},
		CapabilityInvocation: []*VerificationMethod{},
		KeyAgreement:         []*VerificationMethod{},
	}

}

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

func (msg *MsgSetDidDocument) ValidateBasic() error {
	// TODO

	// if err := msg.DidDocument.Validate(); err != nil {
	// 	return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid DID document", err)
	// }

	return nil
}
