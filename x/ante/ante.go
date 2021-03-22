package ante

import (
	"errors"
	"fmt"

	commerciomintKeeper "github.com/commercionetwork/commercionetwork/x/commerciomint/keeper"
	government "github.com/commercionetwork/commercionetwork/x/government/keeper"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// fixedRequiredFee is the amount of fee we apply/require for each transaction processed.
var fixedRequiredFee = sdk.NewDecWithPrec(1, 2)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(
	ak keeper.AccountKeeper,
	supplyKeeper types.SupplyKeeper,
	govKeeper government.Keeper,
	mintKeeper commerciomintKeeper.Keeper,
	sigGasConsumer cosmosante.SignatureVerificationGasConsumer,
	stableCreditsDemon string,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		cosmosante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		cosmosante.NewMempoolFeeDecorator(),
		cosmosante.NewValidateBasicDecorator(),
		cosmosante.NewValidateMemoDecorator(ak),
		NewMinFeeDecorator(govKeeper, mintKeeper, stableCreditsDemon),
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
// The amount can be specified using stableCreditsDenom.
type MinFeeDecorator struct {
	govk               government.Keeper
	mintk              commerciomintKeeper.Keeper
	stableCreditsDenom string
}

func NewMinFeeDecorator(govKeeper government.Keeper, mintk commerciomintKeeper.Keeper, stableCreditsDenom string) MinFeeDecorator {
	return MinFeeDecorator{
		govk:               govKeeper,
		mintk:              mintk,
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

	// skip block with height 0, otherwise no chain initialization could happen!
	if ctx.BlockHeight() == 0 {
		return next(ctx, tx, simulate)
	}

	// calculate required fees for this transaction as (number of messages * fixed required feees)
	requiredFees := fixedRequiredFee.MulInt64(int64(len(stdTx.Msgs)))

	// Check the minimum fees
	if err := checkMinimumFees(stdTx, ctx, mfd.govk, mfd.mintk, mfd.stableCreditsDenom, requiredFees); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func checkMinimumFees(
	stdTx types.StdTx,
	ctx sdk.Context,
	govk government.Keeper,
	mintk commerciomintKeeper.Keeper,
	stableCreditsDenom string,
	requiredFees sdk.Dec,
) error {

	// No Fees for the Tumbler.
	if stdTx.FeePayer().Equals(govk.GetTumblerAddress(ctx)) {
		return nil
	}

	fiatAmount := sdk.ZeroDec()

	stableRequiredQty := requiredFees.MulInt64(1000000)
	fiatAmount = sdk.NewDecFromInt(stdTx.Fee.Amount.AmountOf(stableCreditsDenom))
	if !stableRequiredQty.IsZero() && stableRequiredQty.LTE(fiatAmount) {
		return nil
	}

	ucccConversionRate := mintk.GetConversionRate(ctx)

	if comAmount := stdTx.Fee.Amount.AmountOf("ucommercio"); comAmount.IsPositive() {
		f := comAmount.ToDec().Mul(ucccConversionRate)
		realQty := f.QuoInt64(1000000)

		fiatAmount = fiatAmount.Add(realQty)
	}

	if !fiatAmount.GTE(requiredFees) {
		msg := fmt.Sprintf("insufficient fees. Expected %s fiat amount, got %s", requiredFees, fiatAmount)
		return sdkErr.Wrap(sdkErr.ErrInsufficientFee, msg)
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
