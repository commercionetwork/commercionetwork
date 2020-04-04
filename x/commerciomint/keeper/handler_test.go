package keeper

import (
	"fmt"
	"testing"

	sdkErr "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
)

var testMsgOpenCdp = types.NewMsgOpenCdp(testCdp.Owner, testCdp.Deposit)
var testMsgCloseCdp = types.NewMsgCloseCdp(testCdp.Owner, testCdp.CreatedAt)

func TestHandler_handleMsgOpenCdp(t *testing.T) {
	ctx, bk, pfk, _, k := SetupTestInput()
	handler := NewHandler(k)

	// Test setup
	_, _ = bk.AddCoins(ctx, testCdp.Owner, sdk.NewCoins(testCdp.Deposit))
	balance := bk.GetCoins(ctx, testCdpOwner)

	// Check balance
	require.Equal(t, "100ucommercio", balance.String())

	// Set credits denom and push a price to pricefeed
	k.SetCreditsDenom(ctx, "uccc")
	pfk.SetCurrentPrice(ctx, pricefeed.NewPrice("ucommercio", sdk.NewDec(10), sdk.NewInt(1000)))

	actual, err := handler(ctx, testMsgOpenCdp)
	require.NoError(t, err)
	require.Equal(t, &sdk.Result{Log: "Cdp opened successfully"}, actual)

	// Check final balance
	balance = bk.GetCoins(ctx, testCdpOwner)
	require.Equal(t, "500uccc", balance.String())
}

func TestHandler_handleMsgCloseCdp(t *testing.T) {
	ctx, bk, _, _, k := SetupTestInput()
	handler := NewHandler(k)

	_, _ = bk.AddCoins(ctx, k.supplyKeeper.GetModuleAddress(types.ModuleName), sdk.NewCoins(testCdp.Deposit))
	_ = bk.SetCoins(ctx, testCdp.Owner, testCdp.Credits)
	require.Equal(t, 0, len(k.GetCdps(ctx)))
	k.SetCdp(ctx, testCdp)
	require.Equal(t, 1, len(k.GetCdps(ctx)))

	expected := &sdk.Result{Log: "Cdp closed successfully"}
	actual, err := handler(ctx, testMsgCloseCdp)
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestHandler_handleMsgSetCdpCollateralRate(t *testing.T) {
	ctx, _, _, gk, k := SetupTestInput()
	govAddr := []byte("governance")
	gk.SetGovernmentAddress(ctx, govAddr)
	handler := NewHandler(k)

	msg := types.NewMsgSetCdpCollateralRate(govAddr, sdk.NewDec(3))

	expected := &sdk.Result{Log: "Cdp collateral rate changed successfully to 3.000000000000000000"}
	actual, err := handler(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	msg = types.NewMsgSetCdpCollateralRate(govAddr, sdk.NewDec(0))
	_, err = handler(ctx, msg)
	require.Error(t, err)
	require.Equal(t, "invalid request: cdp collateral rate must be positive: 0.000000000000000000", err.Error())

	msg = types.NewMsgSetCdpCollateralRate([]byte("invalidAddr"), sdk.NewDec(3))
	_, err = handler(ctx, msg)
	require.Error(t, err)
	require.Equal(t, "unauthorized: cosmos1d9h8vctvd9jyzerywgt84wdv cannot set collateral rate", err.Error())
}

func TestHandler_InvalidMsg(t *testing.T) {
	ctx, _, _, _, k := SetupTestInput()
	handler := NewHandler(k)

	invalidMsg := sdk.NewTestMsg()
	errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, invalidMsg.Type())
	expected := sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)

	_, err := handler(ctx, invalidMsg)
	require.Error(t, err)
	require.Equal(t, expected.Error(), err.Error())
}
