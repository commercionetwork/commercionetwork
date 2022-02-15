package keeper

import (
	"reflect"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/did/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_Identity(t *testing.T) {
	type args struct {
		req *types.QueryResolveIdentityRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryResolveIdentityResponse
		wantErr bool
	}{
		{
			name:    "invalid request",
			args:    args{req: nil},
			wantErr: true,
		},
		{
			name: "empty",
			args: args{
				req: &types.QueryResolveIdentityRequest{
					ID: types.ValidIdentity.DidDocument.ID,
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				req: &types.QueryResolveIdentityRequest{
					ID: types.ValidIdentity.DidDocument.ID,
				},
			},
			want: &types.QueryResolveIdentityResponse{
				Identity: &types.ValidIdentity,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			if tt.want != nil {
				k.SetIdentity(ctx, *tt.want.Identity)
			}

			sdkCtx := sdk.WrapSDKContext(ctx)

			got, err := k.Identity(sdkCtx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keeper.Identity() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.Identity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeeper_IdentityHistory(t *testing.T) {
	type args struct {
		req *types.QueryResolveIdentityHistoryRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryResolveIdentityHistoryResponse
		wantErr bool
	}{
		{
			name: "invalid request",
			args: args{
				req: nil,
			},
			wantErr: true,
		},
		{
			name: "empty",
			args: args{
				req: &types.QueryResolveIdentityHistoryRequest{},
			},
			want: &types.QueryResolveIdentityHistoryResponse{
				Identities: []*types.Identity{},
			},
			wantErr: false,
		},
		{
			name: "many updates",
			args: args{
				req: &types.QueryResolveIdentityHistoryRequest{},
			},
			want: &types.QueryResolveIdentityHistoryResponse{
				Identities: identitiesAtIncreasingMoments(5, 0),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := setupKeeper(t)

			if tt.want != nil {
				for _, identity := range tt.want.Identities {
					k.SetIdentity(ctx, *identity)
				}
			}

			sdkCtx := sdk.WrapSDKContext(ctx)

			got, err := k.IdentityHistory(sdkCtx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keeper.IdentityHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.IdentityHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}
