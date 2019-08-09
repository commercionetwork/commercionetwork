package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is they name of the CommercioID module
const RouterKey = "commercioid"

type MsgSetIdentity struct {
	Owner          sdk.AccAddress `json:"owner"`
	DidDocumentUri string         `json:"did_document_uri"`
}

func NewMsgSetIdentity(owner sdk.AccAddress, didDocumentUri string) MsgSetIdentity {
	return MsgSetIdentity{
		Owner:          owner,
		DidDocumentUri: didDocumentUri,
	}
}

// Route Implements Msg.
func (msg MsgSetIdentity) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSetIdentity) Type() string { return "set_identity" }

// ValidateBasic Implements Msg.
func (msg MsgSetIdentity) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.DidDocumentUri) == 0 {
		return sdk.ErrUnknownRequest("Did Document reference cannot be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSetIdentity) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSetIdentity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
