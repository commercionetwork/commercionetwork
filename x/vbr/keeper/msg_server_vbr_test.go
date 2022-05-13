package keeper

import (
	"reflect"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func Test_msgServer_IncrementBlockRewardsPool(t *testing.T) {
	tests := []struct {
		name        string
		msg         *types.MsgIncrementBlockRewardsPool
		funderFunds sdk.Coins
		want        *types.MsgIncrementBlockRewardsPoolResponse
		wantErr     bool
	}{
		{
			name: "invalid funder",
			msg: &types.MsgIncrementBlockRewardsPool{
				Funder: "",
				Amount: types.ValidMsgIncrementBlockRewardsPool.Amount,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "funder has not enough funds",
			msg:     &types.ValidMsgIncrementBlockRewardsPool,
			wantErr: true,
		},
		{
			name:        "add amount in message to empty rewards pool",
			msg:         &types.ValidMsgIncrementBlockRewardsPool,
			funderFunds: types.ValidMsgIncrementBlockRewardsPool.Amount,
			want:        &types.MsgIncrementBlockRewardsPoolResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, ctx, keeper, bank := setupMsgServer(t)

			if !tt.funderFunds.Empty() {
				funder, err := sdk.AccAddressFromBech32(tt.msg.Funder)
				require.NoError(t, err)
				wctx := sdk.UnwrapSDKContext(ctx)
				bank.MintCoins(wctx, types.ModuleName, tt.funderFunds)
				bank.SendCoinsFromModuleToAccount(wctx, types.ModuleName, funder, tt.funderFunds)
				//bank.AddCoins(sdk.UnwrapSDKContext(ctx), funder, tt.funderFunds)
			}

			got, err := srv.IncrementBlockRewardsPool(ctx, tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgServer.IncrementBlockRewardsPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgServer.IncrementBlockRewardsPool() = %v, want %v", got, tt.want)
			}

			if !tt.wantErr {
				actual := keeper.GetTotalRewardPool(sdk.UnwrapSDKContext(ctx))

				expected := sdk.NewDecCoinsFromCoins(tt.msg.Amount...)

				require.Equal(t, expected, actual)
			}
		})
	}
}

var oldParams = types.NewParams(types.ValidMsgSetParams.DistrEpochIdentifier, types.ValidMsgSetParams.EarnRate.Mul(sdk.NewDec(2)))

func Test_msgServer_SetParams(t *testing.T) {
	tests := []struct {
		name      string
		oldParams *types.Params
		msg       *types.MsgSetParams
		want      *types.MsgSetParamsResponse
		wantErr   bool
	}{
		{
			name: "invalid government",
			msg: &types.MsgSetParams{
				Government:           "",
				DistrEpochIdentifier: types.ValidMsgSetParams.DistrEpochIdentifier,
				EarnRate:             types.ValidMsgSetParams.EarnRate,
			},
			wantErr: true,
		},
		{
			name: "different government",
			msg: &types.MsgSetParams{
				Government:           types.ValidMsgIncrementBlockRewardsPool.Funder,
				DistrEpochIdentifier: types.ValidMsgSetParams.DistrEpochIdentifier,
				EarnRate:             types.ValidMsgSetParams.EarnRate,
			},
			wantErr: true,
		},
		{
			name: "invalid DistrEpochIdentifier",
			msg: &types.MsgSetParams{
				Government:           types.ValidMsgSetParams.Government,
				DistrEpochIdentifier: "",
				EarnRate:             types.ValidMsgSetParams.EarnRate,
			},
			wantErr: true,
		},
		{
			name: "invalid EarnRate",
			msg: &types.MsgSetParams{
				Government:           types.ValidMsgSetParams.Government,
				DistrEpochIdentifier: types.ValidMsgSetParams.DistrEpochIdentifier,
				EarnRate:             types.InvalidEarnRate,
			},
			wantErr: true,
		},
		{
			name:      "ok",
			oldParams: &oldParams,
			msg: &types.MsgSetParams{
				Government:           types.ValidMsgSetParams.Government,
				DistrEpochIdentifier: types.ValidMsgSetParams.DistrEpochIdentifier,
				EarnRate:             types.ValidMsgSetParams.EarnRate,
			},
			want: &types.MsgSetParamsResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, ctx, keeper, _ := setupMsgServer(t)

			expected := types.NewParams(tt.msg.DistrEpochIdentifier, tt.msg.EarnRate)

			if tt.oldParams != nil {
				require.NotEqual(t, expected, tt.oldParams)

				err := keeper.SetParamSet(sdk.UnwrapSDKContext(ctx), *tt.oldParams)
				require.NoError(t, err)
			}

			got, err := srv.SetParams(ctx, tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgServer.SetParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgServer.SetParams() = %v, want %v", got, tt.want)
			}

			if !tt.wantErr {
				actual := keeper.GetParamSet(sdk.UnwrapSDKContext(ctx))
				require.Equal(t, expected, actual)
			}
		})
	}
}
