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
	_, ctx, k := SetupTestInput()
	//ensure that assets array are empty
	assets := k.GetAssets(ctx)
	assert.Nil(t, assets)

	k.AddAsset(ctx, TestAsset)
	expected := ctypes.Strings{TestAsset}
	actual := k.GetAssets(ctx)
	assert.Equal(t, expected, actual)
}

func TestKeeper_AddAsset_AlreadyPresent(t *testing.T) {
	_, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)
	assets := ctypes.Strings{TestAsset}
	store.Set([]byte(types.AssetsStoreKey), k.cdc.MustMarshalBinaryBare(assets))
	expected := len(assets)

	k.AddAsset(ctx, TestAsset)
	actual := len(k.GetAssets(ctx))

	assert.Equal(t, expected, actual)
}

func TestKeeper_GetAssets(t *testing.T) {
	_, ctx, k := SetupTestInput()
	expected := ctypes.Strings{TestAsset, TestAsset2}

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.AssetsStoreKey), k.cdc.MustMarshalBinaryBare(expected))

	actual := k.GetAssets(ctx)
	assert.Equal(t, expected, actual)
}

func TestKeeper_GetAssets_EmptyList(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var expected ctypes.Strings
	actual := k.GetAssets(ctx)
	assert.Equal(t, expected, actual)
}

// -----------------
// --- Raw prices
// -----------------

func TestKeeper_SetRawPrice_withValidSigner_PricesNotAlreadyPresent(t *testing.T) {
	_, ctx, k := SetupTestInput()
	k.AddOracle(ctx, TestOracle1)
	k.AddOracle(ctx, TestOracle2)

	err := k.SetRawPrice(ctx, TestRawPrice1)
	assert.Nil(t, err)
	err = k.SetRawPrice(ctx, TestRawPrice3)
	assert.Nil(t, err)

	actual := k.GetRawPrices(ctx, TestRawPrice1.PriceInfo.AssetName)

	expected := types.RawPrices{TestRawPrice1, TestRawPrice3}

	assert.Equal(t, expected, actual)

}

func TestKeeper_SetRawPrice_withValidSigner_PriceAlreadyPresent(t *testing.T) {
	_, ctx, k := SetupTestInput()

	store := ctx.KVStore(k.StoreKey)

	k.AddOracle(ctx, TestOracle1)

	expected := types.RawPrices{TestRawPrice1}
	store.Set(k.getRawPricesKey(TestRawPrice1.PriceInfo.AssetName), k.cdc.MustMarshalBinaryBare(&expected))

	err := k.SetRawPrice(ctx, TestRawPrice1)
	assert.Error(t, err)
}

func TestKeeper_SetRawPrice_WithInvalidSigner(t *testing.T) {
	_, ctx, k := SetupTestInput()
	err := k.SetRawPrice(ctx, TestRawPrice1)
	assert.Error(t, err)
}

func TestKeeper_GetRawPrices(t *testing.T) {
	_, ctx, k := SetupTestInput()
	k.AddOracle(ctx, TestOracle1)
	k.AddOracle(ctx, TestOracle2)
	_ = k.SetRawPrice(ctx, TestRawPrice1)
	_ = k.SetRawPrice(ctx, TestRawPrice3)
	actual := k.GetRawPrices(ctx, TestRawPrice1.PriceInfo.AssetName)
	expected := types.RawPrices{TestRawPrice1, TestRawPrice3}
	assert.Equal(t, expected, actual)
}

// ---------------------
// --- Current prices
// ---------------------

func TestKeeper_SetCurrentPrices_MoreThanOneNotExpiredPrice(t *testing.T) {
	_, ctx, k := SetupTestInput()

	k.AddOracle(ctx, TestOracle1)
	k.AddOracle(ctx, TestOracle2)
	_ = k.SetRawPrice(ctx, TestRawPrice1)
	_ = k.SetRawPrice(ctx, TestRawPrice3)

	sumPrice := TestRawPrice1.PriceInfo.Price.Add(TestRawPrice3.PriceInfo.Price)
	sumExpiry := TestRawPrice1.PriceInfo.Expiry.Add(TestRawPrice3.PriceInfo.Expiry)
	expectedMedianPrice := sumPrice.Quo(sdk.NewDec(2))
	expectedMedianExpiry := sumExpiry.Quo(sdk.NewInt(2))

	_ = k.SetCurrentPrices(ctx)

	actual, found := k.GetCurrentPrice(ctx, TestRawPrice1.PriceInfo.AssetName)
	assert.True(t, found)
	assert.Equal(t, expectedMedianPrice, actual.Price)
	assert.Equal(t, expectedMedianExpiry, actual.Expiry)
}

func TestKeeper_SetCurrentPrices_AllExpiredRawPrices(t *testing.T) {
	_, ctx, k := SetupTestInput()
	k.AddOracle(ctx, TestOracle1)
	_ = k.SetRawPrice(ctx, TestRawPriceE)
	err := k.SetCurrentPrices(ctx)
	assert.Error(t, err)
}

func TestKeeper_SetCurrentPrice_OneNotExpiredPrice(t *testing.T) {
	_, ctx, k := SetupTestInput()
	k.AddOracle(ctx, TestOracle1)
	_ = k.SetRawPrice(ctx, TestRawPrice1)
	_ = k.SetCurrentPrices(ctx)
	actual, _ := k.GetCurrentPrice(ctx, TestRawPrice1.PriceInfo.AssetName)
	assert.Equal(t, TestRawPrice1.PriceInfo.Price, actual.Price)
	assert.Equal(t, TestRawPrice1.PriceInfo.Expiry, actual.Expiry)
}

func TestKeeper_GetCurrentPrices(t *testing.T) {
	_, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.CurrentPricesPrefix+TestPriceInfo.AssetName),
		k.cdc.MustMarshalBinaryBare(TestPriceInfo))
	actual := k.GetCurrentPrices(ctx)
	expected := types.CurrentPrices{TestPriceInfo}
	assert.Equal(t, expected, actual)
}

func TestKeeper_GetCurrentPrice_Found(t *testing.T) {
	_, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.CurrentPricesPrefix+TestPriceInfo.AssetName),
		k.cdc.MustMarshalBinaryBare(TestPriceInfo))
	actual, _ := k.GetCurrentPrice(ctx, TestPriceInfo.AssetName)
	assert.Equal(t, TestPriceInfo, actual)
}

func TestKeeper_GetCurrentPrice_NotFound(t *testing.T) {
	_, ctx, k := SetupTestInput()
	_, found := k.GetCurrentPrice(ctx, TestPriceInfo.AssetName)
	assert.False(t, found)
}

// ------------------
// --- Oracles
// ------------------

func TestKeeper_AddOracle(t *testing.T) {
	_, ctx, k := SetupTestInput()

	expected := ctypes.Addresses{TestOracle1}

	store := ctx.KVStore(k.StoreKey)
	oracleBz := store.Get([]byte(types.OraclePrefix))
	assert.Nil(t, oracleBz)

	k.AddOracle(ctx, TestOracle1)
	oracleBz = store.Get([]byte(types.OraclePrefix))
	var actual ctypes.Addresses
	k.cdc.MustUnmarshalBinaryBare(oracleBz, &actual)

	assert.Equal(t, expected, actual)
}

func TestKeeper_IsOracle_ValidOracle(t *testing.T) {
	_, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.OraclePrefix), k.cdc.MustMarshalBinaryBare(ctypes.Addresses{TestOracle1}))

	isOracle := k.IsOracle(ctx, TestOracle1)
	assert.True(t, isOracle)
}

func TestKeeper_IsOracle_InvalidOracle(t *testing.T) {
	_, ctx, k := SetupTestInput()
	isOracle := k.IsOracle(ctx, TestOracle1)
	assert.False(t, isOracle)
}

func TestKeeper_GetOracles(t *testing.T) {
	_, ctx, k := SetupTestInput()
	expected := ctypes.Addresses{TestOracle1}
	k.AddOracle(ctx, TestOracle1)
	actual := k.GetOracles(ctx)
	assert.Equal(t, expected, actual)
}
