package mint

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/magiconair/properties/assert"
)

var testMsgOpenCdp = NewMsgOpenCdp(TestCdpRequest)
var testMsgCloseCdp = NewMsgCloseCdp(TestCdp.Owner, TestCdp.Timestamp)

func TestHandler_handleMsgOpenCdp(t *testing.T) {
	_, ctx, bk, pfk, k := TestInput()
	handler := NewHandler(k)

	// Test setup
	_, _ = bk.AddCoins(ctx, TestCdp.Owner, TestCdp.DepositedAmount)
	pfk.SetCurrentPrice(ctx, pricefeed.NewCurrentPrice("ucommercio", 10, 1000))
	k.SetCreditsDenom(ctx, "uccc")

	expected := sdk.Result{Log: "Cdp opened successfully"}
	actual := handler(ctx, testMsgOpenCdp)
	assert.Equal(t, expected, actual)
}

func TestHandler_handleMsgCloseCdp(t *testing.T) {
	_, ctx, bk, _, k := TestInput()
	handler := NewHandler(k)

	k.SetLiquidityPool(ctx, TestCdp.DepositedAmount)
	_ = bk.SetCoins(ctx, TestCdp.Owner, TestCdp.CreditsAmount)
	k.AddCdp(ctx, TestCdp)

	expected := sdk.Result{Log: "Cdp closed successfully"}
	actual := handler(ctx, testMsgCloseCdp)
	assert.Equal(t, expected, actual)
}

func TestHandler_InvalidMsg(t *testing.T) {
	_, ctx, _, _, k := TestInput()
	handler := NewHandler(k)

	invalidMsg := sdk.NewTestMsg()
	errMsg := fmt.Sprintf("Unrecognized %s message type: %v", ModuleName, invalidMsg.Type())
	expected := sdk.ErrUnknownRequest(errMsg).Result()
	actual := handler(ctx, invalidMsg)

	assert.Equal(t, expected, actual)
}
