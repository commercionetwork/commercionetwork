package commercioid

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ---------------------------------------
// --- SetIdentity

type MsgSetIdentity struct {
	DID          string
	DDOReference string
	Owner        sdk.AccAddress
}

func NewMsgSetIdentity(did string, ddoReference string, owner sdk.AccAddress) MsgSetIdentity {
	return MsgSetIdentity{
		DID:          did,
		DDOReference: ddoReference,
		Owner:        owner,
	}
}

// Route Implements Msg.
func (msg MsgSetIdentity) Route() string { return "commercioid" }

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
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgSetIdentity) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Owner}
}
