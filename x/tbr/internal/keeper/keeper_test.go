package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/tbr/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dist "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/stretchr/testify/assert"
)

// -------------
// --- Pool
// -------------

func TestKeeper_SetTotalRewardPool(t *testing.T) {
	cdc, ctx, k, _, _ := SetupTestInput()

	k.SetTotalRewardPool(ctx, TestBlockRewardsPool)

	var pool sdk.DecCoins
	store := ctx.KVStore(k.storeKey)
	cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.PoolStoreKey)), &pool)

	assert.Equal(t, TestBlockRewardsPool, pool)
}

func TestKeeper_GetTotalRewardPool_EmptyPool(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()

	actual := k.GetTotalRewardPool(ctx)
	assert.Empty(t, actual)
}

func TestKeeper_GetTotalRewardPool_ExistingPool(t *testing.T) {
	cdc, ctx, k, _, _ := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.PoolStoreKey), cdc.MustMarshalBinaryBare(&TestBlockRewardsPool))

	actual := k.GetTotalRewardPool(ctx)
	assert.Equal(t, TestBlockRewardsPool, actual)
}

func TestKeeper_IncrementBlockRewardsPool(t *testing.T) {
	_, ctx, k, ak, bk := SetupTestInput()

	k.SetTotalRewardPool(ctx, TestBlockRewardsPool)

	acc := ak.NewAccountWithAddress(ctx, TestFunder)
	ak.SetAccount(ctx, acc)
	accountCoins := sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(1000), Denom: k.stakingKeeper.BondDenom(ctx)})
	_ = bk.SetCoins(ctx, acc.GetAddress(), accountCoins)

	_ = k.IncrementBlockRewardsPool(ctx, TestFunder, TestAmount)
	actual := k.GetTotalRewardPool(ctx)

	blockReward := TestBlockRewardsPool.AmountOf(k.stakingKeeper.BondDenom(ctx))
	actualReward := actual.AmountOf(k.stakingKeeper.BondDenom(ctx))
	assert.True(t, blockReward.LT(actualReward))
}

func TestKeeper_IncrementBlockRewardsPool_InsufficientUserFunds(t *testing.T) {
	_, ctx, k, ak, bk := SetupTestInput()

	k.SetTotalRewardPool(ctx, TestBlockRewardsPool)

	acc := ak.NewAccountWithAddress(ctx, TestFunder)
	ak.SetAccount(ctx, acc)
	accountCoins := sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(1), Denom: k.stakingKeeper.BondDenom(ctx)})
	_ = bk.SetCoins(ctx, acc.GetAddress(), accountCoins)

	err := k.IncrementBlockRewardsPool(ctx, TestFunder, TestAmount)

	assert.Error(t, err)
}

// --------------------------
// --- Yearly reward pool
// --------------------------

func TestKeeper_GetYearlyRewardPool_EmptyPool(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()

	actual := k.GetYearlyRewardPool(ctx)
	assert.Empty(t, actual)
}

func TestKeeper_GetYearlyRewardPool_ExistingPool(t *testing.T) {
	cdc, ctx, k, _, _ := SetupTestInput()

	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.YearlyPoolStoreKey), cdc.MustMarshalBinaryBare(&TestBlockRewardsPool))

	actual := k.GetYearlyRewardPool(ctx)
	assert.Equal(t, TestBlockRewardsPool, actual)
}

func TestKeeper_SetYearlyRewardPool(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()

	k.SetYearlyRewardPool(ctx, TestBlockRewardsPool)

	var actual sdk.DecCoins
	store := ctx.KVStore(k.storeKey)
	k.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.YearlyPoolStoreKey)), &actual)

	assert.Equal(t, TestBlockRewardsPool, actual)
}

// --------------------
// --- Year number
// --------------------

func Test_computeYearFromBlockHeight(t *testing.T) {
	assert.Equal(t, int64(0), computeYearFromBlockHeight(0))
	assert.Equal(t, int64(0), computeYearFromBlockHeight(6311519))
	assert.Equal(t, int64(1), computeYearFromBlockHeight(6311520))
	assert.Equal(t, int64(2), computeYearFromBlockHeight(12624040))
	assert.Equal(t, int64(5), computeYearFromBlockHeight(31557600))
}

func TestKeeper_SetYearNumber(t *testing.T) {
	cdc, ctx, k, _, _ := SetupTestInput()

	k.SetYearNumber(ctx, 5)

	var actual int64
	store := ctx.KVStore(k.storeKey)
	cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.YearNumberStoreKey)), &actual)

	assert.Equal(t, int64(5), actual)
}

func TestKeeper_GetYearNumber_Empty(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()
	assert.Equal(t, int64(0), k.GetYearNumber(ctx))
}

func TestKeeper_GetYearNumber_Set(t *testing.T) {
	cdc, ctx, k, _, _ := SetupTestInput()
	store := ctx.KVStore(k.storeKey)

	store.Set([]byte(types.YearNumberStoreKey), cdc.MustMarshalBinaryBare(0))
	assert.Equal(t, int64(0), k.GetYearNumber(ctx))

	store.Set([]byte(types.YearNumberStoreKey), cdc.MustMarshalBinaryBare(5))
	assert.Equal(t, int64(5), k.GetYearNumber(ctx))

	store.Set([]byte(types.YearNumberStoreKey), cdc.MustMarshalBinaryBare(0))
	assert.Equal(t, int64(0), k.GetYearNumber(ctx))
}

func TestKeeper_UpdateYearlyPool_SameYear(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()

	rewards := sdk.DecCoins{sdk.NewInt64DecCoin("uccc", 100)}
	k.SetYearNumber(ctx, 0)
	k.SetTotalRewardPool(ctx, rewards)

	k.UpdateYearlyPool(ctx, 0)

	assert.Equal(t, rewards, k.GetTotalRewardPool(ctx))
	assert.Empty(t, k.GetYearlyRewardPool(ctx))
}

func TestKeeper_UpdateYearlyPool_DifferentYear(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()

	rewards := sdk.DecCoins{sdk.NewInt64DecCoin("uccc", 100)}
	k.SetYearNumber(ctx, 0)
	k.SetTotalRewardPool(ctx, rewards)

	k.UpdateYearlyPool(ctx, 6311520)

	assert.Equal(t, int64(1), k.GetYearNumber(ctx))
	assert.Equal(t, sdk.DecCoins{sdk.NewInt64DecCoin("uccc", 20)}, k.GetYearlyRewardPool(ctx))
}

// ---------------------------
// --- Reward distribution
// ---------------------------

func TestKeeper_ComputeProposerReward_50ValidatorsBalanced(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()

	k.SetYearlyRewardPool(ctx, sdk.DecCoins{sdk.NewInt64DecCoin("uccc", 2500000)})
	TestValidator = TestValidator.UpdateStatus(sdk.Bonded)
	TestValidator, _ = TestValidator.AddTokensFromDel(sdk.NewInt(1))

	reward := k.ComputeProposerReward(ctx, 50, TestValidator, sdk.NewInt(50))

	expected := sdk.DecCoins{sdk.NewDecCoinFromDec("uccc", sdk.NewDecWithPrec(99025274418840450, 18))}
	assert.Equal(t, expected, reward)
}

func TestKeeper_ComputeProposerReward_100ValidatorsBalanced(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()

	k.SetYearlyRewardPool(ctx, sdk.DecCoins{sdk.NewInt64DecCoin("uccc", 2500000)})
	TestValidator = TestValidator.UpdateStatus(sdk.Bonded)
	TestValidator, _ = TestValidator.AddTokensFromDel(sdk.NewInt(1))

	reward := k.ComputeProposerReward(ctx, 100, TestValidator, sdk.NewInt(100))

	expected := sdk.DecCoins{sdk.NewDecCoinFromDec("uccc", sdk.NewDecWithPrec(3961010976753619, 16))}
	assert.Equal(t, expected, reward)
}

func TestKeeper_DistributeBlockRewards_EnoughPoolFunds(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()

	k.SetTotalRewardPool(ctx, TestBlockRewardsPool)
	k.SetYearlyRewardPool(ctx, TestBlockRewardsPool)

	validatorRewards := dist.ValidatorCurrentRewards{Rewards: sdk.DecCoins{}}
	k.DistributionKeeper.SetValidatorCurrentRewards(ctx, TestValidator.GetOperator(), validatorRewards)

	reward := sdk.DecCoins{sdk.NewDecCoin("stake", sdk.NewInt(1000))}
	_ = k.DistributeBlockRewards(ctx, TestValidator, reward)

	actual := k.DistributionKeeper.GetValidatorCurrentRewards(ctx, TestValidator.OperatorAddress)
	assert.Equal(t, reward, actual.Rewards)
}

func TestKeeper_DistributeBlockRewards_InsufficientPoolFunds(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()

	reward := sdk.DecCoins{sdk.NewDecCoin(k.stakingKeeper.BondDenom(ctx), sdk.NewInt(12000))}
	err := k.DistributeBlockRewards(ctx, TestValidator, reward)

	assert.Error(t, err)
}
