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
var msgSetId = NewMsgSetIdentity(DidDocument{
	Context: "https://www.w3.org/2019/did/v1",
	Id:      testOwnerAddress,
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
			Id:           fmt.Sprintf("%s#keys-1", testOwnerAddress),
			Type:         "Secp256k1VerificationKey2018",
			Controller:   testOwnerAddress,
			PublicKeyHex: "02b97c30de767f084ce3080168ee293053ba33b235d7116a3263d29f1450936b71",
		},
		PubKey{
			Id:           fmt.Sprintf("%s#keys-2", testOwnerAddress),
			Type:         "RsaVerificationKey2018",
			Controller:   testOwnerAddress,
			PublicKeyHex: "04418834f5012c808a11830819f300d06092a864886f70d010101050003818d0030818902818100ccaf757e02ec9cfb3beddaa5fe8e9c24df033e9b60db7cb8e2981cb340321faf348731343c7ab2f4920ebd62c5c7617557f66219291ce4e95370381390252b080dfda319bb84808f04078737ab55f291a9024ef3b72aedcf26067d3cee2a470ed056f4e409b73dd6b4fddffa43dff02bf30a9de29357b606df6f0246be267a910203010001a",
		},
		PubKey{
			Id:           fmt.Sprintf("%s#keys-1", testOwnerAddress),
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
	actual := msgSetId.Route()
	assert.Equal(t, ModuleName, actual)
}

func TestMsgSetIdentity_Type(t *testing.T) {
	actual := msgSetId.Type()
	assert.Equal(t, MsgTypeSetIdentity, actual)
}

func TestMsgSetIdentity_ValidateBasic_AllFieldsCorrect(t *testing.T) {
	actual := msgSetId.ValidateBasic()
	assert.Nil(t, actual)
}

func TestMsgSetIdentity_ValidateBasic_InvalidAddress(t *testing.T) {
	invalidMsg := MsgSetIdentity{}
	actual := invalidMsg.ValidateBasic()
	assert.Error(t, actual)
}

func TestMsgSetIdentity_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgSetIdentity","value":{"@context":"https://www.w3.org/2019/did/v1","authentication":["cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1"],"id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","proof":{"created":"2016-02-08T16:02:20Z","creator":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1","signatureValue":"QNB13Y7Q9...1tzjn4w==","type":"LinkedDataSignature2015"},"publicKey":[{"controller":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1","publicKeyHex":"02b97c30de767f084ce3080168ee293053ba33b235d7116a3263d29f1450936b71","type":"Secp256k1VerificationKey2018"},{"controller":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-2","publicKeyHex":"04418834f5012c808a11830819f300d06092a864886f70d010101050003818d0030818902818100ccaf757e02ec9cfb3beddaa5fe8e9c24df033e9b60db7cb8e2981cb340321faf348731343c7ab2f4920ebd62c5c7617557f66219291ce4e95370381390252b080dfda319bb84808f04078737ab55f291a9024ef3b72aedcf26067d3cee2a470ed056f4e409b73dd6b4fddffa43dff02bf30a9de29357b606df6f0246be267a910203010001a","type":"RsaVerificationKey2018"},{"controller":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","id":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0#keys-1","publicKeyHex":"035AD6810A47F073553FF30D2FCC7E0D3B1C0B74B61A1AAA2582344037151E143A","type":"Secp256k1VerificationKey2018"}],"service":null}}`

	actual := msgSetId.GetSignBytes()
	assert.Equal(t, expected, string(actual))
}

func TestNewMsgSetIdentity_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgSetId.Id}
	actual := msgSetId.GetSigners()
	assert.Equal(t, expected, actual)
}
