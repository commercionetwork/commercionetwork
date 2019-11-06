package keeper

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_handleMsgInviteUser(t *testing.T) {
	testData := []struct {
		name           string
		membershipType string
		invitee        sdk.AccAddress
		invitedUser    sdk.AccAddress
		existingInvite types.Invite
		error          string
	}{
		{
			name:        "Invitee has no membership",
			invitee:     TestUser2,
			invitedUser: testUser,
			error:       "Cannot send an invitation without having a membership",
		},
		{
			name:           "Existing invite returns error",
			membershipType: types.MembershipTypeBronze,
			invitee:        TestUser2,
			invitedUser:    testUser,
			existingInvite: types.Invite{Sender: TestUser2, User: testUser, Rewarded: false},
			error:          fmt.Sprintf("%s has already been invited", testUser),
		},
		{
			name:           "New invite is inserted properly",
			membershipType: types.MembershipTypeBronze,
			invitee:        TestUser2,
			invitedUser:    testUser,
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, govK, k := SetupTestInput()

			if !test.existingInvite.Empty() {
				k.SaveInvite(ctx, test.existingInvite)
			}

			if len(test.membershipType) != 0 {
				_, _ = k.AssignMembership(ctx, test.invitee, test.membershipType)
			}

			handler := NewHandler(k, govK)
			msg := types.NewMsgInviteUser(test.invitee, test.invitedUser)
			res := handler(ctx, msg)

			if len(test.error) != 0 {
				assert.False(t, res.IsOK())
				assert.Contains(t, res.Log, test.error)
			} else {
				assert.True(t, res.IsOK())
			}
		})
	}
}

func Test_handleMsgSetUserVerified(t *testing.T) {
	testData := []struct {
		name            string
		tsp             sdk.AccAddress
		user            sdk.AccAddress
		alreadyVerified bool
		error           string
	}{
		{
			name:            "Invalid signer returns error",
			tsp:             nil,
			user:            testUser,
			alreadyVerified: false,
			error:           " is not a valid TSP",
		},
		{
			name:            "Existing credential",
			tsp:             testTsp,
			user:            testUser,
			alreadyVerified: true,
		},
		{
			name:            "New credential",
			tsp:             testTsp,
			user:            testUser,
			alreadyVerified: false,
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, govK, k := SetupTestInput()

			if test.tsp != nil {
				k.AddTrustedServiceProvider(ctx, test.tsp)
			}

			if test.alreadyVerified {
				credential := types.NewCredential(test.user, test.tsp, ctx.BlockHeight())
				k.SaveCredential(ctx, credential)
			}

			handler := NewHandler(k, govK)
			msg := types.NewMsgSetUserVerified(test.user, test.tsp)
			res := handler(ctx, msg)

			if len(test.error) == 0 {
				assert.True(t, res.IsOK())
			} else {
				assert.False(t, res.IsOK())
				assert.Contains(t, res.Log, test.error)
			}
		})
	}
}

// -----------------------------
// --- handleMsgAddTrustedSigner
// -----------------------------

func Test_handleAddTrustedSigner_InvalidGovernment(t *testing.T) {
	ctx, _, govK, k := SetupTestInput()

	government, _ := sdk.AccAddressFromBech32("cosmos15ne6fy8uukkyyf072qklkeleh2zf39k52mcg2f")

	err := govK.SetGovernmentAddress(ctx, testUser)
	assert.Nil(t, err)

	handler := NewHandler(k, govK)
	msg := types.NewMsgAddTsp(testTsp, government)
	res := handler(ctx, msg)

	assert.False(t, res.IsOK())
	assert.Equal(t, sdk.CodeInvalidAddress, res.Code)
}

func Test_handleAddTrustedSigner_ValidGovernment(t *testing.T) {
	ctx, _, govK, k := SetupTestInput()

	government, _ := sdk.AccAddressFromBech32("cosmos15ne6fy8uukkyyf072qklkeleh2zf39k52mcg2f")

	err := govK.SetGovernmentAddress(ctx, government)
	assert.Nil(t, err)

	handler := NewHandler(k, govK)
	msg := types.NewMsgAddTsp(testTsp, government)
	res := handler(ctx, msg)

	assert.True(t, res.IsOK())
	assert.Equal(t, sdk.CodeOK, res.Code)
}

// ---------------------------------
// --- handleMsgBuyMembership
// ---------------------------------

var testInviteSender, _ = sdk.AccAddressFromBech32("cosmos1005d6lt2wcfuulfpegz656ychljt3k3u4hn5my")
var msgBuyMembership = types.NewMsgBuyMembership(testMembershipType, testUser)

func TestHandler_ValidMsgAssignMembership(t *testing.T) {
	ctx, bankK, govK, k := SetupTestInput()

	// Setup everything
	invite := types.Invite{Sender: testInviteSender, User: testUser, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, testInviteSender, types.MembershipTypeBronze)

	credentials := types.Credential{Timestamp: ctx.BlockHeight(), User: testUser, Verifier: testInviteSender}
	k.SaveCredential(ctx, credentials)

	creditsAmnt := sdk.NewCoins(sdk.NewInt64Coin(testStableCreditsDenom, 1000000000))
	_ = bankK.SetCoins(ctx, testUser, creditsAmnt)
	k.supplyKeeper.SetSupply(ctx, supply.NewSupply(creditsAmnt))

	// Perform the call
	var handler = NewHandler(k, govK)
	res := handler(ctx, msgBuyMembership)
	require.True(t, res.IsOK())
}

func TestHandler_InvalidUnknownType(t *testing.T) {
	ctx, _, govK, k := SetupTestInput()

	var handler = NewHandler(k, govK)
	res := handler(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.Contains(t, res.Log, fmt.Sprintf("Unrecognized %s message type", types.ModuleName))
}

func TestHandler_InvalidMembershipType(t *testing.T) {
	ctx, bankK, govK, k := SetupTestInput()

	// Setup everything
	invite := types.Invite{Sender: testInviteSender, User: testUser, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, testInviteSender, types.MembershipTypeBronze)

	credentials := types.Credential{Timestamp: ctx.BlockHeight(), User: testUser, Verifier: testInviteSender}
	k.SaveCredential(ctx, credentials)

	_ = bankK.SetCoins(ctx, testUser, sdk.NewCoins(sdk.NewInt64Coin(testStableCreditsDenom, 1000000000)))

	var handler = NewHandler(k, govK)
	memTypes := []string{"gren", "bronz", "slver", "gol", "blck"}

	for _, memType := range memTypes {
		msg := types.NewMsgBuyMembership(memType, testUser)
		res := handler(ctx, msg)
		require.False(t, res.IsOK())
		require.Contains(t, res.Log, fmt.Sprintf("Invalid membership type: %s", memType))
	}
}

func TestHandler_MembershipUpgrade(t *testing.T) {
	ctx, bankK, govK, k := SetupTestInput()

	// Setup everything
	invite := types.Invite{Sender: testInviteSender, User: testUser, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, testInviteSender, types.MembershipTypeBronze)

	credentials := types.Credential{Timestamp: ctx.BlockHeight(), User: testUser, Verifier: testInviteSender}
	k.SaveCredential(ctx, credentials)

	creditsAmnt := sdk.NewCoins(sdk.NewInt64Coin(testStableCreditsDenom, 100000000000000))
	_ = bankK.SetCoins(ctx, testUser, creditsAmnt)
	k.supplyKeeper.SetSupply(ctx, supply.NewSupply(creditsAmnt))

	// Perform the calls
	var handler = NewHandler(k, govK)
	memTypes := []string{"bronze", "silver", "gold", "black"}

	for index := 1; index < len(memTypes); index++ {
		beforeType := memTypes[index-1]
		memType := memTypes[index]

		_, _ = k.AssignMembership(ctx, testUser, beforeType)

		msg := types.NewMsgBuyMembership(memType, testUser)
		res := handler(ctx, msg)
		require.True(t, res.IsOK())
	}
}

func TestHandler_InvalidMembershipUpgrade(t *testing.T) {
	ctx, bankK, govK, k := SetupTestInput()

	// Setup everything
	invite := types.Invite{Sender: testInviteSender, User: testUser, Rewarded: false}
	k.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, testInviteSender, types.MembershipTypeBronze)

	credentials := types.Credential{Timestamp: ctx.BlockHeight(), User: testUser, Verifier: testInviteSender}
	k.SaveCredential(ctx, credentials)

	_ = bankK.SetCoins(ctx, testUser, sdk.NewCoins(sdk.NewInt64Coin(testStableCreditsDenom, 1000000000)))

	// Performs the call
	var handler = NewHandler(k, govK)
	memTypes := []string{"bronze", "silver", "gold", "black"}

	for _, memType := range memTypes {
		_, _ = k.AssignMembership(ctx, testUser, memType)

		msg := types.NewMsgBuyMembership(memType, testUser)
		res := handler(ctx, msg)

		require.False(t, res.IsOK())
		require.Contains(t, res.Log, fmt.Sprintf("Cannot upgrade from %s membership to %s", memType, memType))
	}
}
