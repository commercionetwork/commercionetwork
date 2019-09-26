package keeper

import (
	"fmt"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	"github.com/commercionetwork/commercionetwork/x/government"
	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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

func GetRawPricesKey(assetName string, assetCode string) []byte {
	return []byte(types.RawPricesPrefix + assetName + assetCode)
}

//GetAssets retrieves all the assets
func (keeper Keeper) GetAssets(ctx sdk.Context) types.Assets {
	store := ctx.KVStore(keeper.StoreKey)
	assetsBz := store.Get([]byte(types.AssetsPrefix))
	var assets types.Assets
	keeper.cdc.MustUnmarshalBinaryBare(assetsBz, &assets)
	return assets
}

//AddAsset add a new priced asset to the assets list
func (keeper Keeper) AddAsset(ctx sdk.Context, assetName string, assetCode string) {
	store := ctx.KVStore(keeper.StoreKey)
	assets := keeper.GetAssets(ctx)
	assets, found := assets.AppendIfMissing(types.Asset{Name: assetName, Code: assetCode})
	if !found {
		store.Set([]byte(types.AssetsPrefix), keeper.cdc.MustMarshalBinaryBare(&assets))
	}
}

//SetRawPrice sets the raw price for a given token after checking the validity of the signer
//If the signer hasn't the rights to set the price, then function returns error
func (keeper Keeper) SetRawPrice(ctx sdk.Context, price types.RawPrice) sdk.Error {
	store := ctx.KVStore(keeper.StoreKey)
	err := keeper.ValidateSigner(ctx, price.Oracle)
	if err != nil {
		return err
	}

	//add Asset's identifiers if it's the first time that it's been priced
	keeper.AddAsset(ctx, price.PriceInfo.AssetName, price.PriceInfo.AssetCode)

	rawPrices := keeper.GetRawPrices(ctx, price.PriceInfo.AssetName, price.PriceInfo.AssetCode)
	rawPrices, found := rawPrices.UpdatePriceOrAppendIfMissing(price)
	if found {
		return sdk.ErrUnknownRequest(fmt.Sprintf("%s, is already been inserted by %s", price, price.Oracle))
	}

	store.Set(GetRawPricesKey(price.PriceInfo.AssetName, price.PriceInfo.AssetCode),
		keeper.cdc.MustMarshalBinaryBare(rawPrices))

	return nil
}

//GetRawPrices retrieves all the raw prices of the given asset
func (keeper Keeper) GetRawPrices(ctx sdk.Context, assetName string, assetCode string) types.RawPrices {
	store := ctx.KVStore(keeper.StoreKey)
	var rawPrices types.RawPrices
	pricesBz := store.Get(GetRawPricesKey(assetName, assetCode))
	keeper.cdc.MustUnmarshalBinaryBare(pricesBz, &rawPrices)
	return rawPrices
}

func (keeper Keeper) SetCurrentPrices(ctx sdk.Context) sdk.Error {

	//Get all listed assets
	assets := keeper.GetAssets(ctx)

	//For every asset, get all its not expired prices and calculate a median price that will be the current one
	for _, asset := range assets {
		// Get all raw prices posted by oracles
		rawPrices := keeper.GetRawPrices(ctx, asset.Name, asset.Code)
		var notExpiredPrices = types.RawPrices{}
		var rawPricesSum = sdk.NewInt(0)
		var rawExpirySum = sdk.NewInt(0)

		// filter out expired prices
		for index, price := range rawPrices {
			if price.PriceInfo.Expiry.GTE(sdk.NewInt(ctx.BlockHeight())) {
				rawPricesSum = rawPricesSum.Add(rawPrices[index].PriceInfo.Price)
				rawExpirySum = rawExpirySum.Add(rawPrices[index].PriceInfo.Expiry)
				notExpiredPrices, _ = notExpiredPrices.UpdatePriceOrAppendIfMissing(price)
			}
		}

		pricesLength := len(notExpiredPrices)
		var medianPrice sdk.Int
		var expiry sdk.Int
		// TODO KAVA suggestion : make threshold for acceptance (ie. require 51% of oracles to have posted valid prices)
		if pricesLength == 0 {
			// Error if there are no valid prices in the raw prices store
			return sdk.ErrInternal("no valid raw prices to calculate current prices")
		} else if pricesLength == 1 {
			// Return if there's only one price
			medianPrice = notExpiredPrices[0].PriceInfo.Price
			expiry = notExpiredPrices[0].PriceInfo.Expiry
		} else {
			pLength := sdk.NewInt(int64(pricesLength))
			medianPrice = rawPricesSum.Quo(pLength)
			expiry = rawExpirySum.Quo(pLength)
		}
		store := ctx.KVStore(keeper.StoreKey)
		currentPrice := types.CurrentPrice{
			AssetName: asset.Name,
			AssetCode: asset.Code,
			Price:     medianPrice,
			Expiry:    expiry,
		}
		store.Set([]byte(types.CurrentPricesPrefix+asset.Name+asset.Code),
			keeper.cdc.MustMarshalBinaryBare(currentPrice))

	}
	return nil
}

//GetCurrentPrices retrieves all the current prices
func (keeper Keeper) GetCurrentPrices(ctx sdk.Context) types.CurrentPrices {
	var curPrices = types.CurrentPrices{}
	store := ctx.KVStore(keeper.StoreKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.CurrentPricesPrefix))
	for ; iterator.Valid(); iterator.Next() {
		var currentPrice types.CurrentPrice
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &currentPrice)
		curPrices, _ = curPrices.AppendIfMissing(currentPrice)
	}
	return curPrices
}

//GetCurrentPrice retrieves the current price for the given token name and code
func (keeper Keeper) GetCurrentPrice(ctx sdk.Context, tokenName string, tokenCode string) (types.CurrentPrice, sdk.Error) {
	currentPrices := keeper.GetCurrentPrices(ctx)
	price, err := currentPrices.GetPrice(tokenName, tokenCode)
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
