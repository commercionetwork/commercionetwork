package types

/*
import (
	"commercio-network/types"
	"commercio-network/x/commerciodocs"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

// ----------------------------------
// --- Message StoreDocument
// ----------------------------------

func TestNewMsgStoreDocument(t *testing.T) {

	actual := types2.NewMsgStoreDocument(commerciodocs.owner, commerciodocs.ownerIdentity, commerciodocs.reference, commerciodocs.metadata)

	assert.Equal(t, commerciodocs.msgStore, actual)
}

func TestMsgStoreDocument_Route(t *testing.T) {
	key := "commerciodocs"

	actual := commerciodocs.msgStore.Route()

	assert.Equal(t, key, actual)
}

func TestMsgStoreDocument_Type(t *testing.T) {
	dType := "store_document"

	actual := commerciodocs.msgStore.Type()

	assert.Equal(t, dType, actual)
}

func TestMsgStoreDocument_ValidateBasic_AllFieldsCorrect(t *testing.T) {

	actual := commerciodocs.msgStore.ValidateBasic()

	assert.Nil(t, actual)
}

func TestMsgStoreDocument_ValidateBasic_EmptyOwnerField(t *testing.T) {
	message := types2.MsgStoreDocument{
		Owner:     sdk.AccAddress{},
		Identity:  commerciodocs.ownerIdentity,
		Reference: commerciodocs.reference,
		Metadata:  commerciodocs.metadata,
	}

	actual := message.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgStoreDocument_ValidateBasic_EmptyReference(t *testing.T) {
	message := types2.MsgStoreDocument{
		Owner:     commerciodocs.owner,
		Identity:  types.Did(""),
		Reference: commerciodocs.reference,
		Metadata:  commerciodocs.metadata,
	}

	actual := message.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgStoreDocument_GetSignBytes(t *testing.T) {

	expected := sdk.MustSortJSON(commerciodocs.input.cdc.MustMarshalJSON(commerciodocs.msgStore))

	actual := commerciodocs.msgStore.GetSignBytes()

	assert.Equal(t, expected, actual)

}

func TestMsgStoreDocument_GetSigners(t *testing.T) {

	expected := []sdk.AccAddress{commerciodocs.msgStore.Owner}

	actual := commerciodocs.msgStore.GetSigners()

	assert.Equal(t, expected, actual)
}

// ----------------------------------
// --- Message ShareDocument
// ----------------------------------

func TestNewMsgShareDocument(t *testing.T) {

	actual := types2.NewMsgShareDocument(commerciodocs.owner, commerciodocs.reference, commerciodocs.ownerIdentity, commerciodocs.recipient)

	assert.Equal(t, commerciodocs.msgShare, actual)
}

func TestMsgShareDocument_Route(t *testing.T) {
	key := "commerciodocs"

	actual := commerciodocs.msgShare.Route()

	assert.Equal(t, key, actual)
}

func TestMsgShareDocument_Type(t *testing.T) {
	dType := "share_document"

	actual := commerciodocs.msgShare.Type()

	assert.Equal(t, dType, actual)
}

func TestMsgShareDocument_ValidateBasic_AllFieldsCorrect(t *testing.T) {

	actual := commerciodocs.msgShare.ValidateBasic()

	assert.Nil(t, actual)
}

func TestMsgShareDocument_ValidateBasic_EmptyOwnerField(t *testing.T) {
	msg := types2.MsgShareDocument{
		Owner:     sdk.AccAddress{},
		Reference: commerciodocs.reference,
		Sender:    commerciodocs.ownerIdentity,
		Receiver:  commerciodocs.recipient,
	}

	actual := msg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgShareDocument_ValidateBasic_EmptyReference(t *testing.T) {
	msg := types2.MsgShareDocument{
		Owner:     commerciodocs.owner,
		Reference: "",
		Sender:    commerciodocs.ownerIdentity,
		Receiver:  commerciodocs.recipient,
	}

	actual := msg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgShareDocument_GetSignBytes(t *testing.T) {
	expected := sdk.MustSortJSON(commerciodocs.input.cdc.MustMarshalJSON(commerciodocs.msgShare))

	actual := commerciodocs.msgShare.GetSignBytes()

	assert.Equal(t, expected, actual)
}

func TestMsgShareDocument_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{commerciodocs.msgShare.Owner}

	actual := commerciodocs.msgShare.GetSigners()

	assert.Equal(t, expected, actual)
}

*/
