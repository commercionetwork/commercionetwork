package ante

import (
	"errors"
	"fmt"

	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"

	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// txTypesWithFixedFees contains all the messages which must have a fixed fee amount, either uccc or ucommercio.
var txTypesWithFixedFees map[string]sdk.Dec

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(
	ak keeper.AccountKeeper,
	supplyKeeper types.SupplyKeeper,
	priceKeeper pricefeed.Keeper,
	sigGasConsumer cosmosante.SignatureVerificationGasConsumer,
	stableCreditsDemon string,
	ml ...ctypes.MessageFeeBinder,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		cosmosante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		cosmosante.NewMempoolFeeDecorator(),
		cosmosante.NewValidateBasicDecorator(),
		cosmosante.NewValidateMemoDecorator(ak),
		NewMinFeeDecorator(priceKeeper, stableCreditsDemon, ml...),
		cosmosante.NewConsumeGasForTxSizeDecorator(ak),
		cosmosante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		cosmosante.NewValidateSigCountDecorator(ak),
		cosmosante.NewDeductFeeDecorator(ak, supplyKeeper),
		cosmosante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		cosmosante.NewSigVerificationDecorator(ak),
		cosmosante.NewIncrementSequenceDecorator(ak),
	)
}

// MinFeeDecorator checks that each transaction containing a MsgShareDocument
// contains also a minimum fee amount corresponding to 0.01 euro per
// MsgShareDocument included into the transaction itself.
// The amount can be specified either using stableCreditsDenom tokens or
// by using any other token which price is contained inside the pricefeedKeeper.
type MinFeeDecorator struct {
	pfk                pricefeed.Keeper
	stableCreditsDenom string
}

func NewMinFeeDecorator(priceKeeper pricefeed.Keeper, stableCreditsDenom string, ml ...ctypes.MessageFeeBinder) MinFeeDecorator {
	txTypesWithFixedFees = make(map[string]sdk.Dec)

	// cycle through all the MessageListers and append their content to txTypesWithFixedFees
	for _, mlister := range ml {
		for _, message := range mlister.Messages() {
			txTypesWithFixedFees[message.Name] = message.Fee
		}
	}

	return MinFeeDecorator{
		pfk:                priceKeeper,
		stableCreditsDenom: stableCreditsDenom,
	}
}

func (mfd MinFeeDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	// all transactions must be of type auth.StdTx
	stdTx, ok := tx.(types.StdTx)
	if !ok {
		// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
		// during runTx.
		newCtx = setGasMeter(simulate, ctx, 0)
		return newCtx, errors.New("tx must be StdTx")
	}

	// Fet the ShareDocument messages
	requiredFees := sdk.NewDec(0)
	for _, msg := range stdTx.Msgs {
		fmt.Println(msg.Type())
		// if we find msg.Type() in txTypesWithFixedFees, that message must pay our custom fixed fee amount
		if fee, msgTypeFound := txTypesWithFixedFees[msg.Type()]; msgTypeFound {
			requiredFees = requiredFees.Add(fee)
		}

	}

	// Check the minimum fees
	if err := checkMinimumFees(stdTx, ctx, mfd.pfk, mfd.stableCreditsDenom, requiredFees); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
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
	// 1 .Using stable credits worth 1€ (10.000 ustable)
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
