package types_test

import (
	"fmt"
	"testing"
	"time"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/id/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var testZone, _ = time.LoadLocation("UTC")
var testTime = time.Date(2016, 2, 8, 16, 2, 20, 0, testZone)
var testOwnerAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var msgSetIdentity = types.NewMsgSetIdentity(types.DidDocument{
	Context: "https://www.w3.org/ns/did/v1",
	ID:      testOwnerAddress,
	Authentication: []string{
		fmt.Sprintf("%s#keys-1", testOwnerAddress),
	},
	Proof: types.Proof{
		Type:           "LinkedDataSignature2015",
		Created:        testTime,
		Creator:        fmt.Sprintf("%s#keys-1", testOwnerAddress),
		SignatureValue: "QNB13Y7Q9...1tzjn4w==",
	},
	PubKeys: types.PubKeys{
		types.PubKey{
			ID:           fmt.Sprintf("%s#keys-1", testOwnerAddress),
			Type:         "Secp256k1VerificationKey2018",
			Controller:   testOwnerAddress,
			PublicKeyHex: "02b97c30de767f084ce3080168ee293053ba33b235d7116a3263d29f1450936b71",
		},
		types.PubKey{
			ID:           fmt.Sprintf("%s#keys-2", testOwnerAddress),
			Type:         "RsaVerificationKey2018",
			Controller:   testOwnerAddress,
			PublicKeyHex: "04418834f5012c808a11830819f300d06092a864886f70d010101050003818d0030818902818100ccaf757e02ec9cfb3beddaa5fe8e9c24df033e9b60db7cb8e2981cb340321faf348731343c7ab2f4920ebd62c5c7617557f66219291ce4e95370381390252b080dfda319bb84808f04078737ab55f291a9024ef3b72aedcf26067d3cee2a470ed056f4e409b73dd6b4fddffa43dff02bf30a9de29357b606df6f0246be267a910203010001a",
		},
		types.PubKey{
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
	require.Equal(t, types.ModuleName, msgSetIdentity.Route())
}

func TestMsgSetIdentity_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeSetIdentity, msgSetIdentity.Type())
}

func TestMsgSetIdentity_ValidateBasic(t *testing.T) {
	testData := []struct {
		name          string
		message       types.MsgSetIdentity
		shouldBeValid bool
	}{
		{
			name:          "Valid message",
			message:       msgSetIdentity,
			shouldBeValid: true,
		},
		{
			name:          "Invalid address",
			message:       types.MsgSetIdentity{},
			shouldBeValid: false,
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.shouldBeValid {
				require.NoError(t, test.message.ValidateBasic())
			} else {
				require.Error(t, test.message.ValidateBasic())
			}
		})
	}
}

func TestMsgSetIdentity_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgSetIdentity","value":{"@context":"https://www.w3.org/ns/did/v1","authentication":["cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1"],"id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","proof":{"created":"2016-02-08T16:02:20Z","creator":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1","signatureValue":"QNB13Y7Q9...1tzjn4w==","type":"LinkedDataSignature2015"},"publicKey":[{"controller":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1","publicKeyHex":"02b97c30de767f084ce3080168ee293053ba33b235d7116a3263d29f1450936b71","type":"Secp256k1VerificationKey2018"},{"controller":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-2","publicKeyHex":"04418834f5012c808a11830819f300d06092a864886f70d010101050003818d0030818902818100ccaf757e02ec9cfb3beddaa5fe8e9c24df033e9b60db7cb8e2981cb340321faf348731343c7ab2f4920ebd62c5c7617557f66219291ce4e95370381390252b080dfda319bb84808f04078737ab55f291a9024ef3b72aedcf26067d3cee2a470ed056f4e409b73dd6b4fddffa43dff02bf30a9de29357b606df6f0246be267a910203010001a","type":"RsaVerificationKey2018"},{"controller":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1","publicKeyHex":"035AD6810A47F073553FF30D2FCC7E0D3B1C0B74B61A1AAA2582344037151E143A","type":"Secp256k1VerificationKey2018"}],"service":null}}`
	require.Equal(t, expected, string(msgSetIdentity.GetSignBytes()))
}

func TestMsgSetIdentity_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgSetIdentity.ID}
	require.Equal(t, expected, msgSetIdentity.GetSigners())
}

var requestSender, _ = sdk.AccAddressFromBech32("cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6")
var requestRecipient, _ = sdk.AccAddressFromBech32("cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa")
var editor, _ = sdk.AccAddressFromBech32("cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6")

// --------------------------
// --- MsgRequestDidPowerUp
// --------------------------

var msgRequestDidPowerUp = types.MsgRequestDidPowerUp{
	Claimant:      requestSender,
	Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
	Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
	EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
}

func TestMsgRequestDidPowerUp_Route(t *testing.T) {
	require.Equal(t, types.ModuleName, msgRequestDidPowerUp.Route())
}

func TestMsgRequestDidPowerUp_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeRequestDidPowerUp, msgRequestDidPowerUp.Type())
}

func TestMsgRequestDidPowerUp_ValidateBasic(t *testing.T) {
	claimant, _ := sdk.AccAddressFromBech32("cosmos1xt9nqxmermu64te9dr8rkjff8eax496hcasju7")
	amount := sdk.NewCoins(sdk.NewInt64Coin("uatom", 100))

	tests := []struct {
		name  string
		msg   types.MsgRequestDidPowerUp
		error error
	}{
		{
			name:  "Invalid claimant returns error",
			msg:   types.MsgRequestDidPowerUp{Claimant: sdk.AccAddress{}},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid claimant: "),
		},
		{
			name:  "Invalid amount returns error",
			msg:   types.MsgRequestDidPowerUp{Claimant: claimant, Amount: sdk.NewCoins()},
			error: sdkErr.Wrap(sdkErr.ErrInvalidCoins, "Power up amount not valid: "),
		},
		{
			name: "Valid message returns no error",
			msg:  msgRequestDidPowerUp,
		},
		{
			name:  "Invalid proof returns error",
			msg:   types.MsgRequestDidPowerUp{Claimant: claimant, Amount: amount, Proof: "230sd"},
			error: sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid proof: 230sd"),
		},
		{
			name:  "Invalid encryption key returns error",
			msg:   types.MsgRequestDidPowerUp{Claimant: claimant, Amount: amount, Proof: "617364", EncryptionKey: "1230xcv"},
			error: sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid encryption key value: 1230xcv"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.error != nil {
				require.Equal(t, test.error.Error(), test.msg.ValidateBasic().Error())
			} else {
				require.NoError(t, test.msg.ValidateBasic())
			}
		})
	}
}

func TestMsgRequestDidPowerUp_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgRequestDidPowerUp","value":{"amount":[{"amount":"100","denom":"uatom"}],"claimant":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","encryption_key":"333b68743231343b6833346832313468354a40617364617364","proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878"}}`
	require.Equal(t, expected, string(msgRequestDidPowerUp.GetSignBytes()))
}

func TestMsgRequestDidPowerUp_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgRequestDidPowerUp.Claimant}
	require.Equal(t, expected, msgRequestDidPowerUp.GetSigners())
}

func TestMsgRequestDidPowerUp_JSON(t *testing.T) {
	json := `{"type":"commercio/MsgRequestDidPowerUp","value":{"amount":[{"amount":"100","denom":"uatom"}],"claimant":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","encryption_key":"333b68743231343b6833346832313468354a40617364617364","proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878"}}`

	var actual types.MsgRequestDidPowerUp
	types.ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)

	expected := types.MsgRequestDidPowerUp{
		Claimant:      requestSender,
		Amount:        sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
		Proof:         "68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
		EncryptionKey: "333b68743231343b6833346832313468354a40617364617364",
	}
	require.Equal(t, expected, actual)
}

// ------------------------
// --- MsgPowerUpDid
// ------------------------

var msgPowerUpDid = types.MsgPowerUpDid{
	Recipient:           requestRecipient,
	Amount:              sdk.NewCoins(sdk.NewInt64Coin("uatom", 100)),
	ActivationReference: "333b68743231343b6833346832313468354a40617364617364",
	Signer:              requestSender,
}

func TestMsgPowerUpDid_Route(t *testing.T) {
	require.Equal(t, types.ModuleName, msgPowerUpDid.Route())
}

func TestMsgPowerUpDid_Type(t *testing.T) {
	require.Equal(t, types.MsgTypePowerUpDid, msgPowerUpDid.Type())
}

func TestMsgPowerUpDid_ValidateBasic(t *testing.T) {
	claimant, _ := sdk.AccAddressFromBech32("cosmos1xt9nqxmermu64te9dr8rkjff8eax496hcasju7")
	amount := sdk.NewCoins(sdk.NewInt64Coin("uatom", 100))

	tests := []struct {
		name  string
		msg   types.MsgPowerUpDid
		error error
	}{
		{
			name:  "Invalid recipient returns error",
			msg:   types.MsgPowerUpDid{Recipient: sdk.AccAddress{}},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid recipient address: "),
		},
		{
			name:  "Invalid amount returns error",
			msg:   types.MsgPowerUpDid{Recipient: claimant, Amount: sdk.NewCoins()},
			error: sdkErr.Wrap(sdkErr.ErrInvalidCoins, "Invalid power up amount: "),
		},
		{
			name: "Valid message returns no error",
			msg:  msgPowerUpDid,
		},
		{
			name:  "Invalid activation reference returns error",
			msg:   types.MsgPowerUpDid{Recipient: claimant, Amount: amount, ActivationReference: "230sd"},
			error: sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid activation_reference: 230sd"),
		},
		{
			name:  "Invalid signer returns error",
			msg:   types.MsgPowerUpDid{Recipient: claimant, Amount: amount, ActivationReference: "617364", Signer: sdk.AccAddress{}},
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid signer address: "),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.error != nil {
				require.Equal(t, test.error.Error(), test.msg.ValidateBasic().Error())
			} else {
				require.NoError(t, test.msg.ValidateBasic())
			}
		})
	}
}

func TestMsgPowerUpDid_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgPowerUpDid","value":{"activation_reference":"333b68743231343b6833346832313468354a40617364617364","amount":[{"amount":"100","denom":"uatom"}],"recipient":"cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa","signer":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6"}}`
	require.Equal(t, expected, string(msgPowerUpDid.GetSignBytes()))
}

func TestMsgPowerUpDid_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgPowerUpDid.Signer}
	require.Equal(t, expected, msgPowerUpDid.GetSigners())
}

func TestMsgPowerUpDid_JSON(t *testing.T) {
	json := `{"type":"commercio/MsgPowerUpDid","value":{"activation_reference":"333b68743231343b6833346832313468354a40617364617364","amount":[{"amount":"100","denom":"uatom"}],"recipient":"cosmos1yhd6h25ksupyezrajk30n7y99nrcgcnppj2haa","signer":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6"}}`

	var actual types.MsgPowerUpDid
	types.ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)
	require.Equal(t, msgPowerUpDid, actual)
}

// --------------------------
// --- MsgInvalidateDidPowerUpRequest
// --------------------------

var msgInvalidateDidPowerUpRequestStatus = types.NewMsgInvalidateDidPowerUpRequest(
	types.RequestStatus{
		Type:    "canceled",
		Message: "Don't want this anymore",
	},
	"68576d5a7134743777217a25432646294a404e635266556a586e327235753878",
	editor,
)

func TestNewMsgInvalidateDidPowerUpRequest_Route(t *testing.T) {
	require.Equal(t, types.ModuleName, msgInvalidateDidPowerUpRequestStatus.Route())
}

func TestNewMsgInvalidateDidPowerUpRequest_Type(t *testing.T) {
	require.Equal(t, types.MsgTypeInvalidateDidPowerUpRequest, msgInvalidateDidPowerUpRequestStatus.Type())
}

func TestNewMsgInvalidateDidPowerUpRequest_ValidateBasic(t *testing.T) {
	status := types.NewRequestStatus("type", "message")

	tests := []struct {
		name  string
		msg   types.MsgInvalidateDidPowerUpRequest
		error error
	}{
		{
			name:  "Empty editor returns error",
			msg:   types.NewMsgInvalidateDidPowerUpRequest(status, "", sdk.AccAddress{}),
			error: sdkErr.Wrap(sdkErr.ErrInvalidAddress, "Invalid editor address: "),
		},
		{
			name:  "Invalid power up proof returns error",
			msg:   types.NewMsgInvalidateDidPowerUpRequest(status, "", editor),
			error: sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid power_up_proof: "),
		},
		{
			name:  "Invalid status returns error",
			msg:   types.NewMsgInvalidateDidPowerUpRequest(types.RequestStatus{}, "31", editor),
			error: sdkErr.Wrap(sdkErr.ErrUnknownRequest, "Invalid status type: "),
		},
		{
			name: "Valid message returns no error",
			msg:  msgInvalidateDidPowerUpRequestStatus,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			if test.error != nil {
				require.Equal(t, test.error.Error(), test.msg.ValidateBasic().Error())
			} else {
				require.NoError(t, test.msg.ValidateBasic())
			}
		})
	}
}

func TestNewMsgInvalidateDidPowerUpRequest_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgInvalidateDidPowerUpRequest","value":{"editor":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","power_up_proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","status":{"message":"Don't want this anymore","type":"canceled"}}}`
	require.Equal(t, expected, string(msgInvalidateDidPowerUpRequestStatus.GetSignBytes()))
}

func TestNewMsgInvalidateDidPowerUpRequest_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgInvalidateDidPowerUpRequestStatus.Editor}
	require.Equal(t, expected, msgInvalidateDidPowerUpRequestStatus.GetSigners())
}

func TestNewMsgInvalidateDidPowerUpRequest_JSON(t *testing.T) {
	json := `{"type":"commercio/MsgInvalidateDidPowerUpRequest","value":{"editor":"cosmos187pz9tpycrhaes72c77p62zjh6p9zwt9amzpp6","power_up_proof":"68576d5a7134743777217a25432646294a404e635266556a586e327235753878","status":{"type":"canceled","message":"Don't want this anymore"}}}`

	var actual types.MsgInvalidateDidPowerUpRequest
	types.ModuleCdc.MustUnmarshalJSON([]byte(json), &actual)
	require.Equal(t, msgInvalidateDidPowerUpRequestStatus, actual)
}
