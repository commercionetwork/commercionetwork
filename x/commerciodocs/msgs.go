package commerciodocs

import (
	"commercio-network/types"
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const Route = "commerciodocs"

// ----------------------------------
// --- StoreDocument
// ----------------------------------

type MsgStoreDocument struct {
	Owner     sdk.AccAddress
	Identity  types.Did
	Reference string
	Metadata  string
}

func NewMsgStoreDocument(owner sdk.AccAddress, identity types.Did, reference string, metadata string) MsgStoreDocument {
	return MsgStoreDocument{
		Owner:     owner,
		Identity:  identity,
		Reference: reference,
		Metadata:  metadata,
	}
}

// Route Implements Msg.
func (msg MsgStoreDocument) Route() string { return Route }

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
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgStoreDocument) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}

// ----------------------------------
// --- ShareDocument
// ----------------------------------

type MsgShareDocument struct {
	Owner     sdk.AccAddress
	Sender    types.Did
	Receiver  types.Did
	Reference string
}

func NewMsgShareDocument(owner sdk.AccAddress, reference string, sender types.Did, receiver types.Did) MsgShareDocument {
	return MsgShareDocument{
		Owner:     owner,
		Sender:    sender,
		Receiver:  receiver,
		Reference: reference,
	}
}

// Route Implements Msg.
func (msg MsgShareDocument) Route() string { return Route }

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
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgShareDocument) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
