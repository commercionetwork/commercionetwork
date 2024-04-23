package keeper

import (
	"reflect"
	"testing"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func Test_msgServer_InviteUser(t *testing.T) {

	type args struct {
		msg *types.MsgInviteUser
	}
	tests := []struct {
		name           string
		membershipType string
		invitee        sdk.AccAddress
		invitedUser    sdk.AccAddress
		existingInvite types.Invite
		existingUser   bool
		want           *types.MsgInviteUserResponse
		wantErr        bool
	}{
		{
			name:           "user is already present",
			invitee:        testUser2,
			invitedUser:    testUser,
			existingInvite: types.Invite{Sender: testUser2.String(), User: testUser.String(), Status: uint64(types.InviteStatusPending)},
			existingUser:   true,
			wantErr:        true,
		},
		// These controls are implemented in invite keeper method: redundant tests
		/*
			{
				name: "user has a invite",
			},
			{
				name: "invitee user has't a membership",
			},

		*/
		{
			name:           "user correctly invited",
			membershipType: types.MembershipTypeBronze,
			invitee:        testUser2,
			invitedUser:    testUser,
			//existingInvite: types.Invite{Sender: testUser2.String(), User: testUser.String(), Status: uint64(types.InviteStatusPending)},
			existingUser: false,
			want:         &types.MsgInviteUserResponse{Status: "1"},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, bk, _, k := SetupTestInput()
			msg := NewMsgServerImpl(k)

			if tt.existingUser {
				require.NoError(t,
					bk.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(stakeDenom, math.NewInt(1)))),
				)
				require.NoError(t,
					bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, tt.invitedUser, sdk.NewCoins(sdk.NewCoin(stakeDenom, math.NewInt(1)))),
				)
				/*
					require.NoError(t,
						bk.SetBalances(ctx, tt.invitedUser, sdk.NewCoins(sdk.NewCoin(stakeDenom, sdk.NewInt(1)))),
					)
				*/
			}

			if !tt.existingInvite.Empty() {
				k.SaveInvite(ctx, tt.existingInvite)
			}

			if len(tt.membershipType) != 0 {
				_ = k.AssignMembership(ctx, tt.invitee, tt.membershipType, testTsp, testExpiration)
			}

			msgInvite := types.NewMsgInviteUser(tt.invitee.String(), tt.invitedUser.String())

			got, err := msg.InviteUser(sdk.WrapSDKContext(ctx), msgInvite)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgServer.InviteUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgServer.InviteUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
