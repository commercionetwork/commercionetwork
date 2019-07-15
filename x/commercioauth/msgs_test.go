package commercioauth

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"testing"
)

func (this MsgCreateAccount) isEqualTo(msg MsgCreateAccount) bool {
	return this.Signer.String() == msg.Signer.String() && this.Address == msg.Address &&
		this.KeyType == msg.KeyType && this.KeyValue == msg.KeyValue
}

var priv = ed25519.GenPrivKey()
var msg = MsgCreateAccount{
	Signer:   sdk.AccAddress(priv.PubKey().Address()),
	Address:  "test-address",
	KeyType:  "test-key-type",
	KeyValue: "test-key-value",
}

var testCdc = codec.New()

func TestNewMsgCreateAccount(t *testing.T) {

	signer := msg.Signer
	address := "test-address"
	keyType := "test-key-type"
	keyValue := "test-key-value"

	expected := msg

	actual := NewMsgCreateAccount(signer, address, keyType, keyValue)

	if !actual.isEqualTo(expected) {
		t.Errorf("actual created account %s is different than expected %v", actual, expected)
	}
}

func TestMsgCreateAccount_Route(t *testing.T) {
	expected := "commercioauth"

	actual := MsgCreateAccount.Route(msg)

	if actual != expected {
		t.Errorf("actual route %s is different than expected %v", actual, expected)
	}
}

func TestMsgCreateAccount_Type(t *testing.T) {
	expected := "create_account"

	actual := MsgCreateAccount.Type(msg)

	if actual != expected {
		t.Errorf("actual type %s is different than expected %v", actual, expected)
	}
}

func TestMsgCreateAccount_ValidateBasic(t *testing.T) {
	actual := MsgCreateAccount.ValidateBasic(msg)

	if actual != nil {
		t.Errorf("validation of %s is failed", actual)
	}
}

//Failing validation with empty signer
func TestMsgCreateAccount_ValidateBasic2(t *testing.T) {
	var errMsg = MsgCreateAccount{
		Signer:   sdk.AccAddress{},
		Address:  "test-address",
		KeyType:  "test-key-value",
		KeyValue: "test-key-type",
	}

	expected := sdk.ErrInvalidAddress(errMsg.Signer.String()).Error()

	actual := MsgCreateAccount.ValidateBasic(errMsg).Error()

	if actual != expected {
		t.Errorf("validation of %s is failed %v", actual, expected)
	}
}

//Failing validation with empty address
func TestMsgCreateAccount_ValidateBasic3(t *testing.T) {
	priv = ed25519.GenPrivKey()

	var errMsg = MsgCreateAccount{
		Signer:   sdk.AccAddress(priv.PubKey().Address()),
		Address:  "",
		KeyType:  "test-key-value",
		KeyValue: "test-key-type",
	}

	expected := sdk.ErrUnknownRequest("Account address cannot be empty").Error()

	actual := MsgCreateAccount.ValidateBasic(errMsg).Error()

	if actual != expected {
		t.Errorf("validation of %s is failed %v", actual, expected)
	}
}

func TestMsgCreateAccount_GetSignBytes(t *testing.T) {

	expected := sdk.MustSortJSON(testCdc.MustMarshalJSON(msg))

	actual := MsgCreateAccount.GetSignBytes(msg)

	for i := 0; i < len(actual) && i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Errorf("signed bytes of %s are different %v", actual, expected)
		}
	}
}

func TestMsgCreateAccount_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msg.Signer}

	actual := MsgCreateAccount.GetSigners(msg)

	for i := 0; i < len(actual) && i < len(expected); i++ {
		if !actual[i].Equals(expected[i]) {
			t.Errorf("signers of %s are different %v", actual, expected)
		}
	}
}
