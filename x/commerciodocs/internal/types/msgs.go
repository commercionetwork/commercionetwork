package types

import (
	"commercio-network/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = "commerciodocs"

// ----------------------------------
// --- StoreDocument
// ----------------------------------

type MsgStoreDocument struct {
	Owner     sdk.AccAddress `json:"owner"`
	Identity  types.Did      `json:"identity"`
	Reference string         `json:"reference"`
	Metadata  string         `json:"metadata"`
}

func NewMsgStoreDocument(owner sdk.AccAddress, identity types.Did, reference string, metadata string) MsgStoreDocument {
	return MsgStoreDocument{
		Owner:     owner,
		Identity:  identity,
		Reference: reference,
		Metadata:  metadata,
	}
}

// RouterKey Implements Msg.
func (msg MsgStoreDocument) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgStoreDocument) Type() string { return "store_document" }

// ValidateBasic Implements Msg.
func (msg MsgStoreDocument) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Reference) == 0 || len(msg.Identity) == 0 || len(msg.Metadata) == 0 {
		return sdk.ErrUnknownRequest("Identity and/or document reference and/or metadata reference cannot be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgStoreDocument) GetSignBytes() []byte {
	return sdk.MustSortJSON(types2.msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgStoreDocument) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// ----------------------------------
// --- ShareDocument
// ----------------------------------

type MsgShareDocument struct {
	Owner     sdk.AccAddress `json:"owner"`
	Sender    types.Did      `json:"sender"`
	Receiver  types.Did      `json:"receiver"`
	Reference string         `json:"reference"`
}

func NewMsgShareDocument(owner sdk.AccAddress, reference string, sender types.Did, receiver types.Did) MsgShareDocument {
	return MsgShareDocument{
		Owner:     owner,
		Sender:    sender,
		Receiver:  receiver,
		Reference: reference,
	}
}

// RouterKey Implements Msg.
func (msg MsgShareDocument) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgShareDocument) Type() string { return "share_document" }

// ValidateBasic Implements Msg.
func (msg MsgShareDocument) ValidateBasic() sdk.Error {
	if msg.Owner.Empty() {
		return sdk.ErrInvalidAddress(msg.Owner.String())
	}
	if len(msg.Sender) == 0 || len(msg.Receiver) == 0 || len(msg.Reference) == 0 {
		return sdk.ErrUnknownRequest("Sender and/or receiver and/or document reference cannot be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgShareDocument) GetSignBytes() []byte {
	return sdk.MustSortJSON(types2.msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgShareDocument) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
