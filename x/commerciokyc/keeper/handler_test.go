package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
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
				_ = k.AssignMembership(ctx, test.invitee, test.membershipType, testTsp, testHeight)
			}

			handler := NewHandler(k, govK)
			msg := types.NewMsgInviteUser(test.invitee, test.invitedUser)
			_, err := handler(ctx, msg)

			if len(test.error) != 0 {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.error)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TODO: Convert all error string in real error
func Test_handleAddTrustedSigner(t *testing.T) {
	government, _ := sdk.AccAddressFromBech32("cosmos15ne6fy8uukkyyf072qklkeleh2zf39k52mcg2f")
	tests := []struct {
		name       string
		tsp        sdk.AccAddress
		signer     sdk.AccAddress
		membership types.Membership
		error      string
	}{
		{
			name:   "Invalid government returns error",
			tsp:    testTsp,
			signer: nil,
			error:  "Invalid government address: ",
		},
		{
			name:   "Invalid tsp has no membership",
			tsp:    testTsp,
			signer: government,
			error:  "Tsp cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0 has no membership",
		},
		{
			name:       "Invalid tsp has no valid membership",
			tsp:        testTsp,
			membership: types.NewMembership(types.MembershipTypeBronze, testTsp, government, testHeight),
			signer:     government,
			error:      "Membership of Tsp cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0 is bronze but must be black",
		},
		{
			name:       "Valid government adds successfully",
			tsp:        testTsp,
			membership: types.NewMembership(types.MembershipTypeBlack, testTsp, government, testHeight),
			signer:     government,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, _, gk, k := SetupTestInput()
			err := gk.SetGovernmentAddress(ctx, government)
			require.NoError(t, err)
			if test.membership.MembershipType != "" {
				k.AssignMembership(ctx, test.membership.Owner, test.membership.MembershipType, test.membership.TspAddress, test.membership.ExpiryAt)
			}

			handler := NewHandler(k, gk)
			msg := types.NewMsgAddTsp(test.tsp, test.signer)
			_, err = handler(ctx, msg)

			if len(test.error) != 0 {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.error)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

var testInviteSender, _ = sdk.AccAddressFromBech32("cosmos1005d6lt2wcfuulfpegz656ychljt3k3u4hn5my")

// TODO: verify balance of tsp and government
func Test_handleMsgBuyMembership(t *testing.T) {
	government, _ := sdk.AccAddressFromBech32("cosmos15ne6fy8uukkyyf072qklkeleh2zf39k52mcg2f")
	tests := []struct {
		name               string
		msg                sdk.Msg
		existingMembership string
		invite             types.Invite
		bankAmount         sdk.Coins
		error              string
	}{
		{
			name:       "Invalid membership type returns error",
			msg:        types.NewMsgBuyMembership("gren", testUser, testTsp),
			invite:     types.NewInvite(testInviteSender, testUser, "bronze"),
			bankAmount: sdk.NewCoins(sdk.NewInt64Coin(stableCreditDenom, 1000000000)),
			error:      "Invalid membership type: gren",
		},
		{
			name: "Invalid message returns error",
			msg:  sdk.NewTestMsg(),
			// TODO change module name invalidate test
			error: "Unrecognized commerciokyc message type: Test message",
		},
		{
			name:       "Valid membership allows buying",
			msg:        types.NewMsgBuyMembership(types.MembershipTypeBronze, testUser, testTsp),
			invite:     types.NewInvite(testInviteSender, testUser, "bronze"),
			bankAmount: sdk.NewCoins(sdk.NewInt64Coin(stableCreditDenom, 1000000000)),
		},
		{
			name:  "Buying without invite returns error",
			msg:   types.NewMsgBuyMembership(types.MembershipTypeBronze, testUser, testTsp),
			error: "Cannot buy a membership without being invited",
		},
		{
			name: "Buying with invalid invite returns error",
			msg:  types.NewMsgBuyMembership(types.MembershipTypeBronze, testUser, testTsp),
			invite: types.Invite{
				Sender:           testInviteSender,
				SenderMembership: "bronze",
				User:             testUser,
				Status:           types.InviteStatusInvalid,
			},
			error: "invite for account cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae has been marked as invalid previously, cannot continue",
		},
		{
			name:               "Valid upgrade works properly",
			existingMembership: types.MembershipTypeBronze,
			msg:                types.NewMsgBuyMembership(types.MembershipTypeSilver, testUser, testTsp),
			invite:             types.NewInvite(testInviteSender, testUser, "bronze"),
			bankAmount:         sdk.NewCoins(sdk.NewInt64Coin(stableCreditDenom, 1000000000)),
		},
		{
			name:               "Valid downgrade works properly",
			existingMembership: types.MembershipTypeSilver,
			msg:                types.NewMsgBuyMembership(types.MembershipTypeBronze, testUser, testTsp),
			invite:             types.NewInvite(testInviteSender, testUser, "bronze"),
			bankAmount:         sdk.NewCoins(sdk.NewInt64Coin(stableCreditDenom, 1000000000)),
		},
		{
			name:       "Invalid buying memebership with diffrent denom",
			msg:        types.NewMsgBuyMembership(types.MembershipTypeBronze, testUser, testTsp),
			invite:     types.NewInvite(testInviteSender, testUser, "bronze"),
			bankAmount: sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1000000000)),
			error:      "insufficient funds: insufficient account funds",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, bk, gk, k := SetupTestInput()

			_ = gk.SetGovernmentAddress(ctx, government)
			govAmtBefore := bk.GetCoins(ctx, government)

			_ = bk.SetCoins(ctx, testTsp, test.bankAmount)

			if !test.invite.Empty() {
				err := k.AssignMembership(ctx, test.invite.Sender, types.MembershipTypeBlack, testTsp, testHeight)
				require.NoError(t, err)
				k.SaveInvite(ctx, test.invite)
			}

			if msg, ok := test.msg.(types.MsgBuyMembership); ok {
				k.SupplyKeeper.SetSupply(ctx, supply.NewSupply(test.bankAmount))
				err := bk.SetCoins(ctx, msg.Buyer, test.bankAmount)
				require.NoError(t, err)
				k.AddTrustedServiceProvider(ctx, msg.Tsp)

				if len(test.existingMembership) != 0 {
					err = k.AssignMembership(ctx, msg.Buyer, test.existingMembership, msg.Tsp, testHeight)
					require.NoError(t, err)
				}

			}

			handler := NewHandler(k, gk)
			_, err := handler(ctx, test.msg)

			if len(test.error) == 0 {
				require.NoError(t, err)

				//userAmt := bk.GetCoins(ctx, test.msg.GetSigners()[0])
				userAmt := bk.GetCoins(ctx, test.msg.GetSigners()[0])
				require.True(t, userAmt.IsAllLT(test.bankAmount))
				spent := test.bankAmount.Sub(userAmt)
				govAmtAfter := bk.GetCoins(ctx, government)
				require.True(t, govAmtAfter.IsEqual(govAmtBefore.Add(spent[0])))
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.error)
			}
		})
	}
}

func Test_handleMsgSetMembership(t *testing.T) {
	tests := []struct {
		name        string
		message     types.MsgSetMembership
		invite      *types.Invite
		senderIsGov bool
		want        string
	}{
		{
			"invited user gets black membership by government",
			types.MsgSetMembership{
				Government:    testInviteSender,
				Subscriber:    testUser,
				NewMembership: "black",
			},
			&types.Invite{
				Sender:           testInviteSender,
				User:             testUser,
				Status:           types.InviteStatusPending,
				SenderMembership: types.MembershipTypeBlack,
			},
			true,
			"",
		},
		{
			"non-invited user gets black membership by government",
			types.MsgSetMembership{
				Government:    testInviteSender,
				Subscriber:    testUser,
				NewMembership: types.MembershipTypeBlack,
			},
			nil,
			true,
			"",
		},
		// Verified was deteled
		/*
			{
				"invited, non-verified user doesn't get black membership by government",
				types.MsgSetMembership{
					Government: testInviteSender,
					Subscriber:        testUser,
					NewMembership:     types.MembershipTypeBlack,
				},
				&types.Invite{
					Sender:           testInviteSender,
					User:             testUser,
					Status:           types.InviteStatusPending,
					SenderMembership: types.MembershipTypeBlack,
				},
				true,
				"User has not yet been verified by a Trusted Service Provider",
			},*/
		{
			"invited, verified user doesn't get black membership because sender is not government",
			types.MsgSetMembership{
				Government:    testInviteSender,
				Subscriber:    testUser,
				NewMembership: "bronze",
			},
			&types.Invite{
				Sender:           testInviteSender,
				User:             testUser,
				Status:           types.InviteStatusPending,
				SenderMembership: types.MembershipTypeBlack,
			},
			false,
			"cosmos1005d6lt2wcfuulfpegz656ychljt3k3u4hn5my is not a government address",
		},
		{
			"Invalid membershipt type",
			types.MsgSetMembership{
				Government:    testInviteSender,
				Subscriber:    testUser,
				NewMembership: "grn",
			},
			&types.Invite{
				Sender:           testInviteSender,
				User:             testUser,
				Status:           types.InviteStatusPending,
				SenderMembership: types.MembershipTypeBlack,
			},
			true,
			"invalid membership type",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, gk, k := SetupTestInput()

			_ = k.AssignMembership(ctx, tt.message.Government, types.MembershipTypeBlack, testTsp, testHeight)

			if tt.invite != nil {
				k.SaveInvite(ctx, *tt.invite)
			}

			if tt.senderIsGov {
				require.NoError(t, gk.SetGovernmentAddress(ctx, tt.message.Government))
				k.AddTrustedServiceProvider(ctx, tt.message.Government)
			}

			handler := NewHandler(k, gk)
			_, err := handler(ctx, tt.message)

			if tt.want != "" {
				require.Error(t, err)
				if err != nil {
					require.Contains(t, err.Error(), tt.want)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_handleMsgRemoveMembership(t *testing.T) {
	tests := []struct {
		name        string
		message     types.MsgRemoveMembership
		senderIsGov bool
		want        string
	}{
		{
			"Valid: Remover user is government",
			types.MsgRemoveMembership{
				Government: testInviteSender,
				Subscriber: testUser,
			},
			true,
			"",
		},
		{
			"Invalid: Remover user is not government",
			types.MsgRemoveMembership{
				Government: testInviteSender,
				Subscriber: testUser,
			},
			false,
			"cosmos1005d6lt2wcfuulfpegz656ychljt3k3u4hn5my is not a government address",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, gk, k := SetupTestInput()

			_ = k.AssignMembership(ctx, tt.message.Subscriber, types.MembershipTypeGreen, tt.message.Government, testHeight)

			if tt.senderIsGov {
				require.NoError(t, gk.SetGovernmentAddress(ctx, tt.message.Government))
				k.AddTrustedServiceProvider(ctx, tt.message.Government)
			}

			handler := NewHandler(k, gk)
			_, err := handler(ctx, tt.message)

			if tt.want != "" {
				require.Error(t, err)
				if err != nil {
					require.Contains(t, err.Error(), tt.want)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_handleMsgRemoveTrustedSigner(t *testing.T) {
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

			handler := NewHandler(k, gk)
			msg := types.NewMsgRemoveTsp(test.tsp, test.signer)
			_, err = handler(ctx, msg)

			if len(test.error) != 0 {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.error)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Dummy test
func Test_handleMsgDepositIntoPool(t *testing.T) {
	tests := []struct {
		name       string
		signer     sdk.AccAddress
		msg        types.MsgDepositIntoLiquidityPool
		bankAmount sdk.Coins
		error      string
	}{
		/*{
			name:       "Valid deposit in uCCC",
			signer:     testUser,
			msg:        types.NewMsgDepositIntoLiquidityPool(depositStableCoin, testUser),
			bankAmount: sdk.NewCoins(sdk.NewInt64Coin(stableCreditDenom, 1000000000)),
		},*/
		{
			name:       "Valid deposit in uCOM",
			signer:     testUser,
			msg:        types.NewMsgDepositIntoLiquidityPool(depositTestCoin, testUser),
			bankAmount: sdk.NewCoins(sdk.NewInt64Coin(testDenom, 1000000000)),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			ctx, bk, gk, k := SetupTestInput()
			_ = bk.SetCoins(ctx, testUser, test.bankAmount)

			handler := NewHandler(k, gk)
			_, err := handler(ctx, test.msg)
			require.NoError(t, err)

		})
	}

}
