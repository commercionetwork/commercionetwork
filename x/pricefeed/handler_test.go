package pricefeed

import (
	"fmt"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

// Test variables
var msgSetPrice = NewMsgSetPrice(TestRawPrice)
var msgAddOracle = NewMsgAddOracle(TestGovernment, TestOracle1)

func TestValidMsgSetPrice(t *testing.T) {
	_, ctx, govK, k := TestInput()

	k.AddOracle(ctx, TestOracle1)

	handler := NewHandler(k, govK)

	actual := handler(ctx, msgSetPrice)
	assert.True(t, actual.IsOK())
}

func TestValidMsgAddOracle(t *testing.T) {
	_, ctx, govK, k := TestInput()
	handler := NewHandler(k, govK)

	_ = govK.SetGovernmentAddress(ctx, TestGovernment)

	actual := handler(ctx, msgAddOracle)
	assert.True(t, actual.IsOK())
}

func TestInvalidMsg(t *testing.T) {
	invalidMsg := sdk.NewTestMsg()
	_, ctx, govK, k := TestInput()
	handler := NewHandler(k, govK)

	actual := handler(ctx, invalidMsg)

	assert.False(t, actual.IsOK())
	assert.True(t, strings.Contains(actual.Log, fmt.Sprintf("Unrecognized %s message type", ModuleName)))
}
