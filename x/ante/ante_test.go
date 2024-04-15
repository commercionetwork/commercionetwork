package ante_test

import (
	"errors"
	"testing"
	"time"

	commerciomintTypes "github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	ptx "github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"

	cosmosante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	"github.com/commercionetwork/commercionetwork/app"
	"github.com/commercionetwork/commercionetwork/testutil/simapp"

	//sdksimapp "cosmossdk.io/simapp"
	"github.com/cosmos/cosmos-sdk/client"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/commercionetwork/commercionetwork/x/ante"
	ctypes "github.com/commercionetwork/commercionetwork/x/common/types"
	docsTypes "github.com/commercionetwork/commercionetwork/x/documents/types"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

const (
	chainID            = "commercionetwork"
	stakeDenom         = "ucommercio"
	stableCreditsDenom = "uccc"
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
	Metadata: &docsTypes.DocumentMetadata{
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
	Sender:     testSender.String(),
	Recipients: ctypes.Strings{testRecipient.String()},
}

var testDocument2 = docsTypes.Document{
	UUID:       "test-document-uuid-2",
	ContentURI: "https://example.com/document",
	Metadata: &docsTypes.DocumentMetadata{
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
	Sender:     testSender.String(),
	Recipients: ctypes.Strings{testRecipient.String()},
}

type AnteTestSuite struct {
	txBuilder client.TxBuilder
}

func SetBalances(ctx sdk.Context, bk bankKeeper.Keeper, addr sdk.AccAddress, coins sdk.Coins) {
	bk.MintCoins(ctx, authtypes.ModuleName, coins)
	bk.SendCoinsFromModuleToAccount(ctx, authtypes.ModuleName, addr, coins)
}

func TestAnteHandlerFees_MsgShareDoc(t *testing.T) {
	// Setup
	// Conversion rate is 2.0
	app, ctx := createTestApp(true, false)
	wasmConfig, _ := wasm.ReadWasmConfig(sdksimapp.EmptyAppOptions{})

	encodingConfig := sdksimapp.MakeTestEncodingConfig()

	anteHandler := ante.NewAnteHandler(
		app.AccountKeeper,
		app.BankKeeper,
		app.GovernmentKeeper,
		app.CommercioMintKeeper,
		cosmosante.DefaultSigVerificationGasConsumer,
		encodingConfig.TxConfig.SignModeHandler(),
		stakeDenom,
		stableCreditsDenom,
		app.FeeGrantKeeper,
		app.IBCKeeper,
		&wasmConfig,
		app.GetKey(wasm.StoreKey),
	)

	// Keys and addresses
	priv1, _, addr1 := testdata.KeyTestPubAddr()

	// Set the accounts
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	amountAcc1 := sdk.NewCoins(sdk.NewInt64Coin("uccc", 1000000000))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, amountAcc1)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, amountAcc1)
	//_ = app.BankKeeper.SetBalances(ctx, addr1, sdk.NewCoins(sdk.NewInt64Coin("uccc", 1000000000)))

	//_ = acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin("uccc", 1000000000)))
	app.AccountKeeper.SetAccount(ctx, acc1)

	// Msg and signatures
	msg := docsTypes.NewMsgShareDocument(docsTypes.Document{
		UUID:           testDocument.UUID,
		Metadata:       testDocument.Metadata,
		ContentURI:     testDocument.ContentURI,
		Checksum:       testDocument.Checksum,
		EncryptionData: testDocument.EncryptionData,
		Sender:         acc1.GetAddress().String(),
		Recipients:     testDocument.Recipients,
	})
	privs, accnums, seqs := []cryptotypes.PrivKey{priv1}, []uint64{0}, []uint64{0}
	msgs := []sdk.Msg{msg}
	txBuilder := encodingConfig.TxConfig.NewTxBuilder()
	as := AnteTestSuite{txBuilder}

	as.txBuilder.SetMsgs(msgs...)
	as.txBuilder.SetGasLimit(200000)

	// Signer has not specified the fees
	fees := sdk.NewCoins()
	as.txBuilder.SetFeeAmount(fees)

	err := as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)

	tx := as.txBuilder.GetTx()
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdkErr.ErrInsufficientFee)

	// Signer has not specified enough stable credits
	fees = sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 9999))
	as.txBuilder.SetFeeAmount(fees)
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)

	tx = as.txBuilder.GetTx()
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdkErr.ErrInsufficientFee)

	// Signer has specified enough stable credits
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 10000))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)
	//_ = app.BankKeeper.SetBalances(ctx, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)

	tx = as.txBuilder.GetTx()
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has not specified enough token frees
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 1))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)

	//_ = app.BankKeeper.SetBalances(ctx, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdkErr.ErrInsufficientFee)

	// Signer has specified enough token fees
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 10000))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)

	//_ = app.BankKeeper.SetBalances(ctx, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	seqs = []uint64{1}
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has not specified enough stake token fees
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 9999))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)

	//_ = app.BankKeeper.SetBalances(ctx, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdkErr.ErrInsufficientFee)

	// Signer has specified enough stake tokens fees but not enough credit tokens fees
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 10000), sdk.NewInt64Coin(stableCreditsDenom, 9999))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)

	//_ = app.BankKeeper.SetBalances(ctx, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	seqs = []uint64{2}
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has specified not enough stake tokens fees but enough credit tokens fees
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 9999), sdk.NewInt64Coin(stableCreditsDenom, 10000))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)

	//_ = app.BankKeeper.SetBalances(ctx, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	seqs = []uint64{3}
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has specified not enough both stake and credit tokens fees
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 9999), sdk.NewInt64Coin(stableCreditsDenom, 9999))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)

	//_ = app.BankKeeper.SetBalances(ctx, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdkErr.ErrInsufficientFee)

	// Test with multiple messages
	msg2 := docsTypes.NewMsgShareDocument(docsTypes.Document{
		UUID:           testDocument2.UUID,
		Metadata:       testDocument2.Metadata,
		ContentURI:     testDocument2.ContentURI,
		Checksum:       testDocument2.Checksum,
		EncryptionData: testDocument2.EncryptionData,
		Sender:         acc1.GetAddress().String(),
		Recipients:     testDocument2.Recipients,
	})
	msgs = []sdk.Msg{msg, msg2}
	as.txBuilder.SetMsgs(msgs...)

	// Signer has not specified enough stable credits
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 19999))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)

	//_ = app.BankKeeper.SetBalances(ctx, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)

	tx = as.txBuilder.GetTx()
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdkErr.ErrInsufficientFee)

	// Signer has specified enough stable credits
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 20000))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)

	//_ = app.BankKeeper.SetBalances(ctx, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	seqs = []uint64{4}
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has specified enough stake tokens
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 20000))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)

	//_ = app.BankKeeper.SetBalances(ctx, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	seqs = []uint64{5}
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has specified enough stake tokens fees but not enough credit tokens fees
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 20000), sdk.NewInt64Coin(stableCreditsDenom, 19999))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)

	//_ = app.BankKeeper.SetBalances(ctx, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	seqs = []uint64{6}
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has specified not enough stake tokens fees but enough credit tokens fees
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 19999), sdk.NewInt64Coin(stableCreditsDenom, 20000))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)

	//_ = app.BankKeeper.SetBalances(ctx, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	seqs = []uint64{7}
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has specified not enough both stake and credit tokens fees
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 19999), sdk.NewInt64Coin(stableCreditsDenom, 19999))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)

	//_ = app.BankKeeper.SetBalances(ctx, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdkErr.ErrInsufficientFee)

	// Test with wasm store messages with others
	msgstore := &wasmtypes.MsgStoreCode{
		Sender:       acc1.GetAddress().String(),
		WASMByteCode: []byte("1"),
	}
	msgs = []sdk.Msg{msg, msgstore}
	as.txBuilder.SetMsgs(msgs...)

	// Signer has not specified enough stable credits
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 100009999))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)

	tx = as.txBuilder.GetTx()
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdkErr.ErrInsufficientFee)

	// Signer has specified enough stable credits
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stableCreditsDenom, 100010000))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	seqs = []uint64{8}
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has specified enough stake tokens
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 100010000))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	seqs = []uint64{9}
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has specified enough stake tokens fees but not enough credit tokens fees
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 100010000), sdk.NewInt64Coin(stableCreditsDenom, 100009999))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	seqs = []uint64{10}
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has specified not enough stake tokens fees but enough credit tokens fees
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 100009999), sdk.NewInt64Coin(stableCreditsDenom, 100010000))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	seqs = []uint64{11}
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkValidTx(t, anteHandler, ctx, tx, true)

	// Signer has specified not enough both stake and credit tokens fees
	app.BankKeeper.SendCoinsFromAccountToModule(ctx, addr1, commerciomintTypes.ModuleName, app.BankKeeper.GetAllBalances(ctx, addr1))
	fees = sdk.NewCoins(sdk.NewInt64Coin(stakeDenom, 100009999), sdk.NewInt64Coin(stableCreditsDenom, 100009999))
	app.BankKeeper.MintCoins(ctx, commerciomintTypes.ModuleName, fees)
	app.BankKeeper.SendCoinsFromModuleToAccount(ctx, commerciomintTypes.ModuleName, addr1, fees)
	as.txBuilder.SetFeeAmount(fees)
	err = as.setupSignatures(privs, accnums, seqs)
	require.NoError(t, err)
	tx = as.txBuilder.GetTx()
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdkErr.ErrInsufficientFee)

}

func (as AnteTestSuite) setupSignatures(privs []cryptotypes.PrivKey, accnums []uint64, seqs []uint64) error {
	encodingConfig := sdksimapp.MakeTestEncodingConfig()
	var sigsV2 []signing.SignatureV2
	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  encodingConfig.TxConfig.SignModeHandler().DefaultMode(),
				Signature: nil,
			},
			Sequence: seqs[i],
		}

		sigsV2 = append(sigsV2, sigV2)
	}
	err := as.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return err
	}

	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := xauthsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: accnums[i],
			Sequence:      seqs[i],
		}

		sigV2, err := ptx.SignWithPrivKey(
			encodingConfig.TxConfig.SignModeHandler().DefaultMode(), signerData,
			as.txBuilder, priv, encodingConfig.TxConfig, seqs[i])
		if err != nil {
			return err
		}
		sigsV2 = append(sigsV2, sigV2)
	}
	err = as.txBuilder.SetSignatures(sigsV2...)
	if err != nil {
		return err
	}
	return nil
}

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool, isBlockZero bool) (*app.App, sdk.Context) {
	app := simapp.New("")
	header := tmproto.Header{ChainID: chainID}
	if !isBlockZero {
		header.Height = 1
	}

	ctx := app.BaseApp.NewContext(isCheckTx, header).
		WithConsensusParams(&abci.ConsensusParams{
			Block: &abci.BlockParams{MaxGas: 200000},
		})

	app.AccountKeeper.SetParams(ctx, authtypes.DefaultParams())

	// TODO shall we drop the following?
	app.CommercioMintKeeper.UpdateParams(ctx, validCommercioMintParams)
	// app.CommercioMintKeeper.UpdateConversionRate(ctx, sdk.NewDec(2))

	return app, ctx
}

var validConversionRate = sdk.NewDec(2)
var validFreezePeriod time.Duration = 0
var validCommercioMintParams = commerciomintTypes.Params{
	ConversionRate: validConversionRate,
	FreezePeriod:   validFreezePeriod,
}
