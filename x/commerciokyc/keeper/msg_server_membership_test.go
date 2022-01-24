package keeper

import (
	"context"
	"reflect"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func Test_msgServer_BuyMembership(t *testing.T) {
	type fields struct {
		Keeper Keeper
	}
	type args struct {
		goCtx context.Context
		msg   *types.MsgBuyMembership
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.MsgBuyMembershipResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := msgServer{
				Keeper: tt.fields.Keeper,
			}
			got, err := k.BuyMembership(tt.args.goCtx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgServer.BuyMembership() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgServer.BuyMembership() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_msgServer_RemoveMembership(t *testing.T) {
	type fields struct {
		Keeper Keeper
	}
	type args struct {
		goCtx context.Context
		msg   *types.MsgRemoveMembership
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.MsgRemoveMembershipResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := msgServer{
				Keeper: tt.fields.Keeper,
			}
			got, err := k.RemoveMembership(tt.args.goCtx, tt.args.msg)
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
	type fields struct {
		Keeper Keeper
	}
	type args struct {
		goCtx context.Context
		msg   *types.MsgSetMembership
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.MsgSetMembershipResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := msgServer{
				Keeper: tt.fields.Keeper,
			}
			got, err := k.SetMembership(tt.args.goCtx, tt.args.msg)
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
