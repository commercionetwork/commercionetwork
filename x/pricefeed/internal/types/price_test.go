package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestCurrentPrice_Equals_true(t *testing.T) {
	actual := TestPriceInfo.Equals(TestPriceInfo)
	assert.True(t, actual)
}

func TestCurrentPrice_Equals_false(t *testing.T) {
	curPrice := CurrentPrice{
		AssetName: TestPriceInfo.AssetName,
		AssetCode: TestPriceInfo.AssetCode,
		Price:     sdk.NewInt(2),
		Expiry:    sdk.NewInt(5),
	}

	actual := TestPriceInfo.Equals(curPrice)
	assert.False(t, actual)
}

func TestCurrentPrices_GetPrice(t *testing.T) {
	prices := CurrentPrices{TestPriceInfo}
	actual, _ := prices.GetPrice(TestPriceInfo.AssetName, TestPriceInfo.AssetCode)
	assert.Equal(t, TestPriceInfo, actual)
}

func TestCurrentPrices_GetPrice_NotFound(t *testing.T) {
	prices := CurrentPrices{}
	_, err := prices.GetPrice("testName", "testCode")
	assert.Error(t, err)
}

func TestCurrentPrices_AppendIfMissing_Notfound(t *testing.T) {
	prices := CurrentPrices{}
	actual, found := prices.AppendIfMissing(TestPriceInfo)
	expected := CurrentPrices{TestPriceInfo}
	assert.Equal(t, expected, actual)
	assert.False(t, found)
}

func TestCurrentPrices_AppendIfMissing_found(t *testing.T) {
	prices := CurrentPrices{TestPriceInfo}
	actual, found := prices.AppendIfMissing(TestPriceInfo)
	assert.Nil(t, actual)
	assert.True(t, found)
}

func TestRawPrice_Equals_false(t *testing.T) {
	curPrice := CurrentPrice{
		AssetName: TestPriceInfo.AssetName,
		AssetCode: TestPriceInfo.AssetCode,
		Price:     sdk.NewInt(2),
		Expiry:    sdk.NewInt(5),
	}
	rawPrice := RawPrice{
		Oracle:    testOracle1,
		PriceInfo: curPrice,
	}

	actual := rawPrice.Equals(TestRawPrice)
	assert.False(t, actual)
}

func TestRawPrice_Equals_true(t *testing.T) {
	actual := TestRawPrice.Equals(TestRawPrice)
	assert.True(t, actual)
}

func TestRawPrices_UpdatePriceOrAppendIfMissing_appendNewRawPrice(t *testing.T) {
	rawPrices := RawPrices{}
	actual, found := rawPrices.UpdatePriceOrAppendIfMissing(TestRawPrice)
	assert.Equal(t, TestRawPrice, actual[0])
	assert.False(t, found)
}

func TestRawPrices_UpdatePriceOrAppendIfMissing_priceAlreadyInserted(t *testing.T) {
	rawPrices := RawPrices{TestRawPrice}
	actual, found := rawPrices.UpdatePriceOrAppendIfMissing(TestRawPrice)
	assert.Nil(t, actual)
	assert.True(t, found)
}

func TestRawPrices_UpdatePriceOrAppendIfMissing_updatedPrice(t *testing.T) {
	curPrice := CurrentPrice{
		AssetName: TestPriceInfo.AssetName,
		AssetCode: TestPriceInfo.AssetCode,
		Price:     sdk.NewInt(200),
		Expiry:    sdk.NewInt(6000),
	}
	rawPrice := RawPrice{
		Oracle:    TestRawPrice.Oracle,
		PriceInfo: curPrice,
	}
	rawPrices := RawPrices{TestRawPrice}
	actual, found := rawPrices.UpdatePriceOrAppendIfMissing(rawPrice)
	assert.Equal(t, rawPrice, actual[0])
	assert.False(t, found)
}
