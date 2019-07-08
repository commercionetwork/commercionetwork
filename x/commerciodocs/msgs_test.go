package commerciodocs

import (
	"commercio-network/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

// ----------------------------------
// --- Message StoreDocument
// ----------------------------------

func TestNewMsgStoreDocument(t *testing.T) {

	actual := NewMsgStoreDocument(owner, ownerIdentity, reference, metadata)

	assert.Equal(t, msgStore, actual)
}

func TestMsgStoreDocument_Route(t *testing.T) {
	key := "commerciodocs"

	actual := msgStore.Route()

	assert.Equal(t, key, actual)
}

func TestMsgStoreDocument_Type(t *testing.T) {
	dType := "store_document"

	actual := msgStore.Type()

	assert.Equal(t, dType, actual)
}

func TestMsgStoreDocument_ValidateBasic_AllFieldsCorrect(t *testing.T) {

	actual := msgStore.ValidateBasic()

	assert.Nil(t, actual)
}

func TestMsgStoreDocument_ValidateBasic_EmptyOwnerField(t *testing.T) {
	message := MsgStoreDocument{
		Owner:     sdk.AccAddress{},
		Identity:  ownerIdentity,
		Reference: reference,
		Metadata:  metadata,
	}

	actual := message.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgStoreDocument_ValidateBasic_EmptyReference(t *testing.T) {
	message := MsgStoreDocument{
		Owner:     owner,
		Identity:  types.Did(""),
		Reference: reference,
		Metadata:  metadata,
	}

	actual := message.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgStoreDocument_GetSignBytes(t *testing.T) {

	expected := sdk.MustSortJSON(input.cdc.MustMarshalJSON(msgStore))

	actual := msgStore.GetSignBytes()

	assert.Equal(t, expected, actual)

}

func TestMsgStoreDocument_GetSigners(t *testing.T) {

	expected := []sdk.AccAddress{msgStore.Owner}

	actual := msgStore.GetSigners()

	assert.Equal(t, expected, actual)
}

// ----------------------------------
// --- Message ShareDocument
// ----------------------------------

func TestNewMsgShareDocument(t *testing.T) {

	actual := NewMsgShareDocument(owner, reference, ownerIdentity, recipient)

	assert.Equal(t, msgShare, actual)
}

func TestMsgShareDocument_Route(t *testing.T) {
	key := "commerciodocs"

	actual := msgShare.Route()

	assert.Equal(t, key, actual)
}

func TestMsgShareDocument_Type(t *testing.T) {
	dType := "share_document"

	actual := msgShare.Type()

	assert.Equal(t, dType, actual)
}

func TestMsgShareDocument_ValidateBasic_AllFieldsCorrect(t *testing.T) {

	actual := msgShare.ValidateBasic()

	assert.Nil(t, actual)
}

func TestMsgShareDocument_ValidateBasic_EmptyOwnerField(t *testing.T) {
	msg := MsgShareDocument{
		Owner:     sdk.AccAddress{},
		Reference: reference,
		Sender:    ownerIdentity,
		Receiver:  recipient,
	}

	actual := msg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgShareDocument_ValidateBasic_EmptyReference(t *testing.T) {
	msg := MsgShareDocument{
		Owner:     owner,
		Reference: "",
		Sender:    ownerIdentity,
		Receiver:  recipient,
	}

	actual := msg.ValidateBasic()

	assert.Error(t, actual)
}

func TestMsgShareDocument_GetSignBytes(t *testing.T) {
	expected := sdk.MustSortJSON(input.cdc.MustMarshalJSON(msgShare))

	actual := msgShare.GetSignBytes()

	assert.Equal(t, expected, actual)
}

func TestMsgShareDocument_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgShare.Owner}

	actual := msgShare.GetSigners()

	assert.Equal(t, expected, actual)
}
