package keeper

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/government"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/types"
)

type Keeper struct {
	StoreKey  sdk.StoreKey
	cdc       *codec.Codec
	govKeeper government.Keeper
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey, govKeeper government.Keeper) Keeper {
	return Keeper{
		StoreKey:  storeKey,
		cdc:       cdc,
		govKeeper: govKeeper,
	}
}

// --------------
// --- Assets
// --------------

// AddAsset add a new priced asset to the assets list
func (keeper Keeper) AddAsset(ctx sdk.Context, assetName string) {
	store := ctx.KVStore(keeper.StoreKey)

	assets := keeper.GetAssets(ctx)
	if assets, updated := assets.AppendIfMissing(assetName); updated {
		store.Set([]byte(types.AssetsStoreKey), keeper.cdc.MustMarshalBinaryBare(&assets))
	}
}

// GetAssets retrieves all the assets
func (keeper Keeper) GetAssets(ctx sdk.Context) (assets ctypes.Strings) {
	store := ctx.KVStore(keeper.StoreKey)
	keeper.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.AssetsStoreKey)), &assets)
	return assets
}

// -----------------
// --- Raw prices
// -----------------

func (keeper Keeper) getRawPricesKey(assetName string) []byte {
	return []byte(types.RawPricesPrefix + assetName)
}

// AddRawPrice sets the raw price for a given token after checking the validity of the signer.
// If the signer hasn't the rights to set the price, then function returns error.
func (keeper Keeper) AddRawPrice(ctx sdk.Context, oracle sdk.AccAddress, price types.Price) error {
	if !keeper.IsOracle(ctx, oracle) {
		return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("%s is not an oracle", oracle))
	}

	// Add the asset's identifiers if it's the first time that it's been priced
	keeper.AddAsset(ctx, price.AssetName)

	// Update the raw prices
	rawPrice := types.OraclePrice{Oracle: oracle, Price: price, Created: sdk.NewInt(ctx.BlockHeight())}
	rawPrices := keeper.GetRawPricesForAsset(ctx, rawPrice.Price.AssetName)
	if rawPrices, updated := rawPrices.UpdatePriceOrAppendIfMissing(rawPrice); updated {
		store := ctx.KVStore(keeper.StoreKey)
		store.Set(keeper.getRawPricesKey(rawPrice.Price.AssetName), keeper.cdc.MustMarshalBinaryBare(&rawPrices))
		return nil
	}

	return sdkErr.Wrap(sdkErr.ErrUnknownRequest, fmt.Sprintf("Price %s already exists", price))
}

// GetRawPricesForAsset retrieves all the raw prices of the given asset
func (keeper Keeper) GetRawPricesForAsset(ctx sdk.Context, assetName string) types.OraclePrices {
	store := ctx.KVStore(keeper.StoreKey)

	var rawPrices types.OraclePrices
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(keeper.getRawPricesKey(assetName)), &rawPrices)
	return rawPrices
}

// GetRawPrices returns the list of the whole raw prices currently stored
func (keeper Keeper) GetRawPrices(ctx sdk.Context) types.OraclePrices {
	store := ctx.KVStore(keeper.StoreKey)

	prices := types.OraclePrices{}
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.RawPricesPrefix))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var price types.OraclePrices
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &price)
		prices = append(prices, price...)
	}

	return prices
}

// ---------------------
// --- Current prices
// ---------------------

func (keeper Keeper) getCurrentPriceKey(assetName string) []byte {
	return []byte(types.CurrentPricesPrefix + assetName)
}

func (keeper Keeper) ComputeAndUpdateCurrentPrices(ctx sdk.Context) {
	// Get all the listed assets
	assets := keeper.GetAssets(ctx)

	// For every asset, get all its not expired prices and calculate a median price that will be the current one
	for _, asset := range assets {

		// Get all raw prices posted by oracles
		rawPrices := keeper.GetRawPricesForAsset(ctx, asset)

		var notExpiredPrices = types.OraclePrices{}
		var rawPricesSum = sdk.NewDec(0)
		var rawExpirySum = sdk.NewInt(0)

		// Filter out expired prices
		for index, price := range rawPrices {
			if price.Price.Expiry.GTE(sdk.NewInt(ctx.BlockHeight())) {
				rawPricesSum = rawPricesSum.Add(rawPrices[index].Price.Value)
				rawExpirySum = rawExpirySum.Add(rawPrices[index].Price.Expiry)
				notExpiredPrices, _ = notExpiredPrices.UpdatePriceOrAppendIfMissing(price)
			}
		}

		pricesLength := len(notExpiredPrices)
		var medianPrice sdk.Dec
		var expiry sdk.Int

		// TODO KAVA suggestion : make threshold for acceptance (ie. require 51% of oracles to have posted valid prices)
		switch pricesLength {
		case 0:
			// Error if there are no valid prices in the raw prices store
			ctx.Logger().Debug("no valid raw prices to calculate current prices")
			continue

		case 1:
			// Return if there's only one price
			medianPrice = notExpiredPrices[0].Price.Value
			expiry = notExpiredPrices[0].Price.Expiry

		default:
			pLength := int64(pricesLength)
			medianPrice = rawPricesSum.Quo(sdk.NewDec(pLength))
			expiry = rawExpirySum.Quo(sdk.NewInt(pLength))
		}

		// Compute the new current price
		currentPrice := types.Price{
			AssetName: asset,
			Value:     medianPrice,
			Expiry:    expiry,
		}

		// Set the price
		keeper.SetCurrentPrice(ctx, currentPrice)
	}
}

// SetCurrentPrice allows to set the current price of a specific asset.
// WARNING: This method should be used for testing purposes only
func (keeper Keeper) SetCurrentPrice(ctx sdk.Context, currentPrice types.Price) {
	store := ctx.KVStore(keeper.StoreKey)
	store.Set(keeper.getCurrentPriceKey(currentPrice.AssetName), keeper.cdc.MustMarshalBinaryBare(currentPrice))
}

// GetCurrentPrice retrieves the current price for the given asset
func (keeper Keeper) GetCurrentPrice(ctx sdk.Context, asset string) (currentPrice types.Price, found bool) {
	store := ctx.KVStore(keeper.StoreKey)

	currentPrice = types.Price{}
	if !store.Has(keeper.getCurrentPriceKey(asset)) {
		return currentPrice, false
	}

	keeper.cdc.MustUnmarshalBinaryBare(store.Get(keeper.getCurrentPriceKey(asset)), &currentPrice)
	return currentPrice, true
}

// GetCurrentPrices retrieves all the current prices
func (keeper Keeper) GetCurrentPrices(ctx sdk.Context) types.Prices {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.CurrentPricesPrefix))

	var curPrices = types.Prices{}
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var currentPrice types.Price
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &currentPrice)
		curPrices, _ = curPrices.AppendIfMissing(currentPrice)
	}

	return curPrices
}

// ------------------
// --- Oracles
// ------------------

// AddOracle adds an Oracle to the store
func (keeper Keeper) AddOracle(ctx sdk.Context, oracle sdk.AccAddress) {
	store := ctx.KVStore(keeper.StoreKey)

	oracles := keeper.GetOracles(ctx)
	if oracles, success := oracles.AppendIfMissing(oracle); success {
		store.Set([]byte(types.OraclePrefix), keeper.cdc.MustMarshalBinaryBare(&oracles))
	}
}

// IsOracle returns true iif the given address is a valid oracle
func (keeper Keeper) IsOracle(ctx sdk.Context, address sdk.Address) bool {
	oracles := keeper.GetOracles(ctx)
	return oracles.Contains(address)
}

// GetOracles returns the list of all the currently present oracles
func (keeper Keeper) GetOracles(ctx sdk.Context) (oracles ctypes.Addresses) {
	store := ctx.KVStore(keeper.StoreKey)
	keeper.cdc.MustUnmarshalBinaryBare(store.Get([]byte(types.OraclePrefix)), &oracles)
	return oracles
}

// ------------------
// --- DenomBlacklist management
// ------------------

// bdKey returns a byte slice containing the store key for a given denom
func bdKey(denom string) []byte {
	return []byte(types.DenomBlacklistKey + denom)
}

// DenomBlacklistIterator returns a store iterator for all the blacklisted denoms
func (keeper Keeper) DenomBlacklistIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(keeper.StoreKey)
	return sdk.KVStorePrefixIterator(store, []byte(types.DenomBlacklistKey))
}

// BlacklistDenom blacklists a list of denoms.
func (keeper Keeper) BlacklistDenom(ctx sdk.Context, denom ...string) {
	store := ctx.KVStore(keeper.StoreKey)

	for _, d := range denom {
		store.Set(bdKey(d), []byte(d))
	}
}

// DenomBlacklist returns the list of all the blacklisted denoms
func (keeper Keeper) DenomBlacklist(ctx sdk.Context) []string {
	iter := keeper.DenomBlacklistIterator(ctx)
	defer iter.Close()

	var ret []string
	for ; iter.Valid(); iter.Next() {
		ret = append(ret, string(iter.Value()))
	}

	return ret
}
