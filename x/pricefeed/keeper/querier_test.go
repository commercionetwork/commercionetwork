package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/types"
)

var request abci.RequestQuery

var TestPriceInfo2 = types.Price{
	AssetName: "test2",
	Value:     sdk.NewDec(8),
	Expiry:    sdk.NewInt(4000),
}

func TestQuerier_getCurrentPrices(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	// setup the keystore with two current prices
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.CurrentPricesPrefix+testPrice.AssetName), k.cdc.MustMarshalBinaryBare(testPrice))
	store.Set([]byte(types.CurrentPricesPrefix+TestPriceInfo2.AssetName), k.cdc.MustMarshalBinaryBare(TestPriceInfo2))

	path := []string{types.QueryGetCurrentPrices}
	querier := NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual types.Prices
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, types.Prices{testPrice, TestPriceInfo2}, actual)
}

func TestQuerier_getCurrentPrice(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	// setup the keystore with two current prices
	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.CurrentPricesPrefix+testPrice.AssetName), k.cdc.MustMarshalBinaryBare(testPrice))
	store.Set([]byte(types.CurrentPricesPrefix+TestPriceInfo2.AssetName), k.cdc.MustMarshalBinaryBare(TestPriceInfo2))

	path := []string{types.QueryGetCurrentPrice, testPrice.AssetName}
	querier := NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual types.Price
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, testPrice, actual)
}

func TestQuerier_getOracles(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	testOracle, _ := sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
	expected := ctypes.Addresses{testOracle}

	store := ctx.KVStore(k.StoreKey)
	store.Set([]byte(types.OraclePrefix), k.cdc.MustMarshalBinaryBare(expected))

	path := []string{types.QueryGetOracles}
	querier := NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual ctypes.Addresses
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, expected, actual)
}

func TestQuerier_queryGetBlacklistedDenoms(t *testing.T) {
	cdc, ctx, _, k := SetupTestInput()

	blacklisted := "blacklisted"
	k.BlacklistDenom(ctx, blacklisted)

	path := []string{types.QueryGetBlacklistedDenoms}
	querier := NewQuerier(k)
	actualBz, _ := querier(ctx, path, request)

	var actual []string
	cdc.MustUnmarshalJSON(actualBz, &actual)

	require.Equal(t, []string{"blacklisted"}, actual)
}

func TestQuerier_unknownEndpoint(t *testing.T) {
	_, ctx, _, k := SetupTestInput()

	path := []string{"test"}
	querier := NewQuerier(k)
	_, err := querier(ctx, path, request)

	require.Error(t, err)
}
