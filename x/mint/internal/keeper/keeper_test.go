package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestKeeper_SetCreditsDenom(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	denom := "test"
	k.SetCreditsDenom(ctx, denom)
	store := ctx.KVStore(k.StoreKey)
	denomBz := store.Get([]byte(types.CreditsDenom))
	assert.Equal(t, denom, string(denomBz))
}

func TestKeeper_GetCreditsDenom(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	denom := "test"
	k.SetCreditsDenom(ctx, denom)
	actual := k.GetCreditsDenom(ctx)
	assert.Equal(t, denom, actual)
}

func TestKeeper_GetCDPkey(t *testing.T) {
	_, _, _, k := SetupTestInput()
	expected := types.CDPSPrefix + TestOwner.String()
	actual := k.GetCDPkey(TestOwner)
	assert.Equal(t, []byte(expected), actual)
}

func TestKeeper_GetUsersSet(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)
	store.Set(k.GetCDPkey(TestOwner), k.Cdc.MustMarshalBinaryBare(TestCdp))
	users := k.GetUsersSet(ctx)
	assert.Equal(t, TestOwner, users[0])
}

func TestKeeper_AddCDP_notAlreadyAdded(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	var cdps types.CDPs
	k.AddCDP(ctx, TestCdp)

	store := ctx.KVStore(k.StoreKey)
	actualBz := store.Get(k.GetCDPkey(TestOwner))
	k.Cdc.MustUnmarshalBinaryBare(actualBz, &cdps)

	assert.Equal(t, TestCdp, cdps[0])
}

func TestKeeper_AddCDP_AlreadyAdded(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	var cdps types.CDPs
	k.AddCDP(ctx, TestCdp)
	k.AddCDP(ctx, TestCdp)

	store := ctx.KVStore(k.StoreKey)
	actualBz := store.Get(k.GetCDPkey(TestOwner))
	k.Cdc.MustUnmarshalBinaryBare(actualBz, &cdps)

	assert.Len(t, cdps, 1)
}

func TestKeeper_GetCDPs(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	var expected = types.CDPs{TestCdp}
	k.AddCDP(ctx, TestCdp)

	actual := k.GetCDPs(ctx, TestOwner)
	assert.Equal(t, expected, actual)
}

func TestKeeper_GetCDPs_empty(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	var expected = types.CDPs(nil)

	actual := k.GetCDPs(ctx, TestOwner)
	assert.Equal(t, expected, actual)
}

func TestKeeper_DeleteCDP_deleted(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	k.AddCDP(ctx, TestCdp)
	actual := k.DeleteCDP(ctx, TestCdp)
	assert.True(t, actual)
}

func TestKeeper_DeleteCDP_notDeleted(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	actual := k.DeleteCDP(ctx, TestCdp)
	assert.False(t, actual)
}

func TestKeeper_setBlockRewardsPool_UtilityFunction(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	var pool sdk.Coins

	k.SetLiquidityPool(ctx, TestLiquidityPool)
	store := ctx.KVStore(k.StoreKey)
	poolBz := store.Get([]byte(types.LiquidityPoolPrefix))
	k.Cdc.MustUnmarshalBinaryBare(poolBz, &pool)

	assert.Equal(t, pool, TestLiquidityPool)
}

func TestKeeper_GetBlockRewardsPool(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	k.SetLiquidityPool(ctx, TestLiquidityPool)
	actual := k.GetLiquidityPool(ctx)

	assert.Equal(t, TestLiquidityPool, actual)
}

func TestKeeper_OpenCDP_InvalidDepositedAmount(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	invalidReq := types.CDPRequest{
		Signer:          TestOwner,
		DepositedAmount: sdk.NewCoins(sdk.NewInt64Coin("testcoin", 0)),
		Timestamp:       "",
	}
	err := k.OpenCDP(ctx, invalidReq)
	assert.Error(t, err)
	expected := sdk.ErrInvalidCoins(invalidReq.DepositedAmount.String())
	assert.Equal(t, expected, err)
}

func TestKeeper_OpenCDP_TokenPriceNotFound(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	err := k.OpenCDP(ctx, TestCdpRequest)
	expected := sdk.ErrInvalidCoins("no current price for given token: ucommercio")
	assert.Error(t, err)
	assert.Equal(t, expected, err)
}

func TestKeeper_OpenCDP_NotEnoughFundsInUsersWallet(t *testing.T) {
	ctx, _, pfk, k := SetupTestInput()

	//tests setup
	store := ctx.KVStore(pfk.StoreKey)
	store.Set([]byte("pricefeed:currentPrices:ucommercio"), k.Cdc.MustMarshalBinaryBare(TestCurrentPrice))
	k.SetCreditsDenom(ctx, "uccc")

	err := k.OpenCDP(ctx, TestCdpRequest)

	assert.Error(t, err)
}

func TestKeeper_OpenCDP_Successfully(t *testing.T) {
	ctx, bk, pfk, k := SetupTestInput()

	_, _ = bk.AddCoins(ctx, TestOwner, TestDepositedAmount)
	store := ctx.KVStore(pfk.StoreKey)
	store.Set([]byte("pricefeed:currentPrices:ucommercio"), k.Cdc.MustMarshalBinaryBare(TestCurrentPrice))
	k.SetCreditsDenom(ctx, "uccc")

	err := k.OpenCDP(ctx, TestCdpRequest)

	assert.Nil(t, err)
	//check if the correct amount of credits has been transfered to user's wallet
	creditsAmount := TestLiquidityAmount.AmountOf("uccc").Quo(sdk.NewInt(2))
	credits := sdk.NewCoins(sdk.NewCoin("uccc", creditsAmount))
	creditsTransfered := bk.HasCoins(ctx, TestOwner, credits)

	assert.True(t, creditsTransfered)
}

func TestKeeper_CloseCDP_InexistentCDP(t *testing.T) {
	ctx, _, _, k := SetupTestInput()

	err := k.CloseCDP(ctx, TestOwner, TestTimestamp)
	expected := sdk.ErrInternal("cannot close an inexistent cdp")
	assert.Error(t, err)
	assert.Equal(t, expected, err)
}

func TestKeeper_CloseCDP_Successfully(t *testing.T) {
	ctx, bk, _, k := SetupTestInput()

	k.AddCDP(ctx, TestCdp)
	k.SetLiquidityPool(ctx, TestLiquidityPool)

	_, _ = bk.AddCoins(ctx, TestOwner, TestLiquidityAmount)

	err := k.CloseCDP(ctx, TestOwner, TestTimestamp)
	assert.Nil(t, err)
	creditsTransfered := bk.HasCoins(ctx, TestOwner, TestDepositedAmount)
	assert.True(t, creditsTransfered)
}
