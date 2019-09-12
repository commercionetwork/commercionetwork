package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/tbr/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	types2 "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/stretchr/testify/assert"

	comm "github.com/commercionetwork/commercionetwork/x/common/types"
)

func TestKeeper_getFundersStoreKey(t *testing.T) {
	_, _, k, _, _ := SetupTestInput()
	actual := k.getFundersStoreKey()
	expected := []byte(types.BlockRewardsPoolFundersPrefix)

	assert.Equal(t, expected, actual)
}

func TestKeeper_setFunders(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()
	var funders comm.Addresses

	k.setFunders(ctx, TestFunders)

	store := ctx.KVStore(k.StoreKey)
	fundersBz := store.Get([]byte(types.BlockRewardsPoolFundersPrefix))
	k.Cdc.MustUnmarshalBinaryBare(fundersBz, &funders)

	assert.Equal(t, TestFunders, funders)
}

func TestKeeper_AddBlockRewardsPoolFunder_FundersNotFound(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()

	k.AddBlockRewardsPoolFunder(ctx, TestFunder)

	actual := k.GetBlockRewardsPoolFunders(ctx)

	assert.Equal(t, TestFunders, actual)
}

func TestKeeper_AddBlockRewardsPoolFunder_FundersFound(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()
	tstFunder, _ := sdk.AccAddressFromBech32("cosmos1nynns8ex9fq6sjjfj8k79ymkdz4sqth06xexae")

	k.setFunders(ctx, TestFunders)
	k.AddBlockRewardsPoolFunder(ctx, tstFunder)

	expected := TestFunders.AppendIfMissing(tstFunder)

	funders := k.GetBlockRewardsPoolFunders(ctx)

	assert.Equal(t, expected, funders)
}

func TestKeeper_GetBlockRewardsPoolFunders(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()
	k.setFunders(ctx, TestFunders)
	actual := k.GetBlockRewardsPoolFunders(ctx)

	assert.Equal(t, TestFunders, actual)
}

func TestKeeper_setBlockRewardsPool_UtilityFunction(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()
	var pool sdk.DecCoins

	k.setBlockRewardsPool(ctx, TestBlockRewardsPool)
	store := ctx.KVStore(k.StoreKey)
	poolBz := store.Get([]byte(types.BlockRewardsPoolPrefix))
	k.Cdc.MustUnmarshalBinaryBare(poolBz, &pool)

	assert.Equal(t, pool, TestBlockRewardsPool)
}

func TestKeeper_getBrPoolStoreKey(t *testing.T) {
	_, _, k, _, _ := SetupTestInput()
	actual := k.getBrPoolStoreKey()
	expected := []byte(types.BlockRewardsPoolPrefix)

	assert.Equal(t, expected, actual)

}

func TestKeeper_GetBlockRewardsPool(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()

	k.setBlockRewardsPool(ctx, TestBlockRewardsPool)
	actual := k.GetBlockRewardsPool(ctx)

	assert.Equal(t, TestBlockRewardsPool, actual)
}

func TestKeeper_ComputeProposerReward(t *testing.T) {

	_, ctx, k, _, _ := SetupTestInput()

	tpy := sdk.NewDecWithPrec(25000, 0)
	tpy = tpy.Mul(sdk.NewDecWithPrec(1000000, 0))

	dpy, _ := sdk.NewDecFromStr("365.24")
	hpd := sdk.NewDecWithPrec(24, 0)
	mph := sdk.NewDecWithPrec(60, 0)
	bpm := sdk.NewDecWithPrec(12, 0)

	dpy = dpy.Mul(hpd).Mul(mph).Mul(bpm)

	averageValidatorsNumber := sdk.NewDecWithPrec(100, 0)
	vNumber := sdk.NewDecWithPrec(32, 0)

	totalStakedToken := sdk.NewInt(15000)

	TestValidator = TestValidator.UpdateStatus(sdk.Bonded)
	TestValidator, _ = TestValidator.AddTokensFromDel(sdk.NewInt(9000))
	validatorStakedTokens := TestValidator.GetBondedTokens().ToDec()

	firstMember := tpy.Quo(dpy).Mul(averageValidatorsNumber.Quo(vNumber))
	println(firstMember.String())
	secondMember := validatorStakedTokens.Quo(totalStakedToken.ToDec())
	println(secondMember.String())

	concreteReward := firstMember.Mul(secondMember)

	expected := sdk.DecCoins{sdk.DecCoin{Denom: types.DefaultBondDenom, Amount: concreteReward}}

	actual := k.ComputeProposerReward(ctx, 32, TestValidator, totalStakedToken)

	assert.Equal(t, expected, actual)
}

func TestKeeper_DistributeBlockRewards_EnoughPoolFunds(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()

	reward := sdk.DecCoins{sdk.NewDecCoin(types.DefaultBondDenom, sdk.NewInt(1000))}

	k.setBlockRewardsPool(ctx, TestBlockRewardsPool)

	validatorRewards := types2.ValidatorCurrentRewards{Rewards: sdk.DecCoins{}}
	k.DistributionKeeper.SetValidatorCurrentRewards(ctx, TestValidator.GetOperator(), validatorRewards)

	_ = k.DistributeBlockRewards(ctx, TestValidator, reward)

	actual := k.DistributionKeeper.GetValidatorCurrentRewards(ctx, TestValidator.OperatorAddress)

	assert.Equal(t, reward, actual.Rewards)
}

func TestKeeper_DistributeBlockRewards_InsufficientPoolFunds(t *testing.T) {
	_, ctx, k, _, _ := SetupTestInput()

	reward := sdk.DecCoins{sdk.NewDecCoin(types.DefaultBondDenom, sdk.NewInt(12000))}
	brPool := sdk.DecCoins{sdk.NewDecCoin(types.DefaultBondDenom, sdk.NewInt(10000))}

	k.setBlockRewardsPool(ctx, brPool)

	err := k.DistributeBlockRewards(ctx, TestValidator, reward)

	assert.Error(t, err)
}

func TestKeeper_IncrementBlockRewardsPool(t *testing.T) {
	_, ctx, k, ak, bk := SetupTestInput()

	k.setBlockRewardsPool(ctx, TestBlockRewardsPool)

	acc := ak.NewAccountWithAddress(ctx, TestFunder)
	ak.SetAccount(ctx, acc)
	accountCoins := sdk.NewCoins(sdk.Coin{Amount: sdk.NewInt(1000), Denom: types.DefaultBondDenom})
	_ = bk.SetCoins(ctx, acc.GetAddress(), accountCoins)

	k.IncrementBlockRewardsPool(ctx, TestFunder, TestAmount)
	actual := k.GetBlockRewardsPool(ctx)

	var greater bool

	if TestBlockRewardsPool.AmountOf(types.DefaultBondDenom).LT(actual.AmountOf(types.DefaultBondDenom)) {
		greater = true
	}

	assert.True(t, greater)
}
