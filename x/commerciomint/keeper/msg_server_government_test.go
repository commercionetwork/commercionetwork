package keeper

// import (
// 	"reflect"
// 	"testing"

// 	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	"github.com/stretchr/testify/require"
// )

// func Test_msgServer_SetConversionRate(t *testing.T) {
// 	type args struct {
// 		msg *types.MsgSetCCCConversionRate
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *types.MsgSetCCCConversionRateResponse
// 		wantErr bool
// 	}{
// 		{
// 			name: "invalid signer",
// 			args: args{
// 				msg: &types.MsgSetCCCConversionRate{
// 					Signer: "",
// 					Rate:   sdk.NewDec(2),
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "signer is not the government",
// 			args: args{
// 				msg: &types.MsgSetCCCConversionRate{
// 					Signer: testEtp.Owner,
// 					Rate:   sdk.NewDec(2),
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "invalid conversion rate",
// 			args: args{
// 				msg: &types.MsgSetCCCConversionRate{
// 					Signer: government.String(),
// 					Rate:   sdk.NewDec(0),
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "ok",
// 			args: args{
// 				msg: &types.MsgSetCCCConversionRate{
// 					Signer: government.String(),
// 					Rate:   sdk.NewDec(3),
// 				},
// 			},
// 			want: &types.MsgSetCCCConversionRateResponse{Rate: sdk.NewDec(3)},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			wctx, _, gk, k, msgServer := SetupMsgServer()

// 			err := gk.SetGovernmentAddress(sdk.UnwrapSDKContext(wctx), government)
// 			require.NoError(t, err)

// 			got, err := msgServer.SetConversionRate(wctx, tt.args.msg)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("msgServer.SetConversionRate() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("msgServer.SetConversionRate() = %v, want %v", got, tt.want)
// 			}

// 			if !tt.wantErr {
// 				require.Equal(t, tt.args.msg.Rate, k.GetConversionRate(sdk.UnwrapSDKContext(wctx)))
// 			}
// 		})
// 	}
// }

// func Test_msgServer_SetFreezePeriod(t *testing.T) {
// 	type args struct {
// 		msg *types.MsgSetCCCFreezePeriod
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *types.MsgSetCCCFreezePeriodResponse
// 		wantErr bool
// 	}{
// 		{
// 			name: "invalid signer",
// 			args: args{
// 				msg: &types.MsgSetCCCFreezePeriod{
// 					Signer:       "",
// 					FreezePeriod: validFreezePeriod.String(),
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "signer is not the government",
// 			args: args{
// 				msg: &types.MsgSetCCCFreezePeriod{

// 					Signer:       testEtp.Owner,
// 					FreezePeriod: validFreezePeriod.String(),
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "unparseable freeze period",
// 			args: args{
// 				msg: &types.MsgSetCCCFreezePeriod{

// 					Signer:       government.String(),
// 					FreezePeriod: "abcd",
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "invalid freeze period",
// 			args: args{
// 				msg: &types.MsgSetCCCFreezePeriod{

// 					Signer:       government.String(),
// 					FreezePeriod: (invalidFreezePeriod).String(),
// 				},
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "ok",
// 			args: args{
// 				msg: &types.MsgSetCCCFreezePeriod{
// 					Signer:       government.String(),
// 					FreezePeriod: validFreezePeriod.String(),
// 				},
// 			},
// 			want: &types.MsgSetCCCFreezePeriodResponse{FreezePeriod: validFreezePeriod.String()},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			wctx, _, gk, k, msgServer := SetupMsgServer()

// 			err := gk.SetGovernmentAddress(sdk.UnwrapSDKContext(wctx), government)
// 			require.NoError(t, err)

// 			got, err := msgServer.SetFreezePeriod(wctx, tt.args.msg)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("msgServer.SetFreezePeriod() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("msgServer.SetFreezePeriod() = %v, want %v", got, tt.want)
// 			}

// 			if !tt.wantErr {
// 				require.Equal(t, tt.args.msg.FreezePeriod, k.GetFreezePeriod(sdk.UnwrapSDKContext(wctx)).String())
// 			}
// 		})
// 	}
// }
