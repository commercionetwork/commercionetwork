package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

func (currentPrice Price) Equals(cp Price) bool {
	return currentPrice.AssetName == cp.AssetName &&
		currentPrice.Value.Equal(cp.Value) &&
		currentPrice.Expiry.Equal(cp.Expiry)
}

type Prices []Price

func (prices Prices) GetPrice(tokenName string) (Price, sdk.Error) {
	for _, ele := range prices {
		if ele.AssetName == tokenName {
			return ele, nil
		}
	}
	return Price{}, sdk.ErrInternal("price not found")
}

func (prices Prices) AppendIfMissing(cp Price) (Prices, bool) {
	for _, ele := range prices {
		if ele.Equals(cp) {
			return nil, true
		}
	}
	return append(prices, cp), false
}

func (prices Prices) UpdatePriceOrAppendIfMissing(rp Price) (Prices, bool) {
	for index, ele := range prices {
		if ele.Equals(rp) {
			return prices, false
		}
		if ele.AssetName == rp.AssetName && ele.Expiry.LTE(rp.Expiry) && ele.Value != rp.Value {
			prices[index] = rp
			return prices, true
		}
	}
	return append(prices, rp), true
}
