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
				Amount: sdk.NewCoins(sdk.NewCoin(types.BondDenom, sdk.NewInt(1000))),
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
			name:        "ok",
			msg:         &types.ValidMsgIncrementBlockRewardsPool,
			funderFunds: types.ValidMsgIncrementBlockRewardsPool.Amount,
			want:        &types.MsgIncrementBlockRewardsPoolResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, ctx, bank := setupMsgServer(t)

			if !tt.funderFunds.Empty() {
				funder, err := sdk.AccAddressFromBech32(tt.msg.Funder)
				require.NoError(t, err)

				bank.AddCoins(sdk.UnwrapSDKContext(ctx), funder, tt.funderFunds)
			}

			got, err := srv.IncrementBlockRewardsPool(ctx, tt.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("msgServer.IncrementBlockRewardsPool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("msgServer.IncrementBlockRewardsPool() = %v, want %v", got, tt.want)
			}
		})
	}
}
