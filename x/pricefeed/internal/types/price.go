package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CurrentPrice struct {
	AssetName string  `json:"token_name"`
	AssetCode string  `json:"token_code"`
	Price     sdk.Int `json:"price"`
	//Block height after that price is invalid
	Expiry sdk.Int `json:"expiry"`
}

func (currentPrice CurrentPrice) Equals(cp CurrentPrice) bool {
	return currentPrice.AssetName == cp.AssetName &&
		currentPrice.AssetCode == cp.AssetCode &&
		currentPrice.Price.Equal(cp.Price) &&
		currentPrice.Expiry.Equal(cp.Expiry)
}

type CurrentPrices []CurrentPrice

func (currentPrices CurrentPrices) GetPrice(tokenName string, tokenCode string) (CurrentPrice, sdk.Error) {
	for _, ele := range currentPrices {
		if ele.AssetCode == tokenCode && ele.AssetName == tokenName {
			return ele, nil
		}
	}
	return CurrentPrice{}, sdk.ErrInternal("price not found")
}

func (currentPrices CurrentPrices) AppendIfMissing(cp CurrentPrice) CurrentPrices {
	for _, ele := range currentPrices {
		if ele.Equals(cp) {
			return currentPrices
		}
	}
	return append(currentPrices, cp)
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

func (rawPrices RawPrices) FindPrice(price RawPrice) bool {
	for _, ele := range rawPrices {
		if ele.Equals(price) {
			return true
		}
	}
	return false
}

func (rawPrices RawPrices) UpdatePriceOrAppendIfMissing(rp RawPrice) RawPrices {
	for index, ele := range rawPrices {
		if ele.Equals(rp) {
			return rawPrices
		}
		if ele.Oracle.Equals(rp.Oracle) &&
			ele.PriceInfo.AssetName == rp.PriceInfo.AssetName &&
			ele.PriceInfo.AssetCode == rp.PriceInfo.AssetCode &&
			ele.PriceInfo.Expiry.LTE(rp.PriceInfo.Expiry) &&
			ele.PriceInfo.Price != rp.PriceInfo.Price {
			rawPrices[index] = rp
			return rawPrices
		}
	}
	return append(rawPrices, rp)
}
