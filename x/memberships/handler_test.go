package memberships

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/accreditations"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var stableCreditsDenom = "uccc"
var testInviteSender, _ = sdk.AccAddressFromBech32("cosmos1005d6lt2wcfuulfpegz656ychljt3k3u4hn5my")

var msgBuyMembership = NewMsgBuyMembership(TestMembershipType, TestUserAddress)

func TestHandler_ValidMsgAssignMembership(t *testing.T) {
	_, ctx, accKeeper, bankKeeper, k := TestSetup()

	// Setup everything
	invite := accreditations.Invite{Sender: testInviteSender, User: TestUserAddress, Rewarded: false}
	accKeeper.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, testInviteSender, MembershipTypeBronze)

	credentials := accreditations.Credential{Timestamp: accreditations.TestTimestamp, User: TestUserAddress, Verifier: testInviteSender}
	accKeeper.SaveCredential(ctx, credentials)

	_ = bankKeeper.SetCoins(ctx, TestUserAddress, sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 1000000000)))

	// Perform the call
	var handler = NewHandler(stableCreditsDenom, k, accKeeper, bankKeeper)
	res := handler(ctx, msgBuyMembership)
	require.True(t, res.IsOK())
}

func TestHandler_InvalidUnknownType(t *testing.T) {
	_, ctx, accKeeper, bankKeeper, k := TestSetup()

	var handler = NewHandler(stableCreditsDenom, k, accKeeper, bankKeeper)
	res := handler(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.Contains(t, res.Log, fmt.Sprintf("Unrecognized %s message type", ModuleName))
}

func TestHandler_InvalidMembershipType(t *testing.T) {
	_, ctx, accKeeper, bankKeeper, k := TestSetup()

	// Setup everything
	invite := accreditations.Invite{Sender: testInviteSender, User: TestUserAddress, Rewarded: false}
	accKeeper.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, testInviteSender, MembershipTypeBronze)

	credentials := accreditations.Credential{Timestamp: accreditations.TestTimestamp, User: TestUserAddress, Verifier: testInviteSender}
	accKeeper.SaveCredential(ctx, credentials)

	_ = bankKeeper.SetCoins(ctx, TestUserAddress, sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 1000000000)))

	var handler = NewHandler(stableCreditsDenom, k, accKeeper, bankKeeper)
	types := []string{"gren", "bronz", "slver", "gol", "blck"}

	for _, memType := range types {
		msg := NewMsgBuyMembership(memType, TestUserAddress)
		res := handler(ctx, msg)
		require.False(t, res.IsOK())
		require.Contains(t, res.Log, fmt.Sprintf("Invalid membership type: %s", memType))
	}
}

func TestHandler_MembershipUpgrade(t *testing.T) {
	_, ctx, accKeeper, bankKeeper, k := TestSetup()

	// Setup everything
	invite := accreditations.Invite{Sender: testInviteSender, User: TestUserAddress, Rewarded: false}
	accKeeper.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, testInviteSender, MembershipTypeBronze)

	credentials := accreditations.Credential{Timestamp: accreditations.TestTimestamp, User: TestUserAddress, Verifier: testInviteSender}
	accKeeper.SaveCredential(ctx, credentials)

	_ = bankKeeper.SetCoins(ctx, TestUserAddress, sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 100000000000000)))

	// Perform the calls
	var handler = NewHandler(stableCreditsDenom, k, accKeeper, bankKeeper)
	types := []string{"bronze", "silver", "gold", "black"}

	for index := 1; index < len(types); index++ {
		beforeType := types[index-1]
		memType := types[index]

		_, _ = k.AssignMembership(ctx, TestUserAddress, beforeType)

		msg := NewMsgBuyMembership(memType, TestUserAddress)
		res := handler(ctx, msg)
		require.True(t, res.IsOK())
	}
}

func TestHandler_InvalidMembershipUpgrade(t *testing.T) {
	_, ctx, accKeeper, bankKeeper, k := TestSetup()

	// Setup everything
	invite := accreditations.Invite{Sender: testInviteSender, User: TestUserAddress, Rewarded: false}
	accKeeper.SaveInvite(ctx, invite)

	_, _ = k.AssignMembership(ctx, testInviteSender, MembershipTypeBronze)

	credentials := accreditations.Credential{Timestamp: accreditations.TestTimestamp, User: TestUserAddress, Verifier: testInviteSender}
	accKeeper.SaveCredential(ctx, credentials)

	_ = bankKeeper.SetCoins(ctx, TestUserAddress, sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 1000000000)))

	// Performs the call
	var handler = NewHandler(stableCreditsDenom, k, accKeeper, bankKeeper)
	types := []string{"bronze", "silver", "gold", "black"}

	for _, memType := range types {
		_, _ = k.AssignMembership(ctx, TestUserAddress, memType)

		msg := NewMsgBuyMembership(memType, TestUserAddress)
		res := handler(ctx, msg)

		require.False(t, res.IsOK())
		require.Contains(t, res.Log, fmt.Sprintf("Cannot upgrade from %s membership to %s", memType, memType))
	}
}
