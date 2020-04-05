package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

const (
	cdpsForExistingPrice       string = "cdp-existing-price"
	liquidityPoolSumEqualsCdps string = "liquidity-pool-sum-equals-cdps"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, cdpsForExistingPrice,
		CdpsForExistingPrice(k))
	// ir.RegisterRoute(types.ModuleName, liquidityPoolSumEqualsCdps,
	// 	LiquidityPoolAmountEqualsCdps(k))
}

// CdpsForExistingPrice checks that each Cdp currently opened refers to an existing token priced by x/pricefeed.
func CdpsForExistingPrice(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		cdps := k.GetCdps(ctx)

		for _, cdp := range cdps {
			price, ok := k.priceFeedKeeper.GetCurrentPrice(ctx, cdp.Deposit.Denom)
			if !ok || price.Value.IsZero() {
				return sdk.FormatInvariant(
					types.ModuleName,
					cdpsForExistingPrice,
					fmt.Sprintf(
						"found cdp from owner %s which refers to a nonexistent asset %s for %s amount",
						cdp.Owner.String(),
						cdp.Deposit.Denom,
						cdp.Deposit.Amount.String(),
					),
				), true
			}

		}
		return "", false
	}
}

// LiquidityPoolAmountEqualsCdps checks that the value of all the opened cdps equals the liquidity pool amount.
func LiquidityPoolAmountEqualsCdps(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		cdps := k.GetCdps(ctx)

		var sums sdk.Coins
		for _, cdp := range cdps {
			sums.Add(cdp.Deposit)
		}

		pool := k.GetLiquidityPoolAmount(ctx)
		if pool.IsZero() && len(cdps) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName,
				cdpsForExistingPrice,
				fmt.Sprintf(
					"cdps opened and liquidity pool is empty",
				),
			), true
		}

		for _, c := range sums {
			name, sum := c.Denom, c.Amount
			for _, token := range pool {
				if token.Denom == name {
					if !sum.Equal(token.Amount) {
						return sdk.FormatInvariant(types.ModuleName, cdpsForExistingPrice, fmt.Sprintf(
							"pool amount for denom %s doesn't correspond to the sum of all the cdps opened for it, which is %s%s",
							name, sum.String(), name)), true
					}
				}
			}
		}

		return "", false
	}
}
