package keeper

import (
	"reflect"
	"testing"
	"time"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func Test_msgServer_BuyMembership(t *testing.T) {

	type args struct {
		msg    *types.MsgBuyMembership
		invite types.Invite
	}
	tests := []struct {
		name string
		args args
		wantErr string
	}{
		{
			name: "Invalid membership type returns error",
			args: args{
				msg:    types.NewMsgBuyMembership("gren", testUser, testTsp),
				invite: types.NewInvite(testInviteSender, testUser, "bronze"),
			},
			wantErr: sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Invalid membership type: %s", testInvalidMembership)).Error(),
		},
		{
			name: "Downgrade from black membership returns error",
			args: args{
				msg:    types.NewMsgBuyMembership("silver", testUser, testTsp),
				invite: types.NewInvite(testInviteSender, testUser, "silver"),
			},
			wantErr: sdkErr.Wrap(sdkErr.ErrUnauthorized, "cannot downgrade from Black membership").Error(),
		},
		{
			name: "Buying without invite returns error",
			args: args{
				msg:    types.NewMsgBuyMembership("bronze", testUser, testTsp),
				invite: types.Invite{},
			},
			wantErr: sdkErr.Wrap(sdkErr.ErrUnauthorized, "Cannot buy a membership without being invited").Error(),
		},
		{
			name: "Valid membership allows buying",
			args: args{
				msg:    types.NewMsgBuyMembership("silver", testUser, testTsp),
				invite: types.NewInvite(testInviteSender, testUser, "silver"),
			},
			wantErr: "",
		},
		{
			name: "Valid upgrade works properly",
			args: args{
				msg:    types.NewMsgBuyMembership("gold", testUser, testTsp),
				invite: types.NewInvite(testInviteSender, testUser, "gold"),
			},
			wantErr: "",
		},
		{
			name: "Valid downgrade works properly",
			args: args{
				msg:    types.NewMsgBuyMembership("green", testUser, testTsp),
				invite: types.NewInvite(testInviteSender, testUser, "green"),
			},
			wantErr: "",
		},
		{
			name: "Normal user buying for other user as tsp returns error",
			args: args{
				msg:    types.NewMsgBuyMembership("silver", testUser2, testUser),
				invite: types.NewInvite(testInviteSender, testUser2, "silver"),
			},
			wantErr: sdkErr.Wrap(sdkErr.ErrUnauthorized, "since you are not a tsp you can buy membership only for yourself").Error(),
		},
		{
			name: "Normal user buy membership for himself",
			args: args{
				msg:    types.NewMsgBuyMembership("green", testUser, testUser),
				invite: types.NewInvite(testInviteSender, testUser, "green"),
			},
			wantErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, bk, _, k := SetupTestInput()

			msg := NewMsgServerImpl(k)
			_ = bk
			_ = msg
			_ = ctx
			k.AddTrustedServiceProvider(ctx, testTsp)
			if !tt.args.invite.Empty() {
				inviteSender, _ := sdk.AccAddressFromBech32(tt.args.invite.Sender)
				//assign membership to inviter
				err := k.AssignMembership(ctx, inviteSender, types.MembershipTypeBlack, testTsp, testExpiration)
				require.NoError(t, err)
				k.SaveInvite(ctx, tt.args.invite)
				//send coins to tsp
				msgTsp, _ := sdk.AccAddressFromBech32(tt.args.msg.Tsp)
				membershipPrice := membershipCosts[tt.args.msg.MembershipType] * 1000000
				membershipCost := sdk.NewCoins(sdk.NewInt64Coin(stableCreditDenom, membershipPrice))
				bk.MintCoins(ctx, types.ModuleName, membershipCost)
				bk.SendCoinsFromModuleToAccount(ctx,types.ModuleName, msgTsp, membershipCost)
			}
			
			name := "Downgrade from black membership returns error"
			if tt.name == name {
				err := k.AssignMembership(ctx, testUser, types.MembershipTypeBlack, testTsp, testExpiration)
				require.NoError(t, err)
			}

			name = "Valid upgrade works properly"
			if tt.name == name {
				err := k.AssignMembership(ctx, testUser, types.MembershipTypeSilver, testTsp, testExpiration)
				require.NoError(t, err)
			}

			name = "Valid downgrade works properly"
			if tt.name == name {
				err := k.AssignMembership(ctx, testUser, types.MembershipTypeSilver, testTsp, testExpiration)
				require.NoError(t, err)
			}

			_, err := msg.BuyMembership(sdk.WrapSDKContext(ctx), tt.args.msg)
			if err != nil {
				require.Equal(t, tt.wantErr, err.Error())
				return
			}
		})
	}
}

func Test_msgServer_RemoveMembership(t *testing.T) {
	tests := []struct {
		name               string
		msg                *types.MsgRemoveMembership
		membershipToCreate types.Membership
		want               *types.MsgRemoveMembershipResponse
		wantErr            bool
	}{
		{
			name: "Remove membership is not from government doesn't work",
			msg:  types.NewMsgRemoveMembership(testInviteSender.String(), testUser.String()),
			membershipToCreate: types.NewMembership(
				types.MembershipTypeBronze,
				testUser,
				testUser3,
				time.Now(),
			),
			want:    nil,
			wantErr: true,
		},
		{
			name: "Remove membership correctly works",
			msg:  types.NewMsgRemoveMembership(testUser3.String(), testUser.String()),
			membershipToCreate: types.NewMembership(
				types.MembershipTypeBronze,
				testUser,
				testUser3,
				time.Now().Add(time.Hour*1),
			),
			want: &types.MsgRemoveMembershipResponse{
				Subscriber: testUser.String(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()

			msg := NewMsgServerImpl(k)
			user, _ := sdk.AccAddressFromBech32(tt.membershipToCreate.Owner)
			tsp, _ := sdk.AccAddressFromBech32(tt.membershipToCreate.TspAddress)
			k.AssignMembership(ctx, user, tt.membershipToCreate.MembershipType, tsp, *tt.membershipToCreate.ExpiryAt)

			got, err := msg.RemoveMembership(sdk.WrapSDKContext(ctx), tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgServer.RemoveMembership() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgServer.RemoveMembership() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_msgServer_SetMembership(t *testing.T) {

	tests := []struct {
		name        string
		msg         *types.MsgSetMembership
		invite      *types.Invite
		want        *types.MsgSetMembershipResponse
		senderIsGov bool
		wantErr     bool
	}{
		{
			name: "invited user gets black membership by government",
			msg: &types.MsgSetMembership{
				Government:    testUser3.String(),
				Subscriber:    testUser.String(),
				NewMembership: "black",
			},
			invite: &types.Invite{
				Sender:           testUser3.String(),
				User:             testUser.String(),
				Status:           uint64(types.InviteStatusPending),
				SenderMembership: types.MembershipTypeBlack,
			},
			want:        &types.MsgSetMembershipResponse{},
			senderIsGov: true,
			wantErr:     false,
		},
		{
			name: "non-invited user gets black membership by government",
			msg: &types.MsgSetMembership{
				Government:    testUser3.String(),
				Subscriber:    testUser.String(),
				NewMembership: types.MembershipTypeBlack,
			},
			invite:      nil,
			want:        &types.MsgSetMembershipResponse{},
			senderIsGov: true,
			wantErr:     false,
		},
		{
			name: "invited, verified user doesn't get black membership because sender is not government",
			msg: &types.MsgSetMembership{
				Government:    testInviteSender.String(),
				Subscriber:    testUser.String(),
				NewMembership: "bronze",
			},
			invite: &types.Invite{
				Sender:           testInviteSender.String(),
				User:             testUser.String(),
				Status:           uint64(types.InviteStatusPending),
				SenderMembership: types.MembershipTypeBlack,
			},
			want:        nil,
			senderIsGov: false,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			government, _ := sdk.AccAddressFromBech32(tt.msg.Government)
			_ = k.AssignMembership(ctx, government, types.MembershipTypeBlack, testTsp, testExpiration)

			if tt.invite != nil {
				k.SaveInvite(ctx, *tt.invite)
			}
			if tt.senderIsGov {
				k.AddTrustedServiceProvider(ctx, government)
			}
			msg := NewMsgServerImpl(k)

			got, err := msg.SetMembership(sdk.WrapSDKContext(ctx), tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgServer.SetMembership() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgServer.SetMembership() = %v, want %v", got, tt.want)
			}
		})
	}
}
