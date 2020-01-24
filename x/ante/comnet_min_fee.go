package ante

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/docs"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

var messagesCosts = map[string]sdk.Dec{
	docs.MsgTypeShareDocument: sdk.NewDecWithPrec(1, 2),
}

// comnetMinFeesChecker checks that each StdTx has a precise fees amount, expressed in uccc tokens.
func comnetMinFeesChecker(stdTx types.StdTx, ctx sdk.Context, pfk pricefeed.Keeper, stableCreditsDenom string) sdk.Result {
	// Fet the ShareDocument messages
	requiredFees := sdk.NewDec(0)
	for _, msg := range stdTx.Msgs {
		if value, ok := messagesCosts[msg.Type()]; ok {
			requiredFees = requiredFees.Add(value)
		}
	}

	// Check the minimum fees
	if err := checkMinimumFees(stdTx, ctx, pfk, stableCreditsDenom, requiredFees); err != nil {
		return err.Result()
	}

	return sdk.Result{}
}

func checkMinimumFees(
	stdTx types.StdTx,
	ctx sdk.Context,
	pfk pricefeed.Keeper,
	stableCreditsDenom string,
	requiredFees sdk.Dec,
) sdk.Error {

	// ----
	// Each message should cost 0.01€, which can be paid:
	// 1 .Using stable credits worth 1€ (10.000 uccc)
	// 2. Using other tokens (their required quantity is based on their value)
	// ----

	// ----
	// Try using stable credits
	// ----

	// Token quantity is always set as millionth of units
	stableRequiredQty := requiredFees.MulInt64(1000000)
	stableFeeAmount := sdk.NewDecFromInt(stdTx.Fee.Amount.AmountOf(stableCreditsDenom))
	if !stableRequiredQty.IsZero() && stableRequiredQty.LTE(stableFeeAmount) {
		return nil
	}

	// ----
	// Stable credits where not sufficient, fall back to normal ones
	// ----

	fiatAmount := sdk.ZeroDec()
	for _, fee := range stdTx.Fee.Amount {

		// Skip stable credits
		if fee.Denom == stableCreditsDenom {
			continue
		}

		// Search for the token price
		if ctPrice, found := pfk.GetCurrentPrice(ctx, fee.Denom); found {
			// The quantity is always set as millionth of unit
			realQty := fee.Amount.ToDec().QuoInt64(1000000)

			// Fiat amount = price * quantity
			tokenFiatAmount := ctPrice.Value.Mul(realQty)

			// Add up everything
			fiatAmount = fiatAmount.Add(tokenFiatAmount)
		}
	}

	if !fiatAmount.GTE(requiredFees) {
		msg := fmt.Sprintf("Insufficient fees. Expected %s fiat amount, got %s", requiredFees, fiatAmount)
		return sdk.ErrInsufficientFee(msg)
	}

	return nil
}

// setGasMeter returns a new context with a gas meter set from a given context.
func setGasMeter(simulate bool, ctx sdk.Context, gasLimit uint64) sdk.Context {
	// In various cases such as simulation and during the genesis block, we do not
	// meter any gas utilization.
	if simulate || ctx.BlockHeight() == 0 {
		return ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	}

	return ctx.WithGasMeter(sdk.NewGasMeter(gasLimit))
}
