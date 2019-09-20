package types

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CurrentPrice struct {
	TokenName string  `json:"token_name"`
	TokenCode string  `json:"token_code"`
	Price     sdk.Int `json:"price"`
	//Block height after that price is invalid
	Expiry sdk.Int `json:"expiry"`
}

func (currentPrice CurrentPrice) Equals(cp CurrentPrice) bool {
	return currentPrice.TokenName == cp.TokenName &&
		currentPrice.TokenCode == cp.TokenCode &&
		currentPrice.Price.Equal(cp.Price)
}

type CurrentPrices []CurrentPrice

func (currentPrices CurrentPrices) FindPrice(tokenName string, tokenCode string) (CurrentPrice, sdk.Error) {
	for _, ele := range currentPrices {
		if ele.TokenCode == tokenCode && ele.TokenName == tokenName {
			return ele, nil
		}
	}
	return CurrentPrice{}, sdk.ErrInternal("price not found")
}

type RawPrice struct {
	Oracle       sdk.AccAddress `json:"oracle"`
	CurrentPrice CurrentPrice   `json:"price"`
}

func (rawprice RawPrice) Equals(rp RawPrice) bool {
	return rawprice.Oracle.Equals(rp.Oracle) &&
		rawprice.CurrentPrice.Equals(rp.CurrentPrice)
}

type RawPrices []RawPrice

func (rawPrices RawPrices) FindPrice(tokenName string, tokenCode string) (RawPrice, error) {
	for _, ele := range rawPrices {
		if ele.CurrentPrice.TokenCode == tokenCode && ele.CurrentPrice.TokenName == tokenName {
			return ele, nil
		}
	}
	return RawPrice{}, errors.New("price not found")
}

func (rawPrices RawPrices) UpdatePriceOrAppendIfMissing(rp RawPrice) RawPrices {
	index := 0
	for _, ele := range rawPrices {
		if ele.Equals(rp) {
			return rawPrices
		}
		if ele.Oracle.Equals(rp.Oracle) &&
			ele.CurrentPrice.TokenName == rp.CurrentPrice.TokenName &&
			ele.CurrentPrice.TokenCode == rp.CurrentPrice.TokenCode &&
			ele.CurrentPrice.Expiry.LTE(rp.CurrentPrice.Expiry) &&
			ele.CurrentPrice.Price != rp.CurrentPrice.Price {
			rawPrices[index] = rp
			return rawPrices
		}
		index++
	}
	return append(rawPrices, rp)
}
