package keeper

import (
	"testing"
	"time"

	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// --------------
// --- Credits
// --------------

func TestKeeper_SetCreditsDenom(t *testing.T) {
	_, ctx, _, _, k := SetupTestInput()

	denom := "test"
	k.SetCreditsDenom(ctx, denom)

	store := ctx.KVStore(k.storeKey)
	denomBz := store.Get([]byte(types.CreditsDenomStoreKey))
	assert.Equal(t, denom, string(denomBz))
}

func TestKeeper_GetCreditsDenom(t *testing.T) {
	_, ctx, _, _, k := SetupTestInput()
	denom := "test"
	k.SetCreditsDenom(ctx, denom)
	actual := k.GetCreditsDenom(ctx)
	assert.Equal(t, denom, actual)
}

// --------------
// --- CDPs
// --------------

func TestKeeper_AddCdp_Existing(t *testing.T) {
	_, ctx, _, _, k := SetupTestInput()
	var cdps types.Cdps
	k.AddCdp(ctx, TestCdp)
	k.AddCdp(ctx, TestCdp)

	store := ctx.KVStore(k.storeKey)
	actualBz := store.Get(k.getCdpKey(TestOwner))
	k.cdc.MustUnmarshalBinaryBare(actualBz, &cdps)

	assert.Len(t, cdps, 1)
}

func TestKeeper_AddCdp_NotExisting(t *testing.T) {
	_, ctx, _, _, k := SetupTestInput()
	var cdps types.Cdps
	k.AddCdp(ctx, TestCdp)

	store := ctx.KVStore(k.storeKey)
	actualBz := store.Get(k.getCdpKey(TestOwner))
	k.cdc.MustUnmarshalBinaryBare(actualBz, &cdps)

	assert.Equal(t, TestCdp, cdps[0])
}

func TestKeeper_OpenCdp_InvalidDepositedAmount(t *testing.T) {
	_, ctx, _, _, k := SetupTestInput()

	invalidReq := types.CdpRequest{
		Signer:          TestOwner,
		DepositedAmount: sdk.NewCoins(sdk.NewInt64Coin("testcoin", 0)),
		Timestamp:       time.Time{},
	}
	err := k.OpenCdp(ctx, invalidReq)

	assert.Error(t, err)
	assert.Equal(t, sdk.ErrInvalidCoins(invalidReq.DepositedAmount.String()), err)
}

func TestKeeper_OpenCdp_TokenPriceNotFound(t *testing.T) {
	_, ctx, _, _, k := SetupTestInput()

	err := k.OpenCdp(ctx, TestCdpRequest)
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeUnknownRequest, err.Code())
}

func TestKeeper_OpenCdp_NotEnoughFundsInUsersWallet(t *testing.T) {
	_, ctx, _, pfk, k := SetupTestInput()

	pfk.SetCurrentPrice(ctx, pricefeed.NewCurrentPrice("ucommercio", 10, 1000))
	k.SetCreditsDenom(ctx, "uccc")

	err := k.OpenCdp(ctx, TestCdpRequest)
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeInsufficientCoins, err.Code())
}

func TestKeeper_OpenCdp_Successful(t *testing.T) {
	_, ctx, bk, pfk, k := SetupTestInput()

	// Setup
	_ = bk.SetCoins(ctx, TestOwner, TestCdpRequest.DepositedAmount)
	pfk.SetCurrentPrice(ctx, pricefeed.NewCurrentPrice(TestLiquidityDenom, 10, 1000))

	// Cdp opening
	err := k.OpenCdp(ctx, TestCdpRequest)
	assert.Nil(t, err)

	// Check that the correct amount of credits has been transferred to the user's wallet
	expected := sdk.NewCoins(sdk.NewInt64Coin(TestCreditsDenom, 10*50))
	actual := bk.GetCoins(ctx, TestCdpRequest.Signer)
	assert.Equal(t, expected, actual)
}

func TestKeeper_GetCdpsByOwner_EmptyList(t *testing.T) {
	_, ctx, _, _, k := SetupTestInput()

	actual := k.GetCdpsByOwner(ctx, TestOwner)
	assert.Empty(t, actual)
}

func TestKeeper_GetCdpsByOwner_NonEmptyList(t *testing.T) {
	cdc, ctx, _, _, k := SetupTestInput()

	k.AddCdp(ctx, TestCdp)

	store := ctx.KVStore(k.storeKey)
	var cdps types.Cdps
	cdc.MustUnmarshalBinaryBare(store.Get(k.getCdpKey(TestCdp.Owner)), &cdps)
	assert.Equal(t, types.Cdps{TestCdp}, k.GetCdpsByOwner(ctx, TestOwner))
}

func TestKeeper_CloseCdp_NonExistentCdp(t *testing.T) {
	_, ctx, _, _, k := SetupTestInput()

	err := k.CloseCdp(ctx, TestOwner, TestCdpRequest.Timestamp)
	assert.Error(t, err)
	assert.Equal(t, sdk.CodeUnknownRequest, err.Code())
	assert.Contains(t, err.Error(), "does not exist")
}

func TestKeeper_CloseCdp_ExistentCdp(t *testing.T) {
	_, ctx, bk, _, k := SetupTestInput()

	k.AddCdp(ctx, TestCdp)
	_ = k.supplyKeeper.MintCoins(ctx, types.ModuleName, TestLiquidityPool)
	_, _ = bk.AddCoins(ctx, TestOwner, TestCdp.CreditsAmount)

	err := k.CloseCdp(ctx, TestOwner, TestCdp.Timestamp)
	assert.Nil(t, err)

	assert.Equal(t, TestCdp.DepositedAmount, bk.GetCoins(ctx, TestOwner))
}

func TestKeeper_DeleteCdp_ExistentCdp(t *testing.T) {
	_, ctx, _, _, k := SetupTestInput()

	k.AddCdp(ctx, TestCdp)
	k.deleteCdp(ctx, TestCdp)

	assert.NotContains(t, k.GetTotalCdps(ctx), TestCdp)
}

func TestKeeper_DeleteCdp_NonExistentCdp(t *testing.T) {
	_, ctx, _, _, k := SetupTestInput()
	k.deleteCdp(ctx, TestCdp)
	assert.NotContains(t, k.GetTotalCdps(ctx), TestCdp)
}
