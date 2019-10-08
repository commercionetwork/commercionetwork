package ante_test

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/ante"
	"github.com/commercionetwork/commercionetwork/x/docs"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// run the tx through the anteHandler and ensure its valid
func checkValidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, simulate bool) {
	_, result, abort := anteHandler(ctx, tx, simulate)
	require.Equal(t, "", result.Log)
	require.False(t, abort)
	require.Equal(t, sdk.CodeOK, result.Code)
	require.True(t, result.IsOK())
}

// run the tx through the anteHandler and ensure it fails with the given code
func checkInvalidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, simulate bool, code sdk.CodeType) {
	newCtx, result, abort := anteHandler(ctx, tx, simulate)
	require.True(t, abort)

	require.Equal(t, code, result.Code, fmt.Sprintf("Expected %v, got %v", code, result))
	require.Equal(t, sdk.CodespaceRoot, result.Codespace)

	if code == sdk.CodeOutOfGas {
		stdTx, ok := tx.(types.StdTx)
		require.True(t, ok, "tx must be in form auth.types.StdTx")
		// GasWanted set correctly
		require.Equal(t, stdTx.Fee.Gas, result.GasWanted, "Gas wanted not set correctly")
		require.True(t, result.GasUsed > result.GasWanted, "GasUsed not greated than GasWanted")
		// Check that context is set correctly
		require.Equal(t, result.GasUsed, newCtx.GasMeter().GasConsumed(), "Context not updated correctly")
	}
}

// defaultSigVerificationGasConsumer is the default implementation of SignatureVerificationGasConsumer. It consumes gas
// for signature verification based upon the public key type. The cost is fetched from the given params and is matched
// by the concrete type.
func defaultSigVerificationGasConsumer(meter sdk.GasMeter, sig []byte, pubkey crypto.PubKey, params types.Params) sdk.Result {
	switch pubkey := pubkey.(type) {
	case ed25519.PubKeyEd25519:
		meter.ConsumeGas(params.SigVerifyCostED25519, "ante verify: ed25519")
		return sdk.ErrInvalidPubKey("ED25519 public keys are unsupported").Result()

	case secp256k1.PubKeySecp256k1:
		meter.ConsumeGas(params.SigVerifyCostSecp256k1, "ante verify: secp256k1")
		return sdk.Result{}

	case multisig.PubKeyMultisigThreshold:
		var multisignature multisig.Multisignature
		codec.Cdc.MustUnmarshalBinaryBare(sig, &multisignature)

		consumeMultisignatureVerificationGas(meter, multisignature, pubkey, params)
		return sdk.Result{}

	default:
		return sdk.ErrInvalidPubKey(fmt.Sprintf("unrecognized public key type: %T", pubkey)).Result()
	}
}

func consumeMultisignatureVerificationGas(meter sdk.GasMeter,
	sig multisig.Multisignature, pubkey multisig.PubKeyMultisigThreshold,
	params types.Params) {

	size := sig.BitArray.Size()
	sigIndex := 0
	for i := 0; i < size; i++ {
		if sig.BitArray.GetIndex(i) {
			defaultSigVerificationGasConsumer(meter, sig.Sigs[sigIndex], pubkey.PubKeys[i], params)
			sigIndex++
		}
	}
}

func TestAnteHandlerFees_MsgShareDoc(t *testing.T) {
	// setup
	app, ctx, pfk := createTestApp(true)
	anteHandler := ante.NewAnteHandler(
		app.AccountKeeper, app.SupplyKeeper, pfk,
		defaultSigVerificationGasConsumer,
		"uccc",
	)

	// keys and addresses
	priv1, _, addr1 := types.KeyTestPubAddr()

	// set the accounts
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	_ = acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("uccc", 1000000000)))
	app.AccountKeeper.SetAccount(ctx, acc1)

	// msg and signatures
	document := docs.Document{
		Uuid: "asdasdasd",
		Metadata: docs.DocumentMetadata{
			ContentUri: "asdasd",
			SchemaType: "asdasd",
			Schema:     nil,
		},
		ContentUri: "asdasd",
		Checksum:   nil,
	}
	msg := docs.NewMsgShareDocument(acc1.GetAddress(), []sdk.AccAddress{acc1.GetAddress()}, document)
	privs, accnums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}
	msgs := []sdk.Msg{msg}

	// signer has not specified enough fees
	var tx sdk.Tx
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(50000, sdk.NewCoins()))
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInsufficientFee)

	// signer has specified enough fees
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(50000, sdk.NewCoins(sdk.NewInt64Coin("uccc", 10000))))
	checkValidTx(t, anteHandler, ctx, tx, true)
}
