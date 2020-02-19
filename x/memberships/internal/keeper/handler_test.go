package keeper_test

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/memberships/internal/keeper"
	"github.com/commercionetwork/commercionetwork/x/memberships/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"
)

func Test_handleMsgInviteUser(t *testing.T) {
	testData := []struct {
		name           string
		membershipType string
		invitee        sdk.AccAddress
		invitedUser    sdk.AccAddress
		existingInvite types.Invite
		existingUser   bool
		error          string
	}{
		{
			name:         "Invitee has no membership",
			invitee:      testUser2,
			invitedUser:  testUser,
			error:        "Cannot send an invitation without having a membership",
			existingUser: false,
		},
		{
			name:           "Existing invite returns error",
			membershipType: types.MembershipTypeBronze,
			invitee:        testUser2,
			invitedUser:    testUser,
			existingInvite: types.Invite{Sender: testUser2, User: testUser, Status: types.InviteStatusPending},
			error:          fmt.Sprintf("%s has already been invited", testUser),
			existingUser:   false,
		},
		{
			name:           "New invite is inserted properly",
			membershipType: types.MembershipTypeBronze,
			invitee:        testUser2,
			invitedUser:    testUser,
			existingUser:   false,
		},
		{
			name:           "existing user is not invited",
			membershipType: types.MembershipTypeBronze,
			invitee:        testUser2,
			invitedUser:    testUser,
			existingUser:   true,
			error:          "cannot invite existing user",
		},
	}

	for _, test := range testData {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, bk, govK, k := SetupTestInput()

			if test.existingUser {
				require.NoError(t,
					bk.SetCoins(ctx, test.invitedUser, sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1)))),
				)
			}

			if !test.existingInvite.Empty() {
				k.SaveInvite(ctx, test.existingInvite)
			}

			if len(test.membershipType) != 0 {
				_ = k.AssignMembership(ctx, test.invitee, test.membershipType)
			}

			handler := keeper.NewHandler(k, govK)
			msg := types.NewMsgInviteUser(test.invitee, test.invitedUser)
			res := handler(ctx, msg)

			if len(test.error) != 0 {
				require.False(t, res.IsOK())
				require.Contains(t, res.Log, test.error)
			} else {
				require.True(t, res.IsOK())
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

			handler := keeper.NewHandler(k, govK)
			msg := types.NewMsgSetUserVerified(test.user, test.tsp)
			res := handler(ctx, msg)

			if len(test.error) == 0 {
				require.True(t, res.IsOK())
			} else {
				require.False(t, res.IsOK())
				require.Contains(t, res.Log, test.error)
			}
		})
	}
}

func Test_handleAddTrustedSigner(t *testing.T) {
	government, _ := sdk.AccAddressFromBech32("cosmos15ne6fy8uukkyyf072qklkeleh2zf39k52mcg2f")
	tests := []struct {
		name   string
		tsp    sdk.AccAddress
		signer sdk.AccAddress
		error  string
	}{
		{
			name:   "Invalid government returns error",
			tsp:    testTsp,
			signer: nil,
			error:  "Invalid government address: ",
		},
		{
			name:   "Valid government adds successfully",
			tsp:    testTsp,
			signer: government,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, gk, k := SetupTestInput()
			err := gk.SetGovernmentAddress(ctx, government)
			require.NoError(t, err)

			handler := keeper.NewHandler(k, gk)
			msg := types.NewMsgAddTsp(test.tsp, test.signer)
			res := handler(ctx, msg)

			if len(test.error) != 0 {
				require.False(t, res.IsOK())
				require.Contains(t, res.Log, test.error)
			} else {
				require.True(t, res.IsOK())
			}
		})
	}
}

var testInviteSender, _ = sdk.AccAddressFromBech32("cosmos1005d6lt2wcfuulfpegz656ychljt3k3u4hn5my")

func TestHandler_ValidMsgAssignMembership(t *testing.T) {
	tests := []struct {
		name               string
		msg                sdk.Msg
		existingMembership string
		credential         types.Credential
		invite             types.Invite
		bankAmount         sdk.Coins
		error              string
	}{
		{
			name:       "Invalid membership type returns error",
			msg:        types.NewMsgBuyMembership("gren", testUser),
			invite:     types.NewInvite(testInviteSender, testUser, "bronze"),
			credential: types.NewCredential(testUser, testTsp, 0),
			bankAmount: sdk.NewCoins(sdk.NewInt64Coin(stableCreditDenom, 1000000000)),
			error:      "Invalid membership type: gren",
		},
		{
			name:  "Invalid message returns error",
			msg:   sdk.NewTestMsg(),
			error: "Unrecognized accreditations message type: Test message",
		},
		{
			name:       "Valid membership allows buying",
			msg:        types.NewMsgBuyMembership(types.MembershipTypeBronze, testUser),
			invite:     types.NewInvite(testInviteSender, testUser, "bronze"),
			credential: types.NewCredential(testUser, testTsp, 0),
			bankAmount: sdk.NewCoins(sdk.NewInt64Coin(stableCreditDenom, 1000000000)),
		},
		{
			name:  "Buying without invite returns error",
			msg:   types.NewMsgBuyMembership(types.MembershipTypeBronze, testUser),
			error: "Cannot buy a membership without being invited",
		},
		{
			name:   "Buying without verification returns error",
			msg:    types.NewMsgBuyMembership(types.MembershipTypeBronze, testUser),
			invite: types.NewInvite(testInviteSender, testUser, "bronze"),
			error:  "User has not yet been verified by a Trusted Service Provider",
		},
		{
			name:               "Valid upgrade works properly",
			existingMembership: types.MembershipTypeBronze,
			msg:                types.NewMsgBuyMembership(types.MembershipTypeSilver, testUser),
			invite:             types.NewInvite(testInviteSender, testUser, "bronze"),
			credential:         types.NewCredential(testUser, testTsp, 0),
			bankAmount:         sdk.NewCoins(sdk.NewInt64Coin(stableCreditDenom, 1000000000)),
		},
		{
			name:               "Invalid upgrade works properly",
			existingMembership: types.MembershipTypeSilver,
			msg:                types.NewMsgBuyMembership(types.MembershipTypeBronze, testUser),
			invite:             types.NewInvite(testInviteSender, testUser, "bronze"),
			credential:         types.NewCredential(testUser, testTsp, 0),
			bankAmount:         sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1000000000)),
			error:              "Cannot upgrade from silver membership to bronze",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, bk, gk, k := SetupTestInput()

			if !test.invite.Empty() {
				err := k.AssignMembership(ctx, test.invite.Sender, types.MembershipTypeBlack)
				require.NoError(t, err)
				k.SaveInvite(ctx, test.invite)
			}

			if !test.credential.Empty() {
				k.SaveCredential(ctx, test.credential)
			}

			if msg, ok := test.msg.(types.MsgBuyMembership); ok {
				k.SupplyKeeper.SetSupply(ctx, supply.NewSupply(test.bankAmount))
				err := bk.SetCoins(ctx, msg.Buyer, test.bankAmount)
				require.NoError(t, err)

				if len(test.existingMembership) != 0 {
					err = k.AssignMembership(ctx, msg.Buyer, test.existingMembership)
					require.NoError(t, err)
				}
			}

			handler := keeper.NewHandler(k, gk)
			res := handler(ctx, test.msg)

			if len(test.error) == 0 {
				require.True(t, res.IsOK())

				userAmt := bk.GetCoins(ctx, test.msg.GetSigners()[0])
				require.True(t, userAmt.IsAllLT(test.bankAmount))
			} else {
				require.False(t, res.IsOK())
				require.Contains(t, res.Log, test.error)
			}
		})
	}
}

func Test_handleMsgSetBlackMembership(t *testing.T) {
	tests := []struct {
		name        string
		message     types.MsgSetBlackMembership
		invite      *types.Invite
		verify      bool
		senderIsGov bool
		want        sdk.Result
	}{
		{
			"invited user gets black membership by government",
			types.MsgSetBlackMembership{
				GovernmentAddress: testInviteSender,
				Subscriber:        testUser,
			},
			&types.Invite{
				Sender: testInviteSender,
				User:   testUser,
				Status: types.InviteStatusPending,
			},
			true,
			true,
			sdk.Result{},
		},
		{
			"non-invited user doesn't get black membership by government",
			types.MsgSetBlackMembership{
				GovernmentAddress: testInviteSender,
				Subscriber:        testUser,
			},
			nil,
			true,
			true,
			sdk.Result{
				Code:      0x9,
				Codespace: "sdk",
				Log:       "{\"codespace\":\"sdk\",\"code\":9,\"message\":\"no membership invite found for user cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae\"}",
			},
		},
		{
			"invited, non-verified user doesn't get black membership by government",
			types.MsgSetBlackMembership{
				GovernmentAddress: testInviteSender,
				Subscriber:        testUser,
			},
			&types.Invite{
				Sender: testInviteSender,
				User:   testUser,
				Status: types.InviteStatusPending,
			},
			false,
			true,
			sdk.Result{
				Code:      0x6,
				Codespace: "sdk",
				Log:       "{\"codespace\":\"sdk\",\"code\":6,\"message\":\"User has not yet been verified by a Trusted Service Provider\"}",
			},
		},
		{
			"invited, verified user doesn't get black membership because sender is not government",
			types.MsgSetBlackMembership{
				GovernmentAddress: testInviteSender,
				Subscriber:        testUser,
			},
			&types.Invite{
				Sender: testInviteSender,
				User:   testUser,
				Status: types.InviteStatusPending,
			},
			true,
			false,
			sdk.Result{
				Code:      0x9,
				Codespace: "sdk",
				Log:       "{\"codespace\":\"sdk\",\"code\":9,\"message\":\"cosmos1005d6lt2wcfuulfpegz656ychljt3k3u4hn5my is not a government address\"}",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, gk, k := SetupTestInput()

			_ = k.AssignMembership(ctx, tt.message.GovernmentAddress, types.MembershipTypeBlack)

			if tt.invite != nil {
				k.SaveInvite(ctx, *tt.invite)
			}

			if tt.senderIsGov {
				require.NoError(t, gk.SetGovernmentAddress(ctx, tt.message.GovernmentAddress))
				k.AddTrustedServiceProvider(ctx, tt.message.GovernmentAddress)
			}

			if tt.verify {
				credential := types.NewCredential(tt.message.Subscriber, tt.message.GovernmentAddress, ctx.BlockHeight())
				k.SaveCredential(ctx, credential)
			}

			handler := keeper.NewHandler(k, gk)
			res := handler(ctx, tt.message)

			require.Equal(t, tt.want, res)
		})
	}
}
