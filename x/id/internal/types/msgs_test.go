package types

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var testZone, _ = time.LoadLocation("UTC")
var testTime = time.Date(2016, 2, 8, 16, 2, 20, 0, testZone)
var testOwnerAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var msgSetIdentity = NewMsgSetIdentity(DidDocument{
	Context: "https://www.w3.org/2019/did/v1",
	ID:      testOwnerAddress,
	Authentication: []string{
		fmt.Sprintf("%s#keys-1", testOwnerAddress),
	},
	Proof: Proof{
		Type:           "LinkedDataSignature2015",
		Created:        testTime,
		Creator:        fmt.Sprintf("%s#keys-1", testOwnerAddress),
		SignatureValue: "QNB13Y7Q9...1tzjn4w==",
	},
	PubKeys: PubKeys{
		PubKey{
			ID:           fmt.Sprintf("%s#keys-1", testOwnerAddress),
			Type:         "Secp256k1VerificationKey2018",
			Controller:   testOwnerAddress,
			PublicKeyHex: "02b97c30de767f084ce3080168ee293053ba33b235d7116a3263d29f1450936b71",
		},
		PubKey{
			ID:           fmt.Sprintf("%s#keys-2", testOwnerAddress),
			Type:         "RsaVerificationKey2018",
			Controller:   testOwnerAddress,
			PublicKeyHex: "04418834f5012c808a11830819f300d06092a864886f70d010101050003818d0030818902818100ccaf757e02ec9cfb3beddaa5fe8e9c24df033e9b60db7cb8e2981cb340321faf348731343c7ab2f4920ebd62c5c7617557f66219291ce4e95370381390252b080dfda319bb84808f04078737ab55f291a9024ef3b72aedcf26067d3cee2a470ed056f4e409b73dd6b4fddffa43dff02bf30a9de29357b606df6f0246be267a910203010001a",
		},
		PubKey{
			ID:           fmt.Sprintf("%s#keys-1", testOwnerAddress),
			Type:         "Secp256k1VerificationKey2018",
			Controller:   testOwnerAddress,
			PublicKeyHex: "035AD6810A47F073553FF30D2FCC7E0D3B1C0B74B61A1AAA2582344037151E143A",
		},
	},
})

// ----------------------------------
// --- SetIdentity
// ----------------------------------

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
	invalidMsg := MsgSetIdentity{}
	actual := invalidMsg.ValidateBasic()
	assert.Error(t, actual)
}

func TestMsgSetIdentity_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgSetIdentity","value":{"@context":"https://www.w3.org/2019/did/v1","authentication":["cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1"],"id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","proof":{"created":"2016-02-08T16:02:20Z","creator":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1","signatureValue":"QNB13Y7Q9...1tzjn4w==","type":"LinkedDataSignature2015"},"publicKey":[{"controller":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1","publicKeyHex":"02b97c30de767f084ce3080168ee293053ba33b235d7116a3263d29f1450936b71","type":"Secp256k1VerificationKey2018"},{"controller":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-2","publicKeyHex":"04418834f5012c808a11830819f300d06092a864886f70d010101050003818d0030818902818100ccaf757e02ec9cfb3beddaa5fe8e9c24df033e9b60db7cb8e2981cb340321faf348731343c7ab2f4920ebd62c5c7617557f66219291ce4e95370381390252b080dfda319bb84808f04078737ab55f291a9024ef3b72aedcf26067d3cee2a470ed056f4e409b73dd6b4fddffa43dff02bf30a9de29357b606df6f0246be267a910203010001a","type":"RsaVerificationKey2018"},{"controller":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1","publicKeyHex":"035AD6810A47F073553FF30D2FCC7E0D3B1C0B74B61A1AAA2582344037151E143A","type":"Secp256k1VerificationKey2018"}],"service":null}}`

	actual := msgSetIdentity.GetSignBytes()
	assert.Equal(t, expected, string(actual))
}

func TestMsgSetIdentity_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgSetIdentity.ID}
	actual := msgSetIdentity.GetSigners()
	assert.Equal(t, expected, actual)
}

// --------------------------
// --- MsgRequestDidDeposit
// --------------------------

var requestSender, _ = sdk.AccAddressFromBech32("cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6")
var requestRecipient, _ = sdk.AccAddressFromBech32("cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa")
var msgRequestDidDeposit = MsgRequestDidDeposit{
	FromAddress:   requestSender,
	Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
	Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
	EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
	Recipient:     requestRecipient,
}

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
	expected := `{"type":"commercio/MsgRequestDidDeposit","value":{"amount":[{"amount":"100","denom":"uatom"}],"encryption_key":"333b68743231343b6833346832313468354a40617364617364","from_address":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","recipient":"cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa"}}`

	actual := msgRequestDidDeposit.GetSignBytes()
	assert.Equal(t, expected, string(actual))
}

func TestMsgRequestDidDeposit_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgRequestDidDeposit.FromAddress}
	actual := msgRequestDidDeposit.GetSigners()
	assert.Equal(t, expected, actual)
}

func TestMsgRequestDidDeposit_JSON(t *testing.T) {
	json := `{"type":"commercio/MsgRequestDidDeposit","value":{"amount":[{"amount":"100","denom":"uatom"}],"encryption_key":"333b68743231343b6833346832313468354a40617364617364","from_address":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","recipient":"cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa"}}`

	var actual MsgRequestDidDeposit
	ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)

	expected := MsgRequestDidDeposit{
		FromAddress:   requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
		EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
		Recipient:     requestRecipient,
	}
	assert.Equal(t, expected, actual)
}

// --------------------------
// --- MsgInvalidateDidDepositRequest
// --------------------------

var editor, _ = sdk.AccAddressFromBech32("cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6")
var msgChangeDidDepositRequestStatus = NewMsgInvalidateDidDepositRequest(
	RequestStatus{
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
	assert.Equal(t, MsgTypeInvalidateDidDepositRequest, actual)
}

func TestMsgChangeDidDepositRequestStatus_ValidateBasic_AllFieldsCorrect(t *testing.T) {
	actual := msgChangeDidDepositRequestStatus.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgChangeDidDepositRequestStatus_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgInvalidateDidDepositRequest","value":{"deposit_proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","editor":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","status":{"message":"Don't want this anymore","type":"canceled"}}}`

	actual := msgChangeDidDepositRequestStatus.GetSignBytes()
	assert.Equal(t, expected, string(actual))
}

func TestMsgChangeDidDepositRequestStatus_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgChangeDidDepositRequestStatus.Editor}
	actual := msgChangeDidDepositRequestStatus.GetSigners()
	assert.Equal(t, expected, actual)
}

func TestMsgChangeDidDepositRequestStatus_JSON(t *testing.T) {
	json := `{"type":"commercio/MsgInvalidateDidDepositRequest","value":{"editor":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","deposit_proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","status":{"type":"canceled","message":"Don't want this anymore"}}}`

	var actual MsgInvalidateDidDepositRequest
	ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)
	assert.Equal(t, msgChangeDidDepositRequestStatus, actual)
}

// --------------------------
// --- MsgRequestDidPowerUp
// --------------------------

var msgRequestDidPowerUp = MsgRequestDidPowerUp{
	Claimant:      requestSender,
	Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
	Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
	EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
}

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
	expected := `{"type":"commercio/MsgRequestDidPowerUp","value":{"amount":[{"amount":"100","denom":"uatom"}],"claimant":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","encryption_key":"333b68743231343b6833346832313468354a40617364617364","proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878"}}`

	actual := msgRequestDidPowerUp.GetSignBytes()
	assert.Equal(t, expected, string(actual))
}

func TestMsgRequestDidPowerUp_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgRequestDidPowerUp.Claimant}
	actual := msgRequestDidPowerUp.GetSigners()
	assert.Equal(t, expected, actual)
}

func TestMsgRequestDidPowerUp_JSON(t *testing.T) {
	json := `{"type":"commercio/MsgRequestDidPowerUp","value":{"amount":[{"amount":"100","denom":"uatom"}],"claimant":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","encryption_key":"333b68743231343b6833346832313468354a40617364617364","proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878"}}`

	var actual MsgRequestDidPowerUp
	ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)

	expected := MsgRequestDidPowerUp{
		Claimant:      requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
		EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
	}
	assert.Equal(t, expected, actual)
}

// --------------------------
// --- MsgInvalidateDidPowerUpRequest
// --------------------------

var msgChangeDidPowerUpRequestStatus = NewMsgInvalidateDidPowerUpRequest(
	RequestStatus{
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
	assert.Equal(t, MsgTypeInvalidateDidPowerUpRequest, actual)
}

func TestMsgChangeDidPowerUpRequestStatus_ValidateBasic_AllFieldsCorrect(t *testing.T) {
	actual := msgChangeDidPowerUpRequestStatus.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgChangeDidPowerUpRequestStatus_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgInvalidateDidPowerUpRequest","value":{"editor":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","power_up_proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","status":{"message":"Don't want this anymore","type":"canceled"}}}`

	actual := msgChangeDidPowerUpRequestStatus.GetSignBytes()
	assert.Equal(t, expected, string(actual))
}

func TestMsgChangeDidPowerUpRequestStatus_JSON(t *testing.T) {
	json := `{"type":"commercio/MsgInvalidateDidPowerUpRequest","value":{"editor":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","power_up_proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","status":{"type":"canceled","message":"Don't want this anymore"}}}`

	var actual MsgInvalidateDidPowerUpRequest
	ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)
	assert.Equal(t, msgChangeDidPowerUpRequestStatus, actual)
}
