package keeper

import (
	"testing"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// --------------
// --- Assets
// --------------

func TestKeeper_AddAsset(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	k.AddAsset(ctx, "ucommercio")

	expected := ctypes.Strings{"ucommercio"}
	actual := k.GetAssets(ctx)

	assert.Equal(t, expected, actual)
}

func TestKeeper_AddAsset_AlreadyPresent(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	assets := ctypes.Strings{"ucommercio"}
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.AssetsStoreKey), k.cdc.MustMarshalBinaryBare(assets))

	k.AddAsset(ctx, "ucommercio")
	actual := k.GetAssets(ctx)

	assert.Len(t, actual, 1)
}

func TestKeeper_GetAssets(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	expected := ctypes.Strings{"ucommercio", "uccc"}

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.AssetsStoreKey), k.cdc.MustMarshalBinaryBare(expected))

	actual := k.GetAssets(ctx)
	assert.Equal(t, expected, actual)
}

func TestKeeper_GetAssets_EmptyList(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	var expected ctypes.Strings
	actual := k.GetAssets(ctx)

	assert.Equal(t, expected, actual)
}

// -----------------
// --- Raw prices
// -----------------

func TestKeeper_SetRawPrice_withValidSigner_PricesNotAlreadyPresent(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	// Add prices
	const assetName = "test"
	testPrice1 := types.Price{AssetName: assetName, Value: sdk.NewDec(10), Expiry: sdk.NewInt(5000)}
	assert.NoError(t, k.AddRawPrice(ctx, testPrice1))

	testPrice2 := types.Price{AssetName: assetName, Value: sdk.NewDec(8), Expiry: sdk.NewInt(4000)}
	assert.NoError(t, k.AddRawPrice(ctx, testPrice2))

	// List prices
	actual := k.GetRawPrices(ctx, assetName)
	expected := types.Prices{testPrice1, testPrice2}
	assert.Equal(t, expected, actual)
}

func TestKeeper_SetRawPrice_withValidSigner_PriceAlreadyPresent(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	// Add the price
	testPrice := types.Price{AssetName: "test", Value: sdk.NewDec(10), Expiry: sdk.NewInt(5000)}
	store := ctx.KVStore(k.StoreKey)
	store.Set(k.getRawPricesKey(testPrice.AssetName), k.cdc.MustMarshalBinaryBare(types.Prices{testPrice}))

	// Try adding the price again
	err := k.AddRawPrice(ctx, testPrice)
	assert.Error(t, err)
}

func TestKeeper_GetRawPrices(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	// Add prices
	const assetName = "test"
	testPrice1 := types.Price{AssetName: assetName, Value: sdk.NewDec(10), Expiry: sdk.NewInt(5000)}
	assert.NoError(t, k.AddRawPrice(ctx, testPrice1))

	testPrice2 := types.Price{AssetName: assetName, Value: sdk.NewDec(8), Expiry: sdk.NewInt(4000)}
	assert.NoError(t, k.AddRawPrice(ctx, testPrice2))

	// List prices
	actual := k.GetRawPrices(ctx, assetName)
	expected := types.Prices{testPrice1, testPrice2}
	assert.Equal(t, expected, actual)
}

// ---------------------
// --- Current prices
// ---------------------

func TestKeeper_SetCurrentPrices_MoreThanOneNotExpiredPrice(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	// Add prices
	const assetName = "test"
	testPrice1 := types.Price{AssetName: assetName, Value: sdk.NewDec(10), Expiry: sdk.NewInt(5000)}
	assert.NoError(t, k.AddRawPrice(ctx, testPrice1))

	testPrice2 := types.Price{AssetName: assetName, Value: sdk.NewDec(8), Expiry: sdk.NewInt(4000)}
	assert.NoError(t, k.AddRawPrice(ctx, testPrice2))

	sumPrice := testPrice1.Value.Add(testPrice2.Value)
	sumExpiry := testPrice1.Expiry.Add(testPrice2.Expiry)
	expectedMedianPrice := sumPrice.Quo(sdk.NewDec(2))
	expectedMedianExpiry := sumExpiry.Quo(sdk.NewInt(2))

	_ = k.ComputeAndUpdateCurrentPrices(ctx)

	actual, found := k.GetCurrentPrice(ctx, assetName)

	assert.True(t, found)
	assert.Equal(t, expectedMedianPrice, actual.Value)
	assert.Equal(t, expectedMedianExpiry, actual.Expiry)
}

func TestKeeper_SetCurrentPrices_AllExpiredRawPrices(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	price := types.Price{AssetName: "uccc", Value: sdk.NewDec(20), Expiry: sdk.NewInt(-1)}
	_ = k.AddRawPrice(ctx, price)

	err := k.ComputeAndUpdateCurrentPrices(ctx)
	assert.Error(t, err)
}

func TestKeeper_SetCurrentPrice_OneNotExpiredPrice(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	const assetName = "test"
	testPrice := types.Price{AssetName: assetName, Value: sdk.NewDec(10), Expiry: sdk.NewInt(5000)}
	assert.NoError(t, k.AddRawPrice(ctx, testPrice))

	_ = k.ComputeAndUpdateCurrentPrices(ctx)

	actual, _ := k.GetCurrentPrice(ctx, assetName)
	assert.Equal(t, testPrice.Value, actual.Value)
	assert.Equal(t, testPrice.Expiry, actual.Expiry)
}

func TestKeeper_GetCurrentPrices(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.CurrentPricesPrefix+TestPrice.AssetName), k.cdc.MustMarshalBinaryBare(TestPrice))

	assert.Equal(t, types.Prices{TestPrice}, k.GetCurrentPrices(ctx))
}

func TestKeeper_GetCurrentPrice_Found(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.CurrentPricesPrefix+TestPrice.AssetName), k.cdc.MustMarshalBinaryBare(TestPrice))

	actual, _ := k.GetCurrentPrice(ctx, TestPrice.AssetName)
	assert.Equal(t, TestPrice, actual)
}

func TestKeeper_GetCurrentPrice_NotFound(t *testing.T) {
	_, ctx, _, k := SetupTestInput()
	_, found := k.GetCurrentPrice(ctx, TestPrice.AssetName)
	assert.False(t, found)
}

// ------------------
// --- Oracles
// ------------------

func TestKeeper_AddOracle(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	testOracle, _ := sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	expected := ctypes.Addresses{testOracle}

	store := ctx.KVStore(k.StoreKey)
	oracleBz := store.Get([]byte(types.OraclePrefix))
	assert.Nil(t, oracleBz)

	k.AddOracle(ctx, testOracle)
	oracleBz = store.Get([]byte(types.OraclePrefix))
	var actual ctypes.Addresses
	k.cdc.MustUnmarshalBinaryBare(oracleBz, &actual)

	assert.Equal(t, expected, actual)
}

func TestKeeper_IsOracle_ValidOracle(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	testOracle, _ := sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.OraclePrefix), k.cdc.MustMarshalBinaryBare(ctypes.Addresses{testOracle}))

	isOracle := k.IsOracle(ctx, testOracle)
	assert.True(t, isOracle)
}

func TestKeeper_IsOracle_InvalidOracle(t *testing.T) {
	_, ctx, _, k := SetupTestInput()
	testOracle, _ := sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	isOracle := k.IsOracle(ctx, testOracle)
	assert.False(t, isOracle)
}

func TestKeeper_GetOracles(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	testOracle, _ := sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	expected := ctypes.Addresses{testOracle}

	k.AddOracle(ctx, testOracle)

	assert.Equal(t, expected, k.GetOracles(ctx))
}
