package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/accreditations/internal/types"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestKeeper_SetAccrediter_NoAccrediter(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)
	store.Delete(TestUser)

	err := TestUtils.AcKeeper.SetAccrediter(TestUtils.Ctx, TestUser, TestAccrediter)
	assert.Nil(t, err)

	accreditationBz := store.Get(TestUser)
	assert.NotNil(t, accreditationBz)

	var accreditation types.Accreditation
	TestUtils.Cdc.MustUnmarshalBinaryBare(accreditationBz, &accreditation)

	assert.Equal(t, TestUser, accreditation.User)
	assert.Equal(t, TestAccrediter, accreditation.Accrediter)
	assert.False(t, accreditation.Rewarded)
}

func TestKeeper_SetAccrediter_ExistingAccrediter(t *testing.T) {
	existingAccredit := types.Accreditation{
		Accrediter: TestAccrediter,
		User:       TestUser,
		Rewarded:   false,
	}

	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)
	store.Set(TestUser, TestUtils.Cdc.MustMarshalBinaryBare(existingAccredit))

	err := TestUtils.AcKeeper.SetAccrediter(TestUtils.Ctx, TestUser, TestAccrediter)
	assert.NotNil(t, err)
}

func TestKeeper_GetAccrediter_NoAccrediter(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)
	store.Delete(TestUser)

	accrediter := TestUtils.AcKeeper.GetAccrediter(TestUtils.Ctx, TestUser)
	assert.Nil(t, accrediter)
}

func TestKeeper_GetAccrediter_ExistingAccrediter(t *testing.T) {
	accreditation := types.Accreditation{
		Accrediter: TestAccrediter,
		User:       TestUser,
		Rewarded:   false,
	}

	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)
	store.Set(TestUser, TestUtils.Cdc.MustMarshalBinaryBare(accreditation))

	accrediter := TestUtils.AcKeeper.GetAccrediter(TestUtils.Ctx, TestUser)
	assert.Equal(t, TestAccrediter, accrediter)
}

func TestKeeper_GetAccreditations_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	accreditations := TestUtils.AcKeeper.GetAccreditations(TestUtils.Ctx)
	assert.Empty(t, accreditations)
}

func TestKeeper_GetAccreditations_NonEmptyList(t *testing.T) {
	acc1 := types.Accreditation{Accrediter: TestAccrediter, User: TestUser, Rewarded: false}
	acc2 := types.Accreditation{Accrediter: TestUser, User: TestAccrediter, Rewarded: false}

	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)
	store.Set(TestUser, TestUtils.Cdc.MustMarshalBinaryBare(acc1))
	store.Set(TestAccrediter, TestUtils.Cdc.MustMarshalBinaryBare(acc2))

	accreditations := TestUtils.AcKeeper.GetAccreditations(TestUtils.Ctx)
	assert.Equal(t, 2, len(accreditations))
	assert.Contains(t, accreditations, acc1)
	assert.Contains(t, accreditations, acc2)
}

func TestKeeper_DepositIntoPool_EmptyPool(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)
	store.Delete([]byte(types.LiquidityPoolKey))

	coins := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))
	err := TestUtils.AcKeeper.DepositIntoPool(TestUtils.Ctx, coins)
	assert.Nil(t, err)

	var pool sdk.Coins
	poolBz := store.Get([]byte(types.LiquidityPoolKey))
	TestUtils.Cdc.MustUnmarshalBinaryBare(poolBz, &pool)
	assert.Equal(t, coins, pool)
}

func TestKeeper_DepositIntoPool_ExistingPool(t *testing.T) {
	pool := sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100)))

	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)
	store.Set([]byte(types.LiquidityPoolKey), TestUtils.Cdc.MustMarshalBinaryBare(&pool))

	addition := sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(1000)))
	err := TestUtils.AcKeeper.DepositIntoPool(TestUtils.Ctx, addition)
	assert.Nil(t, err)

	var actual sdk.Coins
	actualBz := store.Get([]byte(types.LiquidityPoolKey))
	TestUtils.Cdc.MustUnmarshalBinaryBare(actualBz, &actual)

	expected := sdk.NewCoins(
		sdk.NewCoin("uatom", sdk.NewInt(100)),
		sdk.NewCoin("ucommercio", sdk.NewInt(1000)),
	)
	assert.Equal(t, expected, actual)
}

func TestKeeper_AddTrustworthySigner_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	TestUtils.AcKeeper.AddTrustworthySigner(TestUtils.Ctx, TestSigner)

	var signers ctypes.Addresses
	signersBz := store.Get([]byte(types.TrustworthySignersKey))
	TestUtils.Cdc.MustUnmarshalBinaryBare(signersBz, &signers)

	assert.Equal(t, 1, len(signers))
	assert.Contains(t, signers, TestSigner)
}

func TestKeeper_AddTrustworthySigner_ExistingList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	signers := ctypes.Addresses{TestSigner}
	store.Set([]byte(types.TrustworthySignersKey), TestUtils.Cdc.MustMarshalBinaryBare(&signers))

	TestUtils.AcKeeper.AddTrustworthySigner(TestUtils.Ctx, TestUser)

	var actual ctypes.Addresses
	actualBz := store.Get([]byte(types.TrustworthySignersKey))
	TestUtils.Cdc.MustUnmarshalBinaryBare(actualBz, &actual)

	assert.Equal(t, 2, len(actual))
	assert.Contains(t, actual, TestSigner)
	assert.Contains(t, actual, TestUser)
}

func TestKeeper_GetTrustworthySigners_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)

	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	signers := TestUtils.AcKeeper.GetTrustworthySigners(TestUtils.Ctx)

	assert.Empty(t, signers)
}

func TestKeeper_GetTrustworthySigners_ExistingList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	signers := ctypes.Addresses{TestSigner, TestUser, TestAccrediter}
	store.Set([]byte(types.TrustworthySignersKey), TestUtils.Cdc.MustMarshalBinaryBare(&signers))

	actual := TestUtils.AcKeeper.GetTrustworthySigners(TestUtils.Ctx)
	assert.Equal(t, 3, len(actual))
	assert.Contains(t, actual, TestSigner)
	assert.Contains(t, actual, TestUser)
	assert.Contains(t, actual, TestAccrediter)
}

func TestKeeper_IsTrustworthySigner_EmptyList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	assert.False(t, TestUtils.AcKeeper.IsTrustworthySigner(TestUtils.Ctx, TestSigner))
	assert.False(t, TestUtils.AcKeeper.IsTrustworthySigner(TestUtils.Ctx, TestUser))
	assert.False(t, TestUtils.AcKeeper.IsTrustworthySigner(TestUtils.Ctx, TestAccrediter))
}

func TestKeeper_IsTrustworthySigner_ExistingList(t *testing.T) {
	store := TestUtils.Ctx.KVStore(TestUtils.AcKeeper.StoreKey)
	iterator := store.Iterator(nil, nil)
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}

	signers := ctypes.Addresses{TestUser, TestSigner}
	store.Set([]byte(types.TrustworthySignersKey), TestUtils.Cdc.MustMarshalBinaryBare(&signers))

	assert.True(t, TestUtils.AcKeeper.IsTrustworthySigner(TestUtils.Ctx, TestUser))
	assert.True(t, TestUtils.AcKeeper.IsTrustworthySigner(TestUtils.Ctx, TestSigner))
	assert.False(t, TestUtils.AcKeeper.IsTrustworthySigner(TestUtils.Ctx, TestAccrediter))
}
