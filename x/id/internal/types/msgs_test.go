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
	assert.Equal(t, ModuleName, msgSetIdentity.Route())
}

func TestMsgSetIdentity_Type(t *testing.T) {
	assert.Equal(t, MsgTypeSetIdentity, msgSetIdentity.Type())
}

func TestMsgSetIdentity_ValidateBasic(t *testing.T) {
	testData := []struct {
		name          string
		message       MsgSetIdentity
		shouldBeValid bool
	}{
		{
			name:          "Valid message",
			message:       msgSetIdentity,
			shouldBeValid: true,
		},
		{
			name:          "Invalid address",
			message:       MsgSetIdentity{},
			shouldBeValid: false,
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.shouldBeValid {
				assert.NoError(t, test.message.ValidateBasic())
			} else {
				assert.Error(t, test.message.ValidateBasic())
			}
		})
	}
}

func TestMsgSetIdentity_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgSetIdentity","value":{"@context":"https://www.w3.org/2019/did/v1","authentication":["cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1"],"id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","proof":{"created":"2016-02-08T16:02:20Z","creator":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1","signatureValue":"QNB13Y7Q9...1tzjn4w==","type":"LinkedDataSignature2015"},"publicKey":[{"controller":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1","publicKeyHex":"02b97c30de767f084ce3080168ee293053ba33b235d7116a3263d29f1450936b71","type":"Secp256k1VerificationKey2018"},{"controller":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-2","publicKeyHex":"04418834f5012c808a11830819f300d06092a864886f70d010101050003818d0030818902818100ccaf757e02ec9cfb3beddaa5fe8e9c24df033e9b60db7cb8e2981cb340321faf348731343c7ab2f4920ebd62c5c7617557f66219291ce4e95370381390252b080dfda319bb84808f04078737ab55f291a9024ef3b72aedcf26067d3cee2a470ed056f4e409b73dd6b4fddffa43dff02bf30a9de29357b606df6f0246be267a910203010001a","type":"RsaVerificationKey2018"},{"controller":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1","publicKeyHex":"035AD6810A47F073553FF30D2FCC7E0D3B1C0B74B61A1AAA2582344037151E143A","type":"Secp256k1VerificationKey2018"}],"service":null}}`
	assert.Equal(t, expected, string(msgSetIdentity.GetSignBytes()))
}

func TestMsgSetIdentity_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgSetIdentity.ID}
	assert.Equal(t, expected, msgSetIdentity.GetSigners())
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
	assert.Equal(t, ModuleName, msgRequestDidDeposit.Route())
}

func TestMsgRequestDidDeposit_Type(t *testing.T) {
	assert.Equal(t, MsgTypeRequestDidDeposit, msgRequestDidDeposit.Type())
}

func TestMsgRequestDidDeposit_ValidateBasic_AllFieldsCorrect(t *testing.T) {
	assert.Nil(t, msgRequestDidDeposit.ValidateBasic())
}

func TestMsgRequestDidDeposit_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgRequestDidDeposit","value":{"amount":[{"amount":"100","denom":"uatom"}],"encryption_key":"333b68743231343b6833346832313468354a40617364617364","from_address":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","recipient":"cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa"}}`
	assert.Equal(t, expected, string(msgRequestDidDeposit.GetSignBytes()))
}

func TestMsgRequestDidDeposit_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgRequestDidDeposit.FromAddress}
	assert.Equal(t, expected, msgRequestDidDeposit.GetSigners())
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

// ------------------------
// --- MsgMoveDeposit
// ------------------------

var signer, _ = sdk.AccAddressFromBech32("cosmos1ejra5g9prtanzr3mjqj3suh5g6mffyyspemcm0")
var msgMoveDeposit = NewMsgMoveDeposit("333b68743231343b6833346832313468354a40617364617364", signer)

func TestMsgMoveDeposit_Route(t *testing.T) {
	assert.Equal(t, ModuleName, msgMoveDeposit.Route())
}

func TestMsgMoveDeposit_Type(t *testing.T) {
	assert.Equal(t, MsgTypeMoveDeposit, msgMoveDeposit.Type())
}

func TestMsgMoveDeposit_ValidateBasic(t *testing.T) {
	assert.Nil(t, msgMoveDeposit.ValidateBasic())
}

func TestMsgMoveDeposit_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgMoveDeposit","value":{"deposit_proof":"333b68743231343b6833346832313468354a40617364617364","signer":"cosmos1ejra5g9prtanzr3mjqj3suh5g6mffyyspemcm0"}}`
	assert.Equal(t, expected, string(msgMoveDeposit.GetSignBytes()))
}

func TestMsgMoveDeposit_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgMoveDeposit.Signer}
	assert.Equal(t, expected, msgMoveDeposit.GetSigners())
}

func TestMsgMoveDeposit_JSON(t *testing.T) {
	json := `{"type":"commercio/MsgMoveDeposit","value":{"deposit_proof":"333b68743231343b6833346832313468354a40617364617364","signer":"cosmos1ejra5g9prtanzr3mjqj3suh5g6mffyyspemcm0"}}`

	var actual MsgMoveDeposit
	ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)
	assert.Equal(t, msgMoveDeposit, actual)
}

// --------------------------
// --- MsgInvalidateDidDepositRequest
// --------------------------

var editor, _ = sdk.AccAddressFromBech32("cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6")
var msgInvalidateDidDepositRequestStatus = NewMsgInvalidateDidDepositRequest(
	RequestStatus{
		Type:    "canceled",
		Message: "Don't want this anymore",
	},
	"68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
	editor,
)

func TestMsgInvalidateDidDepositRequest_Route(t *testing.T) {
	assert.Equal(t, ModuleName, msgInvalidateDidDepositRequestStatus.Route())
}

func TestMsgInvalidateDidDepositRequest_Type(t *testing.T) {
	assert.Equal(t, MsgTypeInvalidateDidDepositRequest, msgInvalidateDidDepositRequestStatus.Type())
}

func TestMsgInvalidateDidDepositRequest_ValidateBasic(t *testing.T) {
	assert.Nil(t, msgInvalidateDidDepositRequestStatus.ValidateBasic())
}

func TestMsgInvalidateDidDepositRequest_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgInvalidateDidDepositRequest","value":{"deposit_proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","editor":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","status":{"message":"Don't want this anymore","type":"canceled"}}}`
	assert.Equal(t, expected, string(msgInvalidateDidDepositRequestStatus.GetSignBytes()))
}

func TestMsgInvalidateDidDepositRequest_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgInvalidateDidDepositRequestStatus.Editor}
	assert.Equal(t, expected, msgInvalidateDidDepositRequestStatus.GetSigners())
}

func TestMsgInvalidateDidDepositRequest_JSON(t *testing.T) {
	json := `{"type":"commercio/MsgInvalidateDidDepositRequest","value":{"editor":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","deposit_proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","status":{"type":"canceled","message":"Don't want this anymore"}}}`

	var actual MsgInvalidateDidDepositRequest
	ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)
	assert.Equal(t, msgInvalidateDidDepositRequestStatus, actual)
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
	assert.Equal(t, ModuleName, msgRequestDidPowerUp.Route())
}

func TestMsgRequestDidPowerUp_Type(t *testing.T) {
	assert.Equal(t, MsgTypeRequestDidPowerUp, msgRequestDidPowerUp.Type())
}

func TestMsgRequestDidPowerUp_ValidateBasic(t *testing.T) {
	assert.Nil(t, msgRequestDidPowerUp.ValidateBasic())
}

func TestMsgRequestDidPowerUp_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgRequestDidPowerUp","value":{"amount":[{"amount":"100","denom":"uatom"}],"claimant":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","encryption_key":"333b68743231343b6833346832313468354a40617364617364","proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878"}}`
	assert.Equal(t, expected, string(msgRequestDidPowerUp.GetSignBytes()))
}

func TestMsgRequestDidPowerUp_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgRequestDidPowerUp.Claimant}
	assert.Equal(t, expected, msgRequestDidPowerUp.GetSigners())
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

// ------------------------
// --- MsgPowerUpDid
// ------------------------

var msgPowerUpDid = MsgPowerUpDid{
	Recipient:           requestRecipient,
	Amount:              sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
	ActivationReference: "333b68743231343b6833346832313468354a40617364617364",
	Signer:              requestSender,
}

func TestMsgPowerUpDid_Route(t *testing.T) {
	assert.Equal(t, ModuleName, msgPowerUpDid.Route())
}

func TestMsgPowerUpDid_Type(t *testing.T) {
	assert.Equal(t, MsgTypePowerUpDid, msgPowerUpDid.Type())
}

func TestMsgPowerUpDid_ValidateBasic(t *testing.T) {
	assert.Nil(t, msgPowerUpDid.ValidateBasic())
}

func TestMsgPowerUpDid_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgPowerUpDid","value":{"activation_reference":"333b68743231343b6833346832313468354a40617364617364","amount":[{"amount":"100","denom":"uatom"}],"recipient":"cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa","signer":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6"}}`
	assert.Equal(t, expected, string(msgPowerUpDid.GetSignBytes()))
}

func TestMsgPowerUpDid_JSON(t *testing.T) {
	json := `{"type":"commercio/MsgPowerUpDid","value":{"activation_reference":"333b68743231343b6833346832313468354a40617364617364","amount":[{"amount":"100","denom":"uatom"}],"recipient":"cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa","signer":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6"}}`

	var actual MsgPowerUpDid
	ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)
	assert.Equal(t, msgPowerUpDid, actual)
}

// --------------------------
// --- MsgInvalidateDidPowerUpRequest
// --------------------------

var msgInvalidateDidPowerUpRequestStatus = NewMsgInvalidateDidPowerUpRequest(
	RequestStatus{
		Type:    "canceled",
		Message: "Don't want this anymore",
	},
	"68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
	editor,
)

func TestNewMsgInvalidateDidPowerUpRequest_Route(t *testing.T) {
	assert.Equal(t, ModuleName, msgInvalidateDidPowerUpRequestStatus.Route())
}

func TestNewMsgInvalidateDidPowerUpRequest_Type(t *testing.T) {
	assert.Equal(t, MsgTypeInvalidateDidPowerUpRequest, msgInvalidateDidPowerUpRequestStatus.Type())
}

func TestNewMsgInvalidateDidPowerUpRequest_ValidateBasic(t *testing.T) {
	assert.Nil(t, msgInvalidateDidPowerUpRequestStatus.ValidateBasic())
}

func TestNewMsgInvalidateDidPowerUpRequest_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgInvalidateDidPowerUpRequest","value":{"editor":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","power_up_proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","status":{"message":"Don't want this anymore","type":"canceled"}}}`
	assert.Equal(t, expected, string(msgInvalidateDidPowerUpRequestStatus.GetSignBytes()))
}

func TestNewMsgInvalidateDidPowerUpRequest_JSON(t *testing.T) {
	json := `{"type":"commercio/MsgInvalidateDidPowerUpRequest","value":{"editor":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","power_up_proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","status":{"type":"canceled","message":"Don't want this anymore"}}}`

	var actual MsgInvalidateDidPowerUpRequest
	ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)
	assert.Equal(t, msgInvalidateDidPowerUpRequestStatus, actual)
}
