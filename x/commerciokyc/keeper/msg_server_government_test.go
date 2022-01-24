package keeper

import (
	"context"
	"reflect"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/commerciokyc/types"
)

func Test_msgServer_AddTsp(t *testing.T) {
	type fields struct {
		Keeper Keeper
	}
	type args struct {
		goCtx context.Context
		msg   *types.MsgAddTsp
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.MsgAddTspResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := msgServer{
				Keeper: tt.fields.Keeper,
			}
			got, err := k.AddTsp(tt.args.goCtx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgServer.AddTsp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgServer.AddTsp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_msgServer_RemoveTsp(t *testing.T) {
	type fields struct {
		Keeper Keeper
	}
	type args struct {
		goCtx context.Context
		msg   *types.MsgRemoveTsp
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.MsgRemoveTspResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := msgServer{
				Keeper: tt.fields.Keeper,
			}
			got, err := k.RemoveTsp(tt.args.goCtx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgServer.RemoveTsp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgServer.RemoveTsp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_msgServer_DepositIntoLiquidityPool(t *testing.T) {
	type fields struct {
		Keeper Keeper
	}
	type args struct {
		goCtx context.Context
		msg   *types.MsgDepositIntoLiquidityPool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *types.MsgDepositIntoLiquidityPoolResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := msgServer{
				Keeper: tt.fields.Keeper,
			}
			got, err := k.DepositIntoLiquidityPool(tt.args.goCtx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgServer.DepositIntoLiquidityPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgServer.DepositIntoLiquidityPool() = %v, want %v", got, tt.want)
			}
		})
	}
}
