package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IsAllGTE returns false if for any denom in otherCoins,
// the denom is present at a smaller amount in coins;
// else returns true.
func IsAllGTE(coins sdk.DecCoins, otherCoins sdk.DecCoins) bool {
	if len(otherCoins) == 0 {
		return true
	}

	if len(coins) == 0 {
		return false
	}

	for _, otherCoin := range otherCoins {
		if otherCoin.Amount.GT(coins.AmountOf(otherCoin.Denom)) {
			return false
		}
	}

	return true
}
