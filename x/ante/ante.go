package ante

import (
	"fmt"

	"github.com/commercionetwork/commercionetwork/x/docs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(
	ak keeper.AccountKeeper,
	supplyKeeper types.SupplyKeeper,
	sigGasConsumer cosmosante.SignatureVerificationGasConsumer,
	stableCreditsDemon string,
) sdk.AnteHandler {

	cosmosHandler := cosmosante.NewAnteHandler(ak, supplyKeeper, sigGasConsumer)

	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, abort bool) {

		newCtx, res, abort = cosmosHandler(ctx, tx, simulate)
		if abort {
			return newCtx, res, abort
		}

		// all transactions must be of type auth.StdTx
		stdTx, ok := tx.(types.StdTx)
		if !ok {
			// Set a gas meter with limit 0 as to prevent an infinite gas meter attack
			// during runTx.
			newCtx = SetGasMeter(simulate, ctx, 0)
			return newCtx, sdk.ErrInternal("tx must be StdTx").Result(), true
		}

		// get the ShareDocument messages
		shareDocMsgCount := 0
		for _, msg := range stdTx.Msgs {
			if _, ok := msg.(docs.MsgShareDocument); ok {
				shareDocMsgCount = shareDocMsgCount + 1
			}
		}

		// check the fees for the MsgShareDocument messages
		// each message should be paid with a fee of 10.000 ucredits, which is 0.01 credit
		expectedFee := sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDemon, int64(shareDocMsgCount*10000)))
		if shareDocMsgCount > 0 && stdTx.Fee.Amount.IsAllLT(expectedFee) {
			res := sdk.ErrInsufficientFee(fmt.Sprintf("expected at least %s fee", expectedFee)).Result()
			return newCtx, res, true
		}

		return newCtx, sdk.Result{GasWanted: stdTx.Fee.Gas}, false
	}
}

// SetGasMeter returns a new context with a gas meter set from a given context.
func SetGasMeter(simulate bool, ctx sdk.Context, gasLimit uint64) sdk.Context {
	// In various cases such as simulation and during the genesis block, we do not
	// meter any gas utilization.
	if simulate || ctx.BlockHeight() == 0 {
		return ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	}

	return ctx.WithGasMeter(sdk.NewGasMeter(gasLimit))
}
