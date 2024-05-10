package keeper

// import (
// 	"reflect"
// 	"testing"

// 	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/require"
// )

// func Test_msgServer_SetParams(t *testing.T) {
// 	type args struct {
// 		msg *types.MsgSetParams
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *types.MsgSetParamsResponse
// 		wantErr bool
// 	}{
// 		{
// 			name: "invalid signer",
// 			args: args{
// 				msg: &types.MsgSetParams{
// 					Signer: "",
// 					Params: &types.Params{
// 						ConversionRate: types.DefaultConversionRate,
// 						FreezePeriod:   types.DefaultFreezePeriod,
// 					},
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "signer is not the government",
// 			args: args{
// 				msg: &types.MsgSetParams{
// 					Signer: testEtp.Owner,
// 					Params: &types.Params{
// 						ConversionRate: types.DefaultConversionRate,
// 						FreezePeriod:   types.DefaultFreezePeriod,
// 					},
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "invalid conversion rate",
// 			args: args{
// 				msg: &types.MsgSetParams{
// 					Signer: government.String(),
// 					Params: &types.Params{
// 						ConversionRate: invalidConversionRate,
// 						FreezePeriod:   types.DefaultFreezePeriod,
// 					},
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "invalid freeze period",
// 			args: args{
// 				msg: &types.MsgSetParams{
// 					Signer: government.String(),
// 					Params: &types.Params{
// 						ConversionRate: types.DefaultConversionRate,
// 						FreezePeriod:   invalidFreezePeriod,
// 					},
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "ok",
// 			args: args{
// 				msg: &types.MsgSetParams{
// 					Signer: government.String(),
// 					Params: &types.Params{
// 						ConversionRate: types.DefaultConversionRate,
// 						FreezePeriod:   types.DefaultFreezePeriod,
// 					},
// 				},
// 			},
// 			want: &types.MsgSetParamsResponse{},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			wctx, _, gk, k, msgServer := SetupMsgServer()

// 			err := gk.SetGovernmentAddress(sdk.UnwrapSDKContext(wctx), government)
// 			require.NoError(t, err)

// 			got, err := msgServer.SetParams(wctx, tt.args.msg)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("msgServer.SetParams() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("msgServer.SetParams() = %v, want %v", got, tt.want)
// 			}

// 			if !tt.wantErr {
// 				require.Equal(t, *tt.args.msg.Params, k.GetParams(sdk.UnwrapSDKContext(wctx)))

// 			}
// 		})
// 	}
// }
