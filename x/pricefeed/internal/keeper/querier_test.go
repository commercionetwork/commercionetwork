package keeper

import (
	"testing"

	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery

func TestQuerier_getCurrentPrices(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	//setup the keystore with two current prices
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.CurrentPricesPrefix+TestPriceInfo.AssetName+TestPriceInfo.AssetCode), k.cdc.MustMarshalBinaryBare(TestPriceInfo))
	store.Set([]byte(types.CurrentPricesPrefix+TestPriceInfo2.AssetName+TestPriceInfo2.AssetCode), k.cdc.MustMarshalBinaryBare(TestPriceInfo2))

	querier := NewQuerier(k)
	path := []string{types.QueryGetCurrentPrices}

	var actual types.CurrentPrices
	var expected = types.CurrentPrices{TestPriceInfo, TestPriceInfo2}
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, expected, actual)
}

func TestQuerier_getCurrentPrice(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	//setup the keystore with two current prices
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.CurrentPricesPrefix+TestPriceInfo.AssetName+TestPriceInfo.AssetCode), k.cdc.MustMarshalBinaryBare(TestPriceInfo))
	store.Set([]byte(types.CurrentPricesPrefix+TestPriceInfo2.AssetName+TestPriceInfo2.AssetCode), k.cdc.MustMarshalBinaryBare(TestPriceInfo2))

	querier := NewQuerier(k)
	path := []string{types.QueryGetCurrentPrice, TestPriceInfo.AssetName, TestPriceInfo.AssetCode}

	var actual types.CurrentPrice
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, TestPriceInfo, actual)
}
