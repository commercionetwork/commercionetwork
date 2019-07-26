package types

import (
	"commercio-network/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RouterKey is they name of the CommercioID module
const RouterKey = "commercioid"

// ----------------------------------
// --- SetIdentity
// ----------------------------------

type MsgSetIdentity struct {
	DID          types.Did      `json:"did"`
	DDOReference string         `json:"ddo_reference"`
	Owner        sdk.AccAddress `json:"owner"`
}

func NewMsgSetIdentity(did types.Did, ddoReference string, owner sdk.AccAddress) MsgSetIdentity {
	return MsgSetIdentity{
		DID:          did,
		DDOReference: ddoReference,
		Owner:        owner,
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
	if len(msg.DID) == 0 || len(msg.DDOReference) == 0 {
		return sdk.ErrUnknownRequest("DID and/or DDO reference cannot be empty")
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

// ----------------------------------
// --- CreateConnection
// ----------------------------------

type MsgCreateConnection struct {
	FirstUser  types.Did      `json:"first_user"`
	SecondUser types.Did      `json:"second_user"`
	Signer     sdk.AccAddress `json:"signer"`
}

func NewMsgCreateConnection(firstUser types.Did, secondUser types.Did, signer sdk.AccAddress) MsgCreateConnection {
	return MsgCreateConnection{
		FirstUser:  firstUser,
		SecondUser: secondUser,
		Signer:     signer,
	}
}

// Route Implements Msg.
func (msg MsgCreateConnection) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCreateConnection) Type() string { return "create_connection" }

// ValidateBasic Implements Msg.
func (msg MsgCreateConnection) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if len(msg.FirstUser) == 0 || len(msg.SecondUser) == 0 {
		return sdk.ErrUnknownRequest("First user and second user cannot be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateConnection) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgCreateConnection) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
