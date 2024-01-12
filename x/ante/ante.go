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
	ibcante "github.com/cosmos/ibc-go/v4/modules/core/ante"

	//"github.com/cosmos/cosmos-sdk/x/auth/types"
	//bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	ibcKeeper "github.com/cosmos/ibc-go/v4/modules/core/keeper"
)

// fixedRequiredFee is the amount of fee we apply/require for each transaction processed.
var fixedRequiredFee = sdk.NewDecWithPrec(1, 2)
var storeRequiredFee = sdk.NewDecWithPrec(100, 0)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(
	ak keeper.AccountKeeper,
	bankKeeper bankKeeper.Keeper,
	govKeeper government.Keeper,
	mintKeeper commerciomintKeeper.Keeper,
	sigGasConsumer cosmosante.SignatureVerificationGasConsumer,
	signModeHandler authsigning.SignModeHandler,
	stakeDenom string,
	stableCreditsDemon string,
	feegrantKeeper cosmosante.FeegrantKeeper,
	ibcKeeper *ibcKeeper.Keeper,
	wasmConfig *wasmTypes.WasmConfig,
	txCounterStoreKey sdk.StoreKey,
) sdk.AnteHandler {
	// TODO: add check for nil
	/*
		if bankKeeper == nil {
			return nil, sdkErr.Wrap(sdkErr.ErrLogic, "bank keeper is required for AnteHandler")
		}
		if signModeHandler == nil {
			return nil, sdkErr.Wrap(sdkErr.ErrLogic, "sign mode handler is required for ante builder")
		}
		if wasmConfig == nil {
			return nil, sdkErr.Wrap(sdkErr.ErrLogic, "wasm config is required for ante builder")
		}
	*/
	return sdk.ChainAnteDecorators(
		cosmosante.NewSetUpContextDecorator(),                                    // outermost AnteDecorator. SetUpContext must be called first
		wasmkeeper.NewLimitSimulationGasDecorator(wasmConfig.SimulationGasLimit), // after setup context to enforce limits early
		wasmkeeper.NewCountTXDecorator(txCounterStoreKey),

		cosmosante.NewMempoolFeeDecorator(),
		cosmosante.NewValidateBasicDecorator(),
		cosmosante.NewValidateMemoDecorator(ak),
		cosmosante.NewConsumeGasForTxSizeDecorator(ak),
		NewMinFeeDecorator(govKeeper, mintKeeper, stakeDenom, stableCreditsDemon),
		cosmosante.NewConsumeGasForTxSizeDecorator(ak),
		cosmosante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		cosmosante.NewValidateSigCountDecorator(ak),
		cosmosante.NewDeductFeeDecorator(ak, bankKeeper, feegrantKeeper),
		cosmosante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		cosmosante.NewSigVerificationDecorator(ak, signModeHandler),
		cosmosante.NewIncrementSequenceDecorator(ak),
		ibcante.NewAnteDecorator(ibcKeeper),
	)
}

// MinFeeDecorator checks that each transaction containing a MsgShareDocument
// contains also a minimum fee amount corresponding to 0.01 euro per
// MsgShareDocument included into the transaction itself.
// The amount can be specified using stableCreditsDenom or stakeDenom.
// If stakeDenom used the cost of transaction is always 10000ucommercio
type MinFeeDecorator struct {
	govk               government.Keeper
	mintk              commerciomintKeeper.Keeper
	stakeDenom         string
	stableCreditsDenom string
}

func NewMinFeeDecorator(govKeeper government.Keeper, mintk commerciomintKeeper.Keeper, stakeDenom string, stableCreditsDenom string) MinFeeDecorator {
	return MinFeeDecorator{
		govk:               govKeeper,
		mintk:              mintk,
		stakeDenom:         stakeDenom,
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
	// get the len of messages array and find store messages to apply different fees
	lenMsgs := int64(len(stdTx.GetMsgs()))
	lenStoreMsgs := int64(0)
	for _, value := range stdTx.GetMsgs() {
		if sdk.MsgTypeURL(value) == "/cosmwasm.wasm.v1.MsgStoreCode" {
			lenMsgs--
			lenStoreMsgs++
		}
	}

	//requiredFees := fixedRequiredFee.MulInt64(int64(len(stdTx.GetMsgs())))
	requiredFees := fixedRequiredFee.MulInt64(lenMsgs)
	requiredFees = requiredFees.Add(storeRequiredFee.MulInt64(lenStoreMsgs))

	// Check the minimum fees
	if err := checkMinimumFees(stdTx, ctx, mfd.govk, mfd.mintk, mfd.stakeDenom, mfd.stableCreditsDenom, requiredFees); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

func checkMinimumFees(
	stdTx sdk.Tx,
	ctx sdk.Context,
	govk government.Keeper,
	mintk commerciomintKeeper.Keeper,
	stakeDenom string,
	stableCreditsDenom string,
	requiredFees sdk.Dec,
) error {
	fiatAmount := sdk.ZeroDec()
	// Find required quantity of stable coin = number of msg * 10000
	// Every message need 0.01 ccc
	stableRequiredQty := requiredFees.MulInt64(1000000)
	// Extract amount of stable coin from fees
	feeTx, ok := stdTx.(sdk.FeeTx)
	if !ok {
		return sdkErr.Wrap(sdkErr.ErrTxDecode, "Tx must be a FeeTx")
	}
	fiatAmount = sdk.NewDecFromInt(feeTx.GetFee().AmountOf(stableCreditsDenom))
	// Check if amount of stable coin is enough
	if !stableRequiredQty.IsZero() && stableRequiredQty.LTE(fiatAmount) {
		// If amount of stable coin is enough return without error
		return nil
	}
	// NB: if user pay insufficent fiat amount plus enough stake denom, fiat amount will be withdraw from the wallet anyway.
	//return nil
	// stakeDenom must always equal 10000
	comAmount := sdk.ZeroDec()
	comRequiredQty := requiredFees.MulInt64(1000000)
	comAmount = sdk.NewDecFromInt(feeTx.GetFee().AmountOf(stakeDenom))

	// Check if amount of stake coin is enough
	if !comRequiredQty.IsZero() && comRequiredQty.LTE(comAmount) {
		// If amount of stake coin is enough return without error
		return nil
	}

	msg := fmt.Sprintf("insufficient fees. Expected %s fiat amount, got %s, or %s stake denom amount, got %s", stableRequiredQty, fiatAmount, comRequiredQty, comAmount)
	return sdkErr.Wrap(sdkErr.ErrInsufficientFee, msg)
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
