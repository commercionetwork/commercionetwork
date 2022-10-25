package keeper

import (
	"reflect"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_Params(t *testing.T) {

	type args struct {
		req *types.QueryParams
	}
	tests := []struct {
		name    string
		args    args
		want    *types.QueryParamsResponse
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				req: &types.QueryParams{},
			},
			want: &types.QueryParamsResponse{
				Params: &validParams,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _, _, k := SetupTestInput()
			wctx := sdk.WrapSDKContext(ctx)

			got, err := k.Params(wctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Keeper.Params() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Keeper.Params() = %v, want %v", got, tt.want)
			}
		})
	}
}
