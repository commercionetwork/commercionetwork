package keeper

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
)

func setFunds(keeper *Keeper, ctx sdk.Context, pool sdk.DecCoins) {
	if pool.Empty() {
		return
	}

	keeper.SetTotalRewardPool(ctx, pool)
	moduleAcc := keeper.VbrAccount(ctx)
	keeper.accountKeeper.SetModuleAccount(ctx, moduleAcc)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	coins := GetCoins(*keeper, ctx, moduleAcc)
	if coins.Empty() {
		amount, _ := pool.TruncateDecimal()
		//keeper.bankKeeper.SetBalances(ctx, moduleAcc.GetAddress(), amount)
		keeper.bankKeeper.MintCoins(ctx, types.ModuleName, amount)
	}
}

var testFunds1 sdk.DecCoins = sdk.NewDecCoins(sdk.NewDecCoin(types.BondDenom, sdk.NewInt(100000)))

func TestGetBlockRewardsPoolFunds(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	for _, tt := range []struct {
		desc     string
		request  *types.QueryGetBlockRewardsPoolFundsRequest
		response *types.QueryGetBlockRewardsPoolFundsResponse
		err      error
	}{
		{
			desc:     "funds 100000ucommercio",
			request:  &types.QueryGetBlockRewardsPoolFundsRequest{},
			response: &types.QueryGetBlockRewardsPoolFundsResponse{Funds: testFunds1},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		setFunds(keeper, ctx, testFunds1)
		t.Run(tt.desc, func(t *testing.T) {
			response, err := keeper.GetBlockRewardsPoolFunds(wctx, tt.request)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.Equal(t, tt.response, response)
			}
		})
	}
}

var params = types.NewParams(types.EpochDay, sdk.NewDecWithPrec(5, 1))

func Test_GetParams(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetParamsRequest
		response *types.QueryGetParamsResponse
		err      error
	}{
		{
			desc:     "daily epoch and 0.5 earn rate",
			request:  &types.QueryGetParamsRequest{},
			response: &types.QueryGetParamsResponse{Params: params},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		err := keeper.SetParamSet(ctx, params)
		require.NoError(t, err)

		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.GetParams(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
			}
		})
	}
}
