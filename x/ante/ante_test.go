package ante_test

import (
	"errors"
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	"github.com/commercionetwork/commercionetwork/x/ante"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	docsTypes "github.com/commercionetwork/commercionetwork/x/documents/types"
)

// run the tx through the anteHandler and ensure its valid
func checkValidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, simulate bool) {
	_, err := anteHandler(ctx, tx, simulate)
	require.Nil(t, err)
}

// run the tx through the anteHandler and ensure it fails with the given code
func checkInvalidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, simulate bool, code error) {
	_, err := anteHandler(ctx, tx, simulate)
	require.NotNil(t, err)

	require.True(t, errors.Is(sdkErr.ErrInsufficientFee, code))
}

var testSender, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var testRecipient, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")
var testDocument = docsTypes.Document{
	UUID:       "test-document-uuid",
	ContentURI: "https://example.com/document",
	Metadata: docsTypes.DocumentMetadata{
		ContentURI: "https://example.com/document/metadata",
		Schema: &docsTypes.DocumentMetadataSchema{
			URI:     "https://example.com/document/metadata/schema",
			Version: "1.0.0",
		},
	},
	Checksum: &docsTypes.DocumentChecksum{
		Value:     "93dfcaf3d923ec47edb8580667473987",
		Algorithm: "md5",
	},
	Sender:     testSender,
	Recipients: ctypes.Addresses{testRecipient},
}

var testDocument2 = docsTypes.Document{
	UUID:       "test-document-uuid-2",
	ContentURI: "https://example.com/document",
	Metadata: docsTypes.DocumentMetadata{
		ContentURI: "https://example.com/document/metadata",
		Schema: &docsTypes.DocumentMetadataSchema{
			URI:     "https://example.com/document/metadata/schema",
			Version: "1.0.0",
		},
	},
	Checksum: &docsTypes.DocumentChecksum{
		Value:     "93dfcaf3d923ec47edb8580667473987",
		Algorithm: "md5",
	},
	Sender:     testSender,
	Recipients: ctypes.Addresses{testRecipient},
}

func TestAnteHandlerFees_MsgShareDoc(t *testing.T) {

	// Setup
	// Conversion rate is 2.0
	app, ctx := createTestApp(true, false)

	tokenDenom := "ucommercio"
	stableCreditsDenom := "uccc"

	anteHandler := ante.NewAnteHandler(
		app.AccountKeeper, app.SupplyKeeper,
		app.GovernmentKeeper, app.CommercioMintKeeper,
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

	msg := docsTypes.NewMsgShareDocument(docsTypes.Document{
		UUID:           testDocument.UUID,
		Metadata:       testDocument.Metadata,
		ContentURI:     testDocument.ContentURI,
		Checksum:       testDocument.Checksum,
		EncryptionData: testDocument.EncryptionData,
		Sender:         acc1.GetAddress(),
		Recipients:     testDocument.Recipients,
	})
	privs, accnums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}
	msgs := []sdk.Msg{msg}

	// Signer has not specified the fees
	var tx sdk.Tx
	fees := sdk.NewCoins()
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdkErr.ErrInsufficientFee)

	// Signer has not specified enough stable credits
	fees = sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 9999))
	seqs = []uint64{1}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdkErr.ErrInsufficientFee)

	// Signer has specified enough stable credits
	fees = sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 10000))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{2}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has not specified enough token frees
	fees = sdk.NewCoins(sdk.NewInt64Coin(tokenDenom, 1))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{3}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdkErr.ErrInsufficientFee)

	// Signer has specified enough token fees
	fees = sdk.NewCoins(sdk.NewInt64Coin(tokenDenom, 20000))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{2}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has specified not enough token fees with stable credits and token
	fees = sdk.NewCoins(sdk.NewInt64Coin(tokenDenom, 4999), sdk.NewInt64Coin(stableCreditsDenom, 7500))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{6}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkInvalidTx(t, anteHandler, ctx, tx, true, sdkErr.ErrInsufficientFee)

	// Signer has specified enough token fees with stable credits and token
	fees = sdk.NewCoins(sdk.NewInt64Coin(tokenDenom, 5000), sdk.NewInt64Coin(stableCreditsDenom, 7500))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{6}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Test with multiple messages
	msg2 := docsTypes.NewMsgShareDocument(docsTypes.Document{
		UUID:           testDocument2.UUID,
		Metadata:       testDocument2.Metadata,
		ContentURI:     testDocument2.ContentURI,
		Checksum:       testDocument2.Checksum,
		EncryptionData: testDocument2.EncryptionData,
		Sender:         acc1.GetAddress(),
		Recipients:     testDocument2.Recipients,
	})
	msgs = []sdk.Msg{msg, msg2}

	// Signer has specified enough stable credits
	fees = sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 19999))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{7}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkInvalidTx(t, anteHandler, ctx, tx, true, sdkErr.ErrInsufficientFee)

	// Signer has specified enough stable credits
	fees = sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 20000))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{8}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has specified enough token fees
	fees = sdk.NewCoins(sdk.NewInt64Coin(tokenDenom, 40000))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{2}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has specified not enough token fees with stable credits and token
	fees = sdk.NewCoins(sdk.NewInt64Coin(tokenDenom, 9999), sdk.NewInt64Coin(stableCreditsDenom, 15000))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{6}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkInvalidTx(t, anteHandler, ctx, tx, true, sdkErr.ErrInsufficientFee)

	// Signer has specified enough token fees with stable credits and token
	fees = sdk.NewCoins(sdk.NewInt64Coin(tokenDenom, 10000), sdk.NewInt64Coin(stableCreditsDenom, 15000))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{6}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkValidTx(t, anteHandler, ctx, tx, true)

}

func TestAnteHandlerFees_MsgShareDocFromTumbler(t *testing.T) {

	// Setup
	app, ctx := createTestApp(true, false)

	stableCreditsDenom := "uccc"

	anteHandler := ante.NewAnteHandler(
		app.AccountKeeper, app.SupplyKeeper, app.GovernmentKeeper, app.CommercioMintKeeper,
		cosmosante.DefaultSigVerificationGasConsumer,
		stableCreditsDenom,
	)

	// Keys and addresses
	priv1, _, addr1 := types.KeyTestPubAddr()

	// Set the accounts
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	_ = acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("uccc", 1000000000)))
	app.AccountKeeper.SetAccount(ctx, acc1)
	require.NoError(t, app.GovernmentKeeper.SetTumblerAddress(ctx, addr1))

	// Msg and signatures

	msg := docsTypes.NewMsgShareDocument(docsTypes.Document{
		UUID:           testDocument.UUID,
		Metadata:       testDocument.Metadata,
		ContentURI:     testDocument.ContentURI,
		Checksum:       testDocument.Checksum,
		EncryptionData: testDocument.EncryptionData,
		Sender:         acc1.GetAddress(),
		Recipients:     testDocument.Recipients,
	})
	privs, accnums, seqs := []crypto.PrivKey{priv1}, []uint64{0}, []uint64{0}
	msgs := []sdk.Msg{msg}

	// Signer has not specified the fees
	var tx sdk.Tx
	fees := sdk.NewCoins()
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkValidTx(t, anteHandler, ctx, tx, false)
}
