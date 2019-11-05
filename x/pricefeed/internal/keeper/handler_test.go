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
