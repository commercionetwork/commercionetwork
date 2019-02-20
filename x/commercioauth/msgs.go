package commercioauth

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const Route = "commercioauth"

// ----------------------------------
// --- CreateAccount
// ----------------------------------

type MsgCreateAccount struct {
	Signer   sdk.AccAddress
	Address  string
	KeyType  string
	KeyValue string
}

func NewMsgCreateAccount(signer sdk.AccAddress, address string, keyType string, keyValue string) MsgCreateAccount {
	return MsgCreateAccount{
		Address:  address,
		Signer:   signer,
		KeyType:  keyType,
		KeyValue: keyValue,
	}
}

// Route Implements Msg.
func (msg MsgCreateAccount) Route() string { return Route }

// Type Implements Msg.
func (msg MsgCreateAccount) Type() string { return "create_account" }

// ValidateBasic Implements Msg.
func (msg MsgCreateAccount) ValidateBasic() sdk.Error {
	if msg.Signer.Empty() {
		return sdk.ErrInvalidAddress(msg.Signer.String())
	}
	if len(msg.Address) == 0 {
		return sdk.ErrUnknownRequest("Account address cannot be empty")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgCreateAccount) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

// GetSigners Implements Msg.
func (msg MsgCreateAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
