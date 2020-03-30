package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ------------
// --- Price
// ------------

// Price represents the price of an asset that has been set from an oracle
type Price struct {
	AssetName string  `json:"asset_name"`
	Value     sdk.Dec `json:"value"`
	Expiry    sdk.Int `json:"expiry"` // Block height after which the price is to be considered invalid
}

func NewPrice(assetName string, price sdk.Dec, expiry sdk.Int) Price {
	return Price{
		AssetName: assetName,
		Value:     price,
		Expiry:    expiry,
	}
}

func EmptyPrice() Price {
	return Price{AssetName: "", Value: sdk.ZeroDec(), Expiry: sdk.ZeroInt()}
}

// Implements equatable
func (currentPrice Price) Equals(cp Price) bool {
	return currentPrice.AssetName == cp.AssetName &&
		currentPrice.Value.Equal(cp.Value) &&
		currentPrice.Expiry.Equal(cp.Expiry)
}

// implements Stringer
func (currentPrice Price) String() string {
	return string(ModuleCdc.MustMarshalJSON(currentPrice))
}

// Prices represents a slice of Price objects
type Prices []Price

// AppendIfMissing appends the given price to the prices slice, returning the new
// slice as well as a boolean telling if the appending was successful
func (prices Prices) AppendIfMissing(price Price) (Prices, bool) {
	for _, ele := range prices {
		if ele.Equals(price) {
			return nil, true
		}
	}
	return append(prices, price), false
}

// ---------------
// --- RawPrice
// ---------------

// TODO: rename to OraclePrice
// RawPrice represents a raw price
type RawPrice struct {
	Oracle  sdk.AccAddress `json:"oracle"`
	Price   Price          `json:"price"`
	Created sdk.Int        `json:"created"`
}

func (rawPrice RawPrice) Equals(rp RawPrice) bool {
	return rawPrice.Oracle.Equals(rp.Oracle) && rawPrice.Price.Equals(rp.Price) && rawPrice.Created.Equal(rp.Created)
}

// TODO: rename to OraclePrices
type RawPrices []RawPrice

func (rawPrices RawPrices) UpdatePriceOrAppendIfMissing(rp RawPrice) (RawPrices, bool) {
	for index, ele := range rawPrices {
		if ele.Equals(rp) {
			return rawPrices, false
		}
		if ele.Oracle.Equals(rp.Oracle) &&
			ele.Price.AssetName == rp.Price.AssetName &&
			ele.Price.Expiry.LT(rp.Price.Expiry) &&
			ele.Created.LT(rp.Created) {
			rawPrices[index] = rp
			return rawPrices, true
		}
	}
	return append(rawPrices, rp), true
}
