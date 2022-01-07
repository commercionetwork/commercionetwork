package keeper

import (
	"reflect"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_MintCCC(t *testing.T) {

	type args struct {
		msg *types.MsgMintCCC
	}
	tests := []struct {
		name    string
		args    args
		want    *types.MsgMintCCCResponse
		wantErr bool
	}{
		{
			name: "invalid depositor",
			args: args{
				msg: &types.MsgMintCCC{
					Depositor:     "",
					DepositAmount: []*sdk.Coin{&validDepositCoin},
					ID:            testEtp.ID,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid coins",
			args: args{
				msg: &types.MsgMintCCC{
					Depositor:     testEtpOwner.String(),
					DepositAmount: []*sdk.Coin{&validDepositCoin, &inValidDepositCoin},
					ID:            testEtp.ID,
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				msg: &types.MsgMintCCC{
					Depositor:     testEtpOwner.String(),
					DepositAmount: []*sdk.Coin{&validDepositCoin},
					ID:            testID,
				},
			},
			want: &types.MsgMintCCCResponse{ID: testEtp.ID},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wctx, bk, _, _, msgServer := SetupMsgServer()

			if !tt.wantErr {
				ownerAddr, err := sdk.AccAddressFromBech32(tt.args.msg.Depositor)
				require.NoError(t, err)
				err = bk.AddCoins(sdk.UnwrapSDKContext(wctx), ownerAddr, sdk.NewCoins(sdk.NewInt64Coin("ucommercio", 200)))
				require.NoError(t, err)
			}

			got, err := msgServer.MintCCC(wctx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgServer.MintCCC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgServer.MintCCC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_msgServer_BurnCCC(t *testing.T) {

	type args struct {
		msg *types.MsgBurnCCC
	}
	tests := []struct {
		name    string
		args    args
		want    *types.MsgBurnCCCResponse
		wantErr bool
	}{
		{
			name: "empty signer",
			args: args{
				msg: &types.MsgBurnCCC{
					Signer: "",
					Amount: &validBurnCoin,
					ID:     testID,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid coins",
			args: args{
				msg: &types.MsgBurnCCC{
					Signer: testEtpOwner.String(),
					Amount: &inValidBurnCoin,
					ID:     testID,
				},
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				msg: &types.MsgBurnCCC{
					Signer: testEtpOwner.String(),
					Amount: testEtp.Credits,
					ID:     testID,
				},
			},
			want: &types.MsgBurnCCCResponse{ID: testEtp.ID, Residual: &zeroUCCC}, // TODO check residuals
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wctx, bk, _, k, msgServer := SetupMsgServer()

			if !tt.wantErr {
				k.SetPosition(sdk.UnwrapSDKContext(wctx), testEtp)
				_ = bk.MintCoins(sdk.UnwrapSDKContext(wctx), types.ModuleName, testLiquidityPool)
				_ = bk.AddCoins(sdk.UnwrapSDKContext(wctx), testEtpOwner, sdk.NewCoins(*testEtp.Credits))
			}

			got, err := msgServer.BurnCCC(wctx, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgServer.MintCCC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgServer.MintCCC() = %v, want %v", got, tt.want)
			}
		})
	}
}
