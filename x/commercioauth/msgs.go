package commercioauth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = "commercioauth"

// ----------------------------------
// --- CreateAccount
// ----------------------------------

type MsgCreateAccount struct {
	Signer   sdk.AccAddress `json:"signer"`
	Address  string         `json:"address"`
	KeyType  string         `json:"key_type"`
	KeyValue string         `json:"key_value"`
}

func NewMsgCreateAccount(signer sdk.AccAddress, address string, keyType string, keyValue string) MsgCreateAccount {
	return MsgCreateAccount{
		Address:  address,
		Signer:   signer,
		KeyType:  keyType,
		KeyValue: keyValue,
	}
}

// RouterKey Implements Msg.
func (msg MsgCreateAccount) Route() string { return RouterKey }

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
	return sdk.MustSortJSON(msgCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgCreateAccount) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
