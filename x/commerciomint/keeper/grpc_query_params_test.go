package keeper

import (
	"reflect"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// func TestIdentityQuerySingle(t *testing.T) {
// 	keeper, ctx := setupKeeper(t)
// 	wctx := sdk.WrapSDKContext(ctx)
// 	msgs := createNIdentity(keeper, ctx, 2)
// 	for _, tc := range []struct {
// 		desc     string
// 		request  *types.QueryResolveDidDocumentRequest
// 		response *types.QueryResolveDidDocumentResponse
// 		err      error
// 	}{
// 		{
// 			desc:     "First",
// 			request:  &types.QueryResolveDidDocumentRequest{ID: msgs[0].ID},
// 			response: &types.QueryResolveDidDocumentResponse{DidDocument: &msgs[0]},
// 		},
// 		{
// 			desc:     "Second",
// 			request:  &types.QueryResolveDidDocumentRequest{ID: msgs[1].ID},
// 			response: &types.QueryResolveDidDocumentResponse{DidDocument: &msgs[1]},
// 		},
// 		{
// 			desc:    "KeyNotFound",
// 			request: &types.QueryResolveDidDocumentRequest{ID: "x"},
// 			err:     sdkerrors.ErrKeyNotFound,
// 		},
// 		{
// 			desc: "InvalidRequest",
// 			err:  status.Error(codes.InvalidArgument, "invalid request"),
// 		},
// 	} {
// 		tc := tc
// 		t.Run(tc.desc, func(t *testing.T) {
// 			response, err := keeper.Identity(wctx, tc.request)
// 			if tc.err != nil {
// 				require.ErrorIs(t, err, tc.err)
// 			} else {
// 				require.Equal(t, tc.response, response)
// 			}
// 		})
// 	}
// }

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
