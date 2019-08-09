package types

//import (
//	"github.com/commercionetwork/commercionetwork/types"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
////TEST VARS
//
//var testAddress = "cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"
//var testOwner, _ = sdk.AccAddressFromBech32(testAddress)
//var testOwnerIdentity = types.Did("newReader")
//var testReference = "testReference"
//var testMetadata = "testMetadata"
//var testRecipient = types.Did("testRecipient")
//
//var msgStore = MsgStoreDocument{
//	Identity:  testOwnerIdentity,
//	Reference: testReference,
//	Owner:     testOwner,
//	Metadata:  testMetadata,
//}
//
//var msgShare = MsgShareDocument{
//	Owner:     testOwner,
//	Sender:    testOwnerIdentity,
//	Receiver:  testRecipient,
//	Reference: testReference,
//}
//
//// ----------------------------------
//// --- Message StoreDocument
//// ----------------------------------
//
//func TestNewMsgStoreDocument(t *testing.T) {
//
//	actual := NewMsgStoreDocument(testOwner, testOwnerIdentity, testReference, testMetadata)
//
//	assert.Equal(t, msgStore, actual)
//}
//
//func TestMsgStoreDocument_Route(t *testing.T) {
//	key := "commerciodocs"
//
//	actual := msgStore.Route()
//
//	assert.Equal(t, key, actual)
//}
//
//func TestMsgStoreDocument_Type(t *testing.T) {
//	dType := "store_document"
//
//	actual := msgStore.Type()
//
//	assert.Equal(t, dType, actual)
//}
//
//func TestMsgStoreDocument_ValidateBasic_AllFieldsCorrect(t *testing.T) {
//
//	actual := msgStore.ValidateBasic()
//
//	assert.Nil(t, actual)
//}
//
//func TestMsgStoreDocument_ValidateBasic_EmptyOwnerField(t *testing.T) {
//	message := MsgStoreDocument{
//		Owner:     sdk.AccAddress{},
//		Identity:  testOwnerIdentity,
//		Reference: testReference,
//		Metadata:  testMetadata,
//	}
//
//	actual := message.ValidateBasic()
//
//	assert.Error(t, actual)
//}
//
//func TestMsgStoreDocument_ValidateBasic_EmptyReference(t *testing.T) {
//	message := MsgStoreDocument{
//		Owner:     testOwner,
//		Identity:  types.Did(""),
//		Reference: testReference,
//		Metadata:  testMetadata,
//	}
//
//	actual := message.ValidateBasic()
//
//	assert.Error(t, actual)
//}
//
//func TestMsgStoreDocument_GetSignBytes(t *testing.T) {
//
//	expected := `{"type":"commerciodocs/StoreDocument","value":{"identity":"newReader","metadata":"testMetadata","owner":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","reference":"testReference"}}`
//
//	actual := msgStore.GetSignBytes()
//
//	assert.Equal(t, expected, string(actual))
//
//}
//
//func TestMsgStoreDocument_GetSigners(t *testing.T) {
//
//	expected := []sdk.AccAddress{msgStore.Owner}
//
//	actual := msgStore.GetSigners()
//
//	assert.Equal(t, expected, actual)
//}
//
//// ----------------------------------
//// --- Message ShareDocument
//// ----------------------------------
//
//func TestNewMsgShareDocument(t *testing.T) {
//
//	actual := NewMsgShareDocument(testOwner, testReference, testOwnerIdentity, testRecipient)
//
//	assert.Equal(t, msgShare, actual)
//}
//
//func TestMsgShareDocument_Route(t *testing.T) {
//	key := "commerciodocs"
//
//	actual := msgShare.Route()
//
//	assert.Equal(t, key, actual)
//}
//
//func TestMsgShareDocument_Type(t *testing.T) {
//	dType := "share_document"
//
//	actual := msgShare.Type()
//
//	assert.Equal(t, dType, actual)
//}
//
//func TestMsgShareDocument_ValidateBasic_AllFieldsCorrect(t *testing.T) {
//
//	actual := msgShare.ValidateBasic()
//
//	assert.Nil(t, actual)
//}
//
//func TestMsgShareDocument_ValidateBasic_EmptyOwnerField(t *testing.T) {
//	msg := MsgShareDocument{
//		Owner:     sdk.AccAddress{},
//		Reference: testReference,
//		Sender:    testOwnerIdentity,
//		Receiver:  testRecipient,
//	}
//
//	actual := msg.ValidateBasic()
//
//	assert.Error(t, actual)
//}
//
//func TestMsgShareDocument_ValidateBasic_EmptyReference(t *testing.T) {
//	msg := MsgShareDocument{
//		Owner:     testOwner,
//		Reference: "",
//		Sender:    testOwnerIdentity,
//		Receiver:  testRecipient,
//	}
//
//	actual := msg.ValidateBasic()
//
//	assert.Error(t, actual)
//}
//
//func TestMsgShareDocument_GetSignBytes(t *testing.T) {
//	expected := `{"type":"commerciodocs/ShareDocument","value":{"owner":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","receiver":"testRecipient","reference":"testReference","sender":"newReader"}}`
//
//	actual := msgShare.GetSignBytes()
//
//	assert.Equal(t, expected, string(actual))
//}
//
//func TestMsgShareDocument_GetSigners(t *testing.T) {
//	expected := []sdk.AccAddress{msgShare.Owner}
//
//	actual := msgShare.GetSigners()
//
//	assert.Equal(t, expected, actual)
//}
