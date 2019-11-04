package keeper

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/mint/internal/types"
	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/magiconair/properties/assert"
)

var testMsgOpenCdp = types.NewMsgOpenCdp(TestCdpRequest)
var testMsgCloseCdp = types.NewMsgCloseCdp(TestCdp.Owner, TestCdp.Timestamp)

func TestHandler_handleMsgOpenCdp(t *testing.T) {
	_, ctx, bk, pfk, k := SetupTestInput()
	handler := NewHandler(k)

	// Test setup
	_, _ = bk.AddCoins(ctx, TestCdp.Owner, TestCdp.DepositedAmount)
	pfk.SetCurrentPrice(ctx, pricefeed.NewCurrentPrice("ucommercio", sdk.NewDec(10), sdk.NewInt(1000)))
	k.SetCreditsDenom(ctx, "uccc")

	expected := sdk.Result{Log: "Cdp opened successfully"}
	actual := handler(ctx, testMsgOpenCdp)
	assert.Equal(t, expected, actual)
}

func TestHandler_handleMsgCloseCdp(t *testing.T) {
	_, ctx, bk, _, k := SetupTestInput()
	handler := NewHandler(k)

	_, _ = bk.AddCoins(ctx, k.supplyKeeper.GetModuleAddress(types.ModuleName), TestCdp.DepositedAmount)
	_ = bk.SetCoins(ctx, TestCdp.Owner, TestCdp.CreditsAmount)
	k.AddCdp(ctx, TestCdp)

	expected := sdk.Result{Log: "Cdp closed successfully"}
	actual := handler(ctx, testMsgCloseCdp)
	assert.Equal(t, expected, actual)
}

func TestHandler_InvalidMsg(t *testing.T) {
	_, ctx, _, _, k := SetupTestInput()
	handler := NewHandler(k)

	invalidMsg := sdk.NewTestMsg()
	errMsg := fmt.Sprintf("Unrecognized %s message type: %v", types.ModuleName, invalidMsg.Type())
	expected := sdk.ErrUnknownRequest(errMsg).Result()
	actual := handler(ctx, invalidMsg)

	assert.Equal(t, expected, actual)
}
