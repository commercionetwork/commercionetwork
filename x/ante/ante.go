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
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	//bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	//bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

// fixedRequiredFee is the amount of fee we apply/require for each transaction processed.
var fixedRequiredFee = sdk.NewDecWithPrec(1, 2)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(
	ak keeper.AccountKeeper, bankKeeper types.BankKeeper,
	govKeeper government.Keeper,
	mintKeeper commerciomintKeeper.Keeper,
	sigGasConsumer cosmosante.SignatureVerificationGasConsumer,
	signModeHandler authsigning.SignModeHandler,
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
		cosmosante.NewDeductFeeDecorator(ak, bankKeeper),
		cosmosante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		cosmosante.NewSigVerificationDecorator(ak, signModeHandler),
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
	//stdTx, ok := tx.(types.StdTx)
	stdTx, ok := tx.(authsigning.SigVerifiableTx)
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

	requiredFees := fixedRequiredFee.MulInt64(int64(len(stdTx.GetMsgs())))

	// Check the minimum fees
	if err := checkMinimumFees(stdTx, ctx, mfd.govk, mfd.mintk, mfd.stableCreditsDenom, requiredFees); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func checkMinimumFees(
	stdTx sdk.Tx,
	ctx sdk.Context,
	govk government.Keeper,
	mintk commerciomintKeeper.Keeper,
	stableCreditsDenom string,
	requiredFees sdk.Dec,
) error {

	// No Fees for the Tumbler.

	fiatAmount := sdk.ZeroDec()
	// Find required quantity of stable coin = number of msg * 10000
	// Every message need 0.01 ccc
	stableRequiredQty := requiredFees.MulInt64(1000000)
	// Extract amount of stable coin from fees
	feeTx, ok := stdTx.(sdk.FeeTx)
	if !ok {
		return ctx, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}
	fiatAmount = sdk.NewDecFromInt(feeTx.GetFee().AmountOf(stableCreditsDenom))
	// Check if amount of stable coin is enough
	if !stableRequiredQty.IsZero() && stableRequiredQty.LTE(fiatAmount) {
		// If amount of stable coin is enough return without error
		return nil
	}

	// Retrive stable coin conversion rate
	ucccConversionRate := mintk.GetConversionRate(ctx)
	// Retrive amount of commercio token and calculate equivalent in stable coin
	if comAmount := stdTx.Fee.Amount.AmountOf("ucommercio"); comAmount.IsPositive() {
		//f := comAmount.ToDec().Mul(ucccConversionRate)
		f := comAmount.ToDec().Quo(ucccConversionRate)
		//realQty := f.QuoInt64(1000000)
		//fiatAmount = fiatAmount.Add(realQty)
		fiatAmount = fiatAmount.Add(f)

	}

	// Check if amount of stable coin plus equivalent amount of commercio token are enough
	if !stableRequiredQty.LTE(fiatAmount) {
		msg := fmt.Sprintf("insufficient fees. Expected %s fiat amount, got %s", stableRequiredQty, fiatAmount)
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
