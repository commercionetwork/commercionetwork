package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// ---------------------
// --- Accrediters
// ---------------------

func TestKeeper_SetAccrediter_NoAccrediter(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Delete(TestUser)

	err := accreditationKeeper.SetAccrediter(ctx, TestUser, TestAccrediter)
	assert.Nil(t, err)

	accreditationBz := store.Get(TestUser)
	assert.NotNil(t, accreditationBz)

	var accreditation types.Accreditation
	cdc.MustUnmarshalBinaryBare(accreditationBz, &accreditation)

	assert.Equal(t, TestUser, accreditation.User)
	assert.Equal(t, TestAccrediter, accreditation.Accrediter)
	assert.False(t, accreditation.Rewarded)
}

func TestKeeper_SetAccrediter_ExistingAccrediter(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	existingAccredit := types.Accreditation{
		Accrediter: TestAccrediter,
		User:       TestUser,
		Rewarded:   false,
	}

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Set(TestUser, cdc.MustMarshalBinaryBare(existingAccredit))

	err := accreditationKeeper.SetAccrediter(ctx, TestUser, TestAccrediter)
	assert.NotNil(t, err)
}

func TestKeeper_GetAccrediter_NoAccrediter(t *testing.T) {
	ctx, _, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Delete(TestUser)

	accrediter := accreditationKeeper.GetAccrediter(ctx, TestUser)
	assert.Nil(t, accrediter)
}

func TestKeeper_GetAccrediter_ExistingAccrediter(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	accreditation := types.Accreditation{
		Accrediter: TestAccrediter,
		User:       TestUser,
		Rewarded:   false,
	}

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Set(TestUser, cdc.MustMarshalBinaryBare(accreditation))

	accrediter := accreditationKeeper.GetAccrediter(ctx, TestUser)
	assert.Equal(t, TestAccrediter, accrediter)
}

// ---------------------
// --- Accreditations
// ---------------------

func TestKeeper_GetAccreditations_EmptyList(t *testing.T) {
	ctx, _, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	accreditations := accreditationKeeper.GetAccreditations(ctx)
	assert.Empty(t, accreditations)
}

func TestKeeper_GetAccreditations_NonEmptyList(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	acc1 := types.Accreditation{Accrediter: TestAccrediter, User: TestUser, Rewarded: false}
	acc2 := types.Accreditation{Accrediter: TestUser, User: TestAccrediter, Rewarded: false}

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Set(TestUser, cdc.MustMarshalBinaryBare(acc1))
	store.Set(TestAccrediter, cdc.MustMarshalBinaryBare(acc2))

	accreditations := accreditationKeeper.GetAccreditations(ctx)
	assert.Equal(t, 2, len(accreditations))
	assert.Contains(t, accreditations, acc1)
	assert.Contains(t, accreditations, acc2)
}

// ---------------------
// --- Pool
// ---------------------

func TestKeeper_DepositIntoPool_EmptyPool(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Delete([]byte(types.LiquidityPoolKey))

	coins := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
	err := accreditationKeeper.DepositIntoPool(ctx, coins)
	assert.Nil(t, err)

	var pool sdk.Coins
	poolBz := store.Get([]byte(types.LiquidityPoolKey))
	cdc.MustUnmarshalBinaryBare(poolBz, &pool)
	assert.Equal(t, coins, pool)
}

func TestKeeper_DepositIntoPool_ExistingPool(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	pool := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Set([]byte(types.LiquidityPoolKey), cdc.MustMarshalBinaryBare(&pool))

	addition := sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000)))
	err := accreditationKeeper.DepositIntoPool(ctx, addition)
	assert.Nil(t, err)

	var actual sdk.Coins
	actualBz := store.Get([]byte(types.LiquidityPoolKey))
	cdc.MustUnmarshalBinaryBare(actualBz, &actual)

	expected := sdk.NewCoins(
		sdk.NewCoin("uatom", sdk.NewInt(100)),
		sdk.NewCoin("ucommercio", sdk.NewInt(1000)),
	)
	assert.Equal(t, expected, actual)
}

func TestKeeper_GetPoolFunds_EmptyPool(t *testing.T) {
	ctx, _, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Delete([]byte(types.LiquidityPoolKey))

	pool := accreditationKeeper.GetPoolFunds(ctx)

	assert.Empty(t, pool)
}

func TestKeeper_GetPoolFunds_ExistingPool(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	expected := sdk.NewCoins(
		sdk.NewCoin("uatom", sdk.NewInt(100)),
		sdk.NewCoin("ucommercio", sdk.NewInt(1000)),
	)

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Set([]byte(types.LiquidityPoolKey), cdc.MustMarshalBinaryBare(&expected))

	pool := accreditationKeeper.GetPoolFunds(ctx)

	assert.Equal(t, expected, pool)
}

// ---------------------
// --- Signers
// ---------------------

func TestKeeper_AddTrustworthySigner_EmptyList(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	accreditationKeeper.AddTrustworthySigner(ctx, TestSigner)

	var signers ctypes.Addresses
	signersBz := store.Get([]byte(types.TrustworthySignersKey))
	cdc.MustUnmarshalBinaryBare(signersBz, &signers)

	assert.Equal(t, 1, len(signers))
	assert.Contains(t, signers, TestSigner)
}

func TestKeeper_AddTrustworthySigner_ExistingList(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	signers := ctypes.Addresses{TestSigner}
	store.Set([]byte(types.TrustworthySignersKey), cdc.MustMarshalBinaryBare(&signers))

	accreditationKeeper.AddTrustworthySigner(ctx, TestUser)

	var actual ctypes.Addresses
	actualBz := store.Get([]byte(types.TrustworthySignersKey))
	cdc.MustUnmarshalBinaryBare(actualBz, &actual)

	assert.Equal(t, 2, len(actual))
	assert.Contains(t, actual, TestSigner)
	assert.Contains(t, actual, TestUser)
}

func TestKeeper_GetTrustworthySigners_EmptyList(t *testing.T) {
	ctx, _, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	signers := accreditationKeeper.GetTrustworthySigners(ctx)

	assert.Empty(t, signers)
}

func TestKeeper_GetTrustworthySigners_ExistingList(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	signers := ctypes.Addresses{TestSigner, TestUser, TestAccrediter}
	store.Set([]byte(types.TrustworthySignersKey), cdc.MustMarshalBinaryBare(&signers))

	actual := accreditationKeeper.GetTrustworthySigners(ctx)
	assert.Equal(t, 3, len(actual))
	assert.Contains(t, actual, TestSigner)
	assert.Contains(t, actual, TestUser)
	assert.Contains(t, actual, TestAccrediter)
}

func TestKeeper_IsTrustworthySigner_EmptyList(t *testing.T) {
	ctx, _, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	assert.False(t, accreditationKeeper.IsTrustworthySigner(ctx, TestSigner))
	assert.False(t, accreditationKeeper.IsTrustworthySigner(ctx, TestUser))
	assert.False(t, accreditationKeeper.IsTrustworthySigner(ctx, TestAccrediter))
}

func TestKeeper_IsTrustworthySigner_ExistingList(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	signers := ctypes.Addresses{TestUser, TestSigner}
	store.Set([]byte(types.TrustworthySignersKey), cdc.MustMarshalBinaryBare(&signers))

	assert.True(t, accreditationKeeper.IsTrustworthySigner(ctx, TestUser))
	assert.True(t, accreditationKeeper.IsTrustworthySigner(ctx, TestSigner))
	assert.False(t, accreditationKeeper.IsTrustworthySigner(ctx, TestAccrediter))
}

// --------------------------
// --- Reward distribution
// --------------------------

func TestKeeper_DistributeReward_NilReward(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	accreditation := types.Accreditation{Accrediter: TestAccrediter, User: TestUser, Rewarded: false}
	store.Set(TestUser, cdc.MustMarshalBinaryBare(&accreditation))

	var reward sdk.Coins
	err := accreditationKeeper.DistributeReward(ctx, TestAccrediter, reward, TestUser)

	assert.NotNil(t, err)
	assert.Equal(t, "reward cannot be empty", err.Error())
}

func TestKeeper_DistributeReward_NegativeReward(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	accreditation := types.Accreditation{Accrediter: TestAccrediter, User: TestUser, Rewarded: false}
	store.Set(TestUser, cdc.MustMarshalBinaryBare(&accreditation))

	reward := sdk.Coins{
		sdk.Coin{Denom: "uatom", Amount: sdk.NewInt(100)},
		sdk.Coin{Denom: "ucommercio", Amount: sdk.NewInt(-100)},
	}
	err := accreditationKeeper.DistributeReward(ctx, TestAccrediter, reward, TestUser)

	assert.NotNil(t, err)
	assert.Equal(t, "rewards cannot be negative", err.Error())
}

func TestKeeper_DistributeReward_NoAccrediter(t *testing.T) {
	ctx, _, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Delete(TestUser)

	reward := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
	err := accreditationKeeper.DistributeReward(ctx, TestAccrediter, reward, TestUser)

	assert.NotNil(t, err)
	assert.Equal(t, "user does not have an accrediter", err.Error())
}

func TestKeeper_DistributeReward_AlreadyRewarded(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	accreditation := types.Accreditation{Accrediter: TestAccrediter, User: TestUser, Rewarded: true}
	store.Set(TestUser, cdc.MustMarshalBinaryBare(&accreditation))

	reward := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
	err := accreditationKeeper.DistributeReward(ctx, TestAccrediter, reward, TestUser)

	assert.NotNil(t, err)
	assert.Equal(t, "the accrediter has already been rewarded for this user", err.Error())
}

func TestKeeper_DistributeReward_EmptyPool(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)
	store.Delete([]byte(types.LiquidityPoolKey))

	accreditation := types.Accreditation{Accrediter: TestAccrediter, User: TestUser, Rewarded: false}
	store.Set(TestUser, cdc.MustMarshalBinaryBare(&accreditation))

	reward := sdk.NewCoins(
		sdk.NewCoin("uatom", sdk.NewInt(100)),
	)
	err := accreditationKeeper.DistributeReward(ctx, TestAccrediter, reward, TestUser)

	assert.NotNil(t, err)
	assert.Equal(t, "liquidity pool has not a sufficient amount of tokens for this reward", err.Error())
}

func TestKeeper_DistributeReward_InsufficientPool(t *testing.T) {
	ctx, cdc, _, _, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	pool := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(10)))
	store.Set([]byte(types.LiquidityPoolKey), cdc.MustMarshalBinaryBare(&pool))

	accreditation := types.Accreditation{Accrediter: TestAccrediter, User: TestUser, Rewarded: false}
	store.Set(TestUser, cdc.MustMarshalBinaryBare(&accreditation))

	reward := sdk.NewCoins(
		sdk.NewCoin("uatom", sdk.NewInt(100)),
	)
	err := accreditationKeeper.DistributeReward(ctx, TestAccrediter, reward, TestUser)

	assert.NotNil(t, err)
	assert.Equal(t, "liquidity pool has not a sufficient amount of tokens for this reward", err.Error())
}

func TestKeeper_DistributeReward_ValidData(t *testing.T) {
	ctx, cdc, _, bankKeeper, accreditationKeeper := getTestInput()

	store := ctx.KVStore(accreditationKeeper.StoreKey)

	pool := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(1000)))
	store.Set([]byte(types.LiquidityPoolKey), cdc.MustMarshalBinaryBare(&pool))

	accreditation := types.Accreditation{Accrediter: TestAccrediter, User: TestUser, Rewarded: false}
	store.Set(TestUser, cdc.MustMarshalBinaryBare(&accreditation))

	reward := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
	err := accreditationKeeper.DistributeReward(ctx, TestAccrediter, reward, TestUser)

	assert.Nil(t, err)

	// Check the pool
	var remainingPool sdk.Coins
	cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.LiquidityPoolKey)), &remainingPool)
	assert.Equal(t, 1, len(remainingPool))
	assert.Contains(t, remainingPool, sdk.NewCoin("uatom", sdk.NewInt(900)))

	// Check the user funds
	accountCoins := bankKeeper.GetCoins(ctx, TestAccrediter)
	assert.Equal(t, 1, len(accountCoins))
	assert.Contains(t, accountCoins, sdk.NewCoin("uatom", sdk.NewInt(100)))

	// Check the accreditation
	var newAccreditation types.Accreditation
	cdc.MustUnmarshalBinaryBare(store.Get(TestUser), &newAccreditation)
	assert.Equal(t, TestUser, newAccreditation.User)
	assert.Equal(t, TestAccrediter, newAccreditation.Accrediter)
	assert.True(t, newAccreditation.Rewarded)
}
