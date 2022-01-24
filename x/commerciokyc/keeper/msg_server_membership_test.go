package keeper_test

import (
	"context"
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/keeper"

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
		//want    *types.MsgBuyMembershipResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Invalid membership type returns error",
			args: args{
				msg:    types.NewMsgBuyMembership("gren", testUser, testTsp),
				invite: types.NewInvite(testInviteSender, testUser, "bronze"),
			},
			//want: types.MsgBuyMembershipResponse{ExpiryAt: }
			wantErr: true,
		},
		/*{
			name:    "Invalid message returns error",
			wantErr: true,
		},
		{
			name: "Valid membership allows buying",
		},
		{
			name: "Buying without invite returns error",
		},
		{
			name: "Buying with invalid invite returns error",
		},
		{
			name: "Valid upgrade works properly",
		},
		{
			name: "Valid downgrade works properly",
		},
		{
			name: "Invalid buying memebership with diffrent denom",
		},*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, bk, _, k := SetupTestInput()

			msg := keeper.NewMsgServerImpl(k)
			_ = bk
			_ = msg
			_ = ctx
			if !tt.args.invite.Empty() {
				inviteSender, _ := sdk.AccAddressFromBech32(tt.args.invite.Sender)
				err := k.AssignMembership(ctx, inviteSender, types.MembershipTypeBlack, testTsp, testExpiration)
				require.NoError(t, err)
				k.SaveInvite(ctx, tt.args.invite)
			}

			got, err := msg.BuyMembership(sdk.WrapSDKContext(ctx), tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgServer.BuyMembership() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = got
			/*if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgServer.BuyMembership() = %v, want %v", got, tt.want)
			}*/
		})
	}
}

func Test_msgServer_RemoveMembership(t *testing.T) {

	type args struct {
		goCtx context.Context
		msg   *types.MsgRemoveMembership
	}
	tests := []struct {
		name    string
		args    args
		want    *types.MsgRemoveMembershipResponse
		wantErr bool
	}{
		// TODO: Add test cases.

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, bk, _, k := SetupTestInput()

			msg := keeper.NewMsgServerImpl(k)
			_ = bk

			got, err := msg.RemoveMembership(sdk.WrapSDKContext(ctx), tt.args.msg)
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

	type args struct {
		goCtx context.Context
		msg   *types.MsgSetMembership
	}
	tests := []struct {
		name    string
		args    args
		want    *types.MsgSetMembershipResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, bk, _, k := SetupTestInput()

			msg := keeper.NewMsgServerImpl(k)
			_ = bk
			got, err := msg.SetMembership(sdk.WrapSDKContext(ctx), tt.args.msg)
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
