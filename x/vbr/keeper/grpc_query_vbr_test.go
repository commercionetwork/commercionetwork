package keeper

import (
	//"fmt"
	//"sort"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
//	"github.com/cosmos/cosmos-sdk/types/query"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
)

func setFunds (keeper *Keeper, ctx sdk.Context, pool sdk.DecCoins) {
	if pool.Empty(){
		return
	}

	keeper.SetTotalRewardPool(ctx, pool)
	moduleAcc := keeper.VbrAccount(ctx)
	keeper.accountKeeper.SetModuleAccount(ctx, moduleAcc)
	/*if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	coins := GetCoins(*keeper, ctx, moduleAcc)
	if coins.Empty() {
		amount, _ := pool.TruncateDecimal()
		keeper.bankKeeper.SetBalances(ctx,moduleAcc.GetAddress(), amount)
	}*/
}

var testFunds1 sdk.DecCoins = sdk.NewDecCoins(sdk.NewDecCoin("ucommercio",sdk.NewInt(1000)))
var testFunds2 sdk.DecCoins = sdk.NewDecCoins(sdk.NewDecCoin("ucommercio", sdk.NewInt(10)))
func TestGetBlockRewardsPoolFunds(t *testing.T) {
	keeper, ctx := setupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetBlockRewardsPoolFundsRequest
		response *types.QueryGetBlockRewardsPoolFundsResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetBlockRewardsPoolFundsRequest{},
			response: &types.QueryGetBlockRewardsPoolFundsResponse{Funds: testFunds1},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetBlockRewardsPoolFundsRequest{},
			response: &types.QueryGetBlockRewardsPoolFundsResponse{Funds: testFunds2},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetBlockRewardsPoolFundsRequest{},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			setFunds(keeper, ctx, tc.response.Funds)
			response, err := keeper.GetBlockRewardsPoolFunds(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}

func TesGetVbrParams(t *testing.T) {

}