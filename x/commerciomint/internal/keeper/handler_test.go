package keeper

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/commerciomint/internal/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

var testMsgOpenCdp = types.NewMsgOpenCdp(testCdp.DepositedAmount, testCdp.Owner)
var testMsgCloseCdp = types.NewMsgCloseCdp(testCdp.Owner, testCdp.Timestamp)

func TestHandler_handleMsgOpenCdp(t *testing.T) {
	ctx, bk, pfk, k := SetupTestInput()
	handler := NewHandler(k)

	// Test setup
	_, _ = bk.AddCoins(ctx, testCdp.Owner, testCdp.DepositedAmount)
	pfk.SetCurrentPrice(ctx, pricefeed.NewPrice("ucommercio", sdk.NewDec(10), sdk.NewInt(1000)))
	k.SetCreditsDenom(ctx, "uccc")

	expected := sdk.Result{Log: "Cdp opened successfully"}
	actual := handler(ctx, testMsgOpenCdp)
	require.Equal(t, expected, actual)
}

func TestHandler_handleMsgCloseCdp(t *testing.T) {
	ctx, bk, _, k := SetupTestInput()
	handler := NewHandler(k)

	_, _ = bk.AddCoins(ctx, k.supplyKeeper.GetModuleAddress(types.ModuleName), testCdp.DepositedAmount)
	_ = bk.SetCoins(ctx, testCdp.Owner, testCdp.CreditsAmount)
	k.AddCdp(ctx, testCdp)

	expected := sdk.Result{Log: "Cdp closed successfully"}
	actual := handler(ctx, testMsgCloseCdp)
	require.Equal(t, expected, actual)
}

func TestHandler_InvalidMsg(t *testing.T) {
	ctx, _, _, k := SetupTestInput()
	handler := NewHandler(k)

	invalidMsg := sdk.NewTestMsg()
	errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, invalidMsg.Type())
	expected := sdkErr.Wrap(sdkErr.ErrUnknownRequest, errMsg)
	actual := handler(ctx, invalidMsg)

	require.Equal(t, expected, actual)
}
