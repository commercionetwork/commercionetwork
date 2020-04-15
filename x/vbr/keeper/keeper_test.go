package keeper

import (
	"testing"

	dist "github.com/cosmos/cosmos-sdk/x/distribution"

	"github.com/commercionetwork/commercionetwork/x/vbr/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// -------------
// --- Pool
// -------------

func TestKeeper_SetTotalRewardPool(t *testing.T) {
	cdc, ctx, k, _, _ := SetupTestInput(true)

	k.SetTotalRewardPool(ctx, TestBlockRewardsPool)

	var pool sdk.DecCoins
	store := ctx.KVStore(k.storeKey)
	cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.PoolStoreKey)), &pool)

	require.Equal(t, TestBlockRewardsPool, pool)
}

func TestKeeper_GetTotalRewardPool_EmptyPool(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput(true)

	actual := k.GetTotalRewardPool(ctx)
	require.Empty(t, actual)
}

func TestKeeper_GetTotalRewardPool_ExistingPool(t *testing.T) {
	cdc, ctx, k, _, _ := SetupTestInput(false)

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.PoolStoreKey), cdc.MustMarshalBinaryBare(&TestBlockRewardsPool))

	actual := k.GetTotalRewardPool(ctx)
	require.Equal(t, TestBlockRewardsPool, actual)
}

// --------------------------
// --- Yearly reward pool
// --------------------------

func TestKeeper_GetYearlyRewardPool_EmptyPool(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput(true)

	actual := k.GetYearlyRewardPool(ctx)
	require.Empty(t, actual)
}

func TestKeeper_GetYearlyRewardPool_ExistingPool(t *testing.T) {
	cdc, ctx, k, _, _ := SetupTestInput(false)

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.YearlyPoolStoreKey), cdc.MustMarshalBinaryBare(&TestBlockRewardsPool))

	actual := k.GetYearlyRewardPool(ctx)
	require.Equal(t, TestBlockRewardsPool, actual)
}

func TestKeeper_SetYearlyRewardPool(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput(false)

	k.SetYearlyRewardPool(ctx, TestBlockRewardsPool)

	var actual sdk.DecCoins
	store := ctx.KVStore(k.storeKey)
	k.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.YearlyPoolStoreKey)), &actual)

	require.Equal(t, TestBlockRewardsPool, actual)
}

// --------------------
// --- Year number
// --------------------

func Test_computeYearFromBlockHeight(t *testing.T) {
	require.Equal(t, int64(0), computeYearFromBlockHeight(0))
	require.Equal(t, int64(0), computeYearFromBlockHeight(6311519))
	require.Equal(t, int64(1), computeYearFromBlockHeight(6311520))
	require.Equal(t, int64(2), computeYearFromBlockHeight(12624040))
	require.Equal(t, int64(5), computeYearFromBlockHeight(31557600))
}

func TestKeeper_SetYearNumber(t *testing.T) {
	cdc, ctx, k, _, _ := SetupTestInput(false)

	k.SetYearNumber(ctx, 5)

	var actual int64
	store := ctx.KVStore(k.storeKey)
	cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.YearNumberStoreKey)), &actual)

	require.Equal(t, int64(5), actual)
}

func TestKeeper_GetYearNumber_Empty(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput(false)
	require.Equal(t, int64(0), k.GetYearNumber(ctx))
}

func TestKeeper_GetYearNumber_Set(t *testing.T) {
	cdc, ctx, k, _, _ := SetupTestInput(false)
	store := ctx.KVStore(k.storeKey)

	store.Set([]byte(types.YearNumberStoreKey), cdc.MustMarshalBinaryBare(0))
	require.Equal(t, int64(0), k.GetYearNumber(ctx))

	store.Set([]byte(types.YearNumberStoreKey), cdc.MustMarshalBinaryBare(5))
	require.Equal(t, int64(5), k.GetYearNumber(ctx))

	store.Set([]byte(types.YearNumberStoreKey), cdc.MustMarshalBinaryBare(0))
	require.Equal(t, int64(0), k.GetYearNumber(ctx))
}

func TestKeeper_UpdateYearlyPool_SameYear(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput(false)

	rewards := sdk.DecCoins{sdk.NewInt64DecCoin("stake", 100000)}
	k.SetYearNumber(ctx, 0)
	k.SetTotalRewardPool(ctx, rewards)

	k.UpdateYearlyPool(ctx, 0)

	require.Equal(t, rewards, k.GetTotalRewardPool(ctx))
	require.Empty(t, k.GetYearlyRewardPool(ctx))
}

func TestKeeper_UpdateYearlyPool_DifferentYear(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput(false)

	rewards := sdk.DecCoins{sdk.NewInt64DecCoin("stake", 20000)}
	k.SetYearNumber(ctx, 0)
	k.SetTotalRewardPool(ctx, rewards)

	k.UpdateYearlyPool(ctx, 6311520)

	require.Equal(t, int64(1), k.GetYearNumber(ctx))
	require.Equal(t, sdk.DecCoins{sdk.NewInt64DecCoin("stake", 20000)}, k.GetYearlyRewardPool(ctx))
}

// ---------------------------
// --- Reward distribution
// ---------------------------

func TestKeeper_ComputeProposerReward_50ValidatorsBalanced(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput(false)

	k.SetYearlyRewardPool(ctx, sdk.DecCoins{sdk.NewInt64DecCoin("uccc", 2500000)})
	TestValidator = TestValidator.UpdateStatus(sdk.Bonded)
	TestValidator, _ = TestValidator.AddTokensFromDel(sdk.NewInt(1))

	reward := k.ComputeProposerReward(ctx, 50, TestValidator, sdk.NewInt(50))

	expected := sdk.DecCoins{sdk.NewDecCoinFromDec("uccc", sdk.NewDecWithPrec(99025274418840450, 18))}
	require.Equal(t, expected, reward)
}

func TestKeeper_ComputeProposerReward_100ValidatorsBalanced(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput(false)

	k.SetYearlyRewardPool(ctx, sdk.DecCoins{sdk.NewInt64DecCoin("uccc", 2500000)})
	TestValidator = TestValidator.UpdateStatus(sdk.Bonded)
	TestValidator, _ = TestValidator.AddTokensFromDel(sdk.NewInt(1))

	reward := k.ComputeProposerReward(ctx, 100, TestValidator, sdk.NewInt(100))

	expected := sdk.DecCoins{sdk.NewDecCoinFromDec("uccc", sdk.NewDecWithPrec(3961010976753619, 16))}
	require.Equal(t, expected, reward)
}

func TestKeeper_DistributeBlockRewards_EnoughPoolFunds(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput(false)

	pool := sdk.DecCoins{sdk.NewInt64DecCoin("stake", 100000)}
	k.SetTotalRewardPool(ctx, pool)
	k.SetYearlyRewardPool(ctx, pool)

	validatorRewards := dist.ValidatorCurrentRewards{Rewards: sdk.DecCoins{}}
	k.distKeeper.SetValidatorCurrentRewards(ctx, TestValidator.GetOperator(), validatorRewards)

	validatorOutstandingRewards := sdk.DecCoins{}
	k.distKeeper.SetValidatorOutstandingRewards(ctx, TestValidator.GetOperator(), validatorOutstandingRewards)

	reward := sdk.DecCoins{sdk.NewDecCoin("stake", sdk.NewInt(1000))}
	err := k.DistributeBlockRewards(ctx, TestValidator, reward)
	require.NoError(t, err)
	actual := k.distKeeper.GetValidatorOutstandingRewards(ctx, TestValidator.OperatorAddress)
	require.Equal(t, reward, actual)

	remaining := sdk.DecCoins{sdk.NewInt64DecCoin("stake", 99000)}
	require.Equal(t, remaining, k.GetTotalRewardPool(ctx))
	require.Equal(t, remaining, k.GetYearlyRewardPool(ctx))
}

func TestKeeper_DistributeBlockRewards_InsufficientPoolFunds(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput(false)

	reward := sdk.DecCoins{sdk.NewDecCoin("stake", sdk.NewInt(12000))}
	err := k.DistributeBlockRewards(ctx, TestValidator, reward)

	require.Error(t, err)
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
			sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(100000), Denom: "stake"}),
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k, _, _ := SetupTestInput(tt.emptyPool)
			macc := k.VbrAccount(ctx)

			require.Equal(t, macc.GetName(), tt.wantModName)
			require.True(t, macc.GetCoins().IsEqual(tt.wantModAccBalance))
		})
	}
}

func TestKeeper_MintVBRTokens(t *testing.T) {
	tests := []struct {
		name       string
		wantAmount sdk.Coins
	}{
		{
			"add 10stake",
			sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(10), Denom: "stake"}),
		},
		{
			"add no stake",
			sdk.NewCoins(),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, k, _, _ := SetupTestInput(true)
			k.MintVBRTokens(ctx, tt.wantAmount)
			macc := k.VbrAccount(ctx)
			require.True(t, macc.GetCoins().IsEqual(tt.wantAmount))
		})
	}
}
