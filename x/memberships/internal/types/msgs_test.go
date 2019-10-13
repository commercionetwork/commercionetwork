package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var TestBuyer, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestMembershipType = "bronze"

var msgSetId = MsgBuyMembership{
	Buyer:          TestBuyer,
	MembershipType: TestMembershipType,
}

// ----------------------------------
// --- SetIdentity
// ----------------------------------

func TestMsgAssignMembership_Route(t *testing.T) {
	assert.Equal(t, ModuleName, msgSetId.Route())
}

func TestMsgAssignMembership_Type(t *testing.T) {
	assert.Equal(t, MsgTypeBuyMembership, msgSetId.Type())
}

func TestMsgAssignMembership_ValidateBasic_AllFieldsCorrect(t *testing.T) {
	assert.Nil(t, msgSetId.ValidateBasic())
}

func TestMsgAssignMembership_ValidateBasic_InvalidBuyer(t *testing.T) {
	invalidMsg := NewMsgBuyMembership(TestMembershipType, nil)
	assert.Error(t, invalidMsg.ValidateBasic())
}

func TestMsgAssignMembership_ValidateBasic_InvalidTypes(t *testing.T) {
	types := []string{"green", "bronz", "slver", "gld", "blck"}
	for _, memType := range types {
		invalidMsg := NewMsgBuyMembership(memType, TestBuyer)
		assert.Error(t, invalidMsg.ValidateBasic())
	}
}

func TestMsgAssignMembership_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/MsgBuyMembership","value":{"buyer":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0","membership_type":"bronze"}}`
	assert.Equal(t, expected, string(msgSetId.GetSignBytes()))
}

func TestNewMsgSetIdentity_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgSetId.Buyer}
	assert.Equal(t, expected, msgSetId.GetSigners())
}
