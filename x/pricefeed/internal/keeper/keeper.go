package keeper

import (
	"fmt"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"
)

type Keeper struct {
	StoreKey         sdk.StoreKey
	GovernmentKeeper government.Keeper
	cdc              *codec.Codec
}

func NewKeeper(storekey sdk.StoreKey, govK government.Keeper, cdc *codec.Codec) Keeper {
	return Keeper{
		StoreKey:         storekey,
		GovernmentKeeper: govK,
		cdc:              cdc,
	}
}

//getAssets retrieve all the assets
func (keeper Keeper) getAssets(ctx sdk.Context) types.Assets {
	store := ctx.KVStore(keeper.StoreKey)
	assetsBz := store.Get([]byte(types.AssetsPrefix))
	var assets types.Assets
	keeper.cdc.MustUnmarshalBinaryBare(assetsBz, &assets)
	return assets
}

//setAssets add a new priced assets to the assets list
func (keeper Keeper) setAssets(ctx sdk.Context, assetName string, assetCode string) {
	assets := keeper.getAssets(ctx)
	assets = assets.AppendIfMissing(types.Asset{Name: assetName, Code: assetCode})
}

//SetRawPrice sets the raw price for a given token after checking the validity of the signer
//If the signer hasn't the rights to set the price, then function returns error
func (keeper Keeper) SetRawPrice(ctx sdk.Context, price types.RawPrice) sdk.Error {
	store := ctx.KVStore(keeper.StoreKey)
	err := keeper.ValidateSigner(ctx, price.Oracle)
	if err != nil {
		return err
	}
	rawPrices := keeper.GetRawPrices(ctx, price.PriceInfo.AssetName, price.PriceInfo.AssetCode)
	keeper.setAssets(ctx, price.PriceInfo.AssetName, price.PriceInfo.AssetCode)
	rawPrices = rawPrices.UpdatePriceOrAppendIfMissing(price)
	store.Set([]byte(types.RawPricesPrefix + price.PriceInfo.AssetName + price.PriceInfo.AssetCode),
		keeper.cdc.MustMarshalBinaryBare(rawPrices))
	return nil
}

//GetRawPrices retrieves all the raw prices of the given asset if it given, all the raw prices instead
func (keeper Keeper) GetRawPrices(ctx sdk.Context, assetName string, assetCode string) types.RawPrices {
	store := ctx.KVStore(keeper.StoreKey)
	pricesBz := store.Get([]byte(types.RawPricesPrefix + assetName + assetCode))
	var rawPrices types.RawPrices
	keeper.cdc.MustUnmarshalBinaryBare(pricesBz, &rawPrices)
	return rawPrices
}

func (keeper Keeper) SetCurrentPrices(ctx sdk.Context) sdk.Error {

	//Get all listed assets
	assets := keeper.getAssets(ctx)

	//For every asset, get all its not expired prices and calculate a median price that will be the current one
	for _, asset := range assets {
		// Get all raw prices posted by oracles
		rawPrices := keeper.GetRawPrices(ctx, asset.Name, asset.Code)
		var notExpiredPrices = types.RawPrices{}

		// filter out expired prices
		for _, price := range rawPrices {
			if price.PriceInfo.Expiry.GTE(sdk.NewInt(ctx.BlockHeight())) {
				notExpiredPrices = notExpiredPrices.UpdatePriceOrAppendIfMissing(price)
			}
		}

		pricesLength := len(notExpiredPrices)
		var medianPrice sdk.Int
		var expiry 		sdk.Int
		// TODO make threshold for acceptance (ie. require 51% of oracles to have posted valid prices)
		if pricesLength == 0 {
			// Error if there are no valid prices in the raw prices store
			return sdk.ErrInternal("no valid raw prices to calculate current prices")
		} else if pricesLength == 1 {
			// Return if there's only one price
			medianPrice = notExpiredPrices[0].PriceInfo.Price
			expiry = notExpiredPrices[0].PriceInfo.Expiry
		} else {
			// sort the prices
			sort.Slice(notExpiredPrices, func(i, j int) bool {
				return notExpiredPrices[i].PriceInfo.Price.LT(notExpiredPrices[j].PriceInfo.Price)
			})
			// If there's an even number of prices
			if pricesLength%2 == 0 {
				// TODO make sure this is safe.
				// Since it's a price and not a balance, division with precision loss is OK.
				price1 := notExpiredPrices[pricesLength/2-1].PriceInfo.Price
				price2 := notExpiredPrices[pricesLength/2].PriceInfo.Price
				sum := price1.Add(price2)
				medianPrice = sum.Quo(sdk.NewInt(2))
				// TODO Check if safe, makes sense
				// Takes the average of the two expires rounded down to the nearest Int.
				expiry = notExpiredPrices[pricesLength/2-1].PriceInfo.Expiry.
					Add(notExpiredPrices[pricesLength/2].PriceInfo.Expiry).
					Quo(sdk.NewInt(2))
			} else {
				// integer division, so we'll get an integer back, rounded down
				medianPrice = notExpiredPrices[pricesLength/2].PriceInfo.Price
				expiry = notExpiredPrices[pricesLength/2].PriceInfo.Expiry
			}
		}
		store := ctx.KVStore(keeper.StoreKey)


		store.Set([]byte(types.CurrentPricesPrefix+))


	}



	/*
		GetRawPrices
		filter not expire ones
		filter for tokenName+tokeCode
		curPrice = avg(rawPrices)
	*/
}

//GetCurrentPrices retrieves all the current prices
func (keeper Keeper) GetCurrentPrices(ctx sdk.Context) types.CurrentPrices {
	store := ctx.KVStore(keeper.StoreKey)
	pricesBz := store.Get([]byte(types.CurrentPricesPrefix))
	var curPrices types.CurrentPrices
	keeper.cdc.MustUnmarshalBinaryBare(pricesBz, &curPrices)
	return curPrices
}

//GetCurrentPrice retrieves the current price for the given token name and code
func (keeper Keeper) GetCurrentPrice(ctx sdk.Context, tokenName string, tokenCode string) (types.CurrentPrice, sdk.Error) {
	currentPrices := keeper.GetCurrentPrices(ctx)
	price, err := currentPrices.FindPrice(tokenName, tokenCode)
	if err != nil {
		return types.CurrentPrice{}, err
	}
	return price, nil
}

//ValidateSigner makes sure the signer posting the price is an oracle
func (keeper Keeper) ValidateSigner(ctx sdk.Context, signer sdk.AccAddress) sdk.Error {
	oracles := keeper.GetOracles(ctx)
	isOracle := oracles.Contains(signer)
	if !isOracle {
		return sdk.ErrInvalidAddress(fmt.Sprintf("%s isn't an Oracle", signer))
	}
	return nil
}

// AddOracle adds an Oracle to the store
func (keeper Keeper) AddOracle(ctx sdk.Context, oracle sdk.AccAddress) {
	oracles := keeper.GetOracles(ctx)
	oracles = oracles.AppendIfMissing(oracle)
	store := ctx.KVStore(keeper.StoreKey)
	store.Set([]byte(types.OraclePrefix), keeper.cdc.MustMarshalBinaryBare(&oracles))
}

func (keeper Keeper) GetOracles(ctx sdk.Context) ctypes.Addresses {
	store := ctx.KVStore(keeper.StoreKey)
	oraclesBz := store.Get([]byte(types.OraclePrefix))
	var oracles ctypes.Addresses
	keeper.cdc.MustUnmarshalBinaryBare(oraclesBz, &oracles)
	return oracles
}

func (keeper Keeper) GetOracle(ctx sdk.Context, oracle sdk.AccAddress) (sdk.AccAddress, error) {
	oracles := keeper.GetOracles(ctx)
	found := oracles.GetAddress(oracle)
	if found == nil {
		return nil, sdk.ErrUnknownAddress("Oracle address not found")
	}
	return found, nil
}
