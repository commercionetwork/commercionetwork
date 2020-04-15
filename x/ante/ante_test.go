package ante_test

import (
	"errors"
	docsTypes "github.com/commercionetwork/commercionetwork/x/docs/types"
	pricefeedTypes "github.com/commercionetwork/commercionetwork/x/pricefeed/types"
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/commercionetwork/commercionetwork/x/ante"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
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

func TestAnteHandlerFees_MsgShareDoc(t *testing.T) {

	// Setup
	app, ctx := createTestApp(true, false)

	tokenDenom := "ucommercio"
	stableCreditsDenom := "uccc"

	anteHandler := ante.NewAnteHandler(
		app.AccountKeeper, app.SupplyKeeper, app.PriceFeedKeeper,
		app.GovernmentKeeper,
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
	app.PriceFeedKeeper.SetCurrentPrice(ctx, pricefeedTypes.NewPrice(tokenDenom, sdk.NewDec(5), sdk.NewInt(1000)))
	fees = sdk.NewCoins(sdk.NewInt64Coin(tokenDenom, 1))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{3}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdkErr.ErrInsufficientFee)

	// Signer has specified enough token fees
	app.PriceFeedKeeper.SetCurrentPrice(ctx, pricefeedTypes.NewPrice(tokenDenom, sdk.NewDec(2), sdk.NewInt(1000)))
	fees = sdk.NewCoins(sdk.NewInt64Coin(tokenDenom, 5000))
	_ = app.BankKeeper.SetCoins(ctx, addr1, fees)
	seqs = []uint64{2}
	tx = types.NewTestTx(ctx, msgs, privs, accnums, seqs, auth.NewStdFee(200000, fees))
	checkValidTx(t, anteHandler, ctx, tx, true)
}

func TestAnteHandlerFees_MsgShareDocFromTumbler(t *testing.T) {

	// Setup
	app, ctx := createTestApp(true, false)

	stableCreditsDenom := "uccc"

	anteHandler := ante.NewAnteHandler(
		app.AccountKeeper, app.SupplyKeeper, app.PriceFeedKeeper, app.GovernmentKeeper,
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
