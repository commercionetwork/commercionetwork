package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CurrentPrice struct {
	AssetName string  `json:"token_name"`
	Price     sdk.Int `json:"price"`
	//Block height after that price is invalid
	Expiry sdk.Int `json:"expiry"`
}

func (currentPrice CurrentPrice) Equals(cp CurrentPrice) bool {
	return currentPrice.AssetName == cp.AssetName &&
		currentPrice.Price.Equal(cp.Price) &&
		currentPrice.Expiry.Equal(cp.Expiry)
}

type CurrentPrices []CurrentPrice

func (currentPrices CurrentPrices) GetPrice(tokenName string) (CurrentPrice, sdk.Error) {
	for _, ele := range currentPrices {
		if ele.AssetName == tokenName {
			return ele, nil
		}
	}
	return CurrentPrice{}, sdk.ErrInternal("price not found")
}

func (currentPrices CurrentPrices) AppendIfMissing(cp CurrentPrice) (CurrentPrices, bool) {
	for _, ele := range currentPrices {
		if ele.Equals(cp) {
			return nil, true
		}
	}
	return append(currentPrices, cp), false
}

type RawPrice struct {
	Oracle    sdk.AccAddress `json:"oracle"`
	PriceInfo CurrentPrice   `json:"price"`
}

func (rawprice RawPrice) Equals(rp RawPrice) bool {
	return rawprice.Oracle.Equals(rp.Oracle) &&
		rawprice.PriceInfo.Equals(rp.PriceInfo)
}

type RawPrices []RawPrice

func (rawPrices RawPrices) UpdatePriceOrAppendIfMissing(rp RawPrice) (RawPrices, bool) {
	for index, ele := range rawPrices {
		if ele.Equals(rp) {
			return nil, true
		}
		if ele.Oracle.Equals(rp.Oracle) &&
			ele.PriceInfo.AssetName == rp.PriceInfo.AssetName &&
			ele.PriceInfo.Expiry.LTE(rp.PriceInfo.Expiry) &&
			ele.PriceInfo.Price != rp.PriceInfo.Price {
			rawPrices[index] = rp
			return rawPrices, false
		}
	}
	return append(rawPrices, rp), false
}
