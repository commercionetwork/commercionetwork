package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestCurrentPrice_Equals_true(t *testing.T) {
	actual := testPrice.Equals(testPrice)
	assert.True(t, actual)
}

func TestCurrentPrice_Equals_false(t *testing.T) {
	curPrice := Price{
		AssetName: testPrice.AssetName,
		Value:     sdk.NewDec(2),
		Expiry:    sdk.NewInt(5),
	}

	actual := testPrice.Equals(curPrice)
	assert.False(t, actual)
}

func TestCurrentPrices_GetPrice(t *testing.T) {
	prices := Prices{testPrice}
	actual, _ := prices.GetPrice(testPrice.AssetName)
	assert.Equal(t, testPrice, actual)
}

func TestCurrentPrices_GetPrice_NotFound(t *testing.T) {
	prices := Prices{}
	_, err := prices.GetPrice("testName")
	assert.Error(t, err)
}

func TestCurrentPrices_AppendIfMissing_Notfound(t *testing.T) {
	prices := Prices{}
	actual, found := prices.AppendIfMissing(testPrice)
	expected := Prices{testPrice}
	assert.Equal(t, expected, actual)
	assert.False(t, found)
}

func TestCurrentPrices_AppendIfMissing_found(t *testing.T) {
	prices := Prices{testPrice}
	actual, found := prices.AppendIfMissing(testPrice)
	assert.Nil(t, actual)
	assert.True(t, found)
}

func TestPrice_Equals_true(t *testing.T) {
	actual := testPrice.Equals(testPrice)
	assert.True(t, actual)
}

func TestPrices_UpdatePriceOrAppendIfMissing_appendNewPrice(t *testing.T) {
	rawPrices := Prices{}
	actual, updated := rawPrices.UpdatePriceOrAppendIfMissing(testPrice)
	assert.True(t, updated)
	assert.Equal(t, testPrice, actual[0])
}

func TestPrices_UpdatePriceOrAppendIfMissing_priceAlreadyInserted(t *testing.T) {
	rawPrices := Prices{testPrice}
	actual, updated := rawPrices.UpdatePriceOrAppendIfMissing(testPrice)
	assert.False(t, updated)
	assert.Equal(t, rawPrices, actual)
}

func TestPrices_UpdatePriceOrAppendIfMissing_updatedPrice(t *testing.T) {
	curPrice := Price{
		AssetName: testPrice.AssetName,
		Value:     sdk.NewDec(200),
		Expiry:    sdk.NewInt(6000),
	}
	rawPrices := Prices{testPrice}
	actual, updated := rawPrices.UpdatePriceOrAppendIfMissing(curPrice)
	assert.Equal(t, curPrice, actual[0])
	assert.True(t, updated)
}
