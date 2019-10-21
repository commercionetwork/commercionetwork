package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// IsAllGTE returns false if for any denom in coinsB,
// the denom is present at a smaller amount in coins;
// else returns true.
func IsAllGTE(coins sdk.DecCoins, coinsB sdk.DecCoins) bool {
	if len(coinsB) == 0 {
		return true
	}

	if len(coins) == 0 {
		return false
	}

	for _, coinB := range coinsB {
		if coinB.Amount.GT(coins.AmountOf(coinB.Denom)) {
			return false
		}
	}

	return true
}
