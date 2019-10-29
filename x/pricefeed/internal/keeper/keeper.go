package keeper

import (
	"errors"
	"fmt"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	StoreKey sdk.StoreKey
	cdc      *codec.Codec
}

func NewKeeper(cdc *codec.Codec, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		StoreKey: storeKey,
		cdc:      cdc,
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

// SetRawPrice sets the raw price for a given token after checking the validity of the signer.
// If the signer hasn't the rights to set the price, then function returns error.
func (keeper Keeper) SetRawPrice(ctx sdk.Context, price types.RawPrice) error {
	store := ctx.KVStore(keeper.StoreKey)

	// Check the signer
	if !keeper.IsOracle(ctx, price.Oracle) {
		return fmt.Errorf("%s is not an oracle", price.Oracle.String())
	}

	// Add the asset's identifiers if it's the first time that it's been priced
	keeper.AddAsset(ctx, price.PriceInfo.AssetName)

	// Update the raw prices
	rawPrices := keeper.GetRawPrices(ctx, price.PriceInfo.AssetName)
	if rawPrices, updated := rawPrices.UpdatePriceOrAppendIfMissing(price); updated {
		store.Set(keeper.getRawPricesKey(price.PriceInfo.AssetName), keeper.cdc.MustMarshalBinaryBare(rawPrices))
	} else {
		return errors.New("price already present")
	}

	return nil
}

// GetRawPrices retrieves all the raw prices of the given asset
func (keeper Keeper) GetRawPrices(ctx sdk.Context, assetName string) (rawPrices types.RawPrices) {
	store := ctx.KVStore(keeper.StoreKey)
	keeper.cdc.MustUnmarshalBinaryBare(store.Get(keeper.getRawPricesKey(assetName)), &rawPrices)
	return rawPrices
}

// ---------------------
// --- Current prices
// ---------------------

func (keeper Keeper) getCurrentPriceKey(assetName string) []byte {
	return []byte(types.CurrentPricesPrefix + assetName)
}

func (keeper Keeper) ComputeAndUpdateCurrentPrices(ctx sdk.Context) error {
	// Get all the listed assets
	assets := keeper.GetAssets(ctx)

	// For every asset, get all its not expired prices and calculate a median price that will be the current one
	for _, asset := range assets {

		// Get all raw prices posted by oracles
		rawPrices := keeper.GetRawPrices(ctx, asset)

		var notExpiredPrices = types.RawPrices{}
		var rawPricesSum = sdk.NewDec(0)
		var rawExpirySum = sdk.NewInt(0)

		// Filter out expired prices
		for index, price := range rawPrices {
			if price.PriceInfo.Expiry.GTE(sdk.NewInt(ctx.BlockHeight())) {
				rawPricesSum = rawPricesSum.Add(rawPrices[index].PriceInfo.Price)
				rawExpirySum = rawExpirySum.Add(rawPrices[index].PriceInfo.Expiry)
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
			return errors.New("no valid raw prices to calculate current prices")

		case 1:
			// Return if there's only one price
			medianPrice = notExpiredPrices[0].PriceInfo.Price
			expiry = notExpiredPrices[0].PriceInfo.Expiry

		default:
			pLength := int64(pricesLength)
			medianPrice = rawPricesSum.Quo(sdk.NewDec(pLength))
			expiry = rawExpirySum.Quo(sdk.NewInt(pLength))
		}

		// Compute the new current price
		currentPrice := types.CurrentPrice{
			AssetName: asset,
			Price:     medianPrice,
			Expiry:    expiry,
		}

		// Set the price
		keeper.SetCurrentPrice(ctx, currentPrice)

	}
	return nil
}

// SetCurrentPrice allows to set the current price of a specific asset.
// WARNING: This method should be used for testing purposes only
func (keeper Keeper) SetCurrentPrice(ctx sdk.Context, currentPrice types.CurrentPrice) {
	store := ctx.KVStore(keeper.StoreKey)
	store.Set(keeper.getCurrentPriceKey(currentPrice.AssetName), keeper.cdc.MustMarshalBinaryBare(currentPrice))
}

// GetCurrentPrice retrieves the current price for the given asset
func (keeper Keeper) GetCurrentPrice(ctx sdk.Context, asset string) (currentPrice types.CurrentPrice, found bool) {
	store := ctx.KVStore(keeper.StoreKey)

	currentPrice = types.CurrentPrice{}
	if !store.Has(keeper.getCurrentPriceKey(asset)) {
		return currentPrice, false
	}

	keeper.cdc.MustUnmarshalBinaryBare(store.Get(keeper.getCurrentPriceKey(asset)), &currentPrice)
	return currentPrice, true
}

// GetCurrentPrices retrieves all the current prices
func (keeper Keeper) GetCurrentPrices(ctx sdk.Context) types.CurrentPrices {
	store := ctx.KVStore(keeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.CurrentPricesPrefix))

	var curPrices = types.CurrentPrices{}
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var currentPrice types.CurrentPrice
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
