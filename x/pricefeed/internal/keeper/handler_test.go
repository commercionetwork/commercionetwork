package keeper

import (
	"fmt"
	"strings"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// -------------------
// --- MsgSetPrice
// -------------------
var testOracle, _ = sdk.AccAddressFromBech32("cosmos1lwmppctrr6ssnrmuyzu554dzf50apkfvd53jx0")
var price = types.Price{AssetName: "test", Value: sdk.NewDec(10), Expiry: sdk.NewInt(5000)}
var msgSetPrice = types.NewMsgSetPrice(price, testOracle)

func TestValidMsgSetPrice(t *testing.T) {
	_, ctx, govK, k := SetupTestInput()
	k.AddOracle(ctx, testOracle)

	handler := NewHandler(k, govK)

	actual := handler(ctx, msgSetPrice)
	assert.True(t, actual.IsOK())
}

// ---------------------
// --- MsgAddOracle
// ---------------------
var testGovernment, _ = sdk.AccAddressFromBech32("cosmos1tupew4x3rhh0lpqha9wvzmzxjr4e37mfy3qefm")
var msgAddOracle = types.NewMsgAddOracle(testGovernment, testOracle)

func TestValidMsgAddOracle(t *testing.T) {
	_, ctx, govK, k := SetupTestInput()
	handler := NewHandler(k, govK)

	_ = govK.SetGovernmentAddress(ctx, testGovernment)

	actual := handler(ctx, msgAddOracle)
	assert.True(t, actual.IsOK())
}

func TestInvalidMsg(t *testing.T) {
	invalidMsg := sdk.NewTestMsg()
	_, ctx, govK, k := SetupTestInput()
	handler := NewHandler(k, govK)

	actual := handler(ctx, invalidMsg)

	assert.False(t, actual.IsOK())
	assert.True(t, strings.Contains(actual.Log, fmt.Sprintf("Unrecognized %s message type", types.ModuleName)))
}
