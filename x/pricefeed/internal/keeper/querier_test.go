package keeper

import (
	"testing"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
)

var request abci.RequestQuery

func TestQuerier_getCurrentPrices(t *testing.T) {
	cdc, ctx, k := SetupTestInput()

	//setup the keystore with two current prices
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.CurrentPricesPrefix+TestPriceInfo.AssetName), k.cdc.MustMarshalBinaryBare(TestPriceInfo))
	store.Set([]byte(types.CurrentPricesPrefix+TestPriceInfo2.AssetName), k.cdc.MustMarshalBinaryBare(TestPriceInfo2))

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
	store.Set([]byte(types.CurrentPricesPrefix+TestPriceInfo.AssetName), k.cdc.MustMarshalBinaryBare(TestPriceInfo))
	store.Set([]byte(types.CurrentPricesPrefix+TestPriceInfo2.AssetName), k.cdc.MustMarshalBinaryBare(TestPriceInfo2))

	querier := NewQuerier(k)
	path := []string{types.QueryGetCurrentPrice, TestPriceInfo.AssetName}

	var actual types.CurrentPrice
	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)

	assert.Equal(t, TestPriceInfo, actual)
}

func TestQuerier_getOracles(t *testing.T) {
	cdc, ctx, k := SetupTestInput()
	var expected = ctypes.Addresses{TestOracle1}

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.OraclePrefix), k.cdc.MustMarshalBinaryBare(expected))

	querier := NewQuerier(k)
	path := []string{types.QueryGetOracles}

	var actual ctypes.Addresses

	actualBz, _ := querier(ctx, path, request)
	cdc.MustUnmarshalJSON(actualBz, &actual)
	assert.Equal(t, expected, actual)
}

func TestQuerier_unknownEndpoint(t *testing.T) {
	_, ctx, k := SetupTestInput()
	querier := NewQuerier(k)

	path := []string{"test"}
	_, actual := querier(ctx, path, request)

	assert.Error(t, actual)
}
