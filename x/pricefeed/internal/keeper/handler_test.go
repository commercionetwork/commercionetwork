package keeper

import (
	"fmt"
	"strings"
	"testing"

	"github.com/commercionetwork/commercionetwork/x/pricefeed/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// -------------------
// --- MsgSetPrice
// -------------------

func TestValidMsgSetPrice(t *testing.T) {
	_, ctx, govK, k := SetupTestInput()
	k.AddOracle(ctx, testOracle)

	handler := NewHandler(k, govK)

	_, err := handler(ctx, msgSetPrice)
	require.NoError(t, err)
}

// ---------------------
// --- MsgAddOracle
// ---------------------

func TestValidMsgAddOracle(t *testing.T) {
	_, ctx, govK, k := SetupTestInput()
	handler := NewHandler(k, govK)

	_ = govK.SetGovernmentAddress(ctx, testGovernment)

	_, err := handler(ctx, msgAddOracle)
	require.NoError(t, err)
}

func TestInvalidMsg(t *testing.T) {
	invalidMsg := sdk.NewTestMsg()
	_, ctx, govK, k := SetupTestInput()
	handler := NewHandler(k, govK)

	_, err := handler(ctx, invalidMsg)

	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), fmt.Sprintf("Unrecognized %s message type", types.ModuleName)))
}
