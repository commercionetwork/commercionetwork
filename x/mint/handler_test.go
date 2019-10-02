package mint

import (
	"fmt"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/pricefeed"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/magiconair/properties/assert"
)

var TestOwner, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var TestTimestamp = "timestamp-test"
var TestDepositedAmount = sdk.NewCoins(sdk.NewCoin("ucommercio", sdk.NewInt(100)))

var TestCdpRequest = CDPRequest{
	Signer:          TestOwner,
	DepositedAmount: TestDepositedAmount,
	Timestamp:       TestTimestamp,
}
var testMsgOpenCDP = NewMsgOpenCDP(TestCdpRequest)
var testMsgCloseCDP = NewMsgCloseCDP(TestOwner, TestTimestamp)

var TestCurrentPrice = pricefeed.CurrentPrice{
	AssetName: "ucommercio",
	Price:     sdk.NewDecFromInt(sdk.NewInt(10)),
	Expiry:    sdk.NewInt(1000),
}

func TestHandler_handleMsgOpenCDP(t *testing.T) {
	cdc, ctx, bk, pfk, k := TestInput()
	handler := NewHandler(k)

	_, _ = bk.AddCoins(ctx, TestOwner, TestDepositedAmount)
	store := ctx.KVStore(pfk.StoreKey)
	store.Set([]byte("pricefeed:currentPrices:ucommercio"), cdc.MustMarshalBinaryBare(TestCurrentPrice))
	k.SetCreditsDenom(ctx, "uccc")

	expected := sdk.Result{Log: "CDP opened successfully"}
	actual := handler(ctx, testMsgOpenCDP)
	assert.Equal(t, expected, actual)
}

func TestHandler_handleMsgCloseCDP(t *testing.T) {
	_, ctx, _, _, k := TestInput()
	handler := NewHandler(k)

	expected := sdk.Result{Log: "CDP closed successfully"}
	actual := handler(ctx, testMsgCloseCDP)
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
