package ante_test

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/ante"
	"github.com/commercionetwork/commercionetwork/x/docs"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
)

// run the tx through the anteHandler and ensure its valid
func checkValidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, simulate bool) {
	_, err := anteHandler(ctx, tx, simulate)
	require.Nil(t, err)
}

// run the tx through the anteHandler and ensure it fails with the given code
func checkInvalidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, simulate bool, code sdk.CodeType) {
	_, err := anteHandler(ctx, tx, simulate)
	require.NotNil(t, err)

	result := sdk.ResultFromError(err)

	require.Equal(t, code, result.Code, fmt.Sprintf("Expected %v, got %v", code, result))
	require.Equal(t, sdk.CodespaceRoot, result.Codespace)
}

func TestAnteHandlerFees_MsgShareDoc(t *testing.T) {

	// Setup
	app, ctx := createTestApp(true)

	tokenDenom := "ucommercio"
	stableCreditsDenom := "uccc"

	anteHandler := ante.NewAnteHandler(
		app.AccountKeeper, app.SupplyKeeper, app.PriceFeedKeeper,
		cosmosante.DefaultSigVerificationGasConsumer,
		stableCreditsDenom,
	)

	// Keys and addresses
	priv1, _, addr1 := types.KeyTestPubAddr()

	// Set the accounts
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	_ = acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("uccc", 1000000000)))
	app.AccountKeeper.SetAccount(ctx, acc1)

	// Msg and signatures

	msg := docs.NewMsgShareDocument(docs.Document{
		UUID:           docs.TestingDocument.UUID,
		Metadata:       docs.TestingDocument.Metadata,
		ContentURI:     docs.TestingDocument.ContentURI,
		Checksum:       docs.TestingDocument.Checksum,
		EncryptionData: docs.TestingDocument.EncryptionData,
		Sender:         acc1.GetAddress(),
		Recipients:     docs.TestingDocument.Recipients,
	})
	privs, accnums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}
	msgs := []sdk.Msg{msg}

	// Signer has not specified the fees
	var tx sdk.Tx
	fees := sdk.NewCoins()
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInsufficientFee)

	// Signer has not specified enough stable credits
	fees = sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 9999))
	seqs = []uint64{1}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInsufficientFee)

	// Signer has specified enough stable credits
	fees = sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 10000))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{2}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has not specified enough token frees
	app.PriceFeedKeeper.SetCurrentPrice(ctx, pricefeed.NewCurrentPrice(tokenDenom, sdk.NewDec(5), sdk.NewInt(1000)))
	fees = sdk.NewCoins(sdk.NewInt64Coin(tokenDenom, 1))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{3}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInsufficientFee)

	// Signer has specified enough token fees
	app.PriceFeedKeeper.SetCurrentPrice(ctx, pricefeed.NewCurrentPrice(tokenDenom, sdk.NewDec(2), sdk.NewInt(1000)))
	fees = sdk.NewCoins(sdk.NewInt64Coin(tokenDenom, 5000))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{2}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkValidTx(t, anteHandler, ctx, tx, true)
}
