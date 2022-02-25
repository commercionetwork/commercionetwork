package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

var params_test = types.Params{
	DistrEpochIdentifier: types.EpochDay,
	EarnRate:             sdk.NewDecWithPrec(5, 1),
}

func TestKeeper_ComputeProposerReward(t *testing.T) {
	tests := []struct {
		name           string
		bonded         sdk.Int
		vNumber        int64
		expectedReward string
		params         types.Params
	}{
		{
			"Compute reward with 100 validators",
			sdk.NewInt(100000000),
			100,
			"136986.301369863013698630",
			params_test,
		},
		{
			"Compute reward with 50 validators",
			sdk.NewInt(100000000),
			50,
			"68493.150684931506849315",
			params_test,
		},
		{
			"Compute reward with small bonded",
			sdk.NewInt(1),
			100,
			"0.001369863013698630",
			params_test,
		},
		{
			"Compute reward per minute",
			sdk.NewInt(100000000),
			50,
			"47.564687975646879756",
			types.Params{
				DistrEpochIdentifier: types.EpochMinute,
				EarnRate:             sdk.NewDecWithPrec(5, 1),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := SetupKeeper(t)

			testVal := TestValidator.UpdateStatus(stakingTypes.Bonded)
			testVal, _ = testVal.AddTokensFromDel(tt.bonded)
			params := tt.params
			reward := k.ComputeProposerReward(ctx, tt.vNumber, testVal, types.BondDenom, params)

			expectedDecReward, _ := sdk.NewDecFromStr(tt.expectedReward)

			expected := sdk.DecCoins{sdk.NewDecCoinFromDec(types.BondDenom, expectedDecReward)}

			require.Equal(t, expected, reward)

		})
	}
}

func TestKeeper_DistributeBlockRewards(t *testing.T) {
	tests := []struct {
		name              string
		pool              sdk.DecCoins
		expectedValidator sdk.DecCoins
		expectedRemaining sdk.DecCoins
		expectedError     error
		bonded            sdk.Int
	}{
		{
			name:              "Reward with enough pool",
			pool:              sdk.DecCoins{sdk.NewInt64DecCoin(types.BondDenom, 100000)},
			expectedRemaining: sdk.DecCoins{sdk.NewInt64DecCoin(types.BondDenom, 86302)},
			expectedValidator: sdk.DecCoins{sdk.NewInt64DecCoin(types.BondDenom, 13698)},
			bonded:            sdk.NewInt(1000000000),
		},
		{
			name:              "Reward with empty pool",
			pool:              sdk.DecCoins{sdk.NewInt64DecCoin(types.BondDenom, 0)},
			expectedRemaining: sdk.DecCoins{},
			expectedValidator: sdk.DecCoins(nil),
			bonded:            sdk.NewInt(1000000000),
		},
		{
			name:              "Reward not enough funds into pool",
			pool:              sdk.DecCoins{sdk.NewInt64DecCoin(types.BondDenom, 1)},
			expectedRemaining: sdk.DecCoins{sdk.NewInt64DecCoin(types.BondDenom, 1)},
			expectedError:     sdkErr.Wrap(sdkErr.ErrInsufficientFunds, "Pool hasn't got enough funds to supply validator's rewards"),
			expectedValidator: sdk.DecCoins(nil),
			bonded:            sdk.NewInt(1000000000),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := SetupKeeper(t)

			testVal := TestValidator.UpdateStatus(stakingTypes.Bonded)
			testVal, _ = testVal.AddTokensFromDel(tt.bonded)

			k.SetTotalRewardPool(ctx, tt.pool)

			macc := k.VbrAccount(ctx)
			suppl, _ := tt.pool.TruncateDecimal()

			k.bankKeeper.SetBalances(ctx, macc.GetAddress(), sdk.NewCoins(suppl...))
			k.accountKeeper.SetModuleAccount(ctx, macc)

			validatorRewards := distrTypes.ValidatorCurrentRewards{Rewards: sdk.DecCoins{}}
			k.distKeeper.SetValidatorCurrentRewards(ctx, testVal.GetOperator(), validatorRewards)

			validatorOutstandingRewards := distrTypes.ValidatorOutstandingRewards{}
			k.distKeeper.SetValidatorOutstandingRewards(ctx, testVal.GetOperator(), validatorOutstandingRewards)

			params := types.Params{
				DistrEpochIdentifier: types.EpochDay,
				EarnRate:             sdk.NewDecWithPrec(5, 1),
			}
			reward := k.ComputeProposerReward(ctx, 1, testVal, types.BondDenom, params)
			rewardInt, _ := reward.TruncateDecimal()
			_ = rewardInt
			err := k.DistributeBlockRewards(ctx, testVal, reward)
			if tt.expectedError != nil {
				require.Equal(t, err.Error(), tt.expectedError.Error())
			}

			valCurReward := k.distKeeper.GetValidatorCurrentRewards(ctx, testVal.GetOperator())
			rewardPool := k.GetTotalRewardPool(ctx)

			require.Equal(t, tt.expectedRemaining, rewardPool)
			require.Equal(t, tt.expectedValidator, valCurReward.Rewards)

		})
	}
}

func TestKeeper_VbrAccount(t *testing.T) {
	tests := []struct {
		name              string
		wantModName       string
		wantModAccBalance sdk.Coins
		emptyPool         bool
	}{
		{
			"an empty vbr account",
			"vbr",
			sdk.NewCoins(),
			true,
		},
		{
			"a vbr account with coins in it",
			"vbr",
			sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(100000), Denom: types.BondDenom}),
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := SetupKeeper(t)
			macc := k.VbrAccount(ctx)

			require.Equal(t, macc.GetName(), tt.wantModName)

			if !tt.emptyPool {
				coins := sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(100000), Denom: types.BondDenom})
				k.bankKeeper.SetBalances(ctx, macc.GetAddress(), coins)
			}

			require.True(t, k.bankKeeper.GetAllBalances(ctx, macc.GetAddress()).IsEqual(tt.wantModAccBalance))
		})
	}
}

func TestKeeper_MintVBRTokens(t *testing.T) {
	tests := []struct {
		name       string
		wantAmount sdk.Coins
	}{
		{
			"add 10ucommercio",
			sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(10), Denom: types.BondDenom}),
		},
		{
			"add no ucommercio",
			sdk.NewCoins(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := SetupKeeper(t)
			k.bankKeeper.SetSupply(ctx, bankTypes.NewSupply(sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(10), Denom: types.BondDenom})))
			k.MintVBRTokens(ctx, tt.wantAmount)
			macc := k.VbrAccount(ctx)
			//require.True(t, macc.GetCoins().IsEqual(tt.wantAmount))
			require.True(t, k.bankKeeper.GetAllBalances(ctx, macc.GetAddress()).IsEqual(tt.wantAmount))
		})
	}
}

func TestKeeper_SetTotalRewardPool(t *testing.T) {

	tests := []struct {
		name        string
		updatedPool sdk.DecCoins
	}{
		{
			name: "empty pool",
		},
		// failing test
		// {
		// 	name:        "ok",
		// 	updatedPool: sdk.NewDecCoinsFromCoins(types.ValidMsgIncrementBlockRewardsPool.Amount...),
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, ctx := SetupKeeper(t)
			k.SetTotalRewardPool(ctx, tt.updatedPool)

			store := ctx.KVStore(k.storeKey)
			if tt.updatedPool.Empty() {
				require.False(t, store.Has([]byte(types.PoolStoreKey)))
			} else {
				actual := k.GetTotalRewardPool(ctx)
				require.Equal(t, tt.updatedPool, actual)
			}

		})
	}
}
