package keeper

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --------------------------
// --- handleMsgInviteUser
// --------------------------

func Test_handleMsgInviteUser_ExistingInvite(t *testing.T) {
	_, ctx, _, govK, k := GetTestInput()

	invite := types.Invite{Sender: TestUser2, User: TestUser, Rewarded: false}
	k.SaveInvite(ctx, invite)
	_, _ = k.AssignMembership(ctx, invite.Sender, types.MembershipTypeBronze)

	handler := NewHandler(k, govK)
	msg := types.NewMsgInviteUser(TestUser, TestUser2)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnknownRequest, res.Code)
	assert.Contains(t, res.Log, "has already been invited")
}

func Test_handleMsgInviteUser_NewInvite(t *testing.T) {
	_, ctx, _, govK, k := GetTestInput()

	_, _ = k.AssignMembership(ctx, TestUser2, types.MembershipTypeBronze)

	handler := NewHandler(k, govK)
	msg := types.NewMsgInviteUser(TestUser, TestUser2)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
}

// -----------------------------
// --- handleMsgSetUserVerified
// -----------------------------

func Test_handleMsgSetUserVerified_InvalidSigner(t *testing.T) {
	_, ctx, _, govK, k := GetTestInput()

	handler := NewHandler(k, govK)
	msg := types.NewMsgSetUserVerified(TestUser, TestTimestamp, nil)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeUnauthorized, res.Code)
}

func Test_handleMsgSetUserVerified_ExistingCredential(t *testing.T) {
	_, ctx, _, govK, k := GetTestInput()

	credential := types.Credential{User: TestUser, Verifier: TestTsp, Timestamp: TestTimestamp}
	k.SaveCredential(ctx, credential)

	handler := NewHandler(k, govK)
	msg := types.NewMsgSetUserVerified(TestUser, TestTimestamp, TestTsp)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
}

func Test_handleMsgSetUserVerified_NewCredential(t *testing.T) {
	_, ctx, _, govK, k := GetTestInput()

	handler := NewHandler(k, govK)
	msg := types.NewMsgSetUserVerified(TestUser, TestTimestamp, TestTsp)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
}

// -----------------------------
// --- handleMsgAddTrustedSigner
// -----------------------------

func Test_handleAddTrustedSigner_InvalidGovernment(t *testing.T) {
	_, ctx, _, govK, k := GetTestInput()

	government, _ := sdk.AccAddressFromBech32("cosmos15ne6fy8uukkyyf072qklkeleh2zf39k52mcg2f")

	err := govK.SetGovernmentAddress(ctx, TestUser)
	assert.Nil(t, err)

	handler := NewHandler(k, govK)
	msg := types.NewMsgAddTsp(TestTsp, government)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleAddTrustedSigner_ValidGovernment(t *testing.T) {
	_, ctx, _, govK, k := GetTestInput()

	government, _ := sdk.AccAddressFromBech32("cosmos15ne6fy8uukkyyf072qklkeleh2zf39k52mcg2f")

	err := govK.SetGovernmentAddress(ctx, government)
	assert.Nil(t, err)

	handler := NewHandler(k, govK)
	msg := types.NewMsgAddTsp(TestTsp, government)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
	assert.Equal(t, sdk.CodeOK, res.Code)
}

// ---------------------------------
// --- handleMsgBuyMembership
// ---------------------------------

var testInviteSender, _ = sdk.AccAddressFromBech32("cosmos1005d6lt2wcfuulfpegz656ychljt3k3u4hn5my")
var msgBuyMembership = types.NewMsgBuyMembership(TestMembershipType, TestUser)

func TestHandler_ValidMsgAssignMembership(t *testing.T) {
	_, ctx, bankK, govK, k := GetTestInput()

	// Setup everything
	invite := types.Invite{Sender: testInviteSender, User: TestUser, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, testInviteSender, types.MembershipTypeBronze)

	credentials := types.Credential{Timestamp: TestTimestamp, User: TestUser, Verifier: testInviteSender}
	k.SaveCredential(ctx, credentials)

	_ = bankK.SetCoins(ctx, TestUser, sdk.NewCoins(sdk.NewInt64Coin(TestStableCreditsDenom, 1000000000)))

	// Perform the call
	var handler = NewHandler(k, govK)
	res := handler(ctx, msgBuyMembership)
	require.True(t, res.IsOK())
}

func TestHandler_InvalidUnknownType(t *testing.T) {
	_, ctx, _, govK, k := GetTestInput()

	var handler = NewHandler(k, govK)
	res := handler(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.Contains(t, res.Log, fmt.Sprintf("Unrecognized %s message type", types.ModuleName))
}

func TestHandler_InvalidMembershipType(t *testing.T) {
	_, ctx, bankK, govK, k := GetTestInput()

	// Setup everything
	invite := types.Invite{Sender: testInviteSender, User: TestUser, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, testInviteSender, types.MembershipTypeBronze)

	credentials := types.Credential{Timestamp: TestTimestamp, User: TestUser, Verifier: testInviteSender}
	k.SaveCredential(ctx, credentials)

	_ = bankK.SetCoins(ctx, TestUser, sdk.NewCoins(sdk.NewInt64Coin(TestStableCreditsDenom, 1000000000)))

	var handler = NewHandler(k, govK)
	memTypes := []string{"gren", "bronz", "slver", "gol", "blck"}

	for _, memType := range memTypes {
		msg := types.NewMsgBuyMembership(memType, TestUser)
		res := handler(ctx, msg)
		require.False(t, res.IsOK())
		require.Contains(t, res.Log, fmt.Sprintf("Invalid membership type: %s", memType))
	}
}

func TestHandler_MembershipUpgrade(t *testing.T) {
	_, ctx, bankK, govK, k := GetTestInput()

	// Setup everything
	invite := types.Invite{Sender: testInviteSender, User: TestUser, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, testInviteSender, types.MembershipTypeBronze)

	credentials := types.Credential{Timestamp: TestTimestamp, User: TestUser, Verifier: testInviteSender}
	k.SaveCredential(ctx, credentials)

	_ = bankK.SetCoins(ctx, TestUser, sdk.NewCoins(sdk.NewInt64Coin(TestStableCreditsDenom, 100000000000000)))

	// Perform the calls
	var handler = NewHandler(k, govK)
	memTypes := []string{"bronze", "silver", "gold", "black"}

	for index := 1; index < len(memTypes); index++ {
		beforeType := memTypes[index-1]
		memType := memTypes[index]

		_, _ = k.AssignMembership(ctx, TestUser, beforeType)

		msg := types.NewMsgBuyMembership(memType, TestUser)
		res := handler(ctx, msg)
		require.True(t, res.IsOK())
	}
}

func TestHandler_InvalidMembershipUpgrade(t *testing.T) {
	_, ctx, bankK, govK, k := GetTestInput()

	// Setup everything
	invite := types.Invite{Sender: testInviteSender, User: TestUser, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, testInviteSender, types.MembershipTypeBronze)

	credentials := types.Credential{Timestamp: TestTimestamp, User: TestUser, Verifier: testInviteSender}
	k.SaveCredential(ctx, credentials)

	_ = bankK.SetCoins(ctx, TestUser, sdk.NewCoins(sdk.NewInt64Coin(TestStableCreditsDenom, 1000000000)))

	// Performs the call
	var handler = NewHandler(k, govK)
	memTypes := []string{"bronze", "silver", "gold", "black"}

	for _, memType := range memTypes {
		_, _ = k.AssignMembership(ctx, TestUser, memType)

		msg := types.NewMsgBuyMembership(memType, TestUser)
		res := handler(ctx, msg)

		require.False(t, res.IsOK())
		require.Contains(t, res.Log, fmt.Sprintf("Cannot upgrade from %s membership to %s", memType, memType))
	}
}
