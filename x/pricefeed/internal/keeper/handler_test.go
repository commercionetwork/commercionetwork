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

func Test_handleMsgBlacklistDenom(t *testing.T) {
	tests := []struct {
		name        string
		msg         types.MsgBlacklistDenom
		senderIsGov bool
		wantErr     bool
	}{
		{
			"sender is not gov",
			types.NewMsgBlacklistDenom(testOracle, "denom"),
			false,
			true,
		},
		{
			"sender is gov",
			types.NewMsgBlacklistDenom(testOracle, "denom"),
			true,
			false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			_, ctx, govK, k := SetupTestInput()

			if tt.senderIsGov {
				_ = govK.SetGovernmentAddress(ctx, tt.msg.Signer)
			}

			_, err := handleMsgBlacklistDenom(ctx, k, govK, tt.msg)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.Contains(t, k.DenomBlacklist(ctx), tt.msg.Denom)

		})
	}
}
