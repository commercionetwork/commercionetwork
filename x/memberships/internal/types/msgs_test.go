package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// Test vars
var user, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var sender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

// ----------------------
// --- MsgInviteUser
// ----------------------

var msgInviteUser = NewMsgInviteUser(user, sender)

func TestMsgInviteUser_Route(t *testing.T) {
	assert.Equal(t, QuerierRoute, msgInviteUser.Route())
}

func TestMsgInviteUser_Type(t *testing.T) {
	assert.Equal(t, MsgTypeInviteUser, msgInviteUser.Type())
}

func TestMsgInviteUser_ValidateBasic_ValidMsg(t *testing.T) {
	assert.Nil(t, msgInviteUser.ValidateBasic())
}

func TestMsgInviteUser_ValidateBasic_MissingRecipient(t *testing.T) {
	msg := MsgInviteUser{Recipient: nil, Sender: sender}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgInviteUser_ValidateBasic_MissingSender(t *testing.T) {
	msg := MsgInviteUser{Recipient: user, Sender: nil}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgInviteUser_GetSignBytes(t *testing.T) {
	actual := msgInviteUser.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgInviteUser))
	assert.Equal(t, expected, actual)
}

func TestMsgInviteUser_GetSigners(t *testing.T) {
	actual := msgInviteUser.GetSigners()
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, msgInviteUser.Recipient, actual[0])
}

func TestMsgInviteUser_UnmarshalJson(t *testing.T) {
	json := `{"type":"commercio/MsgInviteUser","value":{"recipient":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","sender":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"}}`

	var msg MsgInviteUser
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	assert.Equal(t, user, msg.Recipient)
	assert.Equal(t, sender, msg.Sender)
}

// ---------------------------
// --- MsgSetUserVerified
// ---------------------------

var msgSetUserVerified = NewMsgSetUserVerified(user, tsp)

func TestMsgSetUserVerified_Route(t *testing.T) {
	assert.Equal(t, QuerierRoute, msgSetUserVerified.Route())
}

func TestMsgSetUserVerified_Type(t *testing.T) {
	assert.Equal(t, MsgTypeSetUserVerified, msgSetUserVerified.Type())
}

func TestMsgSetUserVerified_ValidateBasic_ValidMsg(t *testing.T) {
	assert.Nil(t, msgSetUserVerified.ValidateBasic())
}

func TestMsgSetUserVerified_ValidateBasic_MissingUser(t *testing.T) {
	msg := NewMsgSetUserVerified(nil, tsp)
	assert.Error(t, msg.ValidateBasic())
}

func TestMsgSetUserVerified_ValidateBasic_MissingVerifier(t *testing.T) {
	msg := NewMsgSetUserVerified(user, nil)
	assert.Error(t, msg.ValidateBasic())
}

func TestMsgSetUserVerified_GetSignBytes(t *testing.T) {
	json := `{"type":"commercio/MsgSetUserVerified","value":{"user":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","verifier":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}}`
	assert.Equal(t, json, string(msgSetUserVerified.GetSignBytes()))
}

func TestMsgSetUserVerified_GetSigners(t *testing.T) {
	actual := msgSetUserVerified.GetSigners()
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, msgSetUserVerified.Verifier, actual[0])
}

func TestMsgSetUserVerified_UnmarshalJson(t *testing.T) {
	json := `{"type":"commercio/MsgSetUserVerified","value":{"user":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","verifier":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}}`

	var msg MsgSetUserVerified
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	assert.Equal(t, user, msg.User)
	assert.Equal(t, tsp, msg.Verifier)
}

// --------------------------------
// --- MsgDepositIntoLiquidityPool
// --------------------------------

var amount = sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
var msgDepositIntoLiquidityPool = NewMsgDepositIntoLiquidityPool(amount, user)

func TestMsgDepositIntoLiquidityPool_Route(t *testing.T) {
	assert.Equal(t, QuerierRoute, msgDepositIntoLiquidityPool.Route())
}

func TestMsgDepositIntoLiquidityPool_Type(t *testing.T) {
	assert.Equal(t, MsgTypesDepositIntoLiquidityPool, msgDepositIntoLiquidityPool.Type())
}

func TestMsgDepositIntoLiquidityPool_ValidateBasic_ValidMsg(t *testing.T) {
	assert.Nil(t, msgDepositIntoLiquidityPool.ValidateBasic())
}

func TestMsgDepositIntoLiquidityPool_ValidateBasic_MissingDepositor(t *testing.T) {
	msg := MsgDepositIntoLiquidityPool{Depositor: nil, Amount: amount}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgDepositIntoLiquidityPool_ValidateBasic_MissingAmount(t *testing.T) {
	msg := MsgDepositIntoLiquidityPool{Depositor: user, Amount: nil}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgDepositIntoLiquidityPool_ValidateBasic_NegativeAmount(t *testing.T) {
	amount := sdk.Coins{sdk.Coin{Denom: "uatom", Amount: sdk.NewInt(-100)}}
	msg := MsgDepositIntoLiquidityPool{Depositor: user, Amount: amount}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgDepositIntoLiquidityPool_GetSignBytes(t *testing.T) {
	actual := msgDepositIntoLiquidityPool.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgDepositIntoLiquidityPool))
	assert.Equal(t, expected, actual)
}

func TestMsgDepositIntoLiquidityPool_GetSigners(t *testing.T) {
	actual := msgDepositIntoLiquidityPool.GetSigners()
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, msgDepositIntoLiquidityPool.Depositor, actual[0])
}

func TestMsgDepositIntoLiquidityPool_UnmarshalJson(t *testing.T) {
	json := `{"type":"commercio/MsgDepositIntoLiquidityPool","value":{"depositor":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","amount":[{"denom":"uatom","amount":"100"}]}}`

	var msg MsgDepositIntoLiquidityPool
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	assert.Equal(t, user, msg.Depositor)
	assert.Equal(t, amount, msg.Amount)
}

// --------------------------------
// --- MsgAddTsp
// --------------------------------

var government, _ = sdk.AccAddressFromBech32("cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg")
var tsp, _ = sdk.AccAddressFromBech32("cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee")
var msgAddTrustedSigner = MsgAddTsp{
	Government: government,
	Tsp:        tsp,
}

func TestMsgAddTrustedSigner_Route(t *testing.T) {
	assert.Equal(t, QuerierRoute, msgAddTrustedSigner.Route())
}

func TestMsgAddTrustedSigner_Type(t *testing.T) {
	assert.Equal(t, MsgTypeAddTsp, msgAddTrustedSigner.Type())
}

func TestMsgAddTrustedSigner_ValidateBasic_ValidMsg(t *testing.T) {
	assert.Nil(t, msgAddTrustedSigner.ValidateBasic())
}

func TestMsgAddTrustedSigner_ValidateBasic_MissingGovernment(t *testing.T) {
	msg := MsgAddTsp{Government: nil, Tsp: tsp}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgAddTrustedSigner_ValidateBasic_MissingSigner(t *testing.T) {
	msg := MsgAddTsp{Government: government, Tsp: nil}
	assert.NotNil(t, msg.ValidateBasic())
}

func TestMsgAddTrustedSigner_GetSignBytes(t *testing.T) {
	actual := msgAddTrustedSigner.GetSignBytes()
	expected := sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msgAddTrustedSigner))
	assert.Equal(t, expected, actual)
}

func TestMsgAddTrustedSigner_GetSigners(t *testing.T) {
	actual := msgAddTrustedSigner.GetSigners()
	assert.Equal(t, 1, len(actual))
	assert.Equal(t, msgAddTrustedSigner.Government, actual[0])
}

func TestMsgAddTrustedSigner_UnmarshalJson(t *testing.T) {
	json := `{"type":"commercio/MsgAddTsp","value":{"government":"cosmos1ct4ym78j7ksv9weyua4mzlksgwc9qq7q3wvhqg","tsp":"cosmos152eg5tmgsu65mcytrln4jk5pld7qd4us5pqdee"}}`

	var msg MsgAddTsp
	ModuleCdc.MustUnmarshalJSON([]byte(json), &msg)

	assert.Equal(t, tsp, msg.Tsp)
	assert.Equal(t, government, msg.Government)
}

// ---------------------------
// --- MsgBuyMemberships
// ---------------------------

var TestBuyer, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestMembershipType = "bronze"
var msgBuyMembership = NewMsgBuyMembership(TestMembershipType, TestBuyer)

func TestMsgBuyMembership_Route(t *testing.T) {
	assert.Equal(t, ModuleName, msgBuyMembership.Route())
}

func TestMsgBuyMembership_Type(t *testing.T) {
	assert.Equal(t, MsgTypeBuyMembership, msgBuyMembership.Type())
}

func TestMsgBuyMembership_ValidateBasic_AllFieldsCorrect(t *testing.T) {
	assert.Nil(t, msgBuyMembership.ValidateBasic())
}

func TestMsgBuyMembership_ValidateBasic_InvalidBuyer(t *testing.T) {
	invalidMsg := NewMsgBuyMembership(TestMembershipType, nil)
	assert.Error(t, invalidMsg.ValidateBasic())
}

func TestMsgBuyMembership_ValidateBasic_InvalidTypes(t *testing.T) {
	types := []string{"green", "bronz", "slver", "gld", "blck"}
	for _, memType := range types {
		invalidMsg := NewMsgBuyMembership(memType, TestBuyer)
		assert.Error(t, invalidMsg.ValidateBasic())
	}
}

func TestMsgBuyMembership_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgBuyMembership","value":{"buyer":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","membership_type":"bronze"}}`
	assert.Equal(t, expected, string(msgBuyMembership.GetSignBytes()))
}

func TestMsgBuyMembership_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgBuyMembership.Buyer}
	assert.Equal(t, expected, msgBuyMembership.GetSigners())
}
