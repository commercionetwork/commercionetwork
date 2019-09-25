package keeper

import (
	"testing"

	types2 "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGetRawPricesKey(t *testing.T) {
	expected := []byte(types.RawPricesPrefix + TestRawPrice1.PriceInfo.AssetName + TestRawPrice1.PriceInfo.AssetCode)
	actual := GetRawPricesKey(TestRawPrice1.PriceInfo.AssetName, TestRawPrice1.PriceInfo.AssetCode)
	assert.Equal(t, expected, actual)
}

func TestKeeper_GetAssets(t *testing.T) {
	_, ctx, k := SetupTestInput()
	expected := types.Assets{TestAsset, TestAsset2}

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.AssetsPrefix), k.cdc.MustMarshalBinaryBare(expected))

	actual := k.GetAssets(ctx)
	assert.Equal(t, expected, actual)
}

func TestKeeper_GetAssets_empty(t *testing.T) {
	_, ctx, k := SetupTestInput()
	var expected types.Assets
	actual := k.GetAssets(ctx)
	assert.Equal(t, expected, actual)
}

func TestKeeper_AddAsset(t *testing.T) {
	_, ctx, k := SetupTestInput()
	//ensure that assets array are empty
	assets := k.GetAssets(ctx)
	assert.Nil(t, assets)

	k.AddAsset(ctx, TestAsset.Name, TestAsset.Code)
	expected := types.Assets{TestAsset}
	actual := k.GetAssets(ctx)
	assert.Equal(t, expected, actual)
}

func TestKeeper_AddAsset_alreadyPresent(t *testing.T) {
	_, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)
	assets := types.Assets{TestAsset}
	store.Set([]byte(types.AssetsPrefix), k.cdc.MustMarshalBinaryBare(assets))
	expected := len(assets)

	k.AddAsset(ctx, TestAsset.Name, TestAsset.Code)
	actual := len(k.GetAssets(ctx))

	assert.Equal(t, expected, actual)
}

func TestKeeper_ValidateSigner_validSigner(t *testing.T) {
	_, ctx, k := SetupTestInput()
	oracles := types2.Addresses{TestOracle2, TestOracle1}
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.OraclePrefix), k.cdc.MustMarshalBinaryBare(oracles))

	err := k.ValidateSigner(ctx, TestOracle1)
	assert.Nil(t, err)
}

func TestKeeper_ValidateSigner_invalidSigner(t *testing.T) {
	_, ctx, k := SetupTestInput()
	err := k.ValidateSigner(ctx, TestOracle1)
	assert.Error(t, err)
}

func TestKeeper_SetRawPrice_withValidSigner_pricesNotAlreadyPresent(t *testing.T) {
	_, ctx, k := SetupTestInput()
	k.AddOracle(ctx, TestOracle1)
	k.AddOracle(ctx, TestOracle2)

	err := k.SetRawPrice(ctx, TestRawPrice1)
	assert.Nil(t, err)
	err = k.SetRawPrice(ctx, TestRawPrice3)
	assert.Nil(t, err)

	actual := k.GetRawPrices(ctx, TestRawPrice1.PriceInfo.AssetName, TestRawPrice1.PriceInfo.AssetCode)

	expected := types.RawPrices{TestRawPrice1, TestRawPrice3}

	assert.Equal(t, expected, actual)

}

func TestKeeper_SetRawPrice_withValidSigner_priceAlreadyPresent(t *testing.T) {
	_, ctx, k := SetupTestInput()
	k.AddOracle(ctx, TestOracle1)
	store := ctx.KVStore(k.StoreKey)
	expected := types.RawPrices{TestRawPrice1}
	store.Set(GetRawPricesKey(TestRawPrice1.PriceInfo.AssetName, TestRawPrice1.PriceInfo.AssetCode),
		k.cdc.MustMarshalBinaryBare(expected))

	err := k.SetRawPrice(ctx, TestRawPrice1)
	assert.Error(t, err)
}

func TestKeeper_SetRawPrice_withInvalidSigner(t *testing.T) {
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
	actual := k.GetRawPrices(ctx, TestRawPrice1.PriceInfo.AssetName, TestRawPrice1.PriceInfo.AssetCode)
	expected := types.RawPrices{TestRawPrice1, TestRawPrice3}
	assert.Equal(t, expected, actual)
}

func TestKeeper_SetCurrentPrices_MoreThanOneNotExpiredPrice(t *testing.T) {
	_, ctx, k := SetupTestInput()
	//test setup

	k.AddOracle(ctx, TestOracle1)
	k.AddOracle(ctx, TestOracle2)
	_ = k.SetRawPrice(ctx, TestRawPrice1)
	_ = k.SetRawPrice(ctx, TestRawPrice3)

	sumPrice := TestRawPrice1.PriceInfo.Price.Add(TestRawPrice3.PriceInfo.Price)
	sumExpiry := TestRawPrice1.PriceInfo.Expiry.Add(TestRawPrice3.PriceInfo.Expiry)
	expectedMedianPrice := sumPrice.Quo(sdk.NewInt(2))
	expectedMedianExpiry := sumExpiry.Quo(sdk.NewInt(2))

	_ = k.SetCurrentPrices(ctx)

	actual, _ := k.GetCurrentPrice(ctx, TestRawPrice1.PriceInfo.AssetName, TestRawPrice1.PriceInfo.AssetCode)

	assert.Equal(t, expectedMedianPrice, actual.Price)
	assert.Equal(t, expectedMedianExpiry, actual.Expiry)
}

func TestKeeper_SetCurrentPrices_allExpiredRawPrices(t *testing.T) {
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
	actual, _ := k.GetCurrentPrice(ctx, TestRawPrice1.PriceInfo.AssetName, TestRawPrice1.PriceInfo.AssetCode)
	assert.Equal(t, TestRawPrice1.PriceInfo.Price, actual.Price)
	assert.Equal(t, TestRawPrice1.PriceInfo.Expiry, actual.Expiry)
}

func TestKeeper_GetCurrentPrices(t *testing.T) {
	_, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.CurrentPricesPrefix+TestPriceInfo.AssetName+TestPriceInfo.AssetCode),
		k.cdc.MustMarshalBinaryBare(TestPriceInfo))
	actual := k.GetCurrentPrices(ctx)
	expected := types.CurrentPrices{TestPriceInfo}
	assert.Equal(t, expected, actual)
}

func TestKeeper_GetCurrentPrice_Found(t *testing.T) {
	_, ctx, k := SetupTestInput()
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.CurrentPricesPrefix+TestPriceInfo.AssetName+TestPriceInfo.AssetCode),
		k.cdc.MustMarshalBinaryBare(TestPriceInfo))
	actual, _ := k.GetCurrentPrice(ctx, TestPriceInfo.AssetName, TestPriceInfo.AssetCode)
	assert.Equal(t, TestPriceInfo, actual)
}

func TestKeeper_GetCurrentPrice_NotFound(t *testing.T) {
	_, ctx, k := SetupTestInput()
	_, err := k.GetCurrentPrice(ctx, TestPriceInfo.AssetName, TestPriceInfo.AssetCode)
	assert.Error(t, err)
}

func TestKeeper_AddOracle(t *testing.T) {
	_, ctx, k := SetupTestInput()

	expected := types2.Addresses{TestOracle1}

	store := ctx.KVStore(k.StoreKey)
	oracleBz := store.Get([]byte(types.OraclePrefix))
	assert.Nil(t, oracleBz)

	k.AddOracle(ctx, TestOracle1)
	oracleBz = store.Get([]byte(types.OraclePrefix))
	var actual types2.Addresses
	k.cdc.MustUnmarshalBinaryBare(oracleBz, &actual)

	assert.Equal(t, expected, actual)
}

func TestKeeper_GetOracles(t *testing.T) {
	_, ctx, k := SetupTestInput()
	expected := types2.Addresses{TestOracle1}
	k.AddOracle(ctx, TestOracle1)
	actual := k.GetOracles(ctx)
	assert.Equal(t, expected, actual)
}

func TestKeeper_GetOracle_Found(t *testing.T) {
	_, ctx, k := SetupTestInput()
	k.AddOracle(ctx, TestOracle1)
	actual, _ := k.GetOracle(ctx, TestOracle1)
	assert.Equal(t, TestOracle1, actual)
}

func TestKeeper_GetOracle_NotFound(t *testing.T) {
	_, ctx, k := SetupTestInput()
	_, err := k.GetOracle(ctx, TestOracle1)
	assert.Error(t, err)
}
