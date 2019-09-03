package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

var TestSignerAddress, _ = sdk.AccAddressFromBech32("cosmos1u4zeemkg5pytfr7l3vn7uz3arlfppy5yyxeand")
var TestOwnerAddress, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestMembershipType = "green"

var msgSetId = MsgAssignMembership{
	Signer:         TestSignerAddress,
	User:           TestOwnerAddress,
	MembershipType: TestMembershipType,
}

// ----------------------------------
// --- SetIdentity
// ----------------------------------

func TestMsgAssignMembership_Route(t *testing.T) {
	assert.Equal(t, ModuleName, msgSetId.Route())
}

func TestMsgAssignMembership_Type(t *testing.T) {
	assert.Equal(t, "assign_membership", msgSetId.Type())
}

func TestMsgAssignMembership_ValidateBasic_AllFieldsCorrect(t *testing.T) {
	assert.Nil(t, msgSetId.ValidateBasic())
}

func TestMsgAssignMembership_ValidateBasic_InvalidSigner(t *testing.T) {
	invalidMsg := MsgAssignMembership{
		Signer:         sdk.AccAddress{},
		User:           TestOwnerAddress,
		MembershipType: TestMembershipType,
	}
	assert.Error(t, invalidMsg.ValidateBasic())
}

func TestMsgAssignMembership_ValidateBasic_InvalidUser(t *testing.T) {
	invalidMsg := MsgAssignMembership{
		Signer:         TestSignerAddress,
		User:           sdk.AccAddress{},
		MembershipType: TestMembershipType,
	}
	assert.Error(t, invalidMsg.ValidateBasic())
}

func TestMsgAssignMembership_ValidateBasic_InvalidTypes(t *testing.T) {
	types := []string{"gren", "bronz", "slver", "gld", "blck"}
	for _, memType := range types {
		invalidMsg := MsgAssignMembership{
			Signer:         TestSignerAddress,
			User:           sdk.AccAddress{},
			MembershipType: memType,
		}
		assert.Error(t, invalidMsg.ValidateBasic())
	}
}

func TestMsgAssignMembership_GetSignBytes(t *testing.T) {
	expected := `{"type":"commercio/AssignMembership","value":{"membership_type":"green","signer":"cosmos1u4zeemkg5pytfr7l3vn7uz3arlfppy5yyxeand","user":"cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0"}}`

	assert.Equal(t, expected, string(msgSetId.GetSignBytes()))
}

func TestNewMsgSetIdentity_GetSigners(t *testing.T) {
	expected := []sdk.AccAddress{msgSetId.Signer}
	assert.Equal(t, expected, msgSetId.GetSigners())
}
