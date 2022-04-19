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
		keeper.bankKeeper.SetBalances(ctx, moduleAcc.GetAddress(), amount)
	}
}

var testFunds1 sdk.DecCoins = sdk.NewDecCoins(sdk.NewDecCoin(types.BondDenom, sdk.NewInt(1000)))

func TestGetBlockRewardsPoolFunds(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)

	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetBlockRewardsPoolFundsRequest
		response *types.QueryGetBlockRewardsPoolFundsResponse
		err      error
	}{
		{
			desc:     "funds 1000ucommercio",
			request:  &types.QueryGetBlockRewardsPoolFundsRequest{},
			response: &types.QueryGetBlockRewardsPoolFundsResponse{Funds: testFunds1},
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		tc := tc
		setFunds(keeper, ctx, testFunds1)
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.GetBlockRewardsPoolFunds(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.Equal(t, tc.response, response)
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
		keeper.SetParamSet(ctx, params)

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
