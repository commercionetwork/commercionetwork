package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// --------------------
// --- MsgSetIdentity
// --------------------

var owner, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var didDocument = DidDocument{
	Uri:         "https://test.example.com/did-document#1",
	ContentHash: "ebd7cfd95e67f57b4f5bcb3cf4554741a9b5ea666052106625a93a6f8881dc43",
}
var msgSetIdentity = NewMsgSetIdentity(owner, didDocument)

func TestMsgSetIdentity_Route(t *testing.T) {
	actual := msgSetIdentity.Route()
	assert.Equal(t, ModuleName, actual)
}

func TestMsgSetIdentity_Type(t *testing.T) {
	actual := msgSetIdentity.Type()
	assert.Equal(t, MsgTypeSetIdentity, actual)
}

func TestMsgSetIdentity_ValidateBasic_AllFieldsCorrect(t *testing.T) {
	actual := msgSetIdentity.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgSetIdentity_ValidateBasic_InvalidAddress(t *testing.T) {
	invalidMsg := NewMsgSetIdentity(sdk.AccAddress{}, msgSetIdentity.DidDocument)
	actual := invalidMsg.ValidateBasic()
	assert.Error(t, actual)
}

func TestMsgSetIdentity_ValidateBasic_InvalidDidDocumentUri(t *testing.T) {
	invalidDidDocument := DidDocument{
		Uri:         "",
		ContentHash: "ebd7cfd95e67f57b4f5bcb3cf4554741a9b5ea666052106625a93a6f8881dc43",
	}

	invalidMsg := NewMsgSetIdentity(owner, invalidDidDocument)
	actual := invalidMsg.ValidateBasic()
	assert.Error(t, actual)
}

func TestMsgSetIdentity_ValidateBasic_InvalidDidDocumentContentHash(t *testing.T) {
	invalidDidDocument := DidDocument{
		Uri:         msgSetIdentity.DidDocument.Uri,
		ContentHash: "ebd7cfd95e67f57b4f5bcb3cf4554741a9b5ea6662106625a93a6f8881dc43",
	}

	invalidMsg := NewMsgSetIdentity(owner, invalidDidDocument)
	actual := invalidMsg.ValidateBasic()
	assert.Error(t, actual)
}

func TestMsgSetIdentity_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgSetIdentity","value":{"did_document":{"content_hash":"ebd7cfd95e67f57b4f5bcb3cf4554741a9b5ea666052106625a93a6f8881dc43","uri":"https://test.example.com/did-document#1"},"owner":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"}}`

	actual := msgSetIdentity.GetSignBytes()
	assert.Equal(t, expected, string(actual))
}

func TestMsgSetIdentity_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgSetIdentity.Owner}
	actual := msgSetIdentity.GetSigners()
	assert.Equal(t, expected, actual)
}

// --------------------------
// --- MsgRequestDidDeposit
// --------------------------

var request = DidDepositRequest{
	FromAddress:   requestSender,
	Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
	Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
	EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
	Recipient:     requestRecipient,
}
var msgRequestDidDeposit = NewMsgRequestDidDeposit(request)

func TestMsgRequestDidDeposit_Route(t *testing.T) {
	actual := msgRequestDidDeposit.Route()
	assert.Equal(t, ModuleName, actual)
}

func TestMsgRequestDidDeposit_Type(t *testing.T) {
	actual := msgRequestDidDeposit.Type()
	assert.Equal(t, MsgTypeRequestDidDeposit, actual)
}

func TestMsgRequestDidDeposit_ValidateBasic_AllFieldsCorrect(t *testing.T) {
	actual := msgRequestDidDeposit.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgRequestDidDeposit_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgRequestDidDeposit","value":{"amount":[{"amount":"100","denom":"uatom"}],"encryption_key":"333b68743231343b6833346832313468354a40617364617364","from_address":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","recipient":"cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa","status":null}}`

	actual := msgRequestDidDeposit.GetSignBytes()
	assert.Equal(t, expected, string(actual))
}

func TestMsgRequestDidDeposit_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgRequestDidDeposit.FromAddress}
	actual := msgRequestDidDeposit.GetSigners()
	assert.Equal(t, expected, actual)
}

func TestMsgRequestDidDeposit_JSON_NullStatus(t *testing.T) {
	json := `{"type":"commercio/MsgRequestDidDeposit","value":{"amount":[{"amount":"100","denom":"uatom"}],"encryption_key":"333b68743231343b6833346832313468354a40617364617364","from_address":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","recipient":"cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa"}}`

	var actual MsgRequestDidDeposit
	ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)

	expected := NewMsgRequestDidDeposit(DidDepositRequest{
		FromAddress:   requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
		EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
		Recipient:     requestRecipient,
	})
	assert.Equal(t, expected, actual)
}

func TestMsgRequestDidDeposit_JSON_NonNullStatus(t *testing.T) {
	json := `{"type":"commercio/MsgRequestDidDeposit","value":{"amount":[{"amount":"100","denom":"uatom"}],"encryption_key":"333b68743231343b6833346832313468354a40617364617364","from_address":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","recipient":"cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa","status":{"type":"canceled","message":"Don't want this anymore"}}}`

	var actual MsgRequestDidDeposit
	ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)

	expected := NewMsgRequestDidDeposit(DidDepositRequest{
		FromAddress:   requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
		EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
		Recipient:     requestRecipient,
		Status: &DidDepositRequestStatus{
			Type:    "canceled",
			Message: "Don't want this anymore",
		},
	})
	assert.Equal(t, expected, actual)
}

// --------------------------
// --- MsgChangeDidDepositRequestStatus
// --------------------------

var editor, _ = sdk.AccAddressFromBech32("cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6")
var msgChangeDidDepositRequestStatus = NewMsgChangeDidDepositRequestStatus(
	DidDepositRequestStatus{
		Type:    "canceled",
		Message: "Don't want this anymore",
	},
	"68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
	editor,
)

func TestMsgChangeDidDepositRequestStatus_Route(t *testing.T) {
	actual := msgChangeDidDepositRequestStatus.Route()
	assert.Equal(t, ModuleName, actual)
}

func TestMsgChangeDidDepositRequestStatus_Type(t *testing.T) {
	actual := msgChangeDidDepositRequestStatus.Type()
	assert.Equal(t, MsgTypeChangeDidDepositRequestStatus, actual)
}

func TestMsgChangeDidDepositRequestStatus_ValidateBasic_AllFieldsCorrect(t *testing.T) {
	actual := msgChangeDidDepositRequestStatus.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgChangeDidDepositRequestStatus_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgChangeDidDepositRequestStatus","value":{"deposit_proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","editor":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","status":{"message":"Don't want this anymore","type":"canceled"}}}`

	actual := msgChangeDidDepositRequestStatus.GetSignBytes()
	assert.Equal(t, expected, string(actual))
}

func TestMsgChangeDidDepositRequestStatus_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgChangeDidDepositRequestStatus.Editor}
	actual := msgChangeDidDepositRequestStatus.GetSigners()
	assert.Equal(t, expected, actual)
}

func TestMsgChangeDidDepositRequestStatus_JSON(t *testing.T) {
	json := `{"type":"commercio/MsgChangeDidDepositRequestStatus","value":{"editor":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","deposit_proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","status":{"type":"canceled","message":"Don't want this anymore"}}}`

	var actual MsgChangeDidDepositRequestStatus
	ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)
	assert.Equal(t, msgChangeDidDepositRequestStatus, actual)
}

// --------------------------
// --- MsgRequestDidPowerUp
// --------------------------

var powerUpRequest = DidPowerUpRequest{
	Claimant:      requestSender,
	Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
	Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
	EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
}
var msgRequestDidPowerUp = NewMsgRequestDidPowerUp(powerUpRequest)

func TestMsgRequestDidPowerUp_Route(t *testing.T) {
	actual := msgRequestDidPowerUp.Route()
	assert.Equal(t, ModuleName, actual)
}

func TestMsgRequestDidPowerUp_Type(t *testing.T) {
	actual := msgRequestDidPowerUp.Type()
	assert.Equal(t, MsgTypeRequestDidPowerUp, actual)
}

func TestMsgRequestDidPowerUp_ValidateBasic_AllFieldsCorrect(t *testing.T) {
	actual := msgRequestDidPowerUp.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgRequestDidPowerUp_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgRequestDidPowerUp","value":{"amount":[{"amount":"100","denom":"uatom"}],"claimant":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","encryption_key":"333b68743231343b6833346832313468354a40617364617364","proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","status":null}}`

	actual := msgRequestDidPowerUp.GetSignBytes()
	assert.Equal(t, expected, string(actual))
}

func TestMsgRequestDidPowerUp_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgRequestDidPowerUp.Claimant}
	actual := msgRequestDidPowerUp.GetSigners()
	assert.Equal(t, expected, actual)
}

func TestMsgRequestDidPowerUp_JSON_NullStatus(t *testing.T) {
	json := `{"type":"commercio/MsgRequestDidPowerUp","value":{"amount":[{"amount":"100","denom":"uatom"}],"claimant":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","encryption_key":"333b68743231343b6833346832313468354a40617364617364","proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878"}}`

	var actual MsgRequestDidPowerUp
	ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)

	expected := NewMsgRequestDidPowerUp(DidPowerUpRequest{
		Claimant:      requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
		EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
	})
	assert.Equal(t, expected, actual)
}

func TestMsgRequestDidPowerUp_JSON_NonNullStatus(t *testing.T) {
	json := `{"type":"commercio/MsgRequestDidPowerUp","value":{"amount":[{"amount":"100","denom":"uatom"}],"claimant":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","encryption_key":"333b68743231343b6833346832313468354a40617364617364","proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","status":{"type":"canceled","message":"Don't want this anymore"}}}`

	var actual MsgRequestDidPowerUp
	ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)

	expected := NewMsgRequestDidPowerUp(DidPowerUpRequest{
		Claimant:      requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
		EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
		Status: &DidPowerUpRequestStatus{
			Type:    "canceled",
			Message: "Don't want this anymore",
		},
	})
	assert.Equal(t, expected, actual)
}

// --------------------------
// --- MsgChangeDidPowerUpRequestStatus
// --------------------------

var msgChangeDidPowerUpRequestStatus = NewMsgChangeDidPowerUpRequestStatus(
	DidPowerUpRequestStatus{
		Type:    "canceled",
		Message: "Don't want this anymore",
	},
	"68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
	editor,
)

func TestMsgChangeDidPowerUpRequestStatus_Route(t *testing.T) {
	actual := msgChangeDidPowerUpRequestStatus.Route()
	assert.Equal(t, ModuleName, actual)
}

func TestMsgChangeDidPowerUpRequestStatus_Type(t *testing.T) {
	actual := msgChangeDidPowerUpRequestStatus.Type()
	assert.Equal(t, MsgTypeChangeDidPowerUpRequestStatus, actual)
}

func TestMsgChangeDidPowerUpRequestStatus_ValidateBasic_AllFieldsCorrect(t *testing.T) {
	actual := msgChangeDidPowerUpRequestStatus.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgChangeDidPowerUpRequestStatus_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgChangeDidPowerUpRequestStatus","value":{"editor":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","power_up_proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","status":{"message":"Don't want this anymore","type":"canceled"}}}`

	actual := msgChangeDidPowerUpRequestStatus.GetSignBytes()
	assert.Equal(t, expected, string(actual))
}

func TestMsgChangeDidPowerUpRequestStatus_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgChangeDidPowerUpRequestStatus.Editor}
	actual := msgChangeDidPowerUpRequestStatus.GetSigners()
	assert.Equal(t, expected, actual)
}

func TestMsgChangeDidPowerUpRequestStatus_JSON(t *testing.T) {
	json := `{"type":"commercio/MsgChangeDidPowerUpRequestStatus","value":{"editor":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","power_up_proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","status":{"type":"canceled","message":"Don't want this anymore"}}}`

	var actual MsgChangeDidPowerUpRequestStatus
	ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)
	assert.Equal(t, msgChangeDidPowerUpRequestStatus, actual)
}
