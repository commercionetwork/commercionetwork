package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
	if err := msg.DidDocument.Validate(); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid didDocument", err)
	}

	return nil
}
