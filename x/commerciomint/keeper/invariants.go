package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
)

const (
	validatePositions               string = "validate-positions"
	positionsForExistingPrice       string = "position-existing-price"
	liquidityPoolSumEqualsPositions string = "liquidity-pool-sum-equals-positions"
)

func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	//ir.RegisterRoute(types.ModuleName, validatePositions, ValidateAllPositions(k))
	//ir.RegisterRoute(types.ModuleName, positionsForExistingPrice,
	//	PositionsForExistingPrice(k))
	//ir.RegisterRoute(types.ModuleName, liquidityPoolSumEqualsPositions,
	//	LiquidityPoolAmountEqualsPositions(k))
}

// ValidateAllPositions ensures that all Positions are correct.
func ValidateAllPositions(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		positions := k.GetAllPositions(ctx)
		if len(positions) == 0 {
			return "", false
		}
		for _, position := range positions {
			if err := position.Validate(); err != nil {
				return sdk.FormatInvariant(types.ModuleName, validatePositions,
					fmt.Sprintf("found inconsistent position %+v: %v", position, err)), true
			}
		}
		return "", false
	}
}

// PositionsForExistingPrice checks that each Position currently opened refers to an existing token priced by x/pricefeed.
/*func PositionsForExistingPrice(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		positions := k.GetAllPositions(ctx)

		for _, position := range positions {
			for _, deposit := range position.Collateral {
				price, ok := k.priceFeedKeeper.GetCurrentPrice(ctx, deposit.Denom)
				if !ok || price.Value.IsZero() {
					return sdk.FormatInvariant(
						types.ModuleName,
						positionsForExistingPrice,
						fmt.Sprintf(
							"found position from owner %s which refers to a nonexistent asset %s for %s amount",
							position.Owner.String(),
							deposit.Denom,
							deposit.Amount.String(),
						),
					), true
				}
			}

		}
		return "", false
	}
}*/

// LiquidityPoolAmountEqualsPositions checks that the value of all the opened positions equals the liquidity pool amount.
/*func LiquidityPoolAmountEqualsPositions(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		positions := k.GetAllPositions(ctx)

		var sums sdk.Coins
		for _, position := range positions {
			sums.Add(position.Collateral...)
		}

		pool := k.GetLiquidityPoolAmount(ctx)
		if pool.IsZero() && len(positions) > 0 {
			return sdk.FormatInvariant(
				types.ModuleName, positionsForExistingPrice, "positions opened and liquidity pool is empty",
			), true
		}

		for _, c := range sums {
			name, sum := c.Denom, c.Amount
			for _, token := range pool {
				if token.Denom == name {
					if !sum.Equal(token.Amount) {
						return sdk.FormatInvariant(types.ModuleName, positionsForExistingPrice, fmt.Sprintf(
							"pool amount for denom %s doesn't correspond to the sum of all the positions opened for it, which is %s%s",
							name, sum.String(), name)), true
					}
				}
			}
		}

		return "", false
	}
}*/
